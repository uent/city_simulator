## ADDED Requirements

### Requirement: Character struct with persona fields
The system SHALL define a `Character` struct containing: unique ID, type, name, age, occupation, psychological core fields (motivation, fear, core belief, internal tension), formative events (slice of strings), voice profile, relational defaults, dialogue examples, runtime state (location, goals, emotional state), memory buffer, and an `Inbox` field (`[]world.Event`).

The `Inbox` field holds private events addressed to this character. It is populated by director actions (e.g., `reveal_information`) and is not persisted to YAML.

#### Scenario: Character loaded from YAML config
- **WHEN** the simulator reads a `characters.yaml` file with valid character entries
- **THEN** each entry SHALL be deserialized into a `Character` struct with all YAML fields populated and `Inbox` initialized as an empty (non-nil) slice

#### Scenario: Missing optional fields use defaults
- **WHEN** a character config omits optional fields (e.g., emotional state)
- **THEN** the system SHALL apply a default value and continue loading without error

### Requirement: Character inbox for private events
The system SHALL add an `Inbox []world.Event` field to `character.Character`. The field SHALL be omitted from YAML serialization (tag `yaml:"-"`). It SHALL be initialized as an empty non-nil slice by `LoadCharacters`.

When `internal/llm/prompt.go` builds a character's prompt, it SHALL:
1. Check `character.Inbox` for any pending items
2. If non-empty, append a "Private information you recently learned:" section to the prompt listing each inbox item's description
3. Clear `character.Inbox` (set to empty slice) after reading — flush-on-read semantics

#### Scenario: Inbox flushed after prompt build
- **WHEN** a character has one item in `Inbox` and their prompt is built
- **THEN** the prompt SHALL contain the inbox item's description AND `character.Inbox` SHALL be empty after the call

#### Scenario: Empty inbox produces no private section
- **WHEN** a character's `Inbox` is empty
- **THEN** the built prompt SHALL NOT contain a "Private information" section

#### Scenario: Inbox not persisted to YAML
- **WHEN** a `Character` with a non-empty `Inbox` is marshaled to YAML
- **THEN** the output SHALL NOT contain an `inbox` key

### Requirement: Character type field
The system SHALL add a `Type string` field (YAML tag `type`) to the `Character` struct. When the field is blank or absent, the character is treated as a regular character. The value `"game_director"` designates the entry as a Game Director.

#### Scenario: Type field absent from YAML
- **WHEN** a character entry in `characters.yaml` has no `type` key
- **THEN** the loaded `Character.Type` SHALL be an empty string and the character SHALL be treated as a regular character

#### Scenario: Type field set to game_director
- **WHEN** a character entry has `type: game_director`
- **THEN** the loaded `Character.Type` SHALL equal `"game_director"`

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
