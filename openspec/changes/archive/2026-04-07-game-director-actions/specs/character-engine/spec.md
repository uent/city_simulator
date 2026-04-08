## MODIFIED Requirements

### Requirement: Character struct with persona fields
The system SHALL define a `Character` struct containing: unique ID, type, name, age, occupation, psychological core fields (motivation, fear, core belief, internal tension), formative events (slice of strings), voice profile, relational defaults, dialogue examples, runtime state (location, goals, emotional state), memory buffer, and an `Inbox` field (`[]world.Event`).

The `Inbox` field holds private events addressed to this character. It is populated by director actions (e.g., `reveal_information`) and is not persisted to YAML.

#### Scenario: Character loaded from YAML config
- **WHEN** the simulator reads a `characters.yaml` file with valid character entries
- **THEN** each entry SHALL be deserialized into a `Character` struct with all YAML fields populated and `Inbox` initialized as an empty (non-nil) slice

#### Scenario: Missing optional fields use defaults
- **WHEN** a character config omits optional fields (e.g., emotional state)
- **THEN** the system SHALL apply a default value and continue loading without error

## ADDED Requirements

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
