# Honey Heist Scenario

## Requirements

### Requirement: Honey Heist scenario directory

The system SHALL provide a `simulations/honey-heist/` directory containing `characters.yaml`, `world.yaml`, and `scenario.yaml` that can be loaded by the scenario loader without errors.

#### Scenario: Scenario loads successfully
- **WHEN** `--scenario honey-heist` is passed to the simulator CLI
- **THEN** the scenario loader SHALL resolve `simulations/honey-heist/`, read all three YAML files, and return a valid `Scenario` struct with no errors

#### Scenario: Scenario directory is missing a required file
- **WHEN** `characters.yaml` or `world.yaml` is absent from `simulations/honey-heist/`
- **THEN** the scenario loader SHALL return a non-nil error describing which file is missing

### Requirement: Honey Heist character roster

The `characters.yaml` file SHALL define exactly 6 bear characters, each with a unique `name`, an `occupation` describing their criminal specialty, a `motivation` string, a `fear` string, a `core_belief` string, an `internal_tension` string, a `formative_events` list of 2–3 causal bullets, a `voice` block with at least `formality` and `verbal_tics`, a `relational_defaults` block with `strangers`, `authority`, and `vulnerable`, and a `dialogue_examples` list of 3–4 representative lines. The legacy `personality` list and `backstory` prose fields SHALL NOT be present.

#### Scenario: All characters present with new schema
- **WHEN** the honey-heist scenario is loaded
- **THEN** `Scenario.Characters` SHALL contain 6 entries, each with non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension` fields

#### Scenario: Character names are unique
- **WHEN** the honey-heist scenario is loaded
- **THEN** no two characters in the roster SHALL share the same `Name`

#### Scenario: Each character has formative events
- **WHEN** the honey-heist scenario is loaded
- **THEN** each character's `FormativeEvents` slice SHALL have length 2 or 3

#### Scenario: Each character has dialogue examples
- **WHEN** the honey-heist scenario is loaded
- **THEN** each character's `DialogueExamples` slice SHALL have length 3 or 4

#### Scenario: Legacy fields absent from YAML
- **WHEN** `simulations/honey-heist/characters.yaml` is read as raw text
- **THEN** it SHALL NOT contain the keys `personality` or `backstory` at the top level of any character entry

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
