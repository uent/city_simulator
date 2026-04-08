## MODIFIED Requirements

### Requirement: Tool-call prompt format
The system SHALL provide `BuildDirectorPrompt(state *world.State, chars []*character.Character, tick int, language string) string` in `internal/director/prompt.go`. The prompt SHALL include:
- Current tick, time-of-day, weather, atmosphere, tension level
- All locations with names and descriptions
- All character names, IDs, and current locations
- A `<tools>` block listing every registered action with its parameter schema
- The `spawn_character` tool SHALL appear in the `<tools>` block ONLY when `state.Concept.CharacterSpawnRule` is non-empty; when included, the tool description SHALL show the rule text
- Instructions to respond with a `<tool_calls>` JSON array of `{"name": "...", "args": {...}}` objects
- If language is non-empty, a "Respond in <language>." instruction appended at the end

#### Scenario: Prompt includes tool schema
- **WHEN** `BuildDirectorPrompt` is called
- **THEN** the returned string SHALL contain a `<tools>` block with at least the names of all registered actions (excluding `spawn_character` when no spawn rule is defined)

#### Scenario: Prompt includes world state fields
- **WHEN** `BuildDirectorPrompt` is called and `state.Weather` is `"rain"`
- **THEN** the returned string SHALL contain the word `"rain"`

#### Scenario: spawn_character tool included when rule is set
- **WHEN** `state.Concept.CharacterSpawnRule` is `"All NPCs must be street vendors"`
- **THEN** the prompt SHALL contain `spawn_character` in the `<tools>` block and SHALL contain the rule text `"All NPCs must be street vendors"`

#### Scenario: spawn_character tool excluded when rule is empty
- **WHEN** `state.Concept.CharacterSpawnRule` is empty
- **THEN** the prompt SHALL NOT contain `spawn_character` anywhere in the `<tools>` block
