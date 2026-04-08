## Context

The current Game Director calls the LLM once per tick and asks it to return a raw JSON array of events. `ParseDirectorEvents` extracts that array and calls `world.AppendEvent` for each item. This works for narration, but the GM has no way to directly change `State.TimeOfDay`, reposition a character, add a condition to an NPC, or send a private message — it can only inject text into the public event log.

The new architecture makes the GM an **action-dispatch agent**: the LLM is given a tool schema (a JSON block listing all available actions and their parameter shapes) and responds with a sequence of tool-call objects. The engine parses those calls, looks up each handler in a registry, validates arguments, and executes them in order. This is analogous to function-calling in OpenAI or Anthropic's tool-use API, but implemented in the simulation engine itself so it works with any LLM that can follow a structured-output prompt.

## Goals / Non-Goals

**Goals:**
- GM decisions are expressed as named function calls, not free text
- Each action handler validates its own arguments and returns a descriptive error on failure (logged, not fatal)
- World-level fields (`Weather`, `TimeOfDay`, `Atmosphere`, `Tension`) are first-class state that characters can observe
- Private messaging: `reveal_information` writes an event visible only to a named recipient's inbox
- `narrate` is a zero-mutation action that broadcasts a public event, replacing the old free-text narration
- All actions are registered in a central registry; adding a new action requires only a new handler registration
- Backward-compatible: scenarios without a Game Director work exactly as before

**Non-Goals:**
- No streaming or async action execution — actions run sequentially in the order the GM emits them
- No player-facing UI to trigger actions manually (future work)
- No multi-GM support
- No rollback / undo for actions within a tick

## Decisions

### 1. Internal tool-use via structured prompt + JSON parser, not native API tool-use

The LLM prompt includes a `<tools>` block describing each action as a JSON Schema object. The GM is instructed to respond with a `<tool_calls>` block containing an array of `{"name": "...", "args": {...}}` objects, followed by optional free-form narration. The engine parses the `<tool_calls>` block only; the rest is discarded.

**Why not use the Anthropic/OpenAI native tool-use API?**
The simulation uses a generic `llm.Client` interface with a single `Generate(ctx, prompt) (string, error)` method. Native tool-use would require a different call shape and bind the GM to a specific provider. The structured-prompt approach works with any backend and lets us prototype quickly without changing the client interface.

**Alternative**: Separate tool-use client interface — deferred; can be added later without changing action handlers.

### 2. `internal/director/` package owns the action layer

A new package `internal/director/` holds:
- `action.go` — `Action` interface: `Name() string`, `Execute(args map[string]any, state *world.State, chars []*character.Character) error`
- `registry.go` — `Registry` map + `Dispatch(name, args, state, chars)` method
- `actions_env.go` — environment actions: `set_weather`, `set_time_of_day`, `set_atmosphere`
- `actions_npc.go` — NPC actions: `move_npc`, `introduce_npc`, `add_npc_condition`, `remove_npc_condition`
- `actions_world.go` — world actions: `modify_location`, `trigger_encounter`, `trigger_event`, `reveal_information`, `escalate_tension`, `narrate`
- `prompt.go` — `BuildDirectorPrompt(state, chars, tick) string` (replaces `internal/llm/`'s version)
- `parse.go` — `ParseToolCalls(raw string) ([]ToolCall, error)`

**Why a dedicated package?**
The action layer is self-contained and testable in isolation. Keeping it out of `internal/llm/` avoids mixing prompt engineering with business logic.

### 3. World state gains top-level GM-controlled fields

`world.State` gets: `Weather string`, `Atmosphere string`, `Tension int` (0–10). `TimeOfDay` already exists. `PublicSummary()` is updated to include these fields so all characters observe them.

`world.Event` gets two optional fields: `Target string` (location or character ID the event is "about") and `PrivateRecipient string` (if set, only that character sees it in their inbox).

**Why not a separate `GMState` struct?**
Keeping GM-controlled fields on `world.State` means characters and the engine read from a single source of truth. A separate struct would require merging on every prompt build.

### 4. Character inbox for private events

`character.Character` gains `Inbox []world.Event`. When `reveal_information` runs, it appends an event with `PrivateRecipient = targetID` to both the world log and the target character's `Inbox`. The character prompt builder in `internal/llm/prompt.go` appends inbox items as "Private information you recently learned:" and clears the inbox after reading (flush-on-read).

**Why flush-on-read?**
Prevents inbox accumulation across ticks. If the character misses it (LLM ignores it), the information is gone — this is intentional; private info is time-sensitive.

### 5. `escalate_tension` clamps to 0–10

The tension field is an integer clamped to [0, 10]. Actions `escalate_tension` and `trigger_event` (with severity) can raise or lower it. Characters see tension as a descriptor in `PublicSummary`: "The tension in the city is high (8/10)." This gives the GM a single dial to modulate narrative intensity.

### 6. Remove `ParseDirectorEvents` and old `BuildDirectorPrompt`

`internal/llm/director.go` is deleted entirely. `internal/llm/prompt.go`'s director-related functions are removed. All director logic lives in `internal/director/`. This avoids split ownership and makes the old API impossible to accidentally call.

## Risks / Trade-offs

- **LLM doesn't follow the tool schema** → Mitigation: `ParseToolCalls` is lenient — it scans for the first `<tool_calls>` block, parses JSON, and skips malformed entries. If the entire block is missing, zero actions execute (no crash, simulation continues).
- **Unknown action name in output** → Mitigation: `Registry.Dispatch` logs `"unknown action: X"` and skips. No panic.
- **Invalid args for a valid action** → Mitigation: each handler validates before mutating state and returns a typed error. Engine logs the error and continues to the next action.
- **Inbox flush-on-read loses private info** → Trade-off accepted; private messages are ephemeral by design. Future work could add a `MemoryEntry` for important revelations.
- **Engine now imports `internal/director/`** → Creates a new dependency edge. Acceptable; director is simulation-only, not a general utility.

## Migration Plan

1. Add new fields to `world.State` and `world.Event` (backward-compatible zero values)
2. Add `Inbox` to `character.Character`
3. Create `internal/director/` package with all action handlers
4. Update `internal/llm/prompt.go` to include new world fields and character inbox; remove old director functions
5. Delete `internal/llm/director.go`
6. Rewrite `Engine.runDirector` to use `director.BuildDirectorPrompt` + `director.ParseToolCalls` + `registry.Dispatch`
7. Update scenario YAML files to demonstrate new GM capabilities (optional `weather`, `atmosphere` fields)
8. Existing scenarios without a Game Director: unaffected (no code paths change for them)
