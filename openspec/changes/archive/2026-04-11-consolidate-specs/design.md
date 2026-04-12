## Context

Los specs de `openspec/specs/` crecieron uno por change a lo largo de ~25 changes. Cada change introducía un spec nuevo o modificaba uno existente, pero nunca hubo una reorganización global. El resultado es que el contenido sobre una misma capacidad del sistema está repartido en múltiples archivos, algunos specs son tan pequeños que describen una sola función, y hay duplicación explícita (e.g., `character-initial-state` describe campos que ya están en `character-schema`).

El sistema tiene 24 specs actuales. La propuesta los reduce a 15 mediante fusiones y redistribuciones, sin cambiar ningún comportamiento del código.

## Goals / Non-Goals

**Goals:**
- Un spec por capacidad semántica del sistema (no por change histórica)
- Eliminar contenido duplicado entre specs
- Que el nombre de un spec responda a la pregunta "¿sobre qué parte del sistema busco?"
- Mantener todos los requirements y scenarios existentes — ninguno se pierde

**Non-Goals:**
- Cambiar ningún comportamiento del código
- Añadir nuevos requirements
- Modificar los specs de escenarios (`honey-heist`, `everyday-lives`, `doom-hell-crusade`) — están bien delimitados
- Cambiar los specs de infraestructura pura (`message-bus`, `llm-client`, `project-makefile`, `scenario-loader`)

## Decisions

### Decisión 1: Disolver `character-engine` en lugar de fusionarlo

`character-engine` contiene requirements de 4 dominios distintos: datos del personaje (struct, loader, type field), comportamiento runtime (inbox, memory), mensajería (CharChatReply), y output del motor (rendering, JSONL). Ningún spec existente tiene la misma mezcla. La decisión es redistribuir cada requirement al spec semánticamente correcto en lugar de fusionar el archivo entero en algún destino arbitrario.

Alternativa considerada: fusionar `character-engine` completo en `character-actor`. Descartada porque mezclaría datos de configuración YAML con comportamiento runtime.

### Decisión 2: Nombrar el spec consolidado de datos como `character` (no `character-schema`)

`character-schema` es un nombre que enfatiza la forma YAML, no la capacidad del sistema. El spec resultante describe todo lo que ES un personaje — su struct, sus campos, cómo se carga y cómo genera su prompt. `character` es más directo como punto de búsqueda.

### Decisión 3: Crear `simulation-output` como spec nuevo (no extender `simulation-engine`)

El motor de simulación orquesta el tick loop; la generación del resumen y las character cards son un sistema de salida separado que podría ejecutarse independientemente. Mantenerlos en specs distintos respeta esa separación. El rendering per-tick (líneas de acción/speech y JSONL) sí va a `simulation-engine` porque ocurre dentro del loop.

### Decisión 4: Fusionar `configuration` completamente (env-config + simulation-language)

`simulation-language` es efectivamente una fila más en la tabla de env vars de `env-config`. No tienen suficiente complejidad individual para justificar archivos separados. Un desarrollador que busca "¿cómo configuro el idioma?" debería ir al mismo lugar que "¿cómo configuro el modelo?".

### Decisión 5: No tocar los specs de escenarios ni infraestructura

Los specs de escenarios (`honey-heist`, `everyday-lives`, `doom-hell-crusade`) describen datasets concretos, no capacidades del sistema. Los specs de infraestructura (`message-bus`, `llm-client`, `scenario-loader`, `project-makefile`) ya tienen límites semánticos claros. Reorganizarlos aportaría poco.

## Risks / Trade-offs

**[Riesgo] Pérdida de contenido durante la migración** → Mitigación: Las tasks deben exigir migración explícita de cada requirement y cada scenario de los specs fuente antes de eliminarlos. Nada se elimina sin haber sido absorbido.

**[Trade-off] Los specs fusionados serán más largos** → `character` tendrá ~230 líneas, `world` ~320 líneas, `director` ~280 líneas. Esto es aceptable porque toda esa longitud describe una sola capacidad coherente. La alternativa (archivos cortos fragmentados) era peor para la búsqueda.

**[Trade-off] Los links entre specs (referencias cruzadas) dejarán de ser necesarios** → Cuando `character-schema` referenciaba `character-cover-identity`, era porque el contenido estaba partido. Al fusionarse, esas referencias se eliminan, lo que reduce la carga cognitiva.

## Migration Plan

1. Crear los specs nuevos/expandidos con el contenido consolidado
2. Verificar que ningún requirement o scenario se perdió comparando con los originals
3. Eliminar los 9 specs que quedan vacíos/absorbidos
4. No hay rollback de código necesario — solo archivos `.md`

## Open Questions

- ¿`character-judgment` merece fusionarse con `character-actor`? Es complejo (131 líneas) y describe un sistema relativamente autónomo. La decisión actual es dejarlo separado.
- ¿El spec `scenario-loader` debería renombrarse a `scenario` para consistencia con el patrón de nombres? No incluido en este change — cambio mínimo de valor discutible.
