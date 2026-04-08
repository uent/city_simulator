## Context

Characters in the simulation have a `Location` field (string, name of current zone) that is updated each tick after movement decisions. The engine already pre-builds system prompts (`BuildSystemPrompt`, `BuildMovementPrompt`) and embeds them in message payloads (`CharChatPayload.InitiatorSystem/ResponderSystem`, `MoveDecisionPayload.SystemPrompt`) before dispatching to actors. Characters currently make decisions with no knowledge of who else occupies each zone.

## Goals / Non-Goals

**Goals:**
- Characters can see the names of other characters in every zone when deciding where to move.
- Characters can see who else is at their own location when engaging in dialogue.
- The zone roster is derived at the moment each tick's messages are dispatched (consistent snapshot).

**Non-Goals:**
- Real-time roster updates mid-tick (roster is computed once per tick before message dispatch).
- Persisting roster history or exposing it via the event log.
- Changing how the director sees zone presence (already has `e.chars`).

## Decisions

### 1. Roster computation at the engine level, injected into pre-built prompts

The engine computes `map[string][]string` (location → character names) by iterating `e.chars` once per tick. This map is passed to `BuildMovementPrompt` and to a new `BuildZoneContext` helper that appends a zone-presence section to CharChat system prompts.

**Why not in world.State?** State holds authoritative world data (events, weather, tension). Character positions live on `Character.Location`, not on `State`, so computing the roster from `State` would require threading the char list into world — creating a coupling that doesn't belong there.

**Why not in the payload struct?** Payloads carry pre-built string prompts; the actors don't reconstruct prompts from raw data. Adding a `ZoneRoster` field to `CharChatPayload` or `MoveDecisionPayload` would require every actor to re-render it into text, duplicating prompt-building logic. Instead: build the text in the engine, where all prompt building already happens.

### 2. Zone roster appended to existing prompts, not replacing them

`BuildMovementPrompt` gains a `zoneRoster map[string][]string` parameter and appends a "Who is where" section listing all zones with their occupants. `BuildSystemPrompt` is unchanged (static persona); the engine appends a zone-awareness block after calling it when constructing `CharChatPayload`.

**Why a separate helper instead of modifying `BuildSystemPrompt`?** `BuildSystemPrompt` takes only a `Character` and is used in multiple contexts. Adding a roster parameter would change its signature for all callers. A dedicated `BuildZoneContext(roster map[string][]string) string` function is composable and keeps the character package clean.

### 3. Roster excludes the character itself from its own location listing

When showing "who else is at your current location", the character's own name is omitted to avoid redundancy. The full roster (all zones + all names) is included for movement decisions so characters can reason about where to go.

## Risks / Trade-offs

- [Prompt length increases] → Each prompt grows by ~N lines where N = number of characters. Acceptable for current simulation sizes (≤20 characters).
- [Stale roster within a tick] → Characters dispatched later in a tick may see positions from before earlier characters moved. Mitigation: roster is a snapshot taken before any movement messages are dispatched in a tick; movement decisions are applied after all responses arrive.
- [Ordering sensitivity] → Character names in the roster are in iteration order of `e.chars` slice (deterministic for a given run). No mitigation needed.

## Open Questions

- Should the roster use character `Name` or `ID` in prompts? **Name** is more natural for LLM reasoning and matches how characters reference each other in dialogue.
