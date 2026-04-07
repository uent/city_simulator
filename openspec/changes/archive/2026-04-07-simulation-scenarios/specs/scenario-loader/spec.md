## ADDED Requirements

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
The system SHALL define a `WorldConfig` struct with: `Locations []Location` (each with `Name` and `Description` string fields) and `InitialEvents []world.Event` (optional). The loader SHALL populate this struct from `world.yaml`.

#### Scenario: Locations parsed correctly
- **WHEN** `world.yaml` contains two location entries
- **THEN** `Scenario.World.Locations` SHALL have length 2 with matching names and descriptions

#### Scenario: Initial events parsed
- **WHEN** `world.yaml` contains an `initial_events` list
- **THEN** `Scenario.World.InitialEvents` SHALL be populated and appended to world state on simulation start

#### Scenario: No initial events key
- **WHEN** `world.yaml` omits `initial_events`
- **THEN** `Scenario.World.InitialEvents` SHALL be an empty slice and no error returned

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
