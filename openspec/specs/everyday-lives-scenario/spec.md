# Everyday Lives Scenario

### Requirement: Everyday lives scenario directory

The system SHALL provide a `simulations/vida-cotidiana/` directory containing `characters.yaml`, `world.yaml`, and `scenario.yaml` that can be loaded by the scenario loader without errors.

#### Scenario: Scenario loads successfully
- **WHEN** `--scenario vida-cotidiana` is passed to the simulator CLI
- **THEN** the scenario loader SHALL resolve `simulations/vida-cotidiana/`, read all three YAML files, and return a valid `Scenario` struct with no errors

#### Scenario: Scenario directory is missing a required file
- **WHEN** `characters.yaml` or `world.yaml` is absent from `simulations/vida-cotidiana/`
- **THEN** the scenario loader SHALL return a non-nil error describing which file is missing

### Requirement: Everyday lives character roster

The `characters.yaml` file SHALL define exactly 5 characters representing residents of the same apartment building. Each character SHALL have a unique `id`, `name`, `age`, `gender`, `occupation`, `motivation`, `fear`, `core_belief`, `internal_tension`, a `formative_events` list of 2–3 causal bullets, a `voice` block with `formality`, `verbal_tics`, `response_length`, `humor_type`, and `communication_style`, a `relational_defaults` block with `strangers`, `authority`, and `vulnerable`, a `dialogue_examples` list of 3–4 representative lines, a `goals` list of 1–2 concrete objectives, and an `emotional_state` string. No character SHALL have a `type: game_director` field. The legacy `personality` list and `backstory` prose fields SHALL NOT be present. No character SHALL have a `cover_identity` block.

#### Scenario: All characters present with complete schema
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** `Scenario.Characters` SHALL contain exactly 5 entries, each with non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension` fields

#### Scenario: No game director character
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** no character in `Scenario.Characters` SHALL have `Type` equal to `game_director`

#### Scenario: Character names are unique
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** no two characters in the roster SHALL share the same `Name`

#### Scenario: Each character has formative events
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** each character's `FormativeEvents` slice SHALL have length 2 or 3

#### Scenario: Each character has dialogue examples
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** each character's `DialogueExamples` slice SHALL have length 3 or 4

#### Scenario: Legacy fields absent from YAML
- **WHEN** `simulations/vida-cotidiana/characters.yaml` is read as raw text
- **THEN** it SHALL NOT contain the keys `personality`, `backstory`, or `cover_identity` at the top level of any character entry

### Requirement: Everyday lives world layout

The `world.yaml` file SHALL define exactly 5 locations representing shared spaces in and around the apartment building: Lobby del Edificio, Azotea, Escalera y Pasillos, Cafetería Marga, and Parque del Barrio. Each location SHALL have a non-empty `name` and `description`. Location `details` fields SHALL describe only physical and environmental properties — they SHALL NOT reference the name or current state of any specific character. At least two `initial_events` SHALL be present to seed narrative tension without requiring external conflict. The `world.yaml` SHALL define `initial_location` pointing to a valid location name. The `world.yaml` SHALL define a `concept:` block with a non-empty `premise`, a `flavor` string, and at least 3 `rules`.

#### Scenario: All locations present
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** `Scenario.World.Locations` SHALL contain 5 entries, each with a non-empty `name` and `description`

#### Scenario: Initial events present
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** `Scenario.World.InitialEvents` SHALL contain at least 2 events with non-empty `description` fields

#### Scenario: World concept block present
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** `Scenario.World.Concept.Premise` SHALL be non-empty and `Scenario.World.Concept.Rules` SHALL contain at least 3 entries

#### Scenario: initial_location references an existing location
- **WHEN** the vida-cotidiana scenario is loaded
- **THEN** `Scenario.World.InitialLocation` SHALL equal the `name` of one of the 5 defined locations

#### Scenario: Location details contain no character names
- **WHEN** `simulations/vida-cotidiana/world.yaml` is read as raw text
- **THEN** the `details` fields SHALL NOT contain any of the character names defined in `characters.yaml`

### Requirement: Everyday lives runtime configuration

The `scenario.yaml` file SHALL set `turns` to 30 to allow enough simulation time for slow-burn relational dynamics to develop.

#### Scenario: Turn count override applied
- **WHEN** the scenario is loaded and no `--turns` CLI flag is provided
- **THEN** the simulation SHALL run for 30 turns as specified by `scenario.yaml`

#### Scenario: CLI flag overrides scenario.yaml turns
- **WHEN** the scenario is loaded and `--turns 15` is passed on the CLI
- **THEN** the simulation SHALL run for 15 turns, ignoring the `scenario.yaml` value of 30
