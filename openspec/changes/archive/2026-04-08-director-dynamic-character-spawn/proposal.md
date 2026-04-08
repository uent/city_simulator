## Why

El director de juego actualmente solo puede introducir NPCs simples (`introduce_npc`) con campos mínimos. No existe mecanismo para que el director genere personajes completos y ricos de forma dinámica según las necesidades de la trama, ni forma de controlar qué tipo de personajes son válidos para el mundo de la simulación.

## What Changes

- **Nueva regla de creación de personajes en `WorldConfig`**: campo `character_spawn_rule` en `world.yaml` que describe las restricciones y guía narrativa para generar personajes dinámicos.
- **Nueva acción del director `spawn_character`**: permite al director crear un personaje con ficha completa (motivación, miedo, voz, etc.) guiado por la regla del mundo.
- **Bloqueo condicional**: si el mundo no define `character_spawn_rule`, la acción `spawn_character` no aparece en el prompt del director y no puede ser invocada.
- **Personajes creados dinámicamente participan en la simulación como actores completos**: tienen su propio actor LLM, memoria, bandeja de entrada, y aparecen en el estado del mundo.

## Capabilities

### New Capabilities

- `director-character-spawn`: Capacidad del director de crear personajes completos de forma dinámica, condicionada a la existencia de una regla de creación en el mundo.

### Modified Capabilities

- `director`: El prompt del director expone condicionalmente la acción `spawn_character` solo cuando existe `character_spawn_rule` en el mundo.
- `world-concept`: El `WorldConfig` agrega el campo `character_spawn_rule` para definir cómo deben ser los personajes generados dinámicamente.

## Impact

- `internal/world/state.go`: agregar campo `CharacterSpawnRule` a `WorldConcept` o `WorldConfig`
- `internal/director/prompt.go`: incluir `spawn_character` en el bloque `<tools>` condicionalmente
- `internal/director/actions_npc.go`: implementar `spawnCharacterAction`
- `internal/simulation/engine.go`: registrar la nueva acción y manejar personajes recién creados como actores activos
- `simulations/*/world.yaml`: los escenarios pueden opcionalmente definir `character_spawn_rule`
