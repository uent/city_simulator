## ADDED Requirements

<!-- Absorbe character-expression (Expression type, ParseExpression, FormatExpression, instrucción de formato)
     y los requirements runtime de character-engine (memory buffer, CharChatReply, historial de conversación)
     que corresponden al comportamiento del actor, no al modelo de datos. -->

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
