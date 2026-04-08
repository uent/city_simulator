## Requirements

### Requirement: Ollama HTTP client initialization
The system SHALL provide a `NewClient(baseURL string, model string) *Client` constructor that stores the Ollama base URL and default model name. No network call is made at construction time.

#### Scenario: Client created with valid parameters
- **WHEN** `NewClient("http://localhost:11434", "llama3")` is called
- **THEN** the function SHALL return a non-nil `*Client` with the URL and model stored

### Requirement: Connectivity check on startup
The system SHALL provide a `Ping() error` method that issues a GET request to `<baseURL>/api/tags` and returns nil if Ollama responds with HTTP 200, or a descriptive error otherwise.

#### Scenario: Ollama is running
- **WHEN** `Ping()` is called and Ollama responds with HTTP 200
- **THEN** the method SHALL return nil

#### Scenario: Ollama is not running
- **WHEN** `Ping()` is called and the TCP connection is refused
- **THEN** the method SHALL return a non-nil error containing the base URL

### Requirement: Chat completion request
The system SHALL provide a `Chat(messages []Message, opts ...Option) (string, error)` method that sends a POST to `<baseURL>/api/chat` with `stream: false` and returns the assistant's response text.

A `Message` SHALL have `Role` (string: "system", "user", or "assistant") and `Content` (string) fields.

#### Scenario: Successful response
- **WHEN** Ollama returns HTTP 200 with a valid chat response JSON
- **THEN** the method SHALL return the `message.content` field from the response and nil error

#### Scenario: Non-200 HTTP status
- **WHEN** Ollama returns HTTP 500
- **THEN** the method SHALL return an empty string and a non-nil error containing the status code

#### Scenario: Network timeout
- **WHEN** the HTTP request exceeds the configured timeout (default 120 seconds)
- **THEN** the method SHALL return an empty string and a non-nil error describing the timeout

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

### Requirement: Prompt template builder
The system SHALL provide a `BuildSystemPrompt(c character.Character) string` function in the `llm` package that constructs a system-role message text from a character's name, personality, backstory, goals, and emotional state.

#### Scenario: All fields populated
- **WHEN** a character with name, personality, backstory, goals, and emotional state is provided
- **THEN** the returned string SHALL include all those fields in a coherent instruction format

#### Scenario: Empty goals list
- **WHEN** the character has no goals defined
- **THEN** the returned string SHALL omit the goals section without error
