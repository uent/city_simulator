## Requirements

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

### Requirement: Game Director tick invocation
The system SHALL invoke the Game Director at the beginning of each tick, before scheduling any character exchange, when `Engine.director` is non-nil.

The invocation SHALL:
1. Call `director.BuildDirectorPrompt` with the current world state, regular character list (as `[]*character.Character`), and tick number
2. Send the prompt to the LLM client
3. Parse the response with `director.ParseToolCalls`
4. For each `ToolCall` in the result, call `registry.Dispatch(call.Name, call.Args, state, chars)`; log and skip any dispatch error without stopping
5. After each successful dispatch, print `  [Director] <action.Summary(call.Args)>` to stdout
6. Characters that received items in their `Inbox` during this step will have them flushed when their prompt is built later in the same tick

If the LLM call fails, the error SHALL be logged and the tick SHALL continue without executing any actions (fail-open, not fail-stop).

The `ParseDirectorEvents` function and the old `BuildDirectorPrompt` from `internal/llm/` SHALL NOT be called; those functions are removed.

#### Scenario: Director generates actions before characters act
- **WHEN** a tick begins with a non-nil Game Director
- **THEN** all dispatched director actions SHALL complete before `RunExchange` is called for that tick, so characters observe the updated world state

#### Scenario: LLM call fails during director turn
- **WHEN** the LLM returns an error during the director's turn
- **THEN** the error SHALL be logged, no actions SHALL be dispatched, and the simulation SHALL continue to the character exchange step normally

#### Scenario: Director emits unknown action name
- **WHEN** the director response contains a tool call with a name not in the registry
- **THEN** the dispatch error SHALL be logged, that action SHALL be skipped, and subsequent actions in the same response SHALL still be dispatched

#### Scenario: Director generates zero tool calls
- **WHEN** the director response contains no `<tool_calls>` block or an empty array
- **THEN** no actions are dispatched and the tick proceeds normally without any error

#### Scenario: Successful action prints summary
- **WHEN** the director dispatches `set_weather` with `{"type": "storm"}` successfully
- **THEN** the output SHALL contain `[Director] set_weather: storm`

#### Scenario: Summary includes key argument for NPC action
- **WHEN** the director dispatches `move_npc` with `{"id": "alice", "destination": "market"}` successfully
- **THEN** the output SHALL contain `[Director] move_npc: alice → market`

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

### Requirement: Narrative summary at simulation end
The system SHALL call `summary.GenerateSummary` and `summary.SaveSummary` after the tick loop completes in `Engine.Run()`.

The engine SHALL:
- Pass the LLM client, world state, character list, and scenario to `GenerateSummary`
- On success, call `SaveSummary` with the scenario name and generated text, then print the saved file path to stdout
- On LLM or file error, log the error and return nil (fail-open; the simulation already completed successfully)

#### Scenario: Summary generated and saved after run completes
- **WHEN** `Engine.Run` finishes its last tick without context cancellation
- **THEN** the engine SHALL attempt to generate and save a summary file, print the file path to stdout, and return nil

#### Scenario: Summary generation fails
- **WHEN** `GenerateSummary` returns an error
- **THEN** the engine SHALL log the error, skip saving, and return nil (not propagate the error to the caller)

#### Scenario: Run cancelled via context
- **WHEN** the context is cancelled before the last tick completes
- **THEN** the engine SHALL return `ctx.Err()` immediately and SHALL NOT attempt to generate a summary
