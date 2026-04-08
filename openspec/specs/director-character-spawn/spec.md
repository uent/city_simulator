## Requirements

### Requirement: spawn_character action creates a full character
The system SHALL implement a `spawn_character` director action in `internal/director/actions_npc.go`. The action SHALL accept the following args: `id` (string, required), `name` (string, required), `age` (int, optional), `occupation` (string, optional), `motivation` (string, optional), `fear` (string, optional), `core_belief` (string, optional), `internal_tension` (string, optional), `formative_events` ([]string, optional), `location` (string, optional), `emotional_state` (string, optional, default `"neutral"`), `goals` ([]string, optional). The action SHALL construct a `*character.Character` with those fields, set `MaxMemory` to 20, initialize `Inbox` to an empty slice, and append it to `*chars`.

#### Scenario: spawn_character with full args creates complete character
- **WHEN** the director calls `spawn_character` with all fields populated
- **THEN** a new `*character.Character` SHALL be appended to the engine's character slice with all provided fields set

#### Scenario: spawn_character with only required args uses defaults
- **WHEN** the director calls `spawn_character` with only `id` and `name`
- **THEN** a new character SHALL be created with `EmotionalState = "neutral"`, `MaxMemory = 20`, and all optional fields at zero value

#### Scenario: spawn_character emits a public world event
- **WHEN** `spawn_character` executes successfully
- **THEN** a `world.Event` of type `"spawn"` SHALL be appended to `state.EventLog` with `Visibility = "public"` and a description naming the character

#### Scenario: spawn_character fails when id already exists
- **WHEN** `spawn_character` is called with an `id` that matches an existing character in `*chars`
- **THEN** `Execute` SHALL return a non-nil error and SHALL NOT mutate `state` or `*chars`

#### Scenario: spawn_character fails when id is missing
- **WHEN** `spawn_character` is called without the `id` arg
- **THEN** `Execute` SHALL return a non-nil error

#### Scenario: spawn_character fails when name is missing
- **WHEN** `spawn_character` is called without the `name` arg
- **THEN** `Execute` SHALL return a non-nil error

---

### Requirement: spawn_character is gated by character_spawn_rule
The `spawn_character` action SHALL be registered in the director's action registry at all times. However, it SHALL appear in the `<tools>` block of the director prompt ONLY when `state.Concept.CharacterSpawnRule` is non-empty. If the action is invoked at runtime and no rule is defined, `Execute` SHALL return a non-nil error without mutating state.

#### Scenario: spawn_character visible in prompt when rule is defined
- **WHEN** `state.Concept.CharacterSpawnRule` is non-empty
- **THEN** `BuildDirectorPrompt` SHALL include `spawn_character` in the `<tools>` block along with the rule text

#### Scenario: spawn_character absent from prompt when rule is not defined
- **WHEN** `state.Concept.CharacterSpawnRule` is empty
- **THEN** `BuildDirectorPrompt` SHALL NOT include `spawn_character` in the `<tools>` block

#### Scenario: spawn_character invoked without rule returns error
- **WHEN** `Execute` is called on `spawnCharacterAction` and `state.Concept.CharacterSpawnRule == ""`
- **THEN** `Execute` SHALL return a non-nil error containing the words "no character_spawn_rule defined"

---

### Requirement: spawn_character respects max_spawned_characters cap
When `state.Concept.MaxSpawnedCharacters > 0`, the `spawn_character` action SHALL refuse to create additional characters once `state.SpawnedCharacters >= MaxSpawnedCharacters`. The world state SHALL track the count of spawned characters in `State.SpawnedCharacters`, incremented on each successful spawn.

#### Scenario: spawn blocked when cap is reached
- **WHEN** `MaxSpawnedCharacters` is `3` and `SpawnedCharacters` is already `3`
- **THEN** `Execute` SHALL return a non-nil error and SHALL NOT mutate state or `*chars`

#### Scenario: spawn allowed when under cap
- **WHEN** `MaxSpawnedCharacters` is `3` and `SpawnedCharacters` is `2`
- **THEN** `Execute` SHALL succeed and `SpawnedCharacters` SHALL become `3`

#### Scenario: cap of zero means unlimited
- **WHEN** `MaxSpawnedCharacters` is `0`
- **THEN** the cap check SHALL be skipped and spawn SHALL not be blocked by count alone

---

### Requirement: Dynamically spawned characters participate as active actors
The simulation engine SHALL register newly spawned characters as active actors so they receive LLM turns and world events from the next processing cycle onward. The engine SHALL maintain a set of already-registered character IDs (`registeredChars`) to detect newly added characters after each director turn. For each unregistered character found, the engine SHALL create a `CharacterActor`, register it on the message bus, start it with the current context, and mark its ID as registered.

#### Scenario: Spawned character receives LLM turn after creation
- **WHEN** `spawn_character` adds a character in tick N
- **THEN** that character SHALL appear in the actors list and be eligible to act starting from tick N+1 (or tick N if the engine has not yet processed that character's slot)

#### Scenario: Already-registered characters are not re-registered
- **WHEN** `registerSpawnedChars` runs and all characters in `e.chars` are already in `registeredChars`
- **THEN** no new actors SHALL be created and the bus SHALL not be modified

#### Scenario: Engine logs spawn event
- **WHEN** `spawn_character` executes successfully
- **THEN** the engine SHALL log a message prefixed with `[spawn]` identifying the new character's id and name

---

### Requirement: Scheduler integrates spawned characters into pair rotation
The `Scheduler` SHALL expose an `AddCharacter(newChar *character.Character, known []*character.Character, locations []string)` method. When called, it SHALL create pairs between `newChar` and every character in `known`, append those pairs to the existing rotation, and assign `newChar` a random starting location from `locations` if its `Location` field is currently empty.

#### Scenario: New character paired with all existing characters
- **WHEN** `AddCharacter` is called with a new character and a slice of N existing characters
- **THEN** N new pairs SHALL be appended to the scheduler's pair list

#### Scenario: New character gets a random location when location is empty
- **WHEN** `AddCharacter` is called and `newChar.Location == ""`
- **THEN** `newChar.Location` SHALL be set to one of the provided locations

#### Scenario: New character keeps its location when already set
- **WHEN** `AddCharacter` is called and `newChar.Location` is already non-empty
- **THEN** `newChar.Location` SHALL NOT be overwritten
