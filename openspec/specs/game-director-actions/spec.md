## Requirements

### Requirement: Director action interface
The system SHALL define an `Action` interface in `internal/director/action.go` with three methods: `Name() string`, `Execute(args map[string]any, state *world.State, chars *[]*character.Character) error`, and `Summary(args map[string]any) string`. Every named director action SHALL implement this interface.

#### Scenario: Action returns its name
- **WHEN** `Name()` is called on any registered action
- **THEN** it SHALL return the exact string used to register it in the registry

#### Scenario: Action mutates world state
- **WHEN** `Execute` is called with valid args and a non-nil `*world.State`
- **THEN** the action SHALL mutate `state` in-place and return nil

#### Scenario: Action receives invalid args
- **WHEN** `Execute` is called with args that fail validation (e.g., missing required key, out-of-range value)
- **THEN** the action SHALL return a non-nil error without mutating `state`

#### Scenario: Summary returns name and key arg
- **WHEN** `Summary` is called on `set_weather` with `{"type": "storm"}`
- **THEN** it SHALL return `"set_weather: storm"`

#### Scenario: Summary falls back to name when arg missing
- **WHEN** `Summary` is called with an empty args map
- **THEN** it SHALL return the action name alone (no panic)

---

### Requirement: Director action registry
The system SHALL provide a `Registry` in `internal/director/registry.go` that maps action names to `Action` implementations. The registry SHALL expose a `Dispatch(name string, args map[string]any, state *world.State, chars []*character.Character) error` method.

#### Scenario: Dispatch known action
- **WHEN** `Dispatch` is called with a name that matches a registered action
- **THEN** the corresponding `Action.Execute` SHALL be called and its return value returned

#### Scenario: Dispatch unknown action
- **WHEN** `Dispatch` is called with a name not in the registry
- **THEN** it SHALL return a non-nil error containing the unknown name; state SHALL NOT be mutated

---

### Requirement: Tool-call prompt format
The system SHALL provide `BuildDirectorPrompt(state *world.State, chars []*character.Character, tick int) string` in `internal/director/prompt.go`. The prompt SHALL include:
- Current tick, time-of-day, weather, atmosphere, tension level
- All locations with names and descriptions
- All character names, IDs, and current locations
- A `<tools>` block listing every registered action with its parameter schema
- Instructions to respond with a `<tool_calls>` JSON array of `{"name": "...", "args": {...}}` objects

#### Scenario: Prompt includes tool schema
- **WHEN** `BuildDirectorPrompt` is called
- **THEN** the returned string SHALL contain a `<tools>` block with at least the names of all 13 registered actions

#### Scenario: Prompt includes world state fields
- **WHEN** `BuildDirectorPrompt` is called and `state.Weather` is `"rain"`
- **THEN** the returned string SHALL contain the word `"rain"`

---

### Requirement: Tool-call response parser
The system SHALL provide `ParseToolCalls(raw string) ([]ToolCall, error)` in `internal/director/parse.go`. A `ToolCall` SHALL have `Name string` and `Args map[string]any`. The parser SHALL:
- Scan for the first `<tool_calls>` and `</tool_calls>` tags
- Unmarshal the content as a JSON array of tool-call objects
- Return an empty slice (not an error) if no `<tool_calls>` block is found
- Skip individual entries that are missing the `name` field

#### Scenario: Valid tool_calls block parsed
- **WHEN** the raw string contains a `<tool_calls>[{"name":"set_weather","args":{"type":"rain"}}]</tool_calls>` block
- **THEN** `ParseToolCalls` SHALL return a slice of length 1 with `Name="set_weather"` and `Args={"type":"rain"}`

#### Scenario: No tool_calls block present
- **WHEN** the raw string contains no `<tool_calls>` tag
- **THEN** `ParseToolCalls` SHALL return an empty slice and nil error

#### Scenario: Malformed JSON inside block
- **WHEN** the `<tool_calls>` block contains invalid JSON
- **THEN** `ParseToolCalls` SHALL return an empty slice and nil error (lenient parse)

---

### Requirement: Environment actions
The system SHALL register the following environment action handlers:

- **`set_weather`** — sets `state.Weather` to the provided `type` string (e.g., `"rain"`, `"clear"`, `"fog"`). Required arg: `type string`.
- **`set_time_of_day`** — sets `state.TimeOfDay` to the provided `moment` string (e.g., `"dawn"`, `"noon"`, `"midnight"`). Required arg: `moment string`.
- **`set_atmosphere`** — sets `state.Atmosphere` to the provided `descriptor` string (e.g., `"tense"`, `"calm"`, `"oppressive"`). Required arg: `descriptor string`.

Each action SHALL also append a public event to the world log describing the change.

#### Scenario: set_weather with valid type
- **WHEN** `set_weather` is dispatched with `{"type": "storm"}`
- **THEN** `state.Weather` SHALL equal `"storm"` and a public event describing the weather change SHALL be appended to `state.EventLog`

#### Scenario: set_atmosphere with valid descriptor
- **WHEN** `set_atmosphere` is dispatched with `{"descriptor": "tense"}`
- **THEN** `state.Atmosphere` SHALL equal `"tense"` and a public event SHALL be appended

#### Scenario: set_time_of_day with valid moment
- **WHEN** `set_time_of_day` is dispatched with `{"moment": "midnight"}`
- **THEN** `state.TimeOfDay` SHALL equal `"midnight"` and a public event SHALL be appended

---

### Requirement: NPC actions
The system SHALL register the following NPC action handlers:

- **`move_npc`** — sets the `Location` field of the character matching `id` to `destination`. Required args: `id string`, `destination string`. Optional arg: `reason string` (used in the appended event description).
- **`introduce_npc`** — appends a new `character.Character` to the `chars` slice (passed by pointer indirection via the engine). Required args: `id string`, `name string`, `role string`, `attitude string`, `motivation string`. Optional: `location string`.
- **`add_npc_condition`** — appends `condition` to the `EmotionalState` of the character matching `id`. Required args: `id string`, `condition string`.
- **`remove_npc_condition`** — removes the first occurrence of `condition` from the `EmotionalState` of the character matching `id`. Required args: `id string`, `condition string`.

#### Scenario: move_npc with valid id and destination
- **WHEN** `move_npc` is dispatched with `{"id": "alice", "destination": "market"}`
- **THEN** the character with ID `"alice"` SHALL have `Location == "market"` and a public event SHALL be appended describing the move

#### Scenario: move_npc with unknown id
- **WHEN** `move_npc` is dispatched with an `id` that does not match any character
- **THEN** the action SHALL return a non-nil error and no character SHALL be mutated

#### Scenario: introduce_npc creates new character
- **WHEN** `introduce_npc` is dispatched with valid required args
- **THEN** a new character entry SHALL be appended and a public event SHALL announce their arrival

#### Scenario: add_npc_condition appends to emotional state
- **WHEN** `add_npc_condition` is dispatched with `{"id": "bob", "condition": "frightened"}`
- **THEN** the character `bob`'s `EmotionalState` SHALL contain `"frightened"`

---

### Requirement: World and event actions
The system SHALL register the following world/event action handlers:

- **`modify_location`** — updates the `Description` or `Details` of the location matching `name`. Required arg: `name string`. Optional args: `description string`, `details string`.
- **`trigger_encounter`** — appends a public event of type `"encounter"` describing an interaction between the listed characters. Required args: `participants []string`, `context string`.
- **`trigger_event`** — appends a public event with the provided `type` and `description`. Required args: `type string`, `description string`. Optional: `location string`, `severity int` (1–10; if >5, also increases `state.Tension` by 1, clamped to 10).
- **`reveal_information`** — appends a private event visible only to `recipient`. Required args: `recipient string` (character ID), `content string`. The event SHALL have `PrivateRecipient == recipient` and be added to `recipient`'s `Inbox`.
- **`escalate_tension`** — adjusts `state.Tension` by `delta` (positive or negative), clamped to [0, 10]. Required arg: `delta int`.
- **`narrate`** — appends a public event of type `"narration"` with the provided `text`. Required arg: `text string`. Does not mutate any other state field.

#### Scenario: trigger_event with severity > 5 increases tension
- **WHEN** `trigger_event` is dispatched with `{"type": "explosion", "description": "A building collapses", "severity": 8}`
- **THEN** a public event SHALL be appended AND `state.Tension` SHALL increase by 1 (clamped to 10)

#### Scenario: reveal_information reaches only recipient's inbox
- **WHEN** `reveal_information` is dispatched with `{"recipient": "charlie", "content": "The vault is empty"}`
- **THEN** a private event with `PrivateRecipient == "charlie"` SHALL be appended to the world log AND appended to character `charlie`'s `Inbox`; no other character's `Inbox` SHALL be modified

#### Scenario: escalate_tension clamped at 10
- **WHEN** `state.Tension == 9` and `escalate_tension` is dispatched with `{"delta": 5}`
- **THEN** `state.Tension` SHALL equal `10` (not 14)

#### Scenario: escalate_tension clamped at 0
- **WHEN** `state.Tension == 1` and `escalate_tension` is dispatched with `{"delta": -5}`
- **THEN** `state.Tension` SHALL equal `0` (not -4)

#### Scenario: narrate appends public event without side effects
- **WHEN** `narrate` is dispatched with `{"text": "The crowd grows silent"}`
- **THEN** a public event of type `"narration"` SHALL be appended and no other state field SHALL change
