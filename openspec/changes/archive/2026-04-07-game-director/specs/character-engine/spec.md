## ADDED Requirements

### Requirement: Character type field
The system SHALL add a `Type string` field (YAML tag `type`) to the `Character` struct. When the field is blank or absent, the character is treated as a regular character. The value `"game_director"` designates the entry as a Game Director.

#### Scenario: Type field absent from YAML
- **WHEN** a character entry in `characters.yaml` has no `type` key
- **THEN** the loaded `Character.Type` SHALL be an empty string and the character SHALL be treated as a regular character

#### Scenario: Type field set to game_director
- **WHEN** a character entry has `type: game_director`
- **THEN** the loaded `Character.Type` SHALL equal `"game_director"`

## MODIFIED Requirements

### Requirement: Character loader from YAML file
The system SHALL provide a `LoadCharacters(path string) ([]Character, error)` function that reads a YAML file and returns all defined characters.

The function SHALL return ALL character entries regardless of their `Type` value. Filtering by type (separating Game Director from regular characters) is the responsibility of the scenario loader, not this function.

#### Scenario: File not found
- **WHEN** the provided YAML path does not exist
- **THEN** the function SHALL return a non-nil error with a descriptive message including the path

#### Scenario: Malformed YAML
- **WHEN** the YAML file contains a syntax error
- **THEN** the function SHALL return a non-nil error describing the parse failure

#### Scenario: Valid file with multiple characters
- **WHEN** the YAML file contains two or more character entries
- **THEN** the function SHALL return a slice with the same number of `Character` values and a nil error

#### Scenario: Mixed types returned unfiltered
- **WHEN** the YAML file contains one regular character and one `type: game_director` entry
- **THEN** `LoadCharacters` SHALL return a slice of length 2 containing both entries
