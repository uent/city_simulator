## 1. World State Extensions

- [x] 1.1 Add `Weather string`, `Atmosphere string`, and `Tension int` fields to `world.State` in `internal/world/state.go`
- [x] 1.2 Add `Target string` and `PrivateRecipient string` fields to `world.Event` in `internal/world/state.go`
- [x] 1.3 Update `PublicSummary()` to include weather, atmosphere, and tension when non-zero/non-empty

## 2. Character Inbox

- [x] 2.1 Add `Inbox []world.Event` field (yaml:"-") to `character.Character` in `internal/character/character.go`
- [x] 2.2 Initialize `Inbox` as an empty non-nil slice in `LoadCharacters` (or wherever characters are constructed)

## 3. Director Package — Core

- [x] 3.1 Create `internal/director/action.go` with the `Action` interface (`Name() string`, `Execute(args map[string]any, state *world.State, chars []*character.Character) error`) and `ToolCall` struct (`Name string`, `Args map[string]any`)
- [x] 3.2 Create `internal/director/registry.go` with a `Registry` type (map of name → Action), a `NewRegistry() *Registry` constructor that registers all built-in actions, and `Dispatch(name string, args map[string]any, state *world.State, chars []*character.Character) error`
- [x] 3.3 Create `internal/director/parse.go` with `ParseToolCalls(raw string) ([]ToolCall, error)` that scans for `<tool_calls>...</tool_calls>`, unmarshals JSON, skips entries missing `name`, and returns empty slice on missing block or malformed JSON

## 4. Director Package — Environment Actions

- [x] 4.1 Create `internal/director/actions_env.go` with `setWeatherAction`, `setTimeOfDayAction`, and `setAtmosphereAction`; each mutates the corresponding `State` field and appends a public event

## 5. Director Package — NPC Actions

- [x] 5.1 Create `internal/director/actions_npc.go` with `moveNPCAction` (updates `Character.Location`, appends public event; returns error if ID not found)
- [x] 5.2 Add `introduceNPCAction` to `actions_npc.go` (appends new `character.Character` to chars via pointer; appends public event)
- [x] 5.3 Add `addNPCConditionAction` and `removeNPCConditionAction` to `actions_npc.go` (mutate `Character.EmotionalState`; return error if ID not found)

## 6. Director Package — World and Event Actions

- [x] 6.1 Create `internal/director/actions_world.go` with `modifyLocationAction` (updates `Location.Description` and/or `Details` by name; returns error if location not found)
- [x] 6.2 Add `triggerEncounterAction` to `actions_world.go` (appends public event of type `"encounter"` with participants and context)
- [x] 6.3 Add `triggerEventAction` to `actions_world.go` (appends public event; if `severity > 5`, increments `state.Tension` clamped to 10)
- [x] 6.4 Add `revealInformationAction` to `actions_world.go` (appends private event with `PrivateRecipient` set; also appends to target character's `Inbox`; returns error if recipient not found)
- [x] 6.5 Add `escalateTensionAction` to `actions_world.go` (adjusts `state.Tension` by `delta`, clamped to [0, 10])
- [x] 6.6 Add `narrateAction` to `actions_world.go` (appends public event of type `"narration"`; no other state mutation)

## 7. Director Package — Prompt Builder

- [x] 7.1 Create `internal/director/prompt.go` with `BuildDirectorPrompt(state *world.State, chars []*character.Character, tick int) string`; prompt must include tick, time-of-day, weather, atmosphere, tension, locations, character roster (ID + name + location), `<tools>` block with all 13 action schemas, and `<tool_calls>` output instructions

## 8. LLM Prompt — Character Inbox Integration

- [x] 8.1 Update the character prompt builder in `internal/llm/prompt.go` to append a "Private information you recently learned:" section when `character.Inbox` is non-empty, then clear `character.Inbox` (flush-on-read)
- [x] 8.2 Remove `BuildDirectorPrompt` and `ParseDirectorEvents` from `internal/llm/prompt.go` and delete `internal/llm/director.go`

## 9. Simulation Engine — Tool-Use Dispatch Loop

- [x] 9.1 Update `Engine` in `internal/simulation/engine.go` to hold a `*director.Registry` field; initialize it in `NewEngine` using `director.NewRegistry()`
- [x] 9.2 Rewrite `Engine.runDirector` to: call `director.BuildDirectorPrompt`, send to LLM, call `director.ParseToolCalls`, then loop over tool calls dispatching each via the registry; log and skip errors per call
- [x] 9.3 Change the `chars []*character.Character` slice passed to `runDirector` so NPC mutations (move, introduce, condition) propagate back to the engine's character list

## 10. Scenario and Config Updates

- [x] 10.1 Add an example `game_director` entry to `simulations/honey-heist/characters.yaml` that exercises at least one non-narrate action (e.g., `set_weather`, `trigger_event`) in the director's initial motivation description
- [x] 10.2 (Optional) Add `weather` and `atmosphere` initial values to `simulations/honey-heist/world.yaml` for demonstration
