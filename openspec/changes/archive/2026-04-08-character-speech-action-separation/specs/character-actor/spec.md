## MODIFIED Requirements

### Requirement: CharacterActor processes CharChat messages
The system SHALL provide a `CharacterActor` in `internal/character/actor.go` that wraps a `*character.Character` and an `*llm.Client`. When a `CharChat` message arrives in its inbox, the actor SHALL:
1. Build a system prompt from the character's persona and the world context included in the message payload (including zone presence appended by the engine).
2. Retrieve the per-pair conversation history for the `msg.From` character (up to the last 10 turns).
3. Call the LLM client to generate the initiator's expression (raw text with optional `*action*` markers).
4. Parse the raw initiator response using `ParseExpression` to extract `InitiatorSpeech` and `InitiatorAction`.
5. Format the initiator's expression using `FormatExpression` and use it as the user-role message when calling the LLM for the responder's reply.
6. Parse the raw responder response using `ParseExpression` to extract `ResponderSpeech` and `ResponderAction`.
7. Append the raw texts (including markers) to the per-pair history and character memory.
8. Write a `CharChatReply` containing `InitiatorSpeech`, `InitiatorAction`, `ResponderSpeech`, `ResponderAction` to `msg.ReplyChan`.

#### Scenario: First exchange between two characters
- **WHEN** a CharChat message arrives from a character with whom this actor has no prior history
- **THEN** the actor SHALL generate a response using only the system prompt and the incoming message (no prior turns), parse both responses, and write the reply to ReplyChan

#### Scenario: Subsequent exchange with existing history
- **WHEN** a CharChat message arrives from a character with whom this actor has prior turns
- **THEN** the actor SHALL include the last ≤10 prior turns (raw text, including markers) as LLM message history before generating the reply

#### Scenario: Initiator response contains action marker
- **WHEN** the LLM returns `*glances at the door* We should talk somewhere private.` for the initiator
- **THEN** `CharChatReply.InitiatorAction` SHALL equal `glances at the door` and `CharChatReply.InitiatorSpeech` SHALL equal `We should talk somewhere private.`

#### Scenario: Responder receives initiator action as context
- **WHEN** the initiator's expression includes a non-empty `Action`
- **THEN** the user-role message sent to the responder LLM SHALL be formatted as `*{action}* {speech}` so the responder can react to both the physical action and the words

#### Scenario: LLM produces no action markers (fallback)
- **WHEN** the LLM returns plain text with no `*...*` markers for either character
- **THEN** the corresponding `Action` field SHALL be empty and `Speech` SHALL contain the full response; the simulation SHALL continue normally

#### Scenario: LLM error during CharChat
- **WHEN** the LLM client returns an error
- **THEN** the actor SHALL write an error reply (with empty fields and the error wrapped in the payload) to ReplyChan and SHALL NOT crash or exit its goroutine

#### Scenario: CharChat system prompt includes zone presence
- **WHEN** a CharChat message is dispatched and a third character is co-located with the responder
- **THEN** the responder's system prompt (ResponderSystem) SHALL contain that co-located character's name in the zone-presence block
