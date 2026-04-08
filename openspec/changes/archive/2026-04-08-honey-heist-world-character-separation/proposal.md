## Why

El `world.yaml` de honey-heist mezcla propiedades físicas del espacio con estado inicial de personajes (posiciones, inventario, ventanas tácticas activas), lo que acopla el mundo al equipo específico que lo habita y fuerza la simulación a comenzar in media res. Separar estas responsabilidades hace el mundo reutilizable y permite que la tensión narrativa emerja durante la simulación en lugar de heredarse.

## What Changes

- `world.yaml` define únicamente física y espacio: propiedades de locaciones, rutas, restricciones ambientales, y un único `initial_location` donde todos los personajes comienzan.
- `world.yaml` elimina toda referencia a personajes específicos dentro de los `details` de locaciones y en los `initial_events`.
- `world.yaml` reescribe `initial_events` como setup de escena pre-heist (HoneyCon está por comenzar) en lugar de estado mid-operación.
- `characters.yaml` incorpora `inventory` (objetos que cada personaje lleva) e `initial_state` (estado táctico de arranque) para los personajes que lo requieren.
- La simulación comienza antes del heist: los personajes parten del `initial_location` y se dispersan, los temporizadores se activan cuando llegan a sus posiciones.

## Capabilities

### New Capabilities

- `character-initial-state`: Capacidad de definir inventario y estado inicial táctico directamente en el esquema de personaje.

### Modified Capabilities

- `world-concept`: El concepto de mundo incorpora `initial_location` como campo único de punto de partida compartido; los `details` de locaciones dejan de incluir estado de personajes.
- `honey-heist-scenario`: Los archivos `world.yaml` y `characters.yaml` del escenario se actualizan para reflejar la separación. La simulación arranca antes del golpe.

## Impact

- `simulations/honey-heist/world.yaml`: edición de `details` en todas las locaciones, reescritura de `initial_events`, adición de `initial_location`.
- `simulations/honey-heist/characters.yaml`: adición de `inventory` e `initial_state` a los personajes relevantes (Honeydrop, Patches, Lady Marmalade, Dr. Snuffles).
- Sin impacto en el motor de simulación ni en otros escenarios.
