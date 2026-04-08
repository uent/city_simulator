## Requirements

### Requirement: Character inventory field
Each character in `characters.yaml` MAY define an `inventory` field as an ordered list of strings describing objects the character carries at the start of the simulation. The field is optional; omitting it SHALL leave the character with an empty inventory.

#### Scenario: Character with inventory defined
- **WHEN** a character entry in `characters.yaml` contains an `inventory` list with one or more items
- **THEN** the scenario loader SHALL populate the character's inventory with those items in YAML order

#### Scenario: Character without inventory defined
- **WHEN** a character entry in `characters.yaml` omits the `inventory` key
- **THEN** the character's inventory SHALL be empty and loading SHALL return no error

### Requirement: Character initial_state field
Each character in `characters.yaml` MAY define an `initial_state` field as a string describing the character's tactical or narrative state at the start of the simulation. The field is optional; omitting it SHALL leave `initial_state` as an empty string.

#### Scenario: Character with initial_state defined
- **WHEN** a character entry contains `initial_state: "posando como sommelier, distrayendo a los guardias"`
- **THEN** the scenario loader SHALL expose that string as the character's initial state

#### Scenario: Character without initial_state defined
- **WHEN** a character entry omits the `initial_state` key
- **THEN** the character's initial state SHALL be an empty string and loading SHALL return no error
