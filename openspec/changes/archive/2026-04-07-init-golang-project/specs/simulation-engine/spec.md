## ADDED Requirements

### Requirement: Simulation engine initialization
The system SHALL provide a `NewEngine(cfg Config) *Engine` constructor. `Config` SHALL include: list of characters, LLM client, world state reference, conversation manager reference, total tick count, and random seed (int64, 0 means deterministic).

#### Scenario: Valid config provided
- **WHEN** `NewEngine` is called with a non-empty character list and positive tick count
- **THEN** the function SHALL return a non-nil `*Engine` ready to run

#### Scenario: Empty character list
- **WHEN** `NewEngine` is called with zero characters
- **THEN** the function SHALL return an error stating at least two characters are required

### Requirement: Simulation run loop
The system SHALL provide a `Run(ctx context.Context) error` method that executes the configured number of ticks. Each tick: selects an interacting pair via the scheduler, runs one exchange (initiator speaks, responder replies), appends the exchange to both characters' memories, advances world time by one tick, and writes the exchange to the output log.

#### Scenario: Normal run completion
- **WHEN** `Run` completes all configured ticks without error
- **THEN** the method SHALL return nil and the output log SHALL contain one entry per tick

#### Scenario: Context cancellation
- **WHEN** the provided context is cancelled mid-run
- **THEN** `Run` SHALL stop after the current tick completes and return `ctx.Err()`

#### Scenario: LLM error during a tick
- **WHEN** the LLM client returns an error for a chat request
- **THEN** `Run` SHALL log the error, skip that tick's exchange, and continue to the next tick

### Requirement: Interaction scheduler
The system SHALL provide a `Scheduler` that, given a list of characters, returns the next `(initiator, responder)` pair. Default mode is round-robin over all unique pairs. When a non-zero seed is provided, pairs SHALL be shuffled using that seed.

#### Scenario: Round-robin ordering
- **WHEN** three characters A, B, C are loaded with seed 0
- **THEN** successive calls to `Next()` SHALL return pairs in a consistent, repeating order covering all AB, AC, BC combinations

#### Scenario: Seeded shuffle
- **WHEN** a non-zero seed is provided
- **THEN** the pair order SHALL be shuffled and the same seed SHALL produce the same shuffled order across runs
