## MODIFIED Requirements

### Requirement: WorldConfig loaded from `world.yaml`

The system SHALL define a `WorldConfig` struct with: `Locations []Location` (each with `Name string`, `Description string`, and `Details string` fields) and `InitialEvents []world.Event` (optional). The loader SHALL populate this struct from `world.yaml`, including the new `Details`, `Visibility`, and `Location` fields on their respective types.

#### Scenario: Locations parsed correctly with Details
- **WHEN** `world.yaml` contains two location entries, one with a `details` key
- **THEN** `Scenario.World.Locations` SHALL have length 2, with the first location's `Details` field populated and the second's empty

#### Scenario: Initial events parsed with visibility and location
- **WHEN** `world.yaml` contains an `initial_events` list where one event has `visibility: "local"` and `location: "Tavern"`
- **THEN** `Scenario.World.InitialEvents` SHALL contain that event with `Visibility == "local"` and `Location == "Tavern"`

#### Scenario: No initial events key
- **WHEN** `world.yaml` omits `initial_events`
- **THEN** `Scenario.World.InitialEvents` SHALL be an empty slice and no error returned

#### Scenario: Events without visibility default to public
- **WHEN** `world.yaml` contains an event that omits the `visibility` key
- **THEN** the loaded event's `Visibility` field SHALL equal `"public"`
