## Requirements

### Requirement: Summary generation from world state
The system SHALL provide a `GenerateSummary(ctx context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario) (string, error)` function in a new `internal/summary` package that constructs a prompt from world events and character final states, sends it to the LLM, and returns the narrative text.

The prompt SHALL include:
- The scenario name and description
- The total number of ticks that ran
- All world events (up to the last 100 if the list exceeds 100 entries)
- Each character's name, role, and final emotional state

#### Scenario: Successful summary generation
- **WHEN** `GenerateSummary` is called after a completed simulation with events and characters present
- **THEN** the function SHALL return a non-empty narrative string and a nil error

#### Scenario: LLM call fails
- **WHEN** the LLM client returns an error during summary generation
- **THEN** the function SHALL return an empty string and a non-nil error describing the failure

#### Scenario: Event list exceeds 100 entries
- **WHEN** the world state contains more than 100 events
- **THEN** the function SHALL include only the last 100 events in the prompt without error

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
