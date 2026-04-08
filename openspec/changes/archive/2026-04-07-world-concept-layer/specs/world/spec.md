## MODIFIED Requirements

### Requirement: PublicSummary method

`State` SHALL expose a `PublicSummary() string` method that returns the universal world context available to all characters. It SHALL include: the current time of day, the list of all location names with their `Description` (not `Details`), weather and atmosphere when non-empty, tension level when non-zero (as a descriptor, e.g., "Tension: 7/10"), the last 5 events whose `Visibility` is `"public"`, and — when `State.Concept.Premise` is non-empty — a "World Rules" block containing the premise and, if present, the rules as a bulleted list. If there are no public events, only time of day, location list, and (if set) world rules are included.

#### Scenario: Returns time and location names
- **WHEN** `PublicSummary()` is called on a `State` with two locations and no events
- **THEN** the returned string SHALL contain the time-of-day label and both location names

#### Scenario: Includes only public events
- **WHEN** the event log contains one public event and one local event
- **THEN** `PublicSummary()` SHALL include only the public event

#### Scenario: Limits to 5 most recent public events
- **WHEN** the event log contains 8 public events
- **THEN** `PublicSummary()` SHALL include only the 5 most recent ones

#### Scenario: PublicSummary includes weather when set
- **WHEN** `state.Weather == "storm"`
- **THEN** `PublicSummary()` SHALL contain the word `"storm"`

#### Scenario: PublicSummary includes tension when non-zero
- **WHEN** `state.Tension == 7`
- **THEN** `PublicSummary()` SHALL contain `"7"` in the tension description

#### Scenario: PublicSummary omits weather when empty
- **WHEN** `state.Weather == ""`
- **THEN** `PublicSummary()` SHALL NOT contain a weather line

#### Scenario: World Rules block present when Concept.Premise is set
- **WHEN** `State.Concept.Premise == "Bears disguised as humans at a honey convention"`
- **THEN** `PublicSummary()` SHALL contain that premise string in a "World Rules" section

#### Scenario: World Rules block includes rules list when non-empty
- **WHEN** `State.Concept.Rules` contains `["Do not walk on all fours", "Never eat honey directly from a jar"]`
- **THEN** `PublicSummary()` SHALL contain both rule strings

#### Scenario: World Rules block absent when Concept is zero value
- **WHEN** `WorldConfig.Concept.Premise == ""`
- **THEN** `PublicSummary()` SHALL NOT contain the heading "World Rules"

---

### Requirement: World state struct

The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), weather (string, default `""`), atmosphere (string, default `""`), tension (int, 0–10, default 0), list of locations (each with `Name`, `Description`, and `Details` string fields), a shared event log (slice of `Event`), and a `Concept WorldConcept` field that mirrors `WorldConfig.Concept`.

An `Event` SHALL have: tick number, event type (string), description (string), optional participant character IDs (slice of strings), `Visibility` (string, `"public"` or `"local"`, default `"public"`), `Location` (string, name of the location where the event occurred, optional), `Target` (string, optional character or location ID the event is "about"), and `PrivateRecipient` (string, optional — if set, only that character should receive this event in their inbox).

#### Scenario: New world state creation from WorldConfig
- **WHEN** `NewState(cfg scenario.WorldConfig) *State` is called with a valid `WorldConfig`
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", weather `""`, atmosphere `""`, tension 0, locations from `cfg.Locations`, `Concept` copied from `cfg.Concept`, and an event log pre-populated with `cfg.InitialEvents`

#### Scenario: Concept copied from WorldConfig into State
- **WHEN** `WorldConfig.Concept.Premise == "Bears pretending to be humans"`
- **THEN** `NewState(cfg).Concept.Premise` SHALL equal `"Bears pretending to be humans"`

#### Scenario: Initial events appear in log
- **WHEN** `WorldConfig.InitialEvents` contains one event
- **THEN** `State.EventLog` SHALL contain that event at index 0 immediately after `NewState` returns

#### Scenario: Location Details preserved in State
- **WHEN** `WorldConfig.Locations` contains a location with a non-empty `Details` field
- **THEN** `State.Locations` SHALL contain that location with `Details` intact
