## Context

El director de juego puede actualmente introducir NPCs vía `introduce_npc`, pero solo con campos superficiales (nombre, rol, actitud, motivación, ubicación). Los personajes así creados no tienen ficha completa (psicología, voz, historia formativa) y no participan como actores LLM en la simulación.

La simulación ya tiene un modelo de personaje rico (`character.Character`) con motivación, miedo, creencia central, voz, eventos formativos, etc. El objetivo es que el director pueda generar personajes con esa riqueza, pero solo cuando el mundo define explícitamente cómo deben ser dichos personajes (`character_spawn_rule`).

## Goals / Non-Goals

**Goals:**
- Agregar `character_spawn_rule` a `WorldConcept` (o como campo de nivel superior en `WorldConfig`)
- Nueva acción del director `spawn_character` que genera un personaje completo con todos los campos del struct `Character`
- El bloque `<tools>` del director incluye `spawn_character` únicamente si `world.character_spawn_rule` no está vacío
- Los personajes creados dinámicamente se agregan a la simulación y participan como actores activos en el siguiente tick
- La acción falla en runtime si se invoca sin regla definida (doble protección)

**Non-Goals:**
- No se reemplaza `introduce_npc` (sigue existiendo para introducciones simples sin ficha completa)
- No se usa una llamada LLM separada para generar la ficha; el director recibe el esquema completo como args de la herramienta y lo llena él mismo en la misma respuesta
- No se persisten personajes dinámicos en disco entre sesiones

## Decisions

### Decisión 1: Ubicación de `character_spawn_rule`

**Elegido**: campo en `WorldConcept` como `CharacterSpawnRule string` (YAML: `character_spawn_rule`).

**Alternativa descartada**: campo en `WorldConfig` de nivel superior. Conceptualmente la regla es parte del concepto del mundo (qué tipo de seres/personas existen), no de la configuración de runtime.

**Rationale**: `WorldConcept` ya agrupa las reglas y el premise del mundo. `character_spawn_rule` es otra restricción narrativa, no un parámetro operativo.

### Decisión 2: Formato del args de `spawn_character`

**Elegido**: La acción acepta todos los campos del `Character` struct como args planos en el JSON: `id`, `name`, `age`, `occupation`, `motivation`, `fear`, `core_belief`, `internal_tension`, `formative_events` (array), `voice` (objeto), `location`, `goals` (array), `emotional_state`. Campos mínimos requeridos: `id`, `name`.

**Alternativa descartada**: Solo campos básicos + generar el resto internamente. Esto requeriría una llamada LLM adicional desde el engine, aumentando latencia y complejidad.

**Rationale**: El director ya es un LLM con contexto completo. Puede generar la ficha completa en una sola respuesta siguiendo la `character_spawn_rule`. Mantiene el flujo simple: una respuesta del director → una acción → un personaje.

### Decisión 3: Cuándo el personaje creado puede participar

**Elegido**: El personaje se registra inmediatamente en el slice `chars` del engine. Participa como actor en el **mismo tick** si el engine aún no lo ha procesado, o en el siguiente tick si ya pasó su turno.

**Alternativa descartada**: Cola de personajes pendientes procesada al inicio del tick siguiente. Más compleja sin beneficio real dado que el orden de actores dentro de un tick no tiene dependencias estrictas.

### Decisión 4: Visibilidad condicional en el prompt

**Elegido**: `BuildDirectorPrompt` recibe el `WorldConcept` (ya lo tiene via `state.Concept`) y agrega el bloque `spawn_character` al `<tools>` solo si `state.Concept.CharacterSpawnRule != ""`. La regla se incluye como contexto en el bloque de la herramienta.

**Rationale**: El director no debe ver una herramienta que no puede usar. Esto también evita que el director intente invocarla en mundos que no definen regla.

## Risks / Trade-offs

- **[Riesgo] El director genera una ficha incompleta o inconsistente** → Mitigación: campos mínimos validados en `Execute` (`id` y `name`). Campos opcionales se usan con defaults si faltan.
- **[Riesgo] ID duplicado al hacer spawn** → Mitigación: `Execute` verifica que no exista ya un personaje con ese ID; si existe, retorna error sin mutar el estado.
- **[Trade-off] El director genera la ficha completa en una sola llamada** → El LLM debe ser suficientemente capaz. Para modelos débiles la ficha puede ser genérica. Aceptable dado que `character_spawn_rule` guía la generación.

## Open Questions

- ¿Debe el engine loggear con nivel INFO cuando se crea un personaje dinámicamente? (Recomendado: sí, usando el prefijo `[spawn]`.)
- ¿Se debe agregar un ejemplo de `character_spawn_rule` en el escenario honey-heist para validación manual?
