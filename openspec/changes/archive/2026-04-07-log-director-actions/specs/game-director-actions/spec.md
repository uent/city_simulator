## MODIFIED Requirements

### Requirement: Director action interface
The system SHALL define an `Action` interface in `internal/director/action.go` with three methods: `Name() string`, `Execute(args map[string]any, state *world.State, chars *[]*character.Character) error`, and `Summary(args map[string]any) string`. Every named director action SHALL implement this interface.

#### Scenario: Action returns its name
- **WHEN** `Name()` is called on any registered action
- **THEN** it SHALL return the exact string used to register it in the registry

#### Scenario: Action mutates world state
- **WHEN** `Execute` is called with valid args and a non-nil `*world.State`
- **THEN** the action SHALL mutate `state` in-place and return nil

#### Scenario: Action receives invalid args
- **WHEN** `Execute` is called with args that fail validation (e.g., missing required key, out-of-range value)
- **THEN** the action SHALL return a non-nil error without mutating `state`

#### Scenario: Summary returns name and key arg
- **WHEN** `Summary` is called on `set_weather` with `{"type": "storm"}`
- **THEN** it SHALL return `"set_weather: storm"`

#### Scenario: Summary falls back to name when arg missing
- **WHEN** `Summary` is called with an empty args map
- **THEN** it SHALL return the action name alone (no panic)
