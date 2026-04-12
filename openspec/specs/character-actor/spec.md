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

---

## ADDED Requirements

### Requirement: Expression type
El sistema SHALL definir un struct `Expression` en `internal/character/expression.go` con dos campos string: `Speech` (lo que el personaje dice en voz alta) y `Action` (breve descripción en tercera persona de lo que el personaje hace físicamente; puede estar vacío).

#### Scenario: Expression con ambos campos poblados
- **WHEN** una respuesta LLM raw contiene `*picks up the envelope* I wasn't expecting this.`
- **THEN** `Expression.Action` SHALL igualar `picks up the envelope` y `Expression.Speech` SHALL igualar `I wasn't expecting this.`

#### Scenario: Expression sin marcadores de acción
- **WHEN** una respuesta LLM raw no contiene marcadores `*...*`
- **THEN** `Expression.Action` SHALL ser string vacío y `Expression.Speech` SHALL igualar la respuesta completa trimmed

#### Scenario: Expression con solo marcador de acción
- **WHEN** una respuesta LLM raw es solo `*walks away silently*` sin texto restante
- **THEN** `Expression.Action` SHALL igualar `walks away silently` y `Expression.Speech` SHALL ser string vacío

---

### Requirement: Función ParseExpression
El sistema SHALL proveer `ParseExpression(raw string) Expression` que:
1. Encuentra la primera subcadena `*...*` en `raw` (desde el `*` de apertura hasta el primer `*` de cierre).
2. Asigna el contenido entre los marcadores (trimmed) a `Expression.Action`.
3. Concatena el texto antes y después del bloque marcador (ambos trimmed), separados por un espacio si ambos son no-vacíos, y asigna el resultado a `Expression.Speech`.
4. Si no se encuentra marcador `*...*`, asigna el `raw` completo trimmed a `Expression.Speech` y deja `Expression.Action` vacío.

#### Scenario: Marcador de acción al inicio
- **WHEN** raw es `*looks around nervously* You again?`
- **THEN** `Action` es `looks around nervously` y `Speech` es `You again?`

#### Scenario: Marcador de acción en el medio
- **WHEN** raw es `I told you *slams fist on table* I don't know anything.`
- **THEN** `Action` es `slams fist on table` y `Speech` es `I told you I don't know anything.`

#### Scenario: Marcador de acción al final
- **WHEN** raw es `Fine. Have it your way. *turns and walks out*`
- **THEN** `Action` es `turns and walks out` y `Speech` es `Fine. Have it your way.`

#### Scenario: Input vacío o solo whitespace
- **WHEN** raw es string vacío o contiene solo whitespace
- **THEN** tanto `Speech` como `Action` SHALL ser strings vacíos

---

### Requirement: Instrucción de formato de expresión en system prompt
`BuildSystemPrompt` SHALL agregar una instrucción de formato de expresión a cada system prompt de personaje. La instrucción SHALL explicar la convención `*acción*` y decirle al personaje que la use cuando realiza una acción física junto con el diálogo.

La instrucción SHALL colocarse después de la línea de cierre "Stay in character".

#### Scenario: System prompt incluye instrucción de expresión
- **WHEN** `BuildSystemPrompt` se llama para cualquier personaje
- **THEN** el string retornado SHALL contener una línea instruyendo al personaje a prefijar acciones físicas con marcadores `*...*`

#### Scenario: Instrucción es independiente del idioma
- **WHEN** `BuildSystemPrompt` se llama con un parámetro `language` no-vacío
- **THEN** la instrucción de formato de expresión SHALL aparecer en el prompt independientemente del idioma

---

### Requirement: Helper FormatExpression
El sistema SHALL proveer `FormatExpression(e Expression) string` que recombina un `Expression` al formato wire `*acción* speech` usado al inyectar el turno de un personaje en el contexto de otro:
- Si `Action` es no-vacío: retorna `*{action}* {speech}` (trimmed)
- Si `Action` está vacío: retorna `Speech` tal cual
- Si `Speech` está vacío y `Action` es no-vacío: retorna `*{action}*`

#### Scenario: Ambos campos presentes
- **WHEN** `Action` es `picks up the card` y `Speech` es `Interesting.`
- **THEN** `FormatExpression` retorna `*picks up the card* Interesting.`

#### Scenario: Action vacío
- **WHEN** `Action` está vacío y `Speech` es `What do you want?`
- **THEN** `FormatExpression` retorna `What do you want?`

#### Scenario: Speech vacío
- **WHEN** `Action` es `nods silently` y `Speech` está vacío
- **THEN** `FormatExpression` retorna `*nods silently*`

---

### Requirement: CharChatReply incluye campos de acción y speech
El sistema SHALL definir `CharChatReply` en `internal/messaging/message.go` con los siguientes campos:
- `InitiatorSpeech string` — lo que dijo en voz alta el personaje iniciador
- `InitiatorAction string` — acción física del iniciador (puede estar vacío)
- `ResponderSpeech string` — lo que dijo en voz alta el personaje respondedor
- `ResponderAction string` — acción física del respondedor (puede estar vacío)
- `Err error` — no-nil si el intercambio falló

Los campos anteriores `InitiatorText` y `ResponderText` SHALL eliminarse.

#### Scenario: Reply con acción y speech
- **WHEN** el actor genera una respuesta de iniciador con acción y speech
- **THEN** `CharChatReply.InitiatorAction` SHALL ser no-vacío y `CharChatReply.InitiatorSpeech` SHALL contener solo el texto hablado

#### Scenario: Reply solo con speech
- **WHEN** el actor genera una respuesta sin marcadores `*...*`
- **THEN** el campo `Action` correspondiente SHALL ser string vacío y el campo `Speech` SHALL contener el texto completo de la respuesta
