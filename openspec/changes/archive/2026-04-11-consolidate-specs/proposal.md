## Why

Los 24 specs actuales crecieron incrementalmente (uno por change), lo que resulta en fragmentación y duplicación: contenido sobre el `Character` struct está repartido en 3 specs, el world concept está separado del world state, y specs como `simulation-premise-display` o `simulation-language` describen una sola función en lugar de una capacidad coherente. Esto hace que sea difícil saber dónde buscar algo y difícil mantener los specs sincronizados.

## What Changes

- Disolver `character-engine` redistribuyendo sus requirements a los specs correctos
- Fusionar `character-schema` + `character-cover-identity` + `character-initial-state` → `character`
- Fusionar `character-actor` + `character-expression` + partes de `character-engine` → `character-actor`
- Fusionar `world` + `world-concept` → `world`
- Fusionar `director` + `director-character-spawn` → `director`
- Fusionar `simulation-engine` + `simulation-premise-display` + partes de `character-engine` → `simulation-engine`
- Fusionar `simulation-summary` + `character-summary-cards` → `simulation-output`
- Fusionar `env-config` + `simulation-language` → `configuration`
- Eliminar los 9 specs disueltos: `character-engine`, `character-cover-identity`, `character-initial-state`, `world-concept`, `director-character-spawn`, `simulation-premise-display`, `simulation-language`, `character-summary-cards`, `env-config`
- Sin cambios funcionales al sistema — solo reorganización documental

## Capabilities

### New Capabilities

- `character`: Spec maestro del modelo de datos del personaje — struct, todos los campos YAML (schema + cover identity + initial state), `LoadCharacters`, `BuildSystemPrompt`. Reemplaza y absorbe `character-schema`, `character-cover-identity`, `character-initial-state`, y la parte de datos de `character-engine`.
- `simulation-output`: Generación del resumen narrativo final, persistencia a archivo, y renderizado de character cards. Absorbe `simulation-summary` y `character-summary-cards`.
- `configuration`: Configuración de la simulación vía env vars y CLI flags, incluyendo idioma. Absorbe `env-config` y `simulation-language`.

### Modified Capabilities

- `character-actor`: Absorbe requirements de expresiones, historial de conversación, memory buffer, inbox, y `CharChatReply` que estaban dispersos en `character-engine` y `character-expression`. Los requirements de comportamiento no cambian, solo se consolidan.
- `world`: Absorbe el `WorldConcept` struct y su parsing YAML de `world-concept`. Sin cambios de comportamiento.
- `director`: Absorbe todos los requirements de `director-character-spawn`. Sin cambios de comportamiento.
- `simulation-engine`: Absorbe el rendering per-tick (acción/speech) y el JSONL logging de `character-engine`, y la lógica de premise display de `simulation-premise-display`. Sin cambios de comportamiento.

## Impact

- Solo afecta archivos en `openspec/specs/` — no hay cambios de código
- Los 9 specs a eliminar dejan de existir; su contenido migra a los specs de destino
- El total baja de 24 a 15 specs
- Las referencias cruzadas entre specs (e.g., `character-schema` referenciando `character-cover-identity`) se eliminan ya que todo queda en un solo archivo
