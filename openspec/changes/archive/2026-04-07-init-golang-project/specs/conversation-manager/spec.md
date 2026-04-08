## ADDED Requirements

### Requirement: Conversation thread per character pair
The system SHALL maintain one `Thread` per ordered character-pair key (e.g., "alice→bob"). A `Thread` holds a slice of `Turn` values. Each `Turn` has: speaker character ID, message text, and tick number.

#### Scenario: First message in a new thread
- **WHEN** `manager.AddTurn(fromID, toID string, text string, tick int)` is called for a pair that has no existing thread
- **THEN** a new `Thread` SHALL be created and the turn appended to it

#### Scenario: Subsequent messages in existing thread
- **WHEN** `AddTurn` is called for a pair that already has a thread
- **THEN** the turn SHALL be appended to the existing thread, preserving prior turns

### Requirement: Format conversation history for LLM prompt
The system SHALL provide a `History(fromID, toID string, maxTurns int) []llm.Message` method that returns the last `maxTurns` turns of the thread between the two characters, formatted as LLM message objects (role "user" for the initiator, "assistant" for the responder, relative to the caller).

#### Scenario: History shorter than maxTurns
- **WHEN** the thread has fewer turns than `maxTurns`
- **THEN** all turns SHALL be returned

#### Scenario: History longer than maxTurns
- **WHEN** the thread has more turns than `maxTurns`
- **THEN** only the most recent `maxTurns` turns SHALL be returned in chronological order

#### Scenario: No existing thread
- **WHEN** `History` is called for a pair with no thread
- **THEN** the method SHALL return an empty slice and no error

### Requirement: Run a full exchange between two characters
The system SHALL provide a `RunExchange(ctx context.Context, initiator, responder *character.Character, world *world.State, tick int) error` method on the manager that:
1. Builds the initiator's system prompt from their persona + world summary
2. Assembles the message history for this pair
3. Calls the LLM client to generate the initiator's message
4. Stores the initiator's message in both characters' memories and the thread
5. Builds the responder's system prompt
6. Calls the LLM client to generate the responder's reply (with the initiator message as the last user turn)
7. Stores the responder's reply in both memories and the thread
8. Appends a "conversation" event to world state

#### Scenario: Successful exchange
- **WHEN** both LLM calls succeed
- **THEN** the method SHALL return nil and both characters' memories SHALL contain the new turns

#### Scenario: LLM error on initiator turn
- **WHEN** the first LLM call returns an error
- **THEN** the method SHALL return that error without making a second LLM call
