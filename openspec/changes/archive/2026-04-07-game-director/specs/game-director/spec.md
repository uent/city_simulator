## ADDED Requirements

### Requirement: Game Director prompt builder
The system SHALL provide a `BuildDirectorPrompt(state *world.State, chars []character.Character, tick int) string` function in `internal/llm/prompt.go` that produces a system prompt for the Game Director LLM call.

The prompt SHALL include:
- Current tick number and time-of-day
- Full list of all locations (name + description)
- Current location of each regular character (ID, name, location)
- The last 10 events from the world event log (all visibilities)
- Clear instruction to return a JSON array of 0 to 3 world events

Each event in the requested JSON array SHALL have fields: `event_type` (string), `description` (string), `visibility` ("public" or "local"), and optional `location` (string matching a known location name).

#### Scenario: Prompt includes all characters
- **WHEN** `BuildDirectorPrompt` is called with 3 regular characters
- **THEN** the returned string SHALL mention all 3 characters' names and current locations

#### Scenario: Prompt caps event history at 10
- **WHEN** the world event log has 15 events
- **THEN** the prompt SHALL include only the most recent 10 events

#### Scenario: Prompt instructs JSON output
- **WHEN** `BuildDirectorPrompt` is called
- **THEN** the returned string SHALL contain instructions to respond with a JSON array and SHALL specify the required fields (`event_type`, `description`, `visibility`)

### Requirement: Game Director event parsing
The system SHALL provide a `ParseDirectorEvents(raw string) ([]world.Event, error)` function in `internal/llm` that extracts a JSON array of events from the LLM's raw response string.

The function SHALL:
- Attempt to extract a JSON array from the raw response (which may contain surrounding text)
- Ignore any event with a missing or blank `description`
- Cap the returned slice to at most 3 events
- Return an empty slice (not an error) if no valid JSON array is found

#### Scenario: Valid JSON array returned
- **WHEN** the raw string is `[{"event_type":"weather","description":"Heavy rain begins","visibility":"public"}]`
- **THEN** `ParseDirectorEvents` SHALL return a slice with one `world.Event` and a nil error

#### Scenario: JSON embedded in prose
- **WHEN** the raw string contains surrounding prose before and after the JSON array
- **THEN** `ParseDirectorEvents` SHALL locate and parse the array and return the events

#### Scenario: Invalid or missing JSON
- **WHEN** the raw string contains no JSON array
- **THEN** `ParseDirectorEvents` SHALL return an empty slice and a nil error (not crash)

#### Scenario: More than 3 events in array
- **WHEN** the JSON array has 5 entries
- **THEN** `ParseDirectorEvents` SHALL return only the first 3 entries

### Requirement: Game Director YAML configuration
The system SHALL support defining a Game Director in a scenario's `characters.yaml` by setting `type: game_director` on any character entry.

A Game Director entry MUST have an `id` and `name`. All other fields (motivation, fear, etc.) are optional but MAY be used to flavor the director's behavior prompt.

At most one Game Director per scenario is supported; if multiple `type: game_director` entries exist, the first is used and a warning is logged.

#### Scenario: Game Director loaded from YAML
- **WHEN** `characters.yaml` contains one entry with `type: game_director`
- **THEN** `scenario.Load` SHALL populate `Scenario.GameDirector` with that entry and exclude it from `Scenario.Characters`

#### Scenario: No Game Director in YAML
- **WHEN** `characters.yaml` contains no entry with `type: game_director`
- **THEN** `Scenario.GameDirector` SHALL be nil and all entries are treated as regular characters

#### Scenario: Multiple Game Director entries
- **WHEN** `characters.yaml` contains two entries with `type: game_director`
- **THEN** `Scenario.GameDirector` SHALL be set to the first one, a warning SHALL be logged, and the second SHALL be discarded
