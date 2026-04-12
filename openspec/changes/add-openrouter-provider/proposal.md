## Why

The simulator currently only supports Ollama as its LLM backend. Adding OpenRouter support lets users run simulations against cloud models (GPT-4o, Claude, Gemini, etc.) without running a local GPU server, and lays the groundwork for easy provider extensibility in the future. Cost tracking is also introduced so users know what each simulation costs when using paid APIs.

## What Changes

- Introduce an `LLMProvider` interface (`Generate`, `Chat`, `Ping`) that both the existing Ollama client and new OpenRouter client implement.
- Replace direct `*llm.Client` wiring throughout the codebase with `llm.Provider` (interface).
- Add an `internal/llm/openrouter` client that calls the OpenRouter chat-completions endpoint and tracks per-call token usage.
- Add a provider factory (`llm.NewProvider`) that reads `LLM_PROVIDER`, `OPENROUTER_API_KEY`, and `OPENROUTER_MODEL` from the environment and returns the appropriate implementation.
- Expose token usage metadata on every `Generate`/`Chat` call so the simulation engine can accumulate a running cost total.
- Append a **Cost Report** section to the simulation summary when OpenRouter is used, showing prompt tokens, completion tokens, and estimated USD cost.
- Update `.env.example` and `configuration` spec with the new variables.

## Capabilities

### New Capabilities
- `llm-provider-interface`: A `Provider` interface and factory that abstracts over LLM backends; includes the Ollama adapter wrapping the existing client and the new OpenRouter client with token-usage tracking.
- `simulation-cost-report`: Accumulation of token usage across all LLM calls during a simulation run, and rendering of a cost/usage summary appended to the final simulation summary.

### Modified Capabilities
- `llm-client`: The existing Ollama-specific spec becomes an implementation detail behind the new interface; `Generate` and `Chat` signatures remain compatible.
- `configuration`: Add `LLM_PROVIDER`, `OPENROUTER_API_KEY`, `OPENROUTER_MODEL`, and `OPENROUTER_BASE_URL` env vars; keep all existing Ollama vars for backward compatibility.
- `simulation-output`: Summary output gains an optional Cost Report section when cost data is available.

## Impact

- `internal/llm/`: New files — `provider.go` (interface + factory), `openrouter/client.go`.
- `internal/simulation/engine.go`: Switch from `*llm.Client` to `llm.Provider`; accumulate usage stats.
- `internal/summary/summary.go`: Accept and render cost report when present.
- `cmd/`: Wire `llm.NewProvider` instead of `llm.NewClient`; pass cost accumulator through.
- `.env` / `.env.example`: New variables for provider selection and OpenRouter credentials.
- `go.mod` / `go.sum`: No new external dependencies required (OpenRouter uses standard REST/JSON).
