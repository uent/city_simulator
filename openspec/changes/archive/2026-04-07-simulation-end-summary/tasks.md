## 1. New `internal/summary` package

- [x] 1.1 Create `internal/summary/summary.go` with the `GenerateSummary` function that builds an LLM prompt from world events (capped at 100) and character final states and returns the narrative text
- [x] 1.2 Implement `SaveSummary` in the same file: creates `simulations/<scenarioName>/` if needed and writes `summary-<timestamp>.md` with colons replaced by hyphens in the timestamp
- [x] 1.3 Add the summary prompt template to `internal/llm/prompt.go` (or inline in the summary package) that formats scenario name, tick count, events, and character states into a coherent narration request

## 2. Engine integration

- [x] 2.1 Import `internal/summary` in `internal/simulation/engine.go`
- [x] 2.2 Add a `generateAndSaveSummary(ctx context.Context)` helper method on `Engine` that calls `summary.GenerateSummary` then `summary.SaveSummary`, logging and returning on any error
- [x] 2.3 Call `e.generateAndSaveSummary(ctx)` at the end of `Engine.Run()`, after the tick loop, only when the context was not cancelled (guard with `ctx.Err() == nil`)

## 3. Verification

- [x] 3.1 Run `make build` (or `go build ./...`) and confirm no compilation errors
- [ ] 3.2 Run a short simulation (e.g., `make run` or the honey-heist scenario) and verify a `summary-<timestamp>.md` file appears in the correct `simulations/<name>/` directory
- [ ] 3.3 Run the simulation a second time and confirm a second distinct summary file is created without overwriting the first
