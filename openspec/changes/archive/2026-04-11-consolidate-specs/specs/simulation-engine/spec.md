## ADDED Requirements

<!-- Absorbe simulation-premise-display (impresión del bloque de concepto al inicio)
     y los requirements de rendering de character-engine (formateo per-tick de acción/speech y JSONL logging). -->

### Requirement: Motor imprime el bloque de concepto del mundo al inicio
Al comienzo de `Engine.Run`, antes de ejecutar el primer tick, el motor SHALL imprimir el bloque de concepto del mundo a stdout si `Scenario.World.Concept.Premise` es no-vacío.

El bloque SHALL seguir este formato:
```
=== World Concept ===
Premise: <premise>
Flavor:  <flavor>
Rules:
  - <rule>
=====================
```

La línea `Flavor:` SHALL omitirse si `Concept.Flavor` está vacío. La sección `Rules:` SHALL omitirse si `Concept.Rules` está vacío. Si `Concept.Premise` está vacío, el bloque completo SHALL omitirse silenciosamente.

#### Scenario: Bloque completo de concepto impreso
- **WHEN** `Engine.Run` se llama y `Scenario.World.Concept.Premise` es `"Bears hiding among humans"`
- **THEN** stdout SHALL contener un bloque que comienza con `=== World Concept ===` e incluye la línea de premise antes de cualquier output de ticks

#### Scenario: Línea de Flavor presente cuando está seteada
- **WHEN** `Concept.Flavor` es `"absurdist heist comedy"`
- **THEN** stdout SHALL contener `Flavor:  absurdist heist comedy` dentro del bloque de concepto

#### Scenario: Sección de Rules presente cuando está seteada
- **WHEN** `Concept.Rules` contiene dos entradas
- **THEN** stdout SHALL contener `Rules:` seguido de dos líneas `  - <rule>`

#### Scenario: Línea de Flavor omitida cuando está vacía
- **WHEN** `Concept.Flavor` es string vacío
- **THEN** stdout SHALL NOT contener una línea `Flavor:` en el bloque de concepto

#### Scenario: Sección de Rules omitida cuando está vacía
- **WHEN** `Concept.Rules` es nil o vacío
- **THEN** stdout SHALL NOT contener una línea `Rules:` en el bloque de concepto

#### Scenario: Bloque omitido cuando premise está vacío
- **WHEN** `Concept.Premise` es string vacío
- **THEN** ningún bloque de concepto SHALL imprimirse y no SHALL ocurrir ningún error

---

### Requirement: Motor renderiza acción y speech por separado en cada tick
El motor SHALL mostrar la acción y el speech de los personajes en líneas separadas por turno. El formato SHALL ser:

```
── Tick N ── InitiatorName [location] → ResponderName [location] ──
*initiator action*
InitiatorName: initiator speech
*responder action*
ResponderName: responder speech
```

Las líneas de acción solo se imprimen cuando el campo `Action` es no-vacío. Una acción ausente resulta en que la línea de acción se omite completamente (sin línea en blanco como placeholder).

#### Scenario: Tick con acción y speech para ambos personajes
- **WHEN** `CharChatReply` tiene `Action` y `Speech` no-vacíos para initiador y respondedor
- **THEN** la consola SHALL imprimir cuatro líneas: acción del iniciador, speech del iniciador, acción del respondedor, speech del respondedor

#### Scenario: Tick donde un personaje no tiene acción
- **WHEN** `CharChatReply.InitiatorAction` está vacío pero `ResponderAction` es no-vacío
- **THEN** la consola SHALL omitir la línea de acción del iniciador y aún imprimir la línea de acción del respondedor

#### Scenario: Tick donde ningún personaje tiene acción
- **WHEN** ambos campos `Action` están vacíos
- **THEN** la consola SHALL imprimir solo las dos líneas de speech (una por personaje), sin líneas de acción

---

### Requirement: Entrada del log JSONL incluye campos de acción y speech
El struct `logEntry` escrito en `OutputWriter` SHALL incluir los siguientes campos adicionales:
- `initiator_speech string` — texto hablado del iniciador
- `initiator_action string` — texto de acción del iniciador (string vacío si ninguno)
- `responder_speech string` — texto hablado del respondedor
- `responder_action string` — texto de acción del respondedor (string vacío si ninguno)

#### Scenario: Entrada del log con acción presente
- **WHEN** `CharChatReply.InitiatorAction` es `slams the table`
- **THEN** la línea JSONL para ese tick SHALL contener `"initiator_action":"slams the table"`

#### Scenario: Entrada del log sin acción
- **WHEN** `CharChatReply.ResponderAction` está vacío
- **THEN** la línea JSONL SHALL contener `"responder_action":""` (string vacío, no omitido)
