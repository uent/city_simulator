## 1. Implement Generate method

- [x] 1.1 Add `Generate(systemPrompt, userPrompt string, opts ...Option) (string, error)` to `internal/llm/client.go` — builds `[]Message{{Role:"system",...},{Role:"user",...}}` and calls `Chat`

## 2. Smoke test

- [x] 2.1 Run `go build ./...` and confirm no compilation errors
- [x] 2.2 With Ollama running, call `client.Generate("You are a pirate.", "Say hello.")` and verify a response is returned
