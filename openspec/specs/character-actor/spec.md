## Requirements

### Requirement: Actor interface
The system SHALL define an `Actor` interface in `internal/messaging/` with:
- `ID() string` — returns the actor's unique identifier
- `Inbox() chan Message` — returns the actor's inbox channel
- `Start(ctx context.Context)` — spawns the actor's processing goroutine; calling Start more than once is a no-op
- `Stop()` — signals the actor to stop; idempotent

Every concrete actor type (CharacterActor, DirectorActor) SHALL implement this interface.

#### Scenario: Start and Stop lifecycle
- **WHEN** `Start(ctx)` is called on an actor and then the context is cancelled
- **THEN** the actor's goroutine SHALL exit and the actor SHALL stop processing new messages

#### Scenario: Double Start is safe
- **WHEN** `Start(ctx)` is called twice on the same actor
- **THEN** only one goroutine SHALL be running; the second call SHALL be a no-op

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

### Requirement: CharacterActor owns per-pair conversation history
The `CharacterActor` SHALL maintain a `map[string][]Turn` keyed by peer character ID. Each `Turn` holds speaker ID, text, and tick number. The actor is the sole writer of this map; no other component accesses it directly.

#### Scenario: History is isolated per peer
- **WHEN** the actor exchanges messages with characters A and B separately
- **THEN** the history for A SHALL NOT contain turns from the exchange with B

#### Scenario: History is capped at MaxHistory turns per peer
- **WHEN** the number of stored turns for a peer exceeds the configured `MaxHistory` (default 20)
- **THEN** the oldest turns SHALL be evicted so that only the most recent `MaxHistory` turns are kept
