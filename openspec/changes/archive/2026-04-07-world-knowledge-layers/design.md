## Context

The simulation currently collapses all world knowledge into a single `State.Summary()` string that every character receives verbatim. This string contains the time of day and the last five events from a global log — meaning all characters share identical situational awareness regardless of where they are. There is no notion of character position, no spatial filtering of events, and no distinction between public knowledge (what anyone could know) and local knowledge (what only someone present at a location would know).

The change introduces two layers:
1. **Public layer** — what any character in the city can reasonably know: location names, time of day, events marked as public
2. **Local layer** — what a character knows because they are physically at a location: detailed location description, events marked as local that happened there

Character position is tracked via a `Location string` field on `Character`, initialized from YAML and updated as the simulation progresses (future work — for now it is static per scenario).

## Goals / Non-Goals

**Goals:**
- Introduce `Location.Details` (private, location-specific description) and `Location.Name`/`Description` split (public)
- Add `Event.Visibility` (`public`/`local`) and `Event.Location` (where it happened)
- Add `Character.Location` (current position)
- Split `State.Summary()` into `PublicSummary()` and `LocalContext(locationID string)`
- Update `RunExchange` to build a per-character world context string = `PublicSummary() + LocalContext(character.Location)`
- Update both scenario YAML files to populate the new fields

**Non-Goals:**
- Dynamic character movement between locations (characters start at a location and stay there for this change)
- Lorebook-style retrieval for arbitrarily long location histories
- Per-character event filtering based on investigation actions or relationship networks
- Multiplayer or real-time state updates

## Decisions

### 1. `Visibility` on Event as a string enum: `"public"` | `"local"` — default `"public"`

**Decision**: Events without an explicit `visibility` field default to `"public"` so that existing scenario YAML files load without errors.

**Rationale**: Backwards compatibility with any existing `initial_events` entries that predate this change. The default matches prior behavior (all events were globally visible).

**Alternative considered**: Default to `"local"`. Rejected because it would silently break existing scenarios where `initial_events` are plot-critical and should be known to all characters.

---

### 2. `Location.Details` is a separate string field — not a replacement for `Description`

**Decision**: Keep `Location.Description` as the public one-liner and add `Location.Details` for the richer private text shown only to characters at that location.

**Rationale**: `Description` is already used in the `PublicSummary()` location list. Changing its semantics would break existing scenarios. Two fields with distinct roles are clearer than one field with context-dependent behavior.

**Alternative considered**: Replace `Description` with a `Public`/`Private` sub-struct. Rejected as over-engineering for the current need; the flat field approach is simpler to author in YAML.

---

### 3. Character location is a `string` matching a `Location.Name` — not an ID

**Decision**: `Character.Location` holds the name of the location (e.g., `"Town Square"`) to match against `Location.Name` in the world state.

**Rationale**: Location names are already the human-readable identifiers used in event descriptions and world summaries. Using them directly avoids a separate ID system. The tradeoff is case-sensitivity — mitigated by the YAML being authored by humans who can match names exactly.

**Alternative considered**: Add `id` to `Location` (kebab-case). Rejected as unnecessary complexity at this stage; can be added later if dynamic movement requires stable identifiers.

---

### 4. `RunExchange` builds per-character context inline — no new method on Manager

**Decision**: `RunExchange` computes `PublicSummary() + "\n" + w.LocalContext(character.Location)` for each character before building the system prompt, without extracting this into a separate helper.

**Rationale**: The logic is one line per character. Extracting it is premature abstraction. If the context-building logic grows (e.g., adding memory injection), a helper can be extracted then.

---

### 5. `LocalContext` returns empty string for unknown location IDs — no error

**Decision**: If a character's `Location` does not match any known `Location.Name`, `LocalContext` returns an empty string (no private context added).

**Rationale**: Fail-soft behavior prevents a YAML typo from crashing the simulation. Characters simply receive only the public summary. A warning log is emitted so authors can detect the mismatch.

## Risks / Trade-offs

- **[Risk] Location name matching is case-sensitive** → Mitigation: document this in `CHARACTER_RULES.md` and in YAML comments; scenario authors must match names exactly
- **[Risk] Static character positions make scenarios feel artificial** → Mitigation: accepted tradeoff for this change; dynamic movement is a future capability
- **[Risk] Local events accumulate without bound per location** → Mitigation: `LocalContext` applies the same 5-event window as the old `Summary()`, but filtered to the location
- **[Risk] Public summary length grows with many locations** → Mitigation: `PublicSummary` includes only location names (not details), keeping it compact
