## ADDED Requirements

<!-- Absorbe director-character-spawn: spawn_character action, gating por regla, cap de personajes,
     y registro de actores spawneados. Todo es parte del mismo sistema de director. -->

### Requirement: Acción spawn_character crea un personaje completo
El sistema SHALL implementar una acción de director `spawn_character` en `internal/director/actions_npc.go`. La acción SHALL aceptar los siguientes args: `id` (string, requerido), `name` (string, requerido), `age` (int, opcional), `occupation` (string, opcional), `motivation` (string, opcional), `fear` (string, opcional), `core_belief` (string, opcional), `internal_tension` (string, opcional), `formative_events` ([]string, opcional), `location` (string, opcional), `emotional_state` (string, opcional, default `"neutral"`), `goals` ([]string, opcional). La acción SHALL construir un `*character.Character` con esos campos, setear `MaxMemory` a 20, inicializar `Inbox` como slice vacío, y agregarlo a `*chars`.

#### Scenario: spawn_character con todos los args crea personaje completo
- **WHEN** el director llama `spawn_character` con todos los campos poblados
- **THEN** un nuevo `*character.Character` SHALL agregarse al slice de personajes del motor con todos los campos provistos seteados

#### Scenario: spawn_character con solo args requeridos usa defaults
- **WHEN** el director llama `spawn_character` con solo `id` y `name`
- **THEN** un nuevo personaje SHALL crearse con `EmotionalState = "neutral"`, `MaxMemory = 20`, y todos los campos opcionales en valor cero

#### Scenario: spawn_character emite un evento público de mundo
- **WHEN** `spawn_character` se ejecuta exitosamente
- **THEN** un `world.Event` de tipo `"spawn"` SHALL agregarse a `state.EventLog` con `Visibility = "public"` y una descripción que nombra al personaje

#### Scenario: spawn_character falla cuando id ya existe
- **WHEN** `spawn_character` se llama con un `id` que coincide con un personaje existente en `*chars`
- **THEN** `Execute` SHALL retornar error no-nil y SHALL NOT mutar `state` ni `*chars`

#### Scenario: spawn_character falla cuando id está ausente
- **WHEN** `spawn_character` se llama sin el arg `id`
- **THEN** `Execute` SHALL retornar error no-nil

#### Scenario: spawn_character falla cuando name está ausente
- **WHEN** `spawn_character` se llama sin el arg `name`
- **THEN** `Execute` SHALL retornar error no-nil

---

### Requirement: spawn_character controlado por character_spawn_rule
La acción `spawn_character` SHALL estar registrada en el registry del director en todo momento. Sin embargo, SHALL aparecer en el bloque `<tools>` del prompt del director SOLO cuando `state.Concept.CharacterSpawnRule` es no-vacío. Si la acción se invoca en runtime y no hay regla definida, `Execute` SHALL retornar error no-nil sin mutar el estado.

#### Scenario: spawn_character visible en el prompt cuando la regla está definida
- **WHEN** `state.Concept.CharacterSpawnRule` es no-vacío
- **THEN** `BuildDirectorPrompt` SHALL incluir `spawn_character` en el bloque `<tools>` junto con el texto de la regla

#### Scenario: spawn_character ausente del prompt cuando la regla no está definida
- **WHEN** `state.Concept.CharacterSpawnRule` está vacío
- **THEN** `BuildDirectorPrompt` SHALL NOT incluir `spawn_character` en el bloque `<tools>`

#### Scenario: spawn_character invocado sin regla retorna error
- **WHEN** `Execute` se llama en `spawnCharacterAction` y `state.Concept.CharacterSpawnRule == ""`
- **THEN** `Execute` SHALL retornar error no-nil conteniendo las palabras "no character_spawn_rule defined"

---

### Requirement: spawn_character respeta el cap max_spawned_characters
Cuando `state.Concept.MaxSpawnedCharacters > 0`, la acción `spawn_character` SHALL rechazar crear personajes adicionales una vez que `state.SpawnedCharacters >= MaxSpawnedCharacters`. El world state SHALL rastrear el conteo de personajes spawneados en `State.SpawnedCharacters`, incrementado en cada spawn exitoso.

#### Scenario: Spawn bloqueado cuando se alcanza el cap
- **WHEN** `MaxSpawnedCharacters` es `3` y `SpawnedCharacters` ya es `3`
- **THEN** `Execute` SHALL retornar error no-nil y SHALL NOT mutar state ni `*chars`

#### Scenario: Spawn permitido cuando está bajo el cap
- **WHEN** `MaxSpawnedCharacters` es `3` y `SpawnedCharacters` es `2`
- **THEN** `Execute` SHALL tener éxito y `SpawnedCharacters` SHALL volverse `3`

#### Scenario: Cap de cero significa ilimitado
- **WHEN** `MaxSpawnedCharacters` es `0`
- **THEN** el chequeo de cap SHALL omitirse y el spawn SHALL no ser bloqueado solo por conteo

---

### Requirement: Personajes spawneados dinámicamente participan como actores activos
El motor de simulación SHALL registrar los personajes recién spawneados como actores activos para que reciban turnos LLM y eventos del mundo desde el siguiente ciclo de procesamiento. El motor SHALL mantener un set de IDs de personajes ya registrados (`registeredChars`) para detectar personajes nuevos tras cada turno del director. Por cada personaje no registrado encontrado, el motor SHALL crear un `CharacterActor`, registrarlo en el message bus, iniciarlo con el contexto actual, y marcar su ID como registrado.

#### Scenario: Personaje spawneado recibe turno LLM tras ser creado
- **WHEN** `spawn_character` agrega un personaje en el tick N
- **THEN** ese personaje SHALL aparecer en la lista de actores y ser elegible para actuar a partir del tick N+1 (o tick N si el motor aún no procesó ese slot)

#### Scenario: Personajes ya registrados no se re-registran
- **WHEN** `registerSpawnedChars` corre y todos los personajes en `e.chars` ya están en `registeredChars`
- **THEN** no se crearán nuevos actores y el bus SHALL no ser modificado

#### Scenario: Motor loguea evento de spawn
- **WHEN** `spawn_character` se ejecuta exitosamente
- **THEN** el motor SHALL loguear un mensaje prefijado con `[spawn]` identificando el id y nombre del nuevo personaje

---

### Requirement: Scheduler integra personajes spawneados en la rotación de pares
El `Scheduler` SHALL exponer un método `AddCharacter(newChar *character.Character, known []*character.Character, locations []string)`. Cuando se llama, SHALL crear pares entre `newChar` y cada personaje en `known`, agregar esos pares a la rotación existente, y asignar a `newChar` una ubicación inicial aleatoria de `locations` si su campo `Location` está actualmente vacío.

#### Scenario: Nuevo personaje emparejado con todos los existentes
- **WHEN** `AddCharacter` se llama con un nuevo personaje y un slice de N personajes existentes
- **THEN** N nuevos pares SHALL agregarse a la lista de pares del scheduler

#### Scenario: Nuevo personaje recibe ubicación aleatoria cuando está vacío
- **WHEN** `AddCharacter` se llama y `newChar.Location == ""`
- **THEN** `newChar.Location` SHALL setearse a una de las ubicaciones provistas

#### Scenario: Nuevo personaje mantiene su ubicación cuando ya está seteada
- **WHEN** `AddCharacter` se llama y `newChar.Location` es ya no-vacío
- **THEN** `newChar.Location` SHALL NOT sobreescribirse
