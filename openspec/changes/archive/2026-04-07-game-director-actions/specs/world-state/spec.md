## MODIFIED Requirements

### Requirement: World state struct
The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), weather (string, default `""`), atmosphere (string, default `""`), tension (int, 0–10, default 0), list of locations (each with `Name`, `Description`, and `Details` string fields), and a shared event log (slice of `Event`).

An `Event` SHALL have: tick number, event type (string), description (string), optional participant character IDs (slice of strings), `Visibility` (string, `"public"` or `"local"`, default `"public"`), `Location` (string, name of the location where the event occurred, optional), `Target` (string, optional character or location ID the event is "about"), and `PrivateRecipient` (string, optional — if set, only that character should receive this event in their inbox).

`State` SHALL expose `PublicSummary() string` and `LocalContext(locationID string) string` methods. The `Summary() string` method SHALL NOT exist.

`PublicSummary()` SHALL include weather, atmosphere, and tension level (as a descriptor, e.g., "Tension: 7/10") when those fields are non-zero/non-empty, in addition to time-of-day, locations, and recent public events.

#### Scenario: New world state creation from WorldConfig
- **WHEN** `NewState(cfg scenario.WorldConfig) *State` is called with a valid `WorldConfig`
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", weather `""`, atmosphere `""`, tension 0, locations from `cfg.Locations`, and an event log pre-populated with `cfg.InitialEvents`

#### Scenario: Initial events appear in log
- **WHEN** `WorldConfig.InitialEvents` contains one event
- **THEN** `State.EventLog` SHALL contain that event at index 0 immediately after `NewState` returns

#### Scenario: Location Details preserved in State
- **WHEN** `WorldConfig.Locations` contains a location with a non-empty `Details` field
- **THEN** `State.Locations` SHALL contain that location with `Details` intact

#### Scenario: PublicSummary includes weather when set
- **WHEN** `state.Weather == "storm"`
- **THEN** `PublicSummary()` SHALL contain the word `"storm"`

#### Scenario: PublicSummary includes tension when non-zero
- **WHEN** `state.Tension == 7`
- **THEN** `PublicSummary()` SHALL contain `"7"` in the tension description

#### Scenario: PublicSummary omits weather when empty
- **WHEN** `state.Weather == ""`
- **THEN** `PublicSummary()` SHALL NOT contain a weather line
