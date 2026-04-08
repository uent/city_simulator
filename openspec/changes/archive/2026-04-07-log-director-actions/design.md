## Context

The simulation engine currently logs each successful director action as `[Director] <name>` (engine.go:87). The `Action` interface exposes only `Name()` and `Execute()`. To include meaningful detail in the log without coupling the engine to action internals, each action must describe itself.

## Goals / Non-Goals

**Goals:**
- Log director actions with a concise human-readable summary (e.g., `[Director] set_weather: storm`).
- Keep each action responsible for its own summary — no string-formatting logic in the engine.

**Non-Goals:**
- Structured/machine-parseable log output.
- Logging failed or skipped actions differently (existing error path unchanged).
- Changes to the world event log (that already has descriptive strings).

## Decisions

### Add `Summary(args map[string]any) string` to the `Action` interface

Each action already receives its own `args` and knows which fields matter. Putting `Summary` on the interface keeps the engine lean and avoids a switch/map at the call site.

**Alternative considered:** a standalone `Summarize(name string, args map[string]any) string` helper in `engine.go`. Rejected — it would need to duplicate arg-key knowledge already embedded in each action.

### Summary format per action

| Action | Format |
|---|---|
| `set_weather` | `set_weather: <type>` |
| `set_time_of_day` | `set_time_of_day: <moment>` |
| `set_atmosphere` | `set_atmosphere: <descriptor>` |
| `move_npc` | `move_npc: <id> → <destination>` |
| `introduce_npc` | `introduce_npc: <name> (<role>)` |
| `add_npc_condition` | `add_npc_condition: <id> += <condition>` |
| `remove_npc_condition` | `remove_npc_condition: <id> -= <condition>` |
| `modify_location` | `modify_location: <name>` |
| `trigger_encounter` | `trigger_encounter: <context>` |
| `trigger_event` | `trigger_event: <type> — <description>` |
| `reveal_information` | `reveal_information: <recipient> — <content>` |
| `escalate_tension` | `escalate_tension: <+delta>` |
| `narrate` | `narrate: <text>` |

If a required arg is missing (args already validated by `Execute`), `Summary` falls back to the action name alone.

## Risks / Trade-offs

- [Risk] Interface change requires all 13 action structs to be updated at once → Mitigation: all implementations live in three files; the change is mechanical and contained.
- [Risk] Long `text`/`description` strings make log lines noisy → Accepted trade-off; truncation is not required for now.
