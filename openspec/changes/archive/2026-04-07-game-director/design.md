## Context

The simulation currently drives all characters through the same `conversation.Manager` pipeline — each character has personal goals, memory, and responds to dialogue. Nothing in the engine generates autonomous world events (weather shifts, random encounters, environmental changes) that characters react to but did not cause themselves.

The codebase already has an `Event` / `EventLog` model in `internal/world/state.go` and an `AppendEvent` method. The `PublicSummary()` function already surfaces the last 5 public events to all characters. The missing piece is an autonomous agent that _writes_ those events before each tick.

## Goals / Non-Goals

**Goals:**
- Introduce a `GameDirector` character type that has privileged, read-only access to the full world state
- Game Director generates 0-N world events per tick as structured JSON, which are parsed and appended to the shared event log
- Game Director is invoked at the start of each tick, before any character exchanges, so events are already visible when characters act
- Game Director is optional — scenarios without one work exactly as before
- Game Director does NOT participate in conversations; it has no memory buffer and no dialogue turns

**Non-Goals:**
- Game Director cannot directly modify character memory, location, or emotional state (it can only inject events that characters _react_ to autonomously)
- No multi-Game-Director support in this change; only one per scenario is allowed
- No UI or external API for triggering events manually

## Decisions

### 1. Same YAML file, new `type` field on Character

Add `Type string yaml:"type"` to `character.Character`. Values: `"character"` (default, blank) and `"game_director"`. The scenario loader splits the slice: `Scenario.Characters` (regular) stays as-is; `Scenario.GameDirector` is a new `*character.Character` field.

**Alternatives considered:**
- Separate `game_director.yaml` file — rejected because it adds file-loading complexity; a single field in `characters.yaml` is simpler and consistent
- Separate `GameDirector` struct — rejected because the Character struct's persona fields (name, motivation, etc.) are still useful for the director's persona flavoring in prompts

### 2. Output format: JSON array of events

The Game Director's LLM call asks for a JSON array where each element has `event_type`, `description`, `visibility` (`"public"` or `"local"`), and optional `location`. Structured output is easier to parse reliably than free text.

A `MaxEvents` cap (default 3) per tick prevents prompt flooding.

**Alternatives considered:**
- Free text parsed with heuristics — rejected, fragile
- Single event per tick — rejected, directors often want to generate multiple simultaneous events (e.g., weather + a crowd gathering)

### 3. Engine calls Game Director first in each tick

In `Engine.Run`, the existing loop starts with `e.scheduler.Next()`. We add a `if e.director != nil { e.runDirector(ctx, tick) }` call _before_ `scheduler.Next()`. The director's events are appended to world state immediately, so `PublicSummary()` in the same tick already includes them.

**Alternatives considered:**
- End-of-tick — events would only be visible next tick; less reactive
- Separate goroutine/concurrent — unnecessary complexity

### 4. Game Director validation: warn, don't fail

If a scenario defines more than one Game Director entry, the loader logs a warning and uses the first. This avoids breaking existing simulations on misconfiguration.

The `NewEngine` minimum-character check (≥2) counts only regular characters; the Game Director is excluded from this count.

## Risks / Trade-offs

- **LLM output parsing failure** → Mitigation: wrap JSON parse in a best-effort recovery; if the output isn't valid JSON, log the error and skip events for that tick (don't crash the simulation)
- **Game Director adds latency per tick** → Trade-off accepted; it's one additional LLM call. Can be mitigated later by batching or async execution
- **Prompt flooding with full world state** → Mitigation: cap full character list to IDs + names + locations (not full personas); expose complete persona only for scenarios that need it
- **Game Director using character minimum check** → Resolved by excluding the director from the ≥2 character count

## Migration Plan

1. Add `Type` field to `Character` (backward-compatible; defaults to blank/"character")
2. Update `scenario.Load` to populate `Scenario.GameDirector`
3. Add `BuildDirectorPrompt` to `internal/llm/prompt.go`
4. Add `runDirector` method to `Engine`; wire into `Run`
5. Existing scenarios without a `game_director` entry are unaffected
