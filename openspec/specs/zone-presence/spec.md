## Requirements

### Requirement: Zone roster computation
The system SHALL provide a `BuildZoneRoster(chars []*character.Character) map[string][]string` function in `internal/character/` that returns a map from location name to the list of character **names** (not IDs) currently at that location. Characters with an empty `Location` field SHALL be omitted.

#### Scenario: Multiple characters at the same location
- **WHEN** characters Alice and Bob both have `Location == "Market"`
- **THEN** `BuildZoneRoster` SHALL return a map entry `"Market": ["Alice", "Bob"]`

#### Scenario: Character with empty location is omitted
- **WHEN** a character has `Location == ""`
- **THEN** that character SHALL NOT appear in any roster entry

#### Scenario: Each location appears at most once as a key
- **WHEN** three characters occupy two different locations
- **THEN** the returned map SHALL have exactly two keys

### Requirement: Zone context prompt section
The system SHALL provide a `BuildZoneContext(roster map[string][]string) string` function in `internal/character/` that renders the roster as a human-readable block suitable for appending to LLM system prompts. The output SHALL list every location with its occupants; an empty roster SHALL return an empty string.

#### Scenario: Non-empty roster renders all zones
- **WHEN** roster has two entries: `"Park": ["Carlos"]` and `"Market": ["Alice", "Bob"]`
- **THEN** `BuildZoneContext` SHALL return a string containing "Park", "Carlos", "Market", "Alice", and "Bob"

#### Scenario: Empty roster returns empty string
- **WHEN** roster is an empty map
- **THEN** `BuildZoneContext` SHALL return `""`

### Requirement: Engine computes roster once per tick before message dispatch
The simulation engine SHALL compute the zone roster by calling `BuildZoneRoster` once per tick, before any `MoveDecision` or `CharChat` messages are dispatched in that tick. The same snapshot SHALL be used for all messages within the tick.

#### Scenario: Roster snapshot is consistent across a tick
- **WHEN** the engine dispatches multiple MoveDecision messages in one tick
- **THEN** all those messages SHALL be built using the same zone roster snapshot computed at the start of the tick

### Requirement: Movement prompt includes full zone roster
`BuildMovementPrompt` SHALL accept a `zoneRoster map[string][]string` parameter and include a "Who is where" section in the rendered prompt listing all zones and their occupants. The character's own name SHALL still appear at their current location (they are present).

#### Scenario: Movement prompt shows other characters' locations
- **WHEN** `BuildMovementPrompt` is called with a roster where location "Bar" has ["Maria", "Luis"]
- **THEN** the returned prompt string SHALL contain "Bar" and both "Maria" and "Luis"

#### Scenario: Empty roster renders no zone section
- **WHEN** `BuildMovementPrompt` is called with an empty roster
- **THEN** the returned prompt SHALL NOT contain a "Who is where" section

### Requirement: CharChat system prompt includes zone presence for the character's location
When building `CharChatPayload.InitiatorSystem` and `CharChatPayload.ResponderSystem`, the engine SHALL append a zone-presence block (via `BuildZoneContext`) that shows who is at each location. The character's own name SHALL be excluded from the listing of their current location so the block reads as "who else is here".

#### Scenario: System prompt mentions co-located characters
- **WHEN** initiator and a third character Charlie are both at "Plaza" and the engine builds the initiator's system prompt
- **THEN** `InitiatorSystem` SHALL contain "Charlie" in the zone-presence block

#### Scenario: Character's own name is excluded from their location listing
- **WHEN** the roster at "Plaza" is ["Alice", "Charlie"] and Alice is the initiator
- **THEN** `InitiatorSystem` SHALL NOT list "Alice" under "Plaza" in the zone-presence block
