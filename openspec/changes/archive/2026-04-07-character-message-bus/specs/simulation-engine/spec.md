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

## ADDED Requirements

### Requirement: Engine tick loop uses MessageBus for all communication
In `Engine.Run`, each tick SHALL proceed as follows:
1. If a GameDirector actor is registered: send a `DirectorDirective` broadcast to all character actors via `bus.Broadcast`, drain all replies before proceeding
2. Select the next character pair via `Scheduler.Next()`
3. Send a `CharChat` message from the initiator actor to the responder actor via `bus.Send`; await the reply on `msg.ReplyChan`
4. The reply payload SHALL contain both the initiator's utterance and the responder's reply text (the responder actor generates both via two sequential LLM calls when processing a `CharChat`)
5. Send `MoveDecision` messages to both characters in the pair; await both replies
6. Apply location changes from movement replies to character state
7. Log the tick entry and advance world state

#### Scenario: Successful tick with two characters
- **WHEN** both the CharChat reply and both MoveDecision replies are received without error
- **THEN** the engine SHALL log the exchange, update character locations if changed, and advance to the next tick

#### Scenario: CharChat LLM error
- **WHEN** the CharChat reply contains an error payload
- **THEN** the engine SHALL log the error, skip logging the exchange, and advance the tick without updating character state

#### Scenario: Director broadcast before character exchange
- **WHEN** a GameDirector actor is registered and a tick begins
- **THEN** the director's broadcast reply channel SHALL be fully drained before the engine calls `scheduler.Next()`, ensuring all characters have processed director events before the exchange
