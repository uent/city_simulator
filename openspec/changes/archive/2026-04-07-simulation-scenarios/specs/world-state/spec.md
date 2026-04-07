## MODIFIED Requirements

### Requirement: World state struct
The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), list of locations (each with a name and description), and a shared event log (slice of `Event`).

An `Event` SHALL have: tick number, event type (string), description (string), and optional participant character IDs (slice of strings).

#### Scenario: New world state creation from WorldConfig
- **WHEN** `NewState(cfg scenario.WorldConfig) *State` is called with a valid `WorldConfig`
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", locations from `cfg.Locations`, and an event log pre-populated with `cfg.InitialEvents` (empty slice if none defined)

#### Scenario: Initial events appear in log
- **WHEN** `WorldConfig.InitialEvents` contains one event
- **THEN** `State.EventLog` SHALL contain that event at index 0 immediately after `NewState` returns
