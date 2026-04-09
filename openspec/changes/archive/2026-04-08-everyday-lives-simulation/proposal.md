## Why

Las simulaciones existentes operan en géneros con tensión externa explícita: noir de corrupción, robo con planificación táctica, cruzada demoníaca. No existe ninguna simulación que demuestre lo que el motor puede hacer cuando la única fuente de drama es la vida ordinaria — personas reales en situaciones reales, sin crimen que resolver ni objetivo que robar. Esta simulación prueba que los personajes con arquitectura psicológica profunda generan interacciones significativas incluso cuando no pasa "nada".

## What Changes

- Nuevo directorio `simulations/vida-cotidiana/` con `characters.yaml`, `world.yaml`, y `scenario.yaml`
- 5 personajes: residentes de un edificio de apartamentos en una ciudad cualquiera, un domingo ordinario
- Mundo: espacios compartidos del edificio y el barrio inmediato — lobby, azotea, escalera, cafetería de la esquina, parque
- Sin director de juego, sin misión, sin género: la premisa es que estas personas comparten paredes y a veces se cruzan
- Los personajes tienen vidas interiores complejas, tensiones no resueltas, y formas muy distintas de estar en el mundo

## Capabilities

### New Capabilities

- `everyday-lives-scenario`: Scenario de vida cotidiana — define el mundo, los personajes, y las reglas narrativas para una simulación sin género ni misión externa

### Modified Capabilities

_(ninguna — no se modifican specs existentes)_

## Impact

- Nuevos archivos de contenido únicamente: `simulations/vida-cotidiana/characters.yaml`, `world.yaml`, `scenario.yaml`
- No requiere cambios al motor, al loader, ni a ningún sistema existente
- El scenario loader ya soporta este formato — el cambio es puro contenido
- Sirve como caso de prueba de que el sistema funciona sin game director y sin tensión externa explícita
