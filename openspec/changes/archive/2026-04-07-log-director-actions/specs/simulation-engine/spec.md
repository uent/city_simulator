## MODIFIED Requirements

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
