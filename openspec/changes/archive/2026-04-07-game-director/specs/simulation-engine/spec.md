## ADDED Requirements

### Requirement: Game Director tick invocation
The system SHALL invoke the Game Director at the beginning of each tick, before scheduling any character exchange, when `Engine.director` is non-nil.

The invocation SHALL:
1. Call `BuildDirectorPrompt` with the current world state, regular character list, and tick number
2. Send the prompt to the LLM client (using the same model as regular characters)
3. Parse the response with `ParseDirectorEvents`
4. Set `Event.Tick` to the current tick on each parsed event
5. Append each parsed event to the world state via `AppendEvent`

If the LLM call or JSON parsing fails, the error SHALL be logged and the tick SHALL continue without any director events (fail-open, not fail-stop).

#### Scenario: Director generates events before characters act
- **WHEN** a tick begins with a non-nil Game Director
- **THEN** the director's events SHALL be appended to the world state before `RunExchange` is called for that tick, so characters' `PublicSummary` already includes the new events

#### Scenario: LLM call fails during director turn
- **WHEN** the LLM returns an error during the director's turn
- **THEN** the error SHALL be logged, no events SHALL be added for that tick, and the simulation SHALL continue to the character exchange step normally

#### Scenario: Director generates zero events
- **WHEN** the director responds with an empty JSON array `[]`
- **THEN** no events are appended and the tick proceeds normally without any error

## MODIFIED Requirements

### Requirement: Simulation engine initialization
The system SHALL provide a `NewEngine(cfg Config) (*Engine, error)` constructor. `Config` SHALL include: a `scenario.Scenario` (which provides characters, game director, and world config), LLM client, conversation manager reference, total tick count, random seed (int64, 0 means deterministic), and output writer.

The engine SHALL:
- Derive its regular character list from `cfg.Scenario.Characters` (entries where `Type != "game_director"`)
- Store `cfg.Scenario.GameDirector` as `engine.director` (may be nil if not defined)
- Initialize world state by calling `world.NewState(cfg.Scenario.World)`
- Require at least 2 regular characters (Game Director is excluded from this count)

#### Scenario: Valid config provided
- **WHEN** `NewEngine` is called with a `Scenario` containing at least two regular characters and a positive tick count
- **THEN** the function SHALL return a non-nil `*Engine` ready to run and a nil error

#### Scenario: Scenario with fewer than two regular characters
- **WHEN** `NewEngine` is called with a `Scenario` that has zero or one regular character (excluding Game Director)
- **THEN** the function SHALL return a nil engine and a non-nil error stating at least two characters are required

#### Scenario: Scenario with Game Director and two regular characters
- **WHEN** `NewEngine` is called with one Game Director and two regular characters
- **THEN** the function SHALL return a non-nil `*Engine` with `engine.director` set and the regular character count equal to 2
