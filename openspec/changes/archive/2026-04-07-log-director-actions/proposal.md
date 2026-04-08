## Why

The simulation log currently prints `[Director] set_weather` with no detail about what changed, making it impossible to follow the director's decisions at a glance. Adding the key argument value (e.g., `[Director] set_weather: lluvioso`) gives observers immediate context without having to cross-reference world state.

## What Changes

- Add a `Summary(args map[string]any) string` method to the `Action` interface so each action can describe itself concisely.
- Implement `Summary` on all 13 registered director actions, returning a human-readable one-liner (e.g., `set_weather: storm`, `move_npc: alice → market`, `escalate_tension: +3`).
- Update `engine.go` to call `Summary` when logging each director action, replacing the current name-only output.

## Capabilities

### New Capabilities

*(none — this is a pure enhancement to existing behaviour)*

### Modified Capabilities

- `game-director-actions`: The `Action` interface gains a third method `Summary(args map[string]any) string`; all existing action implementations must satisfy it.
- `simulation-engine`: The director log line format changes from `[Director] <name>` to `[Director] <summary>`.

## Impact

- `internal/director/action.go` — interface definition
- `internal/director/actions_env.go`, `actions_npc.go`, `actions_world.go` — new `Summary` implementations on all 13 action structs
- `internal/simulation/engine.go:87` — log line updated to use `Summary`
- No API, config, or dependency changes
