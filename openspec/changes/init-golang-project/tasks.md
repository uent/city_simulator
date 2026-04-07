## 1. Project Bootstrap

- [x] 1.1 Run `go mod init github.com/<user>/city_simulator` and create the full directory skeleton (`cmd/simulator/`, `internal/character/`, `internal/llm/`, `internal/simulation/`, `internal/world/`, `internal/conversation/`, `configs/`)
- [x] 1.2 Add `gopkg.in/yaml.v3` dependency with `go get`
- [x] 1.3 Create `configs/characters.yaml` with 3 sample characters (name, age, occupation, personality, backstory, goals)

## 2. World State Package (`internal/world`)

- [x] 2.1 Define `Location`, `Event`, and `State` structs in `state.go`
- [x] 2.2 Implement `NewState(locations []Location) *State` constructor
- [x] 2.3 Implement `AdvanceTick()` with time-of-day cycle (morning/afternoon/evening/night)
- [x] 2.4 Implement `AppendEvent(e Event)` method
- [x] 2.5 Implement `Summary() string` method returning time-of-day + last 5 events

## 3. Character Engine Package (`internal/character`)

- [x] 3.1 Define `Character` struct with all persona fields in `character.go`
- [x] 3.2 Implement `LoadCharacters(path string) ([]Character, error)` with YAML unmarshaling and default for empty emotional state
- [x] 3.3 Define `MemoryEntry` struct and add `Memory []MemoryEntry` + `MaxMemory int` fields to `Character` in `memory.go`
- [x] 3.4 Implement `AddMemory(entry MemoryEntry)` with sliding-window eviction
- [x] 3.5 Implement `RecentMemory(n int) []MemoryEntry` returning last n entries

## 4. LLM Client Package (`internal/llm`)

- [x] 4.1 Define `Message`, `ChatRequest`, `ChatResponse` structs and `Client` struct in `client.go`
- [x] 4.2 Implement `NewClient(baseURL, model string) *Client` constructor
- [x] 4.3 Implement `Ping() error` via GET to `/api/tags`
- [x] 4.4 Implement `Chat(messages []Message, opts ...Option) (string, error)` via POST to `/api/chat` with `stream: false`, 120-second timeout, and HTTP error handling
- [x] 4.5 Implement `BuildSystemPrompt(c character.Character) string` in `prompt.go` covering persona, backstory, goals, and emotional state

## 5. Conversation Manager Package (`internal/conversation`)

- [x] 5.1 Define `Turn`, `Thread`, and `Manager` types in `manager.go`
- [x] 5.2 Implement `NewManager(llmClient *llm.Client) *Manager`
- [x] 5.3 Implement `AddTurn(fromID, toID, text string, tick int)` with auto-create of missing threads
- [x] 5.4 Implement `History(fromID, toID string, maxTurns int) []llm.Message` with role assignment and maxTurns truncation
- [x] 5.5 Implement `RunExchange(ctx, initiator, responder, world, tick)` following the 8-step flow from the spec (build prompts → LLM call × 2 → store memories → append world event)

## 6. Simulation Engine Package (`internal/simulation`)

- [x] 6.1 Define `Config` struct and `Engine` struct in `engine.go`
- [x] 6.2 Implement `NewEngine(cfg Config) (*Engine, error)` with validation (at least 2 characters)
- [x] 6.3 Implement `Scheduler` in `scheduler.go`: generate all unique character pairs, round-robin `Next()`, and seeded shuffle when seed != 0
- [x] 6.4 Implement `Run(ctx context.Context) error` tick loop: select pair → call `RunExchange` → advance world tick → write JSONL log entry; handle LLM errors by logging and continuing

## 7. CLI Entry Point (`cmd/simulator`)

- [x] 7.1 Implement `main.go` with flags: `--characters`, `--model`, `--ollama-url`, `--turns`, `--seed`, `--output`
- [x] 7.2 Wire startup sequence: parse flags → load characters → ping Ollama (exit with clear error if unreachable) → create world state with default locations → create manager + engine → call `engine.Run`
- [x] 7.3 Implement JSONL output writer that writes one line per exchange event to the output file
- [x] 7.4 Add formatted stdout printing of each exchange (character name + dialogue text) for human-readable narrative

## 8. Smoke Test

- [x] 8.1 Run `go build ./...` and fix any compilation errors
- [ ] 8.2 Start Ollama locally with `ollama run llama3` and run `go run ./cmd/simulator --turns 5` end-to-end, verifying output log and stdout narrative
