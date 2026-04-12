## ADDED Requirements

### Requirement: Provider interface
The `internal/llm` package SHALL define a `Provider` interface with the following methods:
- `Generate(systemPrompt, userPrompt string, opts ...Option) (string, Usage, error)`
- `Chat(messages []Message, opts ...Option) (string, Usage, error)`
- `Ping() error`

Any type implementing these three methods satisfies `Provider` and can be used anywhere the simulator requires an LLM backend.

#### Scenario: Interface satisfied by both adapters
- **WHEN** `OllamaProvider` and `OpenRouterClient` are compiled
- **THEN** both SHALL satisfy `llm.Provider` at compile time with no type assertion needed

---

### Requirement: Usage value type
The `internal/llm` package SHALL define a `Usage` struct with fields:
- `PromptTokens int`
- `CompletionTokens int`
- `EstimatedCostUSD float64`

`Usage` SHALL be returned from every `Generate` and `Chat` call. When token counts or cost are unavailable (e.g. Ollama), all fields SHALL be zero.

#### Scenario: Zero usage from Ollama
- **WHEN** `OllamaProvider.Generate` is called and succeeds
- **THEN** the returned `Usage` SHALL have all fields equal to zero

#### Scenario: Non-zero usage from OpenRouter
- **WHEN** `OpenRouterClient.Generate` is called and the response contains a `usage` object
- **THEN** `PromptTokens` and `CompletionTokens` SHALL match the response values, and `EstimatedCostUSD` SHALL match `usage.cost` if present

---

### Requirement: OllamaProvider adapter
The `internal/llm` package SHALL provide an `OllamaProvider` struct that wraps the existing `*Client` and implements `Provider`. Its `Generate` and `Chat` methods SHALL delegate to the underlying `*Client` and return `Usage{}` alongside the result.

#### Scenario: Delegation to underlying client
- **WHEN** `OllamaProvider.Generate(sys, user)` is called
- **THEN** the method SHALL forward the call to the wrapped `*Client.Generate` and return its text result with zero `Usage` and nil error on success

#### Scenario: Error propagation
- **WHEN** the underlying Ollama client returns an error
- **THEN** `OllamaProvider.Generate` SHALL return empty string, zero `Usage`, and the original error

---

### Requirement: OpenRouter client
The `internal/llm` package SHALL provide an `OpenRouterClient` struct that calls `https://openrouter.ai/api/v1/chat/completions` (or a configurable base URL) using an OpenAI-compatible JSON payload and returns token usage from the response.

The client SHALL include the `Authorization: Bearer <api_key>` header on every request.

#### Scenario: Successful call returns text and usage
- **WHEN** `OpenRouterClient.Generate` is called and OpenRouter responds HTTP 200 with valid JSON
- **THEN** the method SHALL return the assistant's message content, populated `Usage`, and nil error

#### Scenario: Non-200 response returns error
- **WHEN** OpenRouter responds with HTTP 4xx or 5xx
- **THEN** the method SHALL return empty string, zero `Usage`, and a non-nil error containing the status code and response body

#### Scenario: Missing usage field handled gracefully
- **WHEN** the OpenRouter response does not include a `usage` field
- **THEN** the method SHALL return the text content with zero `Usage` and nil error

#### Scenario: Ping validates credentials
- **WHEN** `OpenRouterClient.Ping()` is called with a valid API key
- **THEN** it SHALL make a GET request to `<baseURL>/models` and return nil on HTTP 200

#### Scenario: Ping fails with invalid API key
- **WHEN** `OpenRouterClient.Ping()` is called with an empty or invalid API key and OpenRouter returns 401
- **THEN** it SHALL return a non-nil error describing the authentication failure

---

### Requirement: Provider factory
The `internal/llm` package SHALL expose `NewProvider(cfg ProviderConfig) (Provider, error)` that inspects `cfg.ProviderName` and returns the appropriate implementation:
- `"ollama"` → `OllamaProvider` wrapping a new `*Client`
- `"openrouter"` → `OpenRouterClient`
- any other value → non-nil error

`ProviderConfig` SHALL carry all fields needed by both adapters: `ProviderName`, `OllamaURL`, `OllamaModel`, `OpenRouterAPIKey`, `OpenRouterModel`, `OpenRouterBaseURL`.

#### Scenario: Factory returns Ollama provider by default
- **WHEN** `NewProvider` is called with `ProviderName: "ollama"`
- **THEN** the returned `Provider` SHALL be an `OllamaProvider` and error SHALL be nil

#### Scenario: Factory returns OpenRouter provider
- **WHEN** `NewProvider` is called with `ProviderName: "openrouter"` and a non-empty `OpenRouterAPIKey`
- **THEN** the returned `Provider` SHALL be an `OpenRouterClient` and error SHALL be nil

#### Scenario: Factory rejects unknown provider
- **WHEN** `NewProvider` is called with `ProviderName: "unknown"`
- **THEN** error SHALL be non-nil and SHALL name the unrecognized provider

#### Scenario: Factory rejects OpenRouter with missing API key
- **WHEN** `NewProvider` is called with `ProviderName: "openrouter"` and empty `OpenRouterAPIKey`
- **THEN** error SHALL be non-nil indicating the API key is required
