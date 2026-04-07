## 1. Scenario Directory and Example Content

- [x] 1.1 Create `simulations/default/characters.yaml` by moving content from `configs/characters.yaml` (3 sample characters)
- [x] 1.2 Create `simulations/default/world.yaml` with 4–5 city locations and 1–2 initial events
- [x] 1.3 Create `simulations/default/scenario.yaml` with default overrides (`model: llama3`, `turns: 10`)
- [x] 1.4 Delete the old `configs/` directory

## 2. Scenario Loader Package (`internal/scenario`)

- [x] 2.1 Define `WorldConfig`, `RuntimeOverrides`, `SimConfig`, and `Scenario` structs in `scenario.go`
- [x] 2.2 Implement `Load(dirOrName string) (Scenario, error)`: resolve path (absolute or under `simulations/`), validate required files exist
- [x] 2.3 Implement YAML loading for `characters.yaml` into `Scenario.Characters` (reuse `character.LoadCharacters` logic or delegate to it)
- [x] 2.4 Implement YAML loading for `world.yaml` into `Scenario.World` including optional `initial_events`
- [x] 2.5 Implement YAML loading for optional `scenario.yaml` into `Scenario.Overrides` (all pointer fields, nil when absent)
- [x] 2.6 Implement `MergeConfig(overrides RuntimeOverrides, flags CLIFlags, defaults SimConfig) SimConfig` with CLI > scenario.yaml > defaults priority

## 3. Update World State (`internal/world`)

- [x] 3.1 Change `NewState` signature from `NewState(locations []Location)` to `NewState(cfg scenario.WorldConfig) *State`
- [x] 3.2 Pre-populate `State.EventLog` with `cfg.InitialEvents` in `NewState`

## 4. Update Simulation Engine (`internal/simulation`)

- [x] 4.1 Replace `Characters []character.Character` field in `Config` with `Scenario scenario.Scenario`
- [x] 4.2 Update `NewEngine` to derive characters from `cfg.Scenario.Characters` and initialize world via `world.NewState(cfg.Scenario.World)`

## 5. Update CLI Entry Point (`cmd/simulator`)

- [x] 5.1 Replace `--characters` flag with `--scenario` flag (string, default `"default"`)
- [x] 5.2 Add startup error message when `--characters` is passed: `"--characters has been removed, use --scenario instead"`
- [x] 5.3 Call `scenario.Load(scenarioFlag)` on startup and propagate errors with clear message
- [x] 5.4 Call `scenario.MergeConfig` to produce final `SimConfig` from loaded overrides, CLI flags, and defaults
- [x] 5.5 Wire `SimConfig` fields into engine `Config` and replace direct `character.LoadCharacters` call

## 6. Smoke Test

- [x] 6.1 Run `go build ./...` and fix any compilation errors from signature changes
- [x] 6.2 Run `go run ./cmd/simulator --scenario default --turns 5` and verify it loads `simulations/default/` correctly
- [x] 6.3 Create a second scenario directory `simulations/test-scenario/` with different characters and verify `--scenario test-scenario` switches to it cleanly
