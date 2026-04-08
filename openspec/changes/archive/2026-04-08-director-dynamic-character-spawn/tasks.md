## 1. World Model

- [x] 1.1 Agregar campo `CharacterSpawnRule string` (YAML: `character_spawn_rule`) a `WorldConcept` en `internal/world/state.go`
- [x] 1.2 Verificar que `WorldConcept` se deserializa correctamente con el nuevo campo desde `world.yaml` (campo opcional, default vacío)

## 2. Director Action

- [x] 2.1 Implementar `spawnCharacterAction` en `internal/director/actions_npc.go` con todos los args del spec (`id`, `name`, `age`, `occupation`, `motivation`, `fear`, `core_belief`, `internal_tension`, `formative_events`, `location`, `emotional_state`, `goals`)
- [x] 2.2 Validar args requeridos (`id`, `name`) y retornar error si faltan
- [x] 2.3 Validar que no exista ya un personaje con el mismo `id`; retornar error sin mutar estado si hay colisión
- [x] 2.4 Verificar que `state.Concept.CharacterSpawnRule != ""`; retornar error `"no character_spawn_rule defined"` si está vacío
- [x] 2.5 Emitir evento `world.Event{Type: "spawn", Visibility: "public"}` con descripción del personaje creado
- [x] 2.6 Registrar `spawnCharacterAction` en el registry en `internal/director/registry.go`

## 3. Director Prompt

- [x] 3.1 Modificar `BuildDirectorPrompt` en `internal/director/prompt.go` para incluir `spawn_character` en el bloque `<tools>` condicionalmente cuando `state.Concept.CharacterSpawnRule != ""`
- [x] 3.2 Incluir el texto de `character_spawn_rule` como contexto en la descripción de la herramienta `spawn_character`

## 4. Simulation Engine

- [x] 4.1 Verificar que el engine en `internal/simulation/engine.go` ya maneja el caso donde `*chars` crece durante la ejecución (personajes nuevos añadidos mid-tick o entre ticks)
- [x] 4.2 Agregar log `[spawn] created character <id> (<name>)` cuando `spawn_character` ejecuta exitosamente (en el engine o en la acción)

## 5. Scenario Example

- [x] 5.1 Agregar `character_spawn_rule` al `world.yaml` del escenario `honey-heist` como ejemplo de validación manual
