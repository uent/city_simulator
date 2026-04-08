## 1. Character Model

- [x] 1.1 Add `Type string` field (YAML tag `type`) to `character.Character` struct in `internal/character/character.go`

## 2. Scenario Loader

- [x] 2.1 Add `GameDirector *character.Character` field to `scenario.Scenario` struct in `internal/scenario/scenario.go`
- [x] 2.2 Update `scenario.Load` to separate entries by type: populate `Scenario.GameDirector` with the first `type: game_director` entry, keep remaining entries in `Scenario.Characters`, and log a warning if more than one director entry is found

## 3. LLM Prompt Builder

- [x] 3.1 Add `BuildDirectorPrompt(state *world.State, chars []character.Character, tick int) string` to `internal/llm/prompt.go` — includes tick, time-of-day, all locations, character positions, last 10 events, and JSON output instructions
- [x] 3.2 Add `ParseDirectorEvents(raw string) ([]world.Event, error)` to `internal/llm/` — extracts a JSON array from raw LLM output, filters blank descriptions, caps at 3 events, returns empty slice on parse failure

## 4. Simulation Engine

- [x] 4.1 Add `director *character.Character` field to `Engine` struct in `internal/simulation/engine.go`
- [x] 4.2 Update `NewEngine` to populate `engine.director` from `cfg.Scenario.GameDirector` and exclude the director from the ≥2 regular character validation count
- [x] 4.3 Add `runDirector(ctx context.Context, tick int) error` method to `Engine` that calls `BuildDirectorPrompt`, invokes the LLM, calls `ParseDirectorEvents`, sets `Event.Tick`, and appends events via `world.AppendEvent`
- [x] 4.4 Call `e.runDirector(ctx, tick)` at the start of each tick loop in `Engine.Run`, before `e.scheduler.Next()`, when `e.director != nil`

## 5. Example Scenario Config

- [x] 5.1 Add a Game Director entry to `simulations/honey-heist/characters.yaml` with `type: game_director` to demonstrate the feature in an existing scenario
