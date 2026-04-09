## MODIFIED Requirements

### Requirement: Simulation engine initialization
The system SHALL provide a `NewEngine(cfg Config) (*Engine, error)` constructor. `Config` SHALL include: a `scenario.Scenario` (which provides characters and world config), a `*messaging.MessageBus` with all character actors already registered, total tick count, random seed (int64, 0 means random), and output writer.

The engine SHALL derive its character list from `cfg.Scenario.Characters`, initialize world state by calling `world.NewState(cfg.Scenario.World)`, and start all registered actors via `bus.StartAll(ctx)` at the beginning of `Engine.Run`.

The engine SHALL NOT hold a reference to `conversation.Manager`. All inter-actor communication SHALL go through the `MessageBus`.

#### Scenario: Valid config provided
- **WHEN** `NewEngine` is called with a `Scenario` containing at least two characters and a positive tick count, and a `MessageBus` with those characters registered
- **THEN** the function SHALL return a non-nil `*Engine` ready to run and a nil error

#### Scenario: Scenario with fewer than two characters
- **WHEN** `NewEngine` is called with a `Scenario` that has zero or one regular character (excluding any GameDirector)
- **THEN** the function SHALL return a nil engine and a non-nil error stating at least two characters are required

### Requirement: World concept printed at run start
At the beginning of `Engine.Run`, before the first tick, the engine SHALL print the world concept block by delegating to the `simulation-premise-display` logic. If `Scenario.World.Concept.Premise` is empty, the engine SHALL skip this step silently.

#### Scenario: Concept block appears before first tick
- **WHEN** `Engine.Run` is called with a scenario whose `Concept.Premise` is non-empty
- **THEN** the concept block SHALL be written to stdout before any tick output (e.g., `[Tick 1]` lines)

#### Scenario: No concept block when premise is absent
- **WHEN** `Engine.Run` is called with a scenario whose `Concept.Premise` is empty
- **THEN** no concept block header SHALL appear in stdout and the simulation SHALL proceed normally
