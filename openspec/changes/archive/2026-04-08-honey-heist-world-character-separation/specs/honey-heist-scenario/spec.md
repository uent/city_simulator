## MODIFIED Requirements

### Requirement: Honey Heist character roster
The `characters.yaml` file SHALL define exactly 6 bear characters, each with a unique `name`, an `occupation` describing their criminal specialty, a `motivation` string, a `fear` string, a `core_belief` string, an `internal_tension` string, a `formative_events` list of 2–3 causal bullets, a `voice` block with at least `formality` and `verbal_tics`, a `relational_defaults` block with `strangers`, `authority`, and `vulnerable`, a `dialogue_examples` list of 3–4 representative lines, and a `cover_identity` block with at minimum `alias` and `role`. Characters that carry objects relevant to their role SHALL define an `inventory` list of strings. Characters with a meaningful tactical state at the start of the simulation SHALL define an `initial_state` string. The legacy `personality` list and `backstory` prose fields SHALL NOT be present.

#### Scenario: All characters present with new schema
- **WHEN** the honey-heist scenario is loaded
- **THEN** `Scenario.Characters` SHALL contain 6 entries, each with non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension` fields

#### Scenario: All characters have a cover identity
- **WHEN** the honey-heist scenario is loaded
- **THEN** every character in `Scenario.Characters` SHALL have a non-nil `CoverIdentity` with non-empty `Alias` and `Role`

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

#### Scenario: Honeydrop has inventory and initial_state
- **WHEN** the honey-heist scenario is loaded
- **THEN** Honeydrop's `Inventory` SHALL include at least `bypass kit` and `peso sustituto`, and `InitialState` SHALL be non-empty

#### Scenario: Patches has initial_state
- **WHEN** the honey-heist scenario is loaded
- **THEN** Patches' `InitialState` SHALL be non-empty describing his getaway readiness

### Requirement: Honey Heist world layout
The `world.yaml` file SHALL define exactly 6 locations representing the HoneyCon convention centre and its surroundings: Convention Lobby, Vendor Hall, Security Office, Vault Antechamber, Vault, and Alley (Exit). Each location SHALL have a non-empty `name` and `description`. Location `details` fields SHALL describe only physical and spatial properties of that location — they SHALL NOT reference the name or current state of any specific character. At least one `initial_event` SHALL be present to seed pre-heist narrative context (scene setup before the operation begins). The `world.yaml` SHALL define `initial_location` pointing to a valid location name where all characters begin. The `world.yaml` SHALL also define a `concept:` block with a non-empty `premise`, at least one `rule`, and a `flavor` string.

#### Scenario: All locations present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.Locations` SHALL contain 6 entries, each with a non-empty `name` and `description`

#### Scenario: Initial events present
- **WHEN** the scenario is loaded
- **THEN** `Scenario.World.InitialEvents` SHALL contain at least one event with a non-empty `description`

#### Scenario: World concept block present
- **WHEN** the honey-heist scenario is loaded
- **THEN** `Scenario.World.Concept.Premise` SHALL be non-empty and `Scenario.World.Concept.Rules` SHALL contain at least one entry

#### Scenario: initial_location references an existing location
- **WHEN** the honey-heist scenario is loaded
- **THEN** `Scenario.World.InitialLocation` SHALL equal the `name` of one of the 6 defined locations

#### Scenario: Location details contain no character names
- **WHEN** `simulations/honey-heist/world.yaml` is read as raw text
- **THEN** the `details` fields SHALL NOT contain any of the character names: Grizwald, Honeydrop, Claws McGee, Lady Marmalade, Patches, Dr. Snuffles
