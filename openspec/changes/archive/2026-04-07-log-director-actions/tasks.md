## 1. Extend Action Interface

- [x] 1.1 Add `Summary(args map[string]any) string` to the `Action` interface in `internal/director/action.go`

## 2. Implement Summary on Environment Actions

- [x] 2.1 Add `Summary` to `setWeatherAction` — returns `"set_weather: <type>"`
- [x] 2.2 Add `Summary` to `setTimeOfDayAction` — returns `"set_time_of_day: <moment>"`
- [x] 2.3 Add `Summary` to `setAtmosphereAction` — returns `"set_atmosphere: <descriptor>"`

## 3. Implement Summary on NPC Actions

- [x] 3.1 Add `Summary` to `moveNPCAction` — returns `"move_npc: <id> → <destination>"`
- [x] 3.2 Add `Summary` to `introduceNPCAction` — returns `"introduce_npc: <name> (<role>)"`
- [x] 3.3 Add `Summary` to `addNPCConditionAction` — returns `"add_npc_condition: <id> += <condition>"`
- [x] 3.4 Add `Summary` to `removeNPCConditionAction` — returns `"remove_npc_condition: <id> -= <condition>"`

## 4. Implement Summary on World/Event Actions

- [x] 4.1 Add `Summary` to `modifyLocationAction` — returns `"modify_location: <name>"`
- [x] 4.2 Add `Summary` to `triggerEncounterAction` — returns `"trigger_encounter: <context>"`
- [x] 4.3 Add `Summary` to `triggerEventAction` — returns `"trigger_event: <type> — <description>"`
- [x] 4.4 Add `Summary` to `revealInformationAction` — returns `"reveal_information: <recipient> — <content>"`
- [x] 4.5 Add `Summary` to `escalateTensionAction` — returns `"escalate_tension: <+/-delta>"`
- [x] 4.6 Add `Summary` to `narrateAction` — returns `"narrate: <text>"`

## 5. Update Engine Log Line

- [x] 5.1 In `internal/simulation/engine.go:87`, replace `fmt.Printf("  [Director] %s\n", call.Name)` with `fmt.Printf("  [Director] %s\n", e.registry.Summarize(call.Name, call.Args))` (or equivalent using the action's `Summary` method directly)
- [x] 5.2 Add a `Summarize(name string, args map[string]any) string` helper to `Registry` that delegates to the action's `Summary`, falling back to `name` if the action is not found
