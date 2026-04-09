## Why

Los personajes actuales no tienen campo de género, lo que obliga al motor a inferirlo del nombre o dejarlo ambiguo — esto produce prompts genéricamente escritos y limita la coherencia narrativa cuando otros sistemas necesitan referirse a un personaje por género.

## What Changes

- Se añade el campo `gender` (string) al struct `Character` con clave YAML `gender`.
- `BuildSystemPrompt` incluye el género en la línea de identidad cuando está presente.
- Todos los archivos `characters.yaml` existentes reciben el campo `gender` con el valor correcto para cada personaje.

## Capabilities

### New Capabilities

*(ninguna — el campo se integra en la capacidad existente)*

### Modified Capabilities

- `character-schema`: añade el campo `gender` a la estructura del personaje y actualiza `BuildSystemPrompt` para incluirlo en la línea de identidad.

## Impact

- `internal/character/character.go` — añadir campo `Gender string` al struct y actualizar `BuildSystemPrompt`
- `simulations/default/characters.yaml` — añadir `gender` a elena, marcus, nadia
- `simulations/honey-heist/characters.yaml` — añadir `gender` a todos los personajes
- `simulations/doom-hell-crusade/characters.yaml` — añadir `gender` a todos los personajes
- `simulations/test-scenario/characters.yaml` — añadir `gender` a todos los personajes
- Sin cambios en APIs externas ni dependencias
