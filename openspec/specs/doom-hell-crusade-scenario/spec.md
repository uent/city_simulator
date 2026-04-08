# DOOM Hell Crusade Scenario

## Requirements

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

### Requirement: Doom Guy character sheet

The Doom Guy entry SHALL have non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension`. It SHALL NOT have a `cover_identity` block. It SHALL have a `FormativeEvents` list of exactly 2–3 entries, a `Voice` block with at least `formality`, `verbal_tics`, `response_length`, and `communication_style`, a `RelationalDefaults` block with `strangers`, `authority`, and `vulnerable`, a `DialogueExamples` list of 3–4 entries reflecting sparse internal-monologue style, a `Goals` list of at least 3 entries, an `Inventory` list of at least 3 items, and a non-empty `InitialState` describing his condition at the moment of portal entry.

#### Scenario: Doom Guy present with required fields
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** exactly one character SHALL have a non-empty `Name` matching the marine protagonist, with non-empty `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension`

#### Scenario: Doom Guy has no cover identity
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** the Doom Guy character SHALL NOT have a `cover_identity` block

#### Scenario: Doom Guy has formative events
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `FormativeEvents` slice SHALL have length 2 or 3

#### Scenario: Doom Guy has voice block
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `Voice` block SHALL have non-empty `formality`, `verbal_tics`, `response_length`, and `communication_style`

#### Scenario: Doom Guy has relational defaults
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `RelationalDefaults` block SHALL have non-empty `strangers`, `authority`, and `vulnerable`

#### Scenario: Doom Guy has dialogue examples
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `DialogueExamples` slice SHALL have length 3 or 4, reflecting sparse internal-monologue style

#### Scenario: Doom Guy has goals
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `Goals` slice SHALL have at least 3 entries

#### Scenario: Doom Guy has inventory
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `Inventory` slice SHALL contain at least 3 items including the Praetor Suit and at least one ranged weapon

#### Scenario: Doom Guy has initial state
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Doom Guy's `InitialState` SHALL be non-empty describing his physical and mental condition at portal entry

### Requirement: Vael character sheet

The Vael entry SHALL have non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension`. It SHALL have a `FormativeEvents` list of exactly 2–3 entries, a `Voice` block with at least `formality`, `verbal_tics`, `response_length`, and `communication_style`, a `RelationalDefaults` block with `strangers`, `authority`, and `vulnerable`, a `DialogueExamples` list of 3–4 entries reflecting his negotiating style, a `Goals` list of at least 2 entries reflecting his personal agenda (not the marine's), and a non-empty `InitialState` describing his position and state when the marine arrives.

#### Scenario: Vael present with required fields
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** exactly one character SHALL represent a renegade demon with non-empty `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension` describing his fractured loyalty to Malphas

#### Scenario: Vael has formative events
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `FormativeEvents` slice SHALL have length 2 or 3

#### Scenario: Vael has voice block
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `Voice` block SHALL have non-empty `formality`, `verbal_tics`, `response_length`, and `communication_style`

#### Scenario: Vael has relational defaults
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `RelationalDefaults` block SHALL have non-empty `strangers`, `authority`, and `vulnerable`

#### Scenario: Vael has dialogue examples
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `DialogueExamples` slice SHALL have length 3 or 4

#### Scenario: Vael has independent goals
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `Goals` slice SHALL have at least 2 entries reflecting his own survival and political agenda, distinct from Doom Guy's mission

#### Scenario: Vael has initial state
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** Vael's `InitialState` SHALL be non-empty describing his position and concealment at The Flesh Gate before approaching the marine

### Requirement: DOOM Hell Crusade director character

The director entry SHALL have `type: game_director` and a `Goals` list of at least 4 entries covering: spawning a trapped human soul early, triggering Hellstone pulse events, introducing a renegade demon with an information cost, and revealing partial ritual cancellation conditions before turn 20.

#### Scenario: Director has narrative goals
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** the `game_director` character's `Goals` slice SHALL have at least 4 entries

#### Scenario: Director goals include early soul spawn
- **WHEN** `simulations/doom-hell-crusade/characters.yaml` is read as raw text
- **THEN** the director's goals SHALL reference spawning a human soul within the first few turns

### Requirement: DOOM Hell Crusade world layout

The `world.yaml` file SHALL define exactly 5 locations representing a Hell campaign across three narrative acts: The Flesh Gate (Act 1 entry), The Lava Wastes (Act 1 exploration), The Cathedral of Bone (Act 2 ascent), The Necropolis Vault (Act 2 approach), and The Throne Sanctum (Act 3 confrontation). Each location SHALL have a non-empty `name`, `description`, and `details`. Location `details` SHALL describe only physical and spatial properties — they SHALL NOT reference the name or current state of any specific character. The `initial_location` SHALL be set to `The Flesh Gate`. At least one `initial_event` SHALL be present to establish the ticking-clock threat.

#### Scenario: All five Hell locations present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.Locations` SHALL contain 5 entries, each with a non-empty `name` and `description`

#### Scenario: All locations have details
- **WHEN** the scenario is loaded
- **THEN** each of the 5 locations SHALL have a non-empty `details` field

#### Scenario: Location details contain no character names
- **WHEN** `simulations/doom-hell-crusade/world.yaml` is read as raw text
- **THEN** the `details` fields SHALL NOT contain any of the pre-loaded character names: Doom Slayer, Vael

#### Scenario: Initial location is The Flesh Gate
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialLocation` SHALL equal `"The Flesh Gate"`

#### Scenario: Initial event establishes threat
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialEvents` SHALL contain at least one event describing the Hellstone Convergence countdown

#### Scenario: World concept block present
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.World.Concept.Premise` SHALL be non-empty and `Scenario.World.Concept.Rules` SHALL contain at least one entry

#### Scenario: World has weather and atmosphere
- **WHEN** the doom-hell-crusade scenario is loaded
- **THEN** `Scenario.World` SHALL have non-empty `weather` and `atmosphere` fields

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
