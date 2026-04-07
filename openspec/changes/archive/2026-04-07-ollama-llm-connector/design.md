## Context

`internal/llm/client.go` already contains `Client`, `Message`, `ChatRequest`, `ChatResponse`, `Option`, `NewClient`, `Ping`, and `Chat`. The file uses only stdlib packages. No other files need to change.

## Goals / Non-Goals

**Goals:**
- Single new method `Generate(systemPrompt, userPrompt string, opts ...Option) (string, error)` on `*Client`
- Constructs `[]Message{{Role:"system", Content:systemPrompt}, {Role:"user", Content:userPrompt}}` and calls `Chat`
- Accepts the same variadic `Option` funcs so callers can still override the model per-call

**Non-Goals:**
- Streaming support (out of scope for this change)
- Changing `Chat` signature or any existing method
- New files or packages

## Decisions

### 1. Thin wrapper over `Chat`

`Generate` builds the two-message slice and calls `c.Chat(messages, opts...)`. All HTTP logic, error handling, and timeout stay in `Chat`. No duplication.

**Rationale:** Keeps the single source of truth for HTTP behavior in `Chat`. `Generate` is purely a convenience shim.

### 2. Keep `Chat` public

`Chat` stays exported so callers that need multi-turn history or a custom message order can still use it directly. `Generate` is additive.

**Rationale:** Backward compatibility; `Chat` is already referenced in existing task 5.4 (`History` returns `[]llm.Message` consumed by `Chat`).

### 3. Same `Option` variadic

Reusing the existing `Option` type (e.g., `WithModel`) means no new API surface beyond the method itself.

## Implementation

Add to `internal/llm/client.go`:

```go
// Generate sends a system prompt and a user prompt to Ollama and returns
// the assistant's reply. It is a convenience wrapper around Chat.
func (c *Client) Generate(systemPrompt, userPrompt string, opts ...Option) (string, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	return c.Chat(messages, opts...)
}
```

No new imports needed.
