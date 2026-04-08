## MODIFIED Requirements

### Requirement: World state struct

The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), list of locations (each with `Name`, `Description`, and `Details` string fields), and a shared event log (slice of `Event`).

An `Event` SHALL have: tick number, event type (string), description (string), optional participant character IDs (slice of strings), `Visibility` (string, `"public"` or `"local"`, default `"public"`), and `Location` (string, name of the location where the event occurred, optional).

`State` SHALL expose `PublicSummary() string` and `LocalContext(locationID string) string` methods. The `Summary() string` method SHALL NOT exist.

#### Scenario: New world state creation from WorldConfig
- **WHEN** `NewState(cfg scenario.WorldConfig) *State` is called with a valid `WorldConfig`
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", locations from `cfg.Locations` (including their `Details` fields), and an event log pre-populated with `cfg.InitialEvents` (empty slice if none defined)

#### Scenario: Initial events appear in log
- **WHEN** `WorldConfig.InitialEvents` contains one event
- **THEN** `State.EventLog` SHALL contain that event at index 0 immediately after `NewState` returns

#### Scenario: Location Details preserved in State
- **WHEN** `WorldConfig.Locations` contains a location with a non-empty `Details` field
- **THEN** `State.Locations` SHALL contain that location with `Details` intact
