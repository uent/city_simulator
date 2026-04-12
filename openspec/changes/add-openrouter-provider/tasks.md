## 1. Provider Interface and Usage Type

- [x] 1.1 Create `internal/llm/provider.go` with the `Provider` interface (`Generate`, `Chat`, `Ping`), the `Usage` struct (`PromptTokens int`, `CompletionTokens int`, `EstimatedCostUSD float64`), and the `ProviderConfig` struct
- [x] 1.2 Add `NewProvider(cfg ProviderConfig) (Provider, error)` factory function in `provider.go` that returns `OllamaProvider` for `"ollama"` and `OpenRouterClient` for `"openrouter"`, and errors on unknown names or missing OpenRouter API key

## 2. Update Ollama Client Signatures

- [x] 2.1 Update `(*Client).Chat` signature in `internal/llm/client.go` to return `(string, Usage, error)` — always returning zero `Usage`
- [x] 2.2 Update `(*Client).Generate` signature in `internal/llm/client.go` to return `(string, Usage, error)` — delegating to `Chat` and propagating the zero `Usage`
- [x] 2.3 Create `OllamaProvider` struct in `internal/llm/provider.go` that wraps `*Client` and implements `Provider`, delegating `Generate`, `Chat`, and `Ping` to the underlying client

## 3. OpenRouter Client

- [x] 3.1 Create `internal/llm/openrouter.go` with `OpenRouterClient` struct holding `apiKey`, `model`, `baseURL`, and an `*http.Client`
- [x] 3.2 Implement `OpenRouterClient.Chat` using OpenAI-compatible JSON payload (`model`, `messages`); parse `choices[0].message.content` for text and `usage.prompt_tokens`, `usage.completion_tokens`, `usage.cost` for `Usage`
- [x] 3.3 Implement `OpenRouterClient.Generate` as a convenience wrapper around `Chat` (same pattern as `*Client.Generate`)
- [x] 3.4 Implement `OpenRouterClient.Ping` with a `GET <baseURL>/models` call using the API key; return nil on HTTP 200, error otherwise

## 4. Cost Accumulator

- [x] 4.1 Create `internal/llm/cost.go` with `CostAccumulator` struct (mutex-protected running totals) and `Add(u Usage)` and `Total() Usage` methods

## 5. Update All LLM Callsites to Use Provider

- [x] 5.1 Update `internal/character/actor.go` to accept `llm.Provider` instead of `*llm.Client`; capture the returned `Usage` from `Generate`/`Chat` calls and add to a passed-in `*llm.CostAccumulator`
- [x] 5.2 Update `internal/director/` (prompt and actions files) to accept `llm.Provider` instead of `*llm.Client`; capture and accumulate `Usage`
- [x] 5.3 Update `internal/simulation/engine.go` to hold `llm.Provider` and `*llm.CostAccumulator` instead of `*llm.Client`; pass accumulator to actors and director
- [x] 5.4 Update `simulation.Config` struct in `engine.go` to replace `LLMClient *llm.Client` with `LLMProvider llm.Provider` and add `CostAccumulator *llm.CostAccumulator`

## 6. Cost Report in Summary

- [x] 6.1 Update `internal/summary/summary.go` `GenerateSummary` signature to accept `llm.Provider` instead of `*llm.Client` and add `acc *llm.CostAccumulator` parameter
- [x] 6.2 Implement cost report rendering in `summary.go`: if `acc != nil && acc.Total().PromptTokens+acc.Total().CompletionTokens > 0`, append a `## Cost Report` section with prompt tokens, completion tokens, total tokens, and cost formatted as `$0.000000`

## 7. Wire Everything in cmd/main.go

- [x] 7.1 Add new flags and env var reads in `main.go`: `--provider` / `LLM_PROVIDER` (default `"ollama"`), `--openrouter-base-url` / `OPENROUTER_BASE_URL`; read `OPENROUTER_API_KEY` and `OPENROUTER_MODEL` from env only (no flags for secrets)
- [x] 7.2 Replace `llm.NewClient` + `client.Ping` with `llm.NewProvider(cfg)` + `provider.Ping()`; exit with descriptive error if factory returns error or Ping fails
- [x] 7.3 Create a `llm.CostAccumulator` in `main.go` and pass it into `simulation.Config` and `summary.GenerateSummary`
- [x] 7.4 Update the startup log line to display the active provider name and model instead of hardcoded "Ollama"

## 8. Configuration Files

- [x] 8.1 Update `.env.example` to add `LLM_PROVIDER`, `OPENROUTER_API_KEY`, `OPENROUTER_MODEL`, and `OPENROUTER_BASE_URL` with descriptive comments, keeping all existing Ollama variables intact
