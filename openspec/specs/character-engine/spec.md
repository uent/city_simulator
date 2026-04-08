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

---

### Requirement: CharChatReply payload fields
The system SHALL define `CharChatReply` in `internal/messaging/message.go` with the following fields:
- `InitiatorSpeech string` — what the initiator character said aloud
- `InitiatorAction string` — physical action performed by the initiator (may be empty)
- `ResponderSpeech string` — what the responder character said aloud
- `ResponderAction string` — physical action performed by the responder (may be empty)
- `Err error` — non-nil if the exchange failed

The previous `InitiatorText` and `ResponderText` fields SHALL be removed.

#### Scenario: Reply with action and speech
- **WHEN** the actor generates an initiator response with an action and speech
- **THEN** `CharChatReply.InitiatorAction` SHALL be non-empty and `CharChatReply.InitiatorSpeech` SHALL contain only the spoken text

#### Scenario: Reply with speech only
- **WHEN** the actor generates a response with no `*...*` markers
- **THEN** the corresponding `Action` field SHALL be an empty string and the `Speech` field SHALL contain the full response text

---

### Requirement: Engine renders action and speech separately
The simulation engine SHALL display character action and speech on separate lines per turn. The format SHALL be:

```
── Tick N ── InitiatorName [location] → ResponderName [location] ──
*initiator action*
InitiatorName: initiator speech
*responder action*
ResponderName: responder speech
```

Action lines are only printed when the `Action` field is non-empty. A missing action results in the action line being skipped entirely (no blank line placeholder).

#### Scenario: Tick with both action and speech for both characters
- **WHEN** `CharChatReply` has non-empty Action and Speech for both initiator and responder
- **THEN** the console SHALL print four lines: initiator action, initiator speech, responder action, responder speech

#### Scenario: Tick where one character has no action
- **WHEN** `CharChatReply.InitiatorAction` is empty but `ResponderAction` is non-empty
- **THEN** the console SHALL omit the initiator action line and still print the responder action line

#### Scenario: Tick where neither character has an action
- **WHEN** both `Action` fields are empty
- **THEN** the console SHALL print only the two speech lines (one per character), with no action lines

---

### Requirement: JSONL log entry includes action and speech fields
The `logEntry` struct written to `OutputWriter` SHALL include the following additional fields:
- `initiator_speech string` — spoken text from the initiator
- `initiator_action string` — action text from the initiator (empty string if none)
- `responder_speech string` — spoken text from the responder
- `responder_action string` — action text from the responder (empty string if none)

#### Scenario: Log entry with action present
- **WHEN** `CharChatReply.InitiatorAction` is `slams the table`
- **THEN** the JSONL line for that tick SHALL contain `"initiator_action":"slams the table"`

#### Scenario: Log entry with no action
- **WHEN** `CharChatReply.ResponderAction` is empty
- **THEN** the JSONL line SHALL contain `"responder_action":""` (empty string, not omitted)
