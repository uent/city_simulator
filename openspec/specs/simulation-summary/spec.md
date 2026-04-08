## Requirements

### Requirement: Summary generation from world state
The system SHALL provide a `GenerateSummary(ctx context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario) (string, error)` function in a new `internal/summary` package that constructs a prompt from world events and character final states, sends it to the LLM, and returns the narrative text.

The prompt SHALL include:
- The scenario name
- The world concept (premise, flavor, and rules) when present in the scenario
- The world atmosphere and weather when present
- The total number of ticks that ran
- All world events (up to the last 200 if the list exceeds 200 entries)
- Each character's name, role, location, motivation, fear, goals, and final emotional state (omitting fields that are empty)

The LLM system prompt SHALL instruct the model to produce a rich narrative of at least six detailed paragraphs, covering the arc of the simulation, key turning points, character development, and final outcome.

#### Scenario: Successful summary generation
- **WHEN** `GenerateSummary` is called after a completed simulation with events and characters present
- **THEN** the function SHALL return a non-empty narrative string and a nil error

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

### Requirement: Timestamped summary file persistence
The system SHALL provide a `SaveSummary(scenarioName string, content string) (string, error)` function that writes the summary to `simulations/<scenarioName>/summary-<timestamp>.md`, where `<timestamp>` is formatted as RFC3339 with colons replaced by hyphens.

The function SHALL:
- Create the target directory if it does not exist
- Return the absolute path of the written file on success

#### Scenario: Summary saved to new file
- **WHEN** `SaveSummary` is called with a non-empty scenario name and content
- **THEN** a file named `summary-<timestamp>.md` SHALL be created inside `simulations/<scenarioName>/` and the function SHALL return its path with a nil error

#### Scenario: Multiple runs do not overwrite each other
- **WHEN** `SaveSummary` is called twice in the same second or across different runs
- **THEN** each call SHALL produce a distinct file (timestamp includes seconds; concurrent calls within the same second are acceptable to produce separate files via unique nanosecond suffix fallback if needed)

#### Scenario: Scenario directory does not exist
- **WHEN** `SaveSummary` is called for a scenario whose directory does not yet exist
- **THEN** the function SHALL create the directory and write the file without returning an error
