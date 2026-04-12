## ADDED Requirements

<!-- Spec consolidado: absorbe simulation-summary y character-summary-cards en una sola capacidad de salida. -->

### Requirement: Generación del resumen narrativo
El sistema SHALL proveer `GenerateSummary(ctx context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario, language string) (string, error)` en el package `internal/summary` que construye un prompt desde los eventos del mundo y los estados finales de los personajes, lo envía al LLM, y retorna el string de resumen completo.

El string retornado SHALL consistir en:
1. El texto narrativo generado por el LLM (sin cambios).
2. El bloque de character cards producido por `renderCharacterCards(chars)` agregado después de la narrativa. Si no hay personajes no-director, no se agrega ningún bloque de cards.

El prompt SHALL incluir:
- El nombre del escenario
- El concepto del mundo (premise, flavor, y rules) cuando presente en el escenario
- La atmósfera y el clima del mundo cuando presentes
- El número total de ticks que corrieron
- Todos los eventos del mundo (hasta los últimos 200 si la lista supera 200 entradas)
- El nombre, rol, ubicación, motivación, miedo, goals, y estado emocional final de cada personaje (omitiendo campos vacíos)

El system prompt del LLM SHALL instruir al modelo a producir una narrativa rica de al menos seis párrafos detallados, cubriendo el arco de la simulación, puntos de inflexión clave, desarrollo de personajes, y resultado final.

#### Scenario: Generación exitosa incluye character cards
- **WHEN** `GenerateSummary` se llama tras una simulación completada con eventos y personajes no-director presentes
- **THEN** la función SHALL retornar un string no-vacío conteniendo la narrativa seguida por el bloque de character cards, y nil error

#### Scenario: Fallo del LLM
- **WHEN** el cliente LLM retorna un error durante la generación del resumen
- **THEN** la función SHALL retornar string vacío y error no-nil describiendo el fallo

#### Scenario: Lista de eventos supera 200 entradas
- **WHEN** el world state contiene más de 200 eventos
- **THEN** la función SHALL incluir solo los últimos 200 eventos en el prompt sin error

#### Scenario: Concepto del mundo ausente
- **WHEN** el world config del escenario tiene un bloque `Concept` vacío
- **THEN** la función SHALL omitir la sección de concepto del prompt sin error

#### Scenario: Personaje sin goals ni motivación
- **WHEN** un personaje tiene `Goals`, `Motivation`, o `Fear` vacíos
- **THEN** esos campos SHALL omitirse de la entrada de ese personaje en el prompt

#### Scenario: Solo personajes game_director presentes
- **WHEN** todos los personajes son de tipo `game_director`
- **THEN** la función SHALL retornar el texto narrativo sin bloque de character cards

---

### Requirement: Persistencia del resumen en archivo con timestamp
El sistema SHALL proveer `SaveSummary(scenarioName string, content string) (string, error)` que escribe el resumen en `simulations/<scenarioName>/summary-<timestamp>.md`, donde `<timestamp>` está formateado como RFC3339 con dos puntos reemplazados por guiones.

La función SHALL:
- Crear el directorio destino si no existe
- Retornar el path absoluto del archivo escrito en caso de éxito

#### Scenario: Resumen guardado en nuevo archivo
- **WHEN** `SaveSummary` se llama con nombre de escenario y contenido no-vacíos
- **THEN** un archivo `summary-<timestamp>.md` SHALL crearse dentro de `simulations/<scenarioName>/` y la función SHALL retornar su path con nil error

#### Scenario: Múltiples ejecuciones no se sobreescriben
- **WHEN** `SaveSummary` se llama dos veces en el mismo segundo o en ejecuciones distintas
- **THEN** cada llamada SHALL producir un archivo distinto

#### Scenario: Directorio del escenario no existe
- **WHEN** `SaveSummary` se llama para un escenario cuyo directorio no existe aún
- **THEN** la función SHALL crear el directorio y escribir el archivo sin retornar error

---

### Requirement: Renderizado de character cards
El package `internal/summary` SHALL exponer `renderCharacterCards(chars []*character.Character) string` que produce un bloque Markdown con una card por personaje no-director.

El bloque SHALL comenzar con una línea horizontal (`---`) seguida de un heading `## Character Cards`, luego una sección `### <Name>` por personaje. Los personajes Game Director (`Type == "game_director"`) SHALL ser excluidos. Los campos que son strings vacíos, valores cero, o nil SHALL omitirse de la card.

La card SHALL mostrar los siguientes campos cuando presentes: Name (como heading de sección), Age, Occupation, Appearance, Motivation, Fear, Core Belief, Internal Tension, Formative Events, Location, Emotional State, Goals, Voice (Formality, Verbal Tics, Response Length, Humor Type, Communication Style), Relational Defaults (Strangers, Authority, Vulnerable), Dialogue Examples, Cover Identity (Alias, Role, Backstory, Weaknesses), y Relationships (una entrada por `CharacterJudgment` en `Character.Judgments`, mostrando nombre del personaje conocido, niveles de trust/interest/threat, e impression).

#### Scenario: Todos los campos poblados
- **WHEN** `renderCharacterCards` se llama con un personaje que tiene todos los campos seteados
- **THEN** el output SHALL contener una sección `### <Name>` con cada campo no-vacío renderizado como ítem de lista Markdown o sub-lista

#### Scenario: Campos vacíos omitidos
- **WHEN** un personaje tiene `Fear`, `CoreBelief`, e `InternalTension` vacíos
- **THEN** esos labels SHALL NOT aparecer en la card de ese personaje

#### Scenario: Game Director excluido
- **WHEN** el slice de personajes contiene uno con `Type == "game_director"`
- **THEN** ninguna card SHALL renderizarse para ese personaje

#### Scenario: Sin personajes no-director
- **WHEN** todos los personajes son game directors o el slice está vacío
- **THEN** `renderCharacterCards` SHALL retornar string vacío

#### Scenario: Cover identity presente en la card
- **WHEN** un personaje tiene `CoverIdentity` no-nil con `Alias`, `Role` y `Backstory` seteados
- **THEN** la card SHALL incluir una sub-sección "Cover Identity" con esos campos

#### Scenario: Cover identity nil omitida
- **WHEN** `Character.CoverIdentity` es nil
- **THEN** ninguna sección "Cover Identity" SHALL aparecer en la card de ese personaje

#### Scenario: Relationships renderizadas cuando existen judgments
- **WHEN** el mapa `Judgments` de un personaje contiene una o más entradas
- **THEN** la card SHALL incluir una sub-sección "Relationships" con una línea por judgment mostrando nombre, trust, interest, y threat, seguido de la impression como blockquote

#### Scenario: Sin relationships cuando el mapa de judgments está vacío
- **WHEN** el mapa `Judgments` de un personaje está vacío
- **THEN** ninguna sección "Relationships" SHALL aparecer en la card de ese personaje
