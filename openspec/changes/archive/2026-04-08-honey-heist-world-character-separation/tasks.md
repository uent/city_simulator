## 1. Actualizar world.yaml

- [x] 1.1 Agregar `initial_location: Convention Lobby` al `world.yaml`
- [x] 1.2 Limpiar `details` del Convention Lobby: eliminar referencia a Lady Marmalade distrayendo guardias
- [x] 1.3 Limpiar `details` del Vendor Hall: eliminar referencia a Claws McGee moviéndose por el lado izquierdo
- [x] 1.4 Limpiar `details` del Security Office: eliminar referencia a Dr. Snuffles teniendo visión del cuarto
- [x] 1.5 Limpiar `details` del Vault Antechamber: eliminar "Honeydrop has been here for six minutes" y referencia al bypass kit
- [x] 1.6 Limpiar `details` del Vault: eliminar referencia al peso sustituto que lleva Honeydrop
- [x] 1.7 Limpiar `details` del Alley (Exit): eliminar referencia a Patches en la van
- [x] 1.8 Reescribir `initial_events` como setup de escena pre-heist (HoneyCon está por abrir, el equipo acaba de llegar, el Golden Comb está en exhibición)

## 2. Actualizar characters.yaml

- [x] 2.1 Agregar `inventory` a Honeydrop: bypass kit, peso sustituto calibrado a ±2g
- [x] 2.2 Agregar `initial_state` a Honeydrop: especialista de cajas fuertes lista para infiltrarse
- [x] 2.3 Agregar `initial_state` a Lady Marmalade: posando como sommelier de los Pirineos, lista para distraer
- [x] 2.4 Agregar `initial_state` a Patches: van posicionada y con motor listo, rutas de escape memorizadas
- [x] 2.5 Agregar `initial_state` a Dr. Snuffles: acceso al sistema de cámaras pendiente de activar
- [x] 2.6 Verificar que ningún personaje tiene `initial_location` como campo propio

## 3. Verificación

- [x] 3.1 Confirmar que ningún `details` de locación contiene los nombres: Grizwald, Honeydrop, Claws McGee, Lady Marmalade, Patches, Dr. Snuffles
- [x] 3.2 Confirmar que `initial_location` referencia una locación existente en el `world.yaml`
- [x] 3.3 Cargar el escenario con el simulator CLI (`--scenario honey-heist`) y verificar que no hay errores de parsing
