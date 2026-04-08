## ADDED Requirements

### Requirement: World state struct
The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), list of locations (each with a name and description), and a shared event log (slice of `Event`).

An `Event` SHALL have: tick number, event type (string), description (string), and optional participant character IDs (slice of strings).

#### Scenario: New world state creation
- **WHEN** `NewState(locations []Location) *State` is called
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", the provided locations, and an empty event log

### Requirement: Tick advancement
The system SHALL provide an `AdvanceTick()` method that increments the tick counter by one and updates the time-of-day label according to a fixed cycle: morning (ticks 0–5), afternoon (6–11), evening (12–17), night (18–23), then repeating.

#### Scenario: Tick advances time of day
- **WHEN** `AdvanceTick()` is called 6 times from tick 0
- **THEN** the time-of-day label SHALL change from "morning" to "afternoon"

#### Scenario: Time of day cycles
- **WHEN** `AdvanceTick()` is called 24 times from tick 0
- **THEN** the tick SHALL be 24 and time-of-day SHALL return to "morning"

### Requirement: Event logging
The system SHALL provide an `AppendEvent(e Event)` method that appends a new event to the world's event log.

#### Scenario: Event appended
- **WHEN** `AppendEvent` is called with a valid event
- **THEN** the event SHALL appear as the last entry in `State.EventLog`

### Requirement: World summary for LLM context
The system SHALL provide a `Summary() string` method that returns a concise human-readable description of the current world state (time of day, recent events up to last 5) suitable for inclusion in an LLM system prompt.

#### Scenario: Summary with recent events
- **WHEN** `Summary()` is called after several events have been appended
- **THEN** the returned string SHALL include the current time-of-day and descriptions of the last 5 events

#### Scenario: Summary with no events
- **WHEN** `Summary()` is called on a freshly created state
- **THEN** the returned string SHALL describe the time-of-day and indicate no events have occurred yet
