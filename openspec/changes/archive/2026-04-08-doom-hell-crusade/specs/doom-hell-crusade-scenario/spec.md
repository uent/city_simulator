# DOOM Hell Crusade Scenario

## ADDED Requirements

### Requirement: DOOM Hell Crusade scenario directory

The system SHALL provide a `simulations/doom-hell-crusade/` directory containing `characters.yaml`, `world.yaml`, and `scenario.yaml` that the scenario loader can load without errors.

#### Scenario: Scenario loads successfully
- **WHEN** `--scenario doom-hell-crusade` is passed to the simulator CLI
- **THEN** the scenario loader SHALL resolve `simulations/doom-hell-crusade/`, read all three YAML files, and return a valid `Scenario` struct with no errors

#### Scenario: Scenario directory is missing a required file
- **WHEN** `characters.yaml` or `world.yaml` is absent from `simulations/doom-hell-crusade/`
- **THEN** the scenario loader SHALL return a non-nil error describing which file is missing

### Requirement: DOOM Hell Crusade character roster

The `characters.yaml` file SHALL define exactly 3 characters: one `game_director` character (the Watcher), one protagonist character (Doom Guy), and one pre-loaded renegade demon (Vael) who is present at The Flesh Gate at the start of the simulation. All other antagonists and neutral parties are spawned at runtime by the director.

#### Scenario: Exactly two regular characters pre-loaded
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.Characters` SHALL contain exactly 2 entries (Doom Guy and Vael)

#### Scenario: Director character present
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** exactly one character SHALL have `type: game_director`

#### Scenario: Doom Guy character present
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** exactly one character SHALL have a non-empty `Name` matching the marine protagonist, with non-empty `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension`

#### Scenario: Doom Guy has no cover identity
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** the Doom Guy character SHALL NOT have a `cover_identity` block

#### Scenario: Doom Guy has dialogue examples
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `DialogueExamples` slice SHALL have length 3 or 4, reflecting sparse internal-monologue style

#### Scenario: Vael renegade demon present
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** exactly one character SHALL represent a renegade demon with non-empty `Occupation`, `Motivation`, and `InternalTension` describing his fractured loyalty to Malphas

### Requirement: DOOM Hell Crusade world layout

The `world.yaml` file SHALL define exactly 5 locations representing a Hell campaign across three narrative acts: The Flesh Gate (Act 1 entry), The Lava Wastes (Act 1 exploration), The Cathedral of Bone (Act 2 ascent), The Necropolis Vault (Act 2 approach), and The Throne Sanctum (Act 3 confrontation). Each location SHALL have a non-empty `name` and `description`. The `initial_location` SHALL be set to `The Flesh Gate`. At least one `initial_event` SHALL be present to establish the ticking-clock threat.

#### Scenario: All five Hell locations present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.Locations` SHALL contain 5 entries, each with a non-empty `name` and `description`

#### Scenario: Initial location is The Flesh Gate
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialLocation` SHALL equal `"The Flesh Gate"`

#### Scenario: Initial event establishes threat
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialEvents` SHALL contain at least one event describing the Hellstone Convergence countdown

#### Scenario: World concept block present
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.World.Concept.Premise` SHALL be non-empty and `Scenario.World.Concept.Rules` SHALL contain at least one entry

### Requirement: DOOM Hell Crusade director spawn config

The `world.yaml` concept block SHALL define a `character_spawn_rule` that restricts spawned characters to three archetypes: demon renegade, trapped human soul, and herald of the Prince. `max_spawned_characters` SHALL be set to `4`.

#### Scenario: Spawn rule restricts to three archetypes
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.World.Concept.CharacterSpawnRule` SHALL be a non-empty string describing the three valid archetypes

#### Scenario: Max spawned characters is 4
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.World.Concept.MaxSpawnedCharacters` SHALL equal `4`

### Requirement: DOOM Hell Crusade runtime overrides

The `scenario.yaml` file SHALL set `turns` to 30 to give the narrative arc room to develop across three acts.

#### Scenario: Turn count override applied
- **WHEN** the scenario is loaded and no `--turns` CLI flag is provided
- **THEN** the simulation SHALL run for 30 turns as specified by `scenario.yaml`

#### Scenario: CLI flag overrides scenario.yaml turns
- **WHEN** the scenario is loaded and `--turns 50` is passed on the CLI
- **THEN** the simulation SHALL run for 50 turns, ignoring the `scenario.yaml` value of 30
