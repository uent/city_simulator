## Why

Every character currently receives the same world context via `w.Summary()` — time of day plus the last five events from a single shared log. This makes all characters effectively omniscient: if something happens at the Market, Elena, Marcus, and Nadia all know about it simultaneously, regardless of where they are. Characters need spatial grounding and information scarcity to behave realistically.

## What Changes

- **Add `Details string` to `Location`**: each location in `world.yaml` gets a private details field — richer descriptions visible only to characters present at that location (not the public one-liner)
- **Add `Visibility` and `Location` fields to `Event`**: events can be `public` (broadcast to all) or `local` (visible only to characters at the event's location)
- **Add `Location string` to `Character`**: characters have a starting location that determines their local context
- **Split `State.Summary()` into two methods**:
  - `PublicSummary()` — time of day + location names + public events; given to all characters
  - `LocalContext(locationID string)` — location details + local events at that location; given per character based on their current position
- **Update `RunExchange`** in the conversation manager to pass character-specific world context (public + local) instead of the shared global summary
- **Update `world.yaml` files** for both scenarios to add `details` to locations and `visibility`/`location` to events
- **Update `characters.yaml` files** to add `location` to each character

## Capabilities

### New Capabilities

- `world-knowledge-layers`: Two-layer world knowledge system — a public layer (universal knowledge visible to all) and a local layer (per-location details and events visible only to characters present). Includes character location tracking and per-character world context generation.

### Modified Capabilities

- `world-state`: `Location` struct gains `Details`; `Event` gains `Visibility` and `Location`; `State.Summary()` is replaced by `PublicSummary()` and `LocalContext(locationID string)`
- `scenario-loader`: `WorldConfig` propagates the new `Location.Details`, `Event.Visibility`, and `Event.Location` fields; `Character` struct gains a `Location string` field loaded from YAML

## Impact

- `internal/world/state.go`: `Location` struct, `Event` struct, `State` methods
- `internal/character/character.go`: `Character` struct gains `Location string`
- `internal/conversation/manager.go`: `RunExchange` builds per-character world context
- `simulations/default/world.yaml`: add `details` to locations, `visibility`/`location` to events
- `simulations/default/characters.yaml`: add `location` to each character
- `simulations/honey-heist/world.yaml`: same
- `simulations/honey-heist/characters.yaml`: same
- `openspec/specs/world-state/spec.md` and `scenario-loader/spec.md`: updated
