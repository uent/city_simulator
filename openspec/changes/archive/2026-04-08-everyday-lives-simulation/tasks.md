## 1. Estructura del directorio

- [x] 1.1 Crear el directorio `simulations/vida-cotidiana/`
- [x] 1.2 Crear `simulations/vida-cotidiana/scenario.yaml` con `turns: 30`

## 2. Mundo y locaciones

- [x] 2.1 Crear `simulations/vida-cotidiana/world.yaml` con el bloque `concept` (premise, flavor, rules)
- [x] 2.2 Definir las 5 locaciones: Lobby del Edificio, Azotea, Escalera y Pasillos, Cafetería Marga, Parque del Barrio
- [x] 2.3 Añadir `description` y `details` a cada locación (details sin nombres de personajes)
- [x] 2.4 Establecer `initial_location: Lobby del Edificio`
- [x] 2.5 Añadir al menos 2 `initial_events` que planten semillas relacionales sin resolver

## 3. Personajes

- [x] 3.1 Crear `simulations/vida-cotidiana/characters.yaml` con los 5 residentes del edificio
- [x] 3.2 Escribir personaje 1 (20–30 años) con schema completo — estudiante o recién graduado con primer trabajo
- [x] 3.3 Escribir personaje 2 (20–30 años) con schema completo — en una relación que no sabe cómo nombrar
- [x] 3.4 Escribir personaje 3 (20–30 años) con schema completo — recién llegado a la ciudad, construyendo desde cero
- [x] 3.5 Escribir personaje 4 (20–30 años) con schema completo — con algo que esconder o no procesar
- [x] 3.6 Escribir personaje 5 (20–30 años) con schema completo — el más funcional en apariencia, el más frágil internamente

## 4. Verificación

- [x] 4.1 Confirmar que ningún personaje tiene campos `personality`, `backstory`, o `cover_identity`
- [x] 4.2 Confirmar que ningún personaje tiene `type: game_director`
- [x] 4.3 Confirmar que los `details` del world.yaml no contienen nombres de personajes
- [x] 4.4 Cargar el escenario con `--scenario vida-cotidiana` y verificar que no hay errores de parsing
