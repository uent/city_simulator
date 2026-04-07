## ADDED Requirements

### Requirement: Generate method on Client
The system SHALL provide a `Generate(systemPrompt, userPrompt string, opts ...Option) (string, error)` method on `*Client` that constructs a two-message slice (`system` + `user`) and delegates to `Chat`.

#### Scenario: Successful generation
- **WHEN** `Generate("You are a helpful assistant.", "What is 2+2?")` is called and Ollama responds with HTTP 200
- **THEN** the method SHALL return the assistant's response text and nil error

#### Scenario: System prompt empty
- **WHEN** `Generate("", "hello")` is called
- **THEN** the method SHALL still issue the request with an empty system message and return whatever Ollama replies

#### Scenario: Option override applied
- **WHEN** `Generate(sys, user, WithModel("mistral"))` is called
- **THEN** the request SHALL use `"mistral"` as the model, not the client's default

#### Scenario: Ollama returns error
- **WHEN** Ollama responds with a non-200 status
- **THEN** the method SHALL return an empty string and a non-nil error (propagated from `Chat`)
