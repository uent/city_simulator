## MODIFIED Requirements

### Requirement: Simulation engine initialization
The system SHALL provide a `NewEngine(cfg Config) (*Engine, error)` constructor. `Config` SHALL include: a `scenario.Scenario` (which provides characters and world config), LLM client, conversation manager reference, total tick count, random seed (int64, 0 means deterministic), and output writer.

The engine SHALL derive its character list from `cfg.Scenario.Characters` and initialize world state by calling `world.NewState(cfg.Scenario.World)`.

#### Scenario: Valid config provided
- **WHEN** `NewEngine` is called with a `Scenario` containing at least two characters and a positive tick count
- **THEN** the function SHALL return a non-nil `*Engine` ready to run and a nil error

#### Scenario: Scenario with fewer than two characters
- **WHEN** `NewEngine` is called with a `Scenario` that has zero or one character
- **THEN** the function SHALL return a nil engine and a non-nil error stating at least two characters are required
