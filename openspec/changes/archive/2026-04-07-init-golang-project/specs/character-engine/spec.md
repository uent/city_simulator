## ADDED Requirements

### Requirement: Character struct with persona fields
The system SHALL define a `Character` struct containing: unique ID, name, age, occupation, personality traits (slice of strings), backstory (string), goals (slice of strings), and current emotional state (string).

#### Scenario: Character loaded from YAML config
- **WHEN** the simulator reads a `characters.yaml` file with valid character entries
- **THEN** each entry SHALL be deserialized into a `Character` struct with all fields populated and no error returned

#### Scenario: Missing optional fields use defaults
- **WHEN** a character config omits optional fields (e.g., emotional state)
- **THEN** the system SHALL apply a default value (e.g., emotional state defaults to "neutral") and continue loading without error

### Requirement: Character loader from YAML file
The system SHALL provide a `LoadCharacters(path string) ([]Character, error)` function that reads a YAML file and returns all defined characters.

#### Scenario: File not found
- **WHEN** the provided YAML path does not exist
- **THEN** the function SHALL return a non-nil error with a descriptive message including the path

#### Scenario: Malformed YAML
- **WHEN** the YAML file contains a syntax error
- **THEN** the function SHALL return a non-nil error describing the parse failure

#### Scenario: Valid file with multiple characters
- **WHEN** the YAML file contains two or more character entries
- **THEN** the function SHALL return a slice with the same number of `Character` values and a nil error

### Requirement: Per-character memory buffer
The system SHALL maintain a sliding memory buffer per character, capped at a configurable `MaxMemory` integer (default 20 entries). Each `MemoryEntry` records: speaker name, message text, and tick number.

#### Scenario: Memory under capacity
- **WHEN** fewer entries than `MaxMemory` have been added
- **THEN** all entries SHALL be retrievable in insertion order

#### Scenario: Memory at capacity receives new entry
- **WHEN** a new entry is added and the buffer is already at `MaxMemory`
- **THEN** the oldest entry SHALL be evicted and the new entry appended, keeping total count at `MaxMemory`

#### Scenario: Retrieve memory as message slice
- **WHEN** `character.RecentMemory(n int)` is called
- **THEN** the system SHALL return up to the last `n` entries in chronological order
