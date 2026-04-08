## Why

Al finalizar una simulación, no hay forma de revisar qué ocurrió durante la historia: los eventos, conversaciones y decisiones de los personajes se pierden. Generar un resumen narrativo persistente al terminar cada ejecución permitiría revisar y comparar historias de distintas simulaciones.

## What Changes

- Al finalizar la simulación, el motor genera automáticamente un resumen narrativo de los eventos ocurridos.
- El resumen se guarda en un archivo con timestamp en el directorio de la simulación activa, evitando sobrescribir resúmenes anteriores.
- El resumen incluye: eventos principales, conversaciones relevantes, decisiones de personajes y estado final del mundo.

## Capabilities

### New Capabilities

- `simulation-summary`: Generación y persistencia de un resumen narrativo al finalizar la simulación. El resumen es producido via LLM a partir del historial de eventos y se guarda en un archivo con nombre único por ejecución.

### Modified Capabilities

- `simulation-engine`: El motor debe invocar la generación del resumen como paso final antes de terminar la ejecución.

## Impact

- `internal/simulation/engine.go`: añadir llamada al generador de resumen al finalizar.
- `internal/llm/prompt.go`: nuevo prompt para sintetizar eventos en narrativa.
- Directorio de simulación activa (`simulations/<name>/`): nuevos archivos `summary-<timestamp>.md` por cada ejecución.
