## Requirements

### Requirement: Scenario directory structure
A scenario SHALL be a directory containing: `characters.yaml` (required), `world.yaml` (required), and `scenario.yaml` (optional). No other files are required. Extra files SHALL be silently ignored.

#### Scenario: Valid scenario directory
- **WHEN** a directory contains both `characters.yaml` and `world.yaml`
- **THEN** `Load` SHALL parse both files and return a populated `Scenario` without error

#### Scenario: Missing required file
- **WHEN** a directory is missing `characters.yaml` or `world.yaml`
- **THEN** `Load` SHALL return a non-nil error naming the missing file and the directory path

### Requirement: Scenario resolution by name or path
The system SHALL provide a `Load(dirOrName string) (Scenario, error)` function. If `dirOrName` is an absolute path it SHALL be used directly. If it is a plain name (no path separators), the function SHALL look it up under `simulations/<dirOrName>` relative to the working directory.

#### Scenario: Load by name
- **WHEN** `Load("default")` is called and `simulations/default/` exists with valid files
- **THEN** the function SHALL return the populated `Scenario` for that directory

#### Scenario: Load by absolute path
- **WHEN** `Load("/some/absolute/path")` is called and the directory exists with valid files
- **THEN** the function SHALL load from that path regardless of the `simulations/` convention

#### Scenario: Name not found
- **WHEN** `Load("nonexistent")` is called and `simulations/nonexistent/` does not exist
- **THEN** the function SHALL return a non-nil error indicating the scenario was not found

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

### Requirement: RuntimeOverrides loaded from `scenario.yaml`
The system SHALL define a `RuntimeOverrides` struct where all fields are pointers (nil means "not set"): `Model *string`, `Turns *int`, `Seed *int64`, `Output *string`. The loader SHALL parse `scenario.yaml` if present and populate only the fields that appear in the file.

#### Scenario: Override file present with partial fields
- **WHEN** `scenario.yaml` defines only `model: mistral`
- **THEN** `Scenario.Overrides.Model` SHALL point to `"mistral"` and all other override fields SHALL be nil

#### Scenario: Override file absent
- **WHEN** no `scenario.yaml` exists in the scenario directory
- **THEN** all `Scenario.Overrides` fields SHALL be nil and no error returned

### Requirement: Override merge priority (CLI > scenario.yaml > defaults)
The system SHALL provide a `MergeConfig(overrides RuntimeOverrides, flags CLIFlags, defaults SimConfig) SimConfig` function that returns a final `SimConfig` applying the priority: explicit CLI flag values beat scenario overrides, which beat compiled defaults.

#### Scenario: CLI flag wins over scenario override
- **WHEN** `scenario.yaml` sets `turns: 30` and the CLI flag `--turns 10` is provided
- **THEN** the resolved `SimConfig.Turns` SHALL be 10

#### Scenario: Scenario override fills unset CLI flag
- **WHEN** `scenario.yaml` sets `model: mistral` and no `--model` flag is given
- **THEN** the resolved `SimConfig.Model` SHALL be `"mistral"`

#### Scenario: Default used when neither CLI nor scenario sets a value
- **WHEN** neither CLI nor `scenario.yaml` specifies `seed`
- **THEN** the resolved `SimConfig.Seed` SHALL equal the compiled default (0)

### Requirement: Game Director separation during load
The system SHALL separate Game Director entries from regular characters when loading `characters.yaml`. `Scenario` SHALL expose a `GameDirector *character.Character` field (nil if none defined) alongside the existing `Characters []character.Character` field (regular characters only).

`Load` SHALL:
- Populate `Scenario.GameDirector` with the first entry whose `Type == "game_director"`
- Exclude all `type: game_director` entries from `Scenario.Characters`
- Log a warning if more than one `type: game_director` entry is found and use only the first

#### Scenario: Single Game Director entry
- **WHEN** `characters.yaml` contains one entry with `type: game_director` and two regular entries
- **THEN** `Scenario.GameDirector` SHALL point to that entry, `Scenario.Characters` SHALL have length 2, and no warning is logged

#### Scenario: No Game Director entry
- **WHEN** `characters.yaml` contains no `type: game_director` entry
- **THEN** `Scenario.GameDirector` SHALL be nil and all entries SHALL appear in `Scenario.Characters`

#### Scenario: Multiple Game Director entries
- **WHEN** `characters.yaml` contains two entries with `type: game_director`
- **THEN** `Scenario.GameDirector` SHALL be set to the first one, a warning SHALL be logged, and `Scenario.Characters` SHALL contain only the non-director entries
