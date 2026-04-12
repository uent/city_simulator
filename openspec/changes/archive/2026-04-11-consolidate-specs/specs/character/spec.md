## ADDED Requirements

<!-- Spec consolidado: absorbe character-schema, character-cover-identity, character-initial-state,
     y los requirements de datos de character-engine (struct, type field, loader, inbox, memory). -->

### Requirement: Character struct completo
El sistema SHALL definir un struct `Character` con los siguientes campos:

**Identidad y tipo:**
- `ID string` (yaml: `id`) — identificador único
- `Type string` (yaml: `type`) — `"game_director"` o vacío para personaje regular
- `Name string` (yaml: `name`)
- `Age int` (yaml: `age`)
- `Gender string` (yaml: `gender`) — campo libre, opcional
- `Occupation string` (yaml: `occupation`)
- `Appearance string` (yaml: `appearance`) — cómo otros perciben al personaje; excluido de su propio system prompt

**Núcleo psicológico:**
- `Motivation string` (yaml: `motivation`)
- `Fear string` (yaml: `fear`)
- `CoreBelief string` (yaml: `core_belief`)
- `InternalTension string` (yaml: `internal_tension`)
- `FormativeEvents []string` (yaml: `formative_events`) — bullets causales "evento → consecuencia"

**Voz y relaciones:**
- `Voice VoiceProfile` (yaml: `voice`) — sub-struct con: `Formality`, `VerbalTics`, `ResponseLength`, `HumorType`, `CommunicationStyle`
- `RelationalDefaults RelationalProfile` (yaml: `relational_defaults`) — sub-struct con: `Strangers`, `Authority`, `Vulnerable`
- `DialogueExamples []string` (yaml: `dialogue_examples`)

**Identidad encubierta:**
- `CoverIdentity *CoverIdentity` (yaml: `cover_identity`) — nil cuando el personaje no tiene cobertura

**Estado runtime (nunca persistido en YAML):**
- `Location string` (yaml: `"-"`) — posición actual en el mundo
- `Goals []string` (yaml: `goals`)
- `EmotionalState string` (yaml: `emotional_state`)
- `Inventory []string` (yaml: `inventory`) — objetos que lleva al inicio de la simulación
- `InitialState string` (yaml: `initial_state`) — estado táctico o narrativo al inicio
- `Inbox []world.Event` (yaml: `"-"`) — eventos privados pendientes; inicializado como slice vacío no-nil por `LoadCharacters`
- `Memory []MemoryEntry` — buffer deslizante, no persistido
- `Judgments map[string]CharacterJudgment` (yaml: `"-"`) — juicios sobre otros personajes

Los campos `personality []string` y `backstory string` NO SHALL existir. YAMLs que los contengan los ignoran silenciosamente.

#### Scenario: Todos los campos YAML se cargan correctamente
- **WHEN** un `characters.yaml` tiene todas las claves del struct
- **THEN** `LoadCharacters` SHALL popular todos los campos y retornar nil error

#### Scenario: Campos opcionales ausentes usan valor cero
- **WHEN** una entrada omite campos opcionales como `gender`, `appearance`, `inventory`, `initial_state`
- **THEN** esos campos quedan como string vacío o slice nil sin error

#### Scenario: Campos legacy `personality` y `backstory` se ignoran
- **WHEN** un `characters.yaml` contiene `personality` y `backstory`
- **THEN** `LoadCharacters` SHALL retornar un `Character` poblado sin error, ignorando esas claves

#### Scenario: `type: game_director` se parsea correctamente
- **WHEN** una entrada tiene `type: game_director`
- **THEN** `Character.Type` SHALL igualar `"game_director"`

#### Scenario: `type` ausente resulta en string vacío
- **WHEN** una entrada omite el campo `type`
- **THEN** `Character.Type` SHALL ser string vacío y el personaje se trata como regular

---

### Requirement: CoverIdentity struct
El sistema SHALL definir un struct `CoverIdentity` con: `Alias string` (yaml: `alias`), `Role string` (yaml: `role`), `Backstory string` (yaml: `backstory`), `Weaknesses []string` (yaml: `weaknesses`). Un puntero nil significa que el personaje no tiene cobertura.

#### Scenario: Bloque cover_identity completo parseado
- **WHEN** una entrada contiene `cover_identity:` con `alias`, `role`, `backstory` y `weaknesses`
- **THEN** `Character.CoverIdentity` SHALL ser no-nil con los cuatro campos poblados

#### Scenario: Bloque cover_identity parcial aceptado
- **WHEN** una entrada contiene `cover_identity: { alias: "Gerald", role: "Honey sommelier" }` sin `backstory` ni `weaknesses`
- **THEN** `Character.CoverIdentity` SHALL ser no-nil con `Alias` y `Role` seteados; `Backstory` vacío; `Weaknesses` nil

#### Scenario: cover_identity ausente deja puntero nil
- **WHEN** una entrada omite `cover_identity`
- **THEN** `Character.CoverIdentity` SHALL ser nil sin error

---

### Requirement: LoadCharacters — cargador desde archivo YAML
El sistema SHALL proveer `LoadCharacters(path string) ([]Character, error)` que lee un archivo YAML y retorna TODOS los personajes (sin filtrar por tipo). La separación de directores es responsabilidad del scenario loader.

#### Scenario: Archivo no encontrado
- **WHEN** el path no existe
- **THEN** la función SHALL retornar error no-nil con el path en el mensaje

#### Scenario: YAML malformado
- **WHEN** el archivo contiene error de sintaxis
- **THEN** la función SHALL retornar error no-nil describiendo el fallo

#### Scenario: Archivo válido con múltiples personajes
- **WHEN** el archivo contiene dos o más entradas
- **THEN** la función SHALL retornar un slice con el mismo número de `Character` y nil error

#### Scenario: Tipos mezclados retornados sin filtrar
- **WHEN** el archivo contiene un personaje regular y uno `type: game_director`
- **THEN** `LoadCharacters` SHALL retornar un slice de longitud 2 con ambas entradas

#### Scenario: Inbox inicializado como slice vacío no-nil
- **WHEN** `LoadCharacters` carga cualquier personaje
- **THEN** `Character.Inbox` SHALL ser un slice vacío no-nil (no nil)

---

### Requirement: Buffer de memoria por personaje
El sistema SHALL mantener un buffer deslizante `Memory []MemoryEntry` por personaje, limitado a `MaxMemory` entradas (default 20). Cada `MemoryEntry` registra: nombre del hablante, texto del mensaje, y número de tick.

#### Scenario: Memoria bajo capacidad
- **WHEN** se han agregado menos entradas que `MaxMemory`
- **THEN** todas las entradas SHALL ser recuperables en orden de inserción

#### Scenario: Memoria al límite recibe nueva entrada
- **WHEN** se agrega una entrada y el buffer ya está en `MaxMemory`
- **THEN** la entrada más antigua SHALL ser eviccionada y la nueva agregada, manteniendo el total en `MaxMemory`

#### Scenario: Recuperar memoria como slice de mensajes
- **WHEN** se llama `character.RecentMemory(n int)`
- **THEN** el sistema SHALL retornar hasta las últimas `n` entradas en orden cronológico

---

### Requirement: Inbox para eventos privados
`Character.Inbox []world.Event` SHALL ser omitido de la serialización YAML (tag `yaml:"-"`). Cuando `BuildSystemPrompt` construye el prompt de un personaje, SHALL:

1. Revisar `character.Inbox` por ítems pendientes
2. Si no está vacío, agregar una sección "Private information you recently learned:" listando la descripción de cada ítem
3. Limpiar `character.Inbox` (set a slice vacío) tras leer — semántica flush-on-read

#### Scenario: Inbox flushed tras construir el prompt
- **WHEN** un personaje tiene un ítem en `Inbox` y se construye su prompt
- **THEN** el prompt SHALL contener la descripción del ítem Y `character.Inbox` SHALL estar vacío tras la llamada

#### Scenario: Inbox vacío no produce sección privada
- **WHEN** `character.Inbox` está vacío
- **THEN** el prompt construido SHALL NOT contener una sección "Private information"

#### Scenario: Inbox no persistido en YAML
- **WHEN** un `Character` con `Inbox` no-vacío se marshalea a YAML
- **THEN** el output SHALL NOT contener la clave `inbox`

---

### Requirement: BuildSystemPrompt con template psicológico estructurado
`BuildSystemPrompt(c character.Character) string` SHALL producir un system prompt con las siguientes secciones (omitiendo secciones cuyos campos estén todos vacíos):

1. Línea de identidad: `"You are {Name}, a {Age}-year-old {Gender} {Occupation}."` cuando `Gender` es no-vacío, o `"You are {Name}, a {Age}-year-old {Occupation}."` cuando es vacío
2. `Motivación:` usando `Motivation`
3. `Miedo:` usando `Fear`
4. `Creencia central:` usando `CoreBelief`
5. `Tensión interna:` usando `InternalTension`
6. `Eventos formativos:` bloque bulleteado usando `FormativeEvents`
7. `Voz:` bloque usando sub-campos de Voice
8. `Relaciones default:` bloque usando RelationalDefaults
9. `Objetivos:` bloque bulleteado usando `Goals`
10. `Estado emocional actual:` usando `EmotionalState`
11. `Ejemplos de diálogo:` bloque con comillas usando `DialogueExamples`
12. Sección de Cover Identity cuando `CoverIdentity` es no-nil: alias, role, backstory (si no vacío), weaknesses como lista bulleteada (si no vacíos)
13. Instrucción de cierre: `"Stay in character at all times. Respond as this person would. Keep responses concise."`
14. Instrucción de formato de expresión: explicar la convención `*acción*` y pedir su uso al realizar acciones físicas (ver capability `character-actor` para el formato exacto)
15. Si `language` es no-vacío: `"Respond in {language}."` al final

`Appearance` SHALL NOT incluirse en el propio prompt del personaje (describe cómo lo ven otros, no cómo se ve a sí mismo).

#### Scenario: Gender incluido en línea de identidad cuando presente
- **WHEN** `BuildSystemPrompt` se llama con `Gender == "femenino"` y `Occupation == "Detective"`
- **THEN** la línea de identidad SHALL contener `"femenino Detective"`

#### Scenario: Gender omitido cuando vacío
- **WHEN** `BuildSystemPrompt` se llama con `Gender == ""`
- **THEN** la línea de identidad SHALL seguir el formato sin token de género

#### Scenario: Cover identity presente en el prompt cuando seteada
- **WHEN** un personaje tiene `CoverIdentity.Alias = "Gerald"` y `Role = "Honey sommelier"`
- **THEN** `BuildSystemPrompt` SHALL contener ambos strings

#### Scenario: Sin sección cover identity cuando nil
- **WHEN** `Character.CoverIdentity` es nil
- **THEN** `BuildSystemPrompt` SHALL NOT contener "Cover Identity" ni "alias"

#### Scenario: Campos vacíos omitidos silenciosamente
- **WHEN** `FormativeEvents` es nil y `InternalTension` es vacío
- **THEN** el prompt SHALL NOT contener headers `"Tensión interna:"` ni `"Eventos formativos:"`

#### Scenario: Personaje mínimo produce prompt válido
- **WHEN** solo `Name`, `Age` y `Occupation` están seteados
- **THEN** el prompt SHALL contener solo la línea de identidad y la instrucción de cierre

---

### Requirement: ObservableSnapshot respeta CoverIdentity
`ObservableSnapshot(c Character) ObservableProfile` SHALL retornar una vista filtrada que otro personaje podría observar en un primer encuentro:
- Si `CoverIdentity` es no-nil: `Name = CoverIdentity.Alias`, `Occupation = CoverIdentity.Role`
- Si `CoverIdentity` es nil: `Name = Character.Name`, `Occupation = Character.Occupation`
- Siempre incluido: `Age`, `EmotionalState`, `Appearance`, `Location`
- Nunca incluido: `Motivation`, `Fear`, `CoreBelief`, `InternalTension`, `FormativeEvents`, `Goals`, internals de `CoverIdentity`

#### Scenario: Personaje con cover identity produce snapshot basado en alias
- **WHEN** `ObservableSnapshot` se llama con `CoverIdentity.Alias = "Don Gregorio"` y `Role = "Coleccionista"`
- **THEN** el perfil retornado SHALL tener `Name = "Don Gregorio"` y `Occupation = "Coleccionista"`, sin `Motivation`

#### Scenario: Personaje sin cover identity usa campos reales
- **WHEN** `ObservableSnapshot` se llama en un personaje sin `CoverIdentity`
- **THEN** el perfil SHALL tener `Name = Character.Name` y `Occupation = Character.Occupation`

#### Scenario: Appearance incluida en el observable snapshot
- **WHEN** `ObservableSnapshot` se llama en un personaje con `Appearance` no-vacío
- **THEN** `ObservableProfile.Appearance` SHALL igualar el valor del personaje

#### Scenario: Appearance ausente del propio system prompt
- **WHEN** `BuildSystemPrompt` se llama con un personaje que tiene `Appearance` no-vacío
- **THEN** el system prompt retornado SHALL NOT contener el texto de `Appearance`

---

### Requirement: CHARACTER_RULES.md como rulebook de creación
El sistema SHALL proveer `simulations/CHARACTER_RULES.md` — un documento Markdown usable como contexto para LLMs que crean nuevos personajes. SHALL incluir: el template YAML completo, descripción campo por campo, guía de anti-patrones, y al menos un personaje de ejemplo completo.

#### Scenario: El archivo existe y es legible
- **WHEN** se clona el repositorio
- **THEN** `simulations/CHARACTER_RULES.md` SHALL existir como archivo Markdown no-vacío

#### Scenario: El rulebook contiene todas las secciones requeridas
- **WHEN** se lee `CHARACTER_RULES.md`
- **THEN** SHALL contener secciones de: template YAML, descripciones de campos, anti-patrones, y un ejemplo completo

---

### Requirement: Campo appearance en todos los personajes de escenarios existentes
Todos los archivos `characters.yaml` en `simulations/` SHALL incluir un campo `appearance` en cada personaje no-director. Los personajes directores (`type: game_director`) SHALL NOT tener campo `appearance`.

#### Scenario: Personajes de escenarios existentes tienen campo appearance
- **WHEN** cualquier `characters.yaml` en `simulations/default/`, `simulations/honey-heist/`, `simulations/doom-hell-crusade/` o `simulations/test-scenario/` se carga
- **THEN** cada personaje no-director SHALL tener `Appearance` no-vacío tras cargar
