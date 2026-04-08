## MODIFIED Requirements

### Requirement: WorldConcept struct in WorldConfig

The system SHALL define a `WorldConcept` struct with four fields: `Premise string` (YAML key `premise`) — a single sentence describing the fundamental nature of this world and the hidden truth characters must conceal; `Rules []string` (YAML key `rules`) — a list of constraints that define what is "normal" in this world and what would expose a character as out of place; `Flavor string` (YAML key `flavor`) — a short tone/mood string (e.g., "absurdist heist comedy"); and `CharacterSpawnRule string` (YAML key `character_spawn_rule`) — a rule describing how dynamically created characters must be designed (e.g., their archetype, tone, constraints). `WorldConfig` SHALL expose a `Concept WorldConcept` field (YAML key `concept`) and an `InitialLocation string` field (YAML key `initial_location`). All sub-fields are optional; omitting the entire `concept:` block SHALL leave `WorldConcept` at its zero value. Omitting `character_spawn_rule` SHALL leave it as an empty string, which disables the `spawn_character` director action.

#### Scenario: Full concept block parsed from world.yaml
- **WHEN** a `world.yaml` contains a `concept:` block with `premise`, `rules`, `flavor`, and `character_spawn_rule` set
- **THEN** all four fields SHALL be populated after loading

#### Scenario: Partial concept block accepted
- **WHEN** a `world.yaml` contains `concept: { premise: "Bears disguised as humans" }` with no other fields
- **THEN** `WorldConfig.Concept.Premise` SHALL be `"Bears disguised as humans"`, all other fields SHALL be zero value, and loading SHALL return no error

#### Scenario: Missing concept block defaults to zero value
- **WHEN** a `world.yaml` omits the `concept:` key entirely
- **THEN** `WorldConfig.Concept` SHALL equal the zero-value `WorldConcept{}` and loading SHALL return no error

#### Scenario: Rules parsed as ordered list
- **WHEN** `world.yaml` contains `concept.rules` with three entries
- **THEN** `WorldConfig.Concept.Rules` SHALL have length 3 with entries in YAML order

#### Scenario: character_spawn_rule parsed from world.yaml
- **WHEN** a `world.yaml` contains `concept.character_spawn_rule: "All characters must be bears in human disguise"`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL be that string after loading

#### Scenario: Missing character_spawn_rule defaults to empty string
- **WHEN** a `world.yaml` omits `character_spawn_rule` under `concept`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL be empty string and loading SHALL return no error

#### Scenario: initial_location parsed from world.yaml
- **WHEN** a `world.yaml` contains `initial_location: "Convention Lobby"`
- **THEN** `WorldConfig.InitialLocation` SHALL be `"Convention Lobby"` after loading

#### Scenario: Missing initial_location defaults to empty string
- **WHEN** a `world.yaml` omits the `initial_location` key
- **THEN** `WorldConfig.InitialLocation` SHALL be an empty string and loading SHALL return no error
