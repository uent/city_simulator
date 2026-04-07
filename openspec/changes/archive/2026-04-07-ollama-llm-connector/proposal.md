## Why

The existing `llm.Client` exposes a `Chat(messages []Message)` method that requires callers to manually build the `[]Message` slice every time. In practice, every call site has the same two-message pattern: one `system` role message and one `user` role message. Forcing callers to assemble that slice adds boilerplate and risks role ordering mistakes.

A dedicated `Generate(systemPrompt, userPrompt string)` method encapsulates that construction, making call sites one line instead of four and keeping the message-building logic in one place.

## What Changes

- Add `Generate(systemPrompt, userPrompt string, opts ...Option) (string, error)` to `internal/llm/client.go`
- The method builds the `[]Message{system, user}` slice internally and delegates to `Chat`

## Capabilities

### Modified Capabilities

- `llm-client`: Extends the existing Ollama HTTP client with a `Generate` convenience method that accepts system and user prompts as plain strings instead of a pre-built message slice

### New Capabilities

<!-- none -->

## Impact

- One new exported method on `*Client` — no breaking changes to existing callers of `Chat` or `NewClient`
- No new dependencies
- All existing usages of `Chat` continue to work unchanged
