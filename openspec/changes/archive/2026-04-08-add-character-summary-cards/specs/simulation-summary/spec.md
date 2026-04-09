## MODIFIED Requirements

### Requirement: Summary generation from world state
The system SHALL provide a `GenerateSummary(ctx context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario, language string) (string, error)` function in the `internal/summary` package that constructs a prompt from world events and character final states, sends it to the LLM, and returns the complete summary string.

The returned string SHALL consist of:
1. The LLM-generated narrative text (unchanged behavior).
2. The character cards block produced by `renderCharacterCards(chars)` appended after the narrative. If there are no non-director characters, no cards block is appended.

The prompt SHALL include:
- The scenario name
- The world concept (premise, flavor, and rules) when present in the scenario
- The world atmosphere and weather when present
- The total number of ticks that ran
- All world events (up to the last 200 if the list exceeds 200 entries)
- Each character's name, role, location, motivation, fear, goals, and final emotional state (omitting fields that are empty)

The LLM system prompt SHALL instruct the model to produce a rich narrative of at least six detailed paragraphs, covering the arc of the simulation, key turning points, character development, and final outcome.

#### Scenario: Successful summary generation includes character cards
- **WHEN** `GenerateSummary` is called after a completed simulation with events and non-director characters present
- **THEN** the function SHALL return a non-empty string containing the narrative followed by the character cards block, and a nil error

#### Scenario: LLM call fails
- **WHEN** the LLM client returns an error during summary generation
- **THEN** the function SHALL return an empty string and a non-nil error describing the failure

#### Scenario: Event list exceeds 200 entries
- **WHEN** the world state contains more than 200 events
- **THEN** the function SHALL include only the last 200 events in the prompt without error

#### Scenario: World concept is absent
- **WHEN** the scenario world config has an empty `Concept` block
- **THEN** the function SHALL omit the concept section from the prompt without error

#### Scenario: Character has no goals or motivation
- **WHEN** a character has empty `Goals`, `Motivation`, or `Fear` fields
- **THEN** those fields SHALL be omitted from that character's entry in the prompt

#### Scenario: Only game director characters present
- **WHEN** all characters are of type `game_director`
- **THEN** the function SHALL return the narrative text with no character cards block appended
