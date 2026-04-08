## Why

Characters currently make movement and dialogue decisions without knowing who else is present in each zone, which breaks social realism — a character might "decide" to go somewhere without any awareness that allies or enemies are already there. Zone presence is a core input for believable autonomous behavior.

## What Changes

- Introduce a `ZoneRoster` query that returns the characters currently in each location.
- Expose zone presence to characters in their movement decision prompts (MoveDecision) and in the system prompt built for CharChat exchanges.
- Include the list of characters at a location in `world.State.LocalContext()` so the director and summary layers also benefit.

## Capabilities

### New Capabilities
- `zone-presence`: Tracks which characters occupy each location and exposes this roster to character prompts (movement decisions and chat system prompts).

### Modified Capabilities
- `character-actor`: MoveDecision and CharChat message payloads must carry zone-presence data so the actor can include it when building LLM prompts.
- `world-state`: `LocalContext(locationID string)` must include the list of characters currently present at that location.

## Impact

- `internal/world/state.go` — add `CharactersAt(locationID string, chars []*character.Character) []string` helper or equivalent method
- `internal/character/actor.go` — consume zone presence in movement and chat prompt builders
- `internal/character/prompt.go` — update prompt templates to render the zone roster
- `internal/messaging/` — message payloads for `MoveDecision` and `CharChat` gain a `ZoneRoster map[string][]string` field
- `internal/simulation/engine.go` — populate `ZoneRoster` before dispatching each message
