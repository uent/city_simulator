## ADDED Requirements

<!-- Absorbe world-concept: el struct WorldConcept y su parsing YAML pertenecen al spec del mundo. -->

### Requirement: WorldConcept struct en WorldConfig
El sistema SHALL definir un struct `WorldConcept` con cinco campos:
- `Premise string` (yaml: `premise`) — una sola oración describiendo la naturaleza fundamental de este mundo y la verdad oculta que los personajes deben preservar
- `Rules []string` (yaml: `rules`) — lista de restricciones que definen qué es "normal" en este mundo
- `Flavor string` (yaml: `flavor`) — string corto de tono/mood (e.g., "absurdist heist comedy")
- `CharacterSpawnRule string` (yaml: `character_spawn_rule`) — regla que describe cómo deben diseñarse los personajes creados dinámicamente; string vacío deshabilita el `spawn_character` del director
- `MaxSpawnedCharacters int` (yaml: `max_spawned_characters`) — número máximo de personajes que el director puede spawnear en runtime (`0` significa ilimitado)

`WorldConfig` SHALL exponer un campo `Concept WorldConcept` (yaml: `concept`) y un campo `InitialLocation string` (yaml: `initial_location`). Todos los sub-campos son opcionales; omitir el bloque `concept:` completo SHALL dejar `WorldConcept` en su valor cero.

#### Scenario: Bloque concept completo parseado de world.yaml
- **WHEN** un `world.yaml` contiene un bloque `concept:` con `premise`, `rules`, `flavor`, `character_spawn_rule` y `max_spawned_characters` seteados
- **THEN** los cinco campos SHALL estar poblados tras la carga

#### Scenario: Bloque concept parcial aceptado
- **WHEN** un `world.yaml` contiene `concept: { premise: "Bears disguised as humans" }` sin otros campos
- **THEN** `WorldConfig.Concept.Premise` SHALL ser `"Bears disguised as humans"`, todos los demás campos SHALL ser valor cero, y la carga SHALL retornar nil error

#### Scenario: Bloque concept ausente resulta en valor cero
- **WHEN** un `world.yaml` omite la clave `concept:` completamente
- **THEN** `WorldConfig.Concept` SHALL igualar el `WorldConcept{}` de valor cero y la carga SHALL retornar nil error

#### Scenario: Rules parseadas como lista ordenada
- **WHEN** `world.yaml` contiene `concept.rules` con tres entradas
- **THEN** `WorldConfig.Concept.Rules` SHALL tener longitud 3 con entradas en orden YAML

#### Scenario: character_spawn_rule parseado de world.yaml
- **WHEN** un `world.yaml` contiene `concept.character_spawn_rule: "All characters must be bears in human disguise"`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL ser ese string tras la carga

#### Scenario: character_spawn_rule ausente resulta en string vacío
- **WHEN** un `world.yaml` omite `character_spawn_rule` bajo `concept`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL ser string vacío y la carga SHALL retornar nil error

#### Scenario: max_spawned_characters parseado de world.yaml
- **WHEN** un `world.yaml` contiene `concept.max_spawned_characters: 3`
- **THEN** `WorldConfig.Concept.MaxSpawnedCharacters` SHALL ser `3` tras la carga

#### Scenario: max_spawned_characters ausente resulta en cero (ilimitado)
- **WHEN** un `world.yaml` omite `max_spawned_characters` bajo `concept`
- **THEN** `WorldConfig.Concept.MaxSpawnedCharacters` SHALL ser `0` y la carga SHALL retornar nil error

#### Scenario: initial_location parseado de world.yaml
- **WHEN** un `world.yaml` contiene `initial_location: "Convention Lobby"`
- **THEN** `WorldConfig.InitialLocation` SHALL ser `"Convention Lobby"` tras la carga

#### Scenario: initial_location ausente resulta en string vacío
- **WHEN** un `world.yaml` omite la clave `initial_location`
- **THEN** `WorldConfig.InitialLocation` SHALL ser string vacío y la carga SHALL retornar nil error
