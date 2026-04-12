## Context

The simulator's LLM layer is currently a single concrete struct (`internal/llm.Client`) that speaks directly to Ollama's `/api/chat` endpoint. Every callsite — character actors, the game director, the summary generator — imports `*llm.Client` by concrete type. There is no seam to swap backends at runtime.

OpenRouter exposes an OpenAI-compatible `/chat/completions` endpoint with token-usage fields in every response, making cost tracking straightforward. The challenge is introducing the provider abstraction without breaking the existing Ollama path and without requiring a large refactor of all callsites.

## Goals / Non-Goals

**Goals:**
- Define a `Provider` interface in `internal/llm` that both Ollama and OpenRouter implement.
- Add an `OpenRouterClient` in `internal/llm` that calls the OpenRouter API and returns usage data.
- Provide a `NewProvider(cfg ProviderConfig) Provider` factory that reads `LLM_PROVIDER` from env and returns the right implementation.
- Accumulate token usage across all LLM calls in a `CostAccumulator` and render a cost summary at the end.
- All configuration driven by `.env` / env vars; zero code changes needed to switch providers.
- Backward-compatible: existing Ollama `.env` setups continue to work unchanged.

**Non-Goals:**
- Streaming responses.
- Supporting more than two providers in this change (interface design makes adding future providers trivial).
- Per-model pricing database — only OpenRouter's response-level `usage` object is used.
- Retry logic or circuit-breaker patterns.

## Decisions

### 1. Interface shape: thin `Provider` over existing method set

**Decision**: Define `Provider` with exactly `Generate(systemPrompt, userPrompt string, opts ...Option) (string, Usage, error)` and `Chat(messages []Message, opts ...Option) (string, Usage, error)`, plus `Ping() error`. Add a `Usage` value type (`PromptTokens`, `CompletionTokens`, `EstimatedCostUSD float64`).

**Why**: All callers today use only `Generate` or `Chat`. Adding a `Usage` return rather than a side-channel (context, callback) keeps callers explicit and avoids global state. The `Option` functional-options pattern is already established.

**Alternative considered**: Context-based side-channel (`context.WithValue`). Rejected: harder to test and invisible to callers.

### 2. Ollama adapter: wrap existing `Client`, return zero `Usage`

**Decision**: Create an `OllamaProvider` thin wrapper around the existing `*Client` that returns `Usage{}` (all zeros) from every call.

**Why**: Ollama's local API does not expose token counts in its response. Returning zero usage is honest and avoids fake estimates. The cost summary section is simply omitted when total usage is zero.

**Alternative considered**: Estimate tokens via a local tokenizer. Rejected: adds complexity and a new dependency for marginal value on a local, free provider.

### 3. OpenRouter client: OpenAI-compatible endpoint

**Decision**: Implement `OpenRouterClient` against `https://openrouter.ai/api/v1/chat/completions` using the standard OpenAI request/response JSON shapes (already familiar from the OpenRouter docs). Parse `usage.prompt_tokens`, `usage.completion_tokens`, and `usage.cost` (OpenRouter-specific field) from the response.

**Why**: OpenRouter documents an OpenAI-compatible API. Using the same JSON structs avoids custom serialization and makes it easy to verify behavior against the OpenAI spec.

**Alternative considered**: Use an OpenAI Go SDK. Rejected: adds a dependency; the structs needed are trivial to define inline.

### 4. Cost accumulation: value type passed through simulation engine

**Decision**: Introduce `llm.CostAccumulator` (a simple struct with a mutex and running totals) created once in `main`, passed into the engine and summary generator. Each `Generate`/`Chat` call returns `Usage` which callers add to the accumulator via `accumulator.Add(usage)`.

**Why**: Keeps cost tracking explicit and testable. No global state. The engine and director already receive the provider as a parameter, so adding the accumulator alongside it is minimal churn.

**Alternative considered**: Attach accumulator to the `Provider` interface. Rejected: mixes transport concerns with accounting; also makes the interface harder to implement for simple providers.

### 5. Provider selection via `LLM_PROVIDER` env var

**Decision**: `LLM_PROVIDER=ollama` (default) or `LLM_PROVIDER=openrouter`. Factory `NewProvider` reads this plus provider-specific vars and returns the implementation.

**Why**: Single toggle, clear intent, backward-compatible default.

### 6. Cost report in summary: appended section, only when non-zero

**Decision**: `summary.GenerateSummary` accepts an optional `*llm.CostAccumulator`. If totals are non-zero, a `## Cost Report` section is appended after the character cards, showing prompt tokens, completion tokens, and estimated USD cost.

**Why**: Non-intrusive — Ollama users see no change. OpenRouter users get actionable spend data.

## Risks / Trade-offs

- **OpenRouter pricing changes** → The `usage.cost` field is provided by OpenRouter per-call; we consume it directly rather than computing it locally, so price changes on their side are reflected automatically. If they remove the field, cost shows as $0.00 with a note.
- **Interface breaks existing tests** → `Generate` and `Chat` gain a `Usage` return value. All existing callers need to accept the extra return. Since there are no test files yet, this is purely a refactor cost.
- **Ollama Ping vs OpenRouter Ping** → OpenRouter has no equivalent of `/api/tags`. The `Ping()` for OpenRouter will make a lightweight `GET /models` call with the API key to validate credentials at startup.

## Migration Plan

1. Add `Provider` interface and `Usage` type to `internal/llm/provider.go`.
2. Create `OllamaProvider` adapter wrapping existing `*Client`.
3. Create `OpenRouterClient` in `internal/llm/openrouter.go`.
4. Add `CostAccumulator` to `internal/llm/cost.go`.
5. Update `NewProvider` factory in `internal/llm/provider.go`.
6. Update all callsites (`engine.go`, `director`, `character/actor.go`, `summary.go`) to use `Provider` and thread `Usage` returns into the accumulator.
7. Update `cmd/` to wire factory and accumulator, pass to summary.
8. Update `.env.example` and `configuration` spec.

Rollback: revert to `*llm.Client` directly — no data migrations involved.

## Open Questions

- Should `OPENROUTER_MODEL` default to a specific model (e.g. `openai/gpt-4o-mini`) or require explicit configuration? → Default to `openai/gpt-4o-mini` for low cost; document in `.env.example`.
- Do we want to surface per-call cost in the JSONL simulation log as well, or only in the summary? → Summary only for this change; JSONL extension is a future concern.
