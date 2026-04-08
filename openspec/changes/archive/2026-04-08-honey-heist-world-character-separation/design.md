## Context

El `world.yaml` de honey-heist actualmente mezcla dos responsabilidades distintas:

1. **Propiedades físicas del espacio**: sensores, delays, dead zones de cámaras, restricciones de radio, rutas.
2. **Estado inicial de personajes**: dónde está cada uno, qué lleva encima, qué acción táctica está ejecutando al arranque.

Esto fuerza la simulación a comenzar in media res (el loop de cámaras ya está activo, Honeydrop lleva 6 minutos en el Vault, Patches ya esperó 4 horas). El resultado es un mundo acoplado al equipo específico que lo habita — no puede reutilizarse con otro elenco ni con otra disposición táctica.

## Goals / Non-Goals

**Goals:**
- `world.yaml` describe únicamente física y espacio + un `initial_location` único donde todos los personajes arrancan.
- `characters.yaml` define `inventory` (objetos que cada personaje lleva) e `initial_state` (estado táctico al inicio).
- Los `initial_events` del mundo se reescriben como setup de escena pre-heist.
- La tensión narrativa emerge durante la simulación, no se hereda del estado inicial.

**Non-Goals:**
- No se modifica el motor de simulación ni el scenario loader.
- No se incorpora lógica de posicionamiento dinámico de personajes.
- No se aplica esta separación a otros escenarios existentes.

## Decisions

### 1. `initial_location` como campo único en el mundo, no por personaje

**Decisión**: el mundo define un único `initial_location` (string que referencia el nombre de una locación existente). Todos los personajes arrancan ahí.

**Alternativa descartada**: `initial_location` como campo por personaje en `characters.yaml`. Permite más flexibilidad pero introduce riesgo de inconsistencia (un personaje puede referenciar una locación que no existe en el mundo). La validación se vuelve más compleja sin ganancia narrativa real para este escenario.

**Rationale**: la consistencia supera la flexibilidad aquí. Si en el futuro se necesita dispersión inicial, puede resolverse con un campo `initial_positions` opcional en `scenario.yaml` que sobreescriba el default.

### 2. `inventory` e `initial_state` como campos opcionales en cada personaje

**Decisión**: agregar `inventory` (lista de strings) e `initial_state` (string descriptivo) como campos opcionales en el esquema de personaje.

**Rationale**: son datos que pertenecen al personaje — lo que lleva y cómo arranca la historia. Mantenerlos en el personaje permite que el Director y los demás personajes tengan acceso a esta información de forma natural durante la simulación.

### 3. Reescritura de `initial_events` como pre-heist

**Decisión**: los `initial_events` del mundo pasan de describir estado mid-operación a describir el contexto de arranque: HoneyCon está por abrir, el equipo acaba de llegar, el Golden Comb está en exhibición.

**Rationale**: si la simulación empieza antes del heist, los eventos iniciales deben reflejar ese punto de partida. Los eventos tácticos (loop de cámaras activo, Reginald distraído) se convierten en oportunidades que los personajes deben crear, no en estado heredado.

## Risks / Trade-offs

- **[Riesgo] El spec de `honey-heist-scenario` requiere actualización** → Los requirements sobre el roster de personajes y el world layout incluyen detalles que cambian con esta separación. Se crean specs delta para reflejar los cambios.
- **[Trade-off] La simulación es más larga** → Al arrancar antes, hay más turns necesarios para llegar al momento de tensión. El valor narrativo (la tensión emerge en vez de heredarse) justifica el trade-off.
- **[Riesgo] El motor no lee `initial_location` ni los nuevos campos de personaje** → Esta propuesta es solo de datos de simulación (YAMLs). Si el motor necesita consumir estos campos, requerirá cambios adicionales fuera del alcance de este cambio.
