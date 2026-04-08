## MODIFIED Requirements

### Requirement: CharacterActor processes movement decision messages
When a `MoveDecision` message arrives, the `CharacterActor` SHALL call the LLM with the character's movement prompt (current location + available locations + zone roster), parse the response, and write back a reply containing the chosen location name or `"stay"` on error.

The `MoveDecisionPayload.SystemPrompt` field SHALL be a fully pre-built string (built by the engine via `BuildMovementPrompt` with the zone roster injected) so the actor does not need to rebuild it.

#### Scenario: Valid location returned by LLM
- **WHEN** the LLM returns a string matching one of the available location names (exact or case-insensitive)
- **THEN** the reply SHALL contain that location name

#### Scenario: LLM returns unknown location
- **WHEN** the LLM returns a string that does not match any available location
- **THEN** the reply SHALL contain `"stay"`

#### Scenario: Movement prompt received by LLM contains zone roster
- **WHEN** the engine dispatches a MoveDecision with a roster listing characters at various locations
- **THEN** the system prompt sent to the LLM SHALL include that zone roster information

### Requirement: CharacterActor processes CharChat messages
The system SHALL provide a `CharacterActor` in `internal/character/actor.go` that wraps a `*character.Character` and an `*llm.Client`. When a `CharChat` message arrives in its inbox, the actor SHALL:
1. Build a system prompt from the character's persona and the world context included in the message payload (including zone presence appended by the engine)
2. Retrieve the per-pair conversation history for the `msg.From` character (up to the last 10 turns)
3. Call the LLM client to generate a dialogue response
4. Append the exchange to the per-pair history
5. Update the character's memory with the new entry
6. Write a reply `Message` containing the generated text to `msg.ReplyChan`

#### Scenario: First exchange between two characters
- **WHEN** a CharChat message arrives from a character with whom this actor has no prior history
- **THEN** the actor SHALL generate a response using only the system prompt and the incoming message (no prior turns), and SHALL write the reply to ReplyChan

#### Scenario: Subsequent exchange with existing history
- **WHEN** a CharChat message arrives from a character with whom this actor has prior turns
- **THEN** the actor SHALL include the last ≤10 prior turns as LLM message history before generating the reply

#### Scenario: LLM error during CharChat
- **WHEN** the LLM client returns an error
- **THEN** the actor SHALL write an error reply (with empty text and the error wrapped in the payload) to ReplyChan and SHALL NOT crash or exit its goroutine

#### Scenario: CharChat system prompt includes zone presence
- **WHEN** a CharChat message is dispatched and a third character is co-located with the responder
- **THEN** the responder's system prompt (ResponderSystem) SHALL contain that co-located character's name in the zone-presence block
