## ADDED Requirements

### Requirement: Honey Heist scenario directory

The system SHALL provide a `simulations/honey-heist/` directory containing `characters.yaml`, `world.yaml`, and `scenario.yaml` that can be loaded by the scenario loader without errors.

#### Scenario: Scenario loads successfully
- **WHEN** `--scenario honey-heist` is passed to the simulator CLI
- **THEN** the scenario loader SHALL resolve `simulations/honey-heist/`, read all three YAML files, and return a valid `Scenario` struct with no errors

#### Scenario: Scenario directory is missing a required file
- **WHEN** `characters.yaml` or `world.yaml` is absent from `simulations/honey-heist/`
- **THEN** the scenario loader SHALL return a non-nil error describing which file is missing

### Requirement: Honey Heist character roster

The `characters.yaml` file SHALL define exactly 6 bear characters, each with a unique `name`, a `role` describing their criminal specialty, and a `personality` string of 1–2 sentences that the LLM uses to stay in character.

#### Scenario: All characters present
- **WHEN** the scenario is loaded
- **THEN** the resulting `Scenario.Characters` slice SHALL contain 6 entries with non-empty `name`, `role`, and `personality` fields

#### Scenario: Character names are unique
- **WHEN** the scenario is loaded
- **THEN** no two characters in the roster SHALL share the same `name`

### Requirement: Honey Heist world layout

The `world.yaml` file SHALL define exactly 6 locations representing the HoneyCon convention centre and its surroundings: Convention Lobby, Vendor Hall, Security Office, Vault Antechamber, Vault, and Alley (Exit). Each location SHALL have a non-empty `name` and `description`. At least one `initial_event` SHALL be present to seed narrative context.

#### Scenario: All locations present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.Locations` SHALL contain 6 entries, each with a non-empty `name` and `description`

#### Scenario: Initial events present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialEvents` SHALL contain at least one event with a non-empty `description`

### Requirement: Honey Heist runtime overrides

The `scenario.yaml` file SHALL set `turns` to 20 so the heist resolves in a short run by default.

#### Scenario: Turn count override applied
- **WHEN** the scenario is loaded and no `--turns` CLI flag is provided
- **THEN** the simulation SHALL run for 20 turns as specified by `scenario.yaml`

#### Scenario: CLI flag overrides scenario.yaml turns
- **WHEN** the scenario is loaded and `--turns 50` is passed on the CLI
- **THEN** the simulation SHALL run for 50 turns, ignoring the `scenario.yaml` value of 20
