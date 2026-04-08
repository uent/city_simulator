## Why

The existing Game Director generates free-form JSON events that get appended to the world log, but it cannot directly mutate world state — it can only hint that something happened. A tool-use architecture replaces this with a structured action layer: the GM calls named functions (`set_weather`, `move_npc`, `trigger_encounter`, etc.) whose implementations validate inputs, update world state, and dispatch messages to affected actors. This makes GM decisions programmatically trackable, replayable, and auditable instead of relying on narrative prose.

## What Changes

- Replace the current JSON-array event output from the Game Director with a **structured action dispatch loop**: the GM receives a tool schema and responds with tool-call JSON; the engine executes each action sequentially
- Introduce a `DirectorAction` interface and a registry of named action handlers in a new `internal/director/` package
- Each action handler: validates its arguments, mutates the relevant slice of `world.State`, and optionally enqueues public or private messages for affected characters
- The `narrate` action is a special no-op that produces a public broadcast without mutating state
- `escalate_tension` adjusts a numeric tension level on the world state that characters can observe
- `reveal_information` places a targeted private message into a specific character's inbox
- Remove `ParseDirectorEvents` and `BuildDirectorPrompt` from `internal/llm/` (superseded)
- **BREAKING**: `world.Event` gains new optional fields (`Target`, `PrivateRecipient`) to support private messaging; existing scenarios are unaffected (fields are optional)

## Capabilities

### New Capabilities

- `game-director-actions`: The action dispatch layer — action interface, handler registry, and all named action implementations (`set_weather`, `set_time_of_day`, `set_atmosphere`, `modify_location`, `move_npc`, `introduce_npc`, `add_npc_condition`, `remove_npc_condition`, `trigger_encounter`, `trigger_event`, `reveal_information`, `escalate_tension`, `narrate`)

### Modified Capabilities

- `world-state`: Add `Tension int`, `Weather string`, `TimeOfDay string`, `Atmosphere string` top-level fields; extend `Event` with optional `Target string` and `PrivateRecipient string` fields for targeted messaging
- `simulation-engine`: Replace `runDirector` free-JSON call with a tool-use dispatch loop; route private events to the correct character's context
- `character-engine`: Characters gain an `Inbox []world.Event` field (or equivalent) so private messages from `reveal_information` are visible only to them in their next prompt

## Impact

- New package `internal/director/` (~6 files): action interface, registry, and one file per action group
- `internal/world/state.go`: new top-level fields + Event field extensions
- `internal/simulation/engine.go`: `runDirector` rewritten to use action dispatch loop
- `internal/character/character.go`: add `Inbox` slice for private events
- `internal/llm/prompt.go`: include private inbox in character prompt; include world-level fields (weather, time, atmosphere, tension) in both character and director prompts; remove old director prompt builder
- `internal/llm/director.go`: replaced by new director prompt builder that emits a tool-schema block
- Scenario YAML files: no changes required (weather/time/atmosphere fields are optional with sensible defaults)
