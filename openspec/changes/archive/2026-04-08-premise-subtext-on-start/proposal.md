## Why

Al iniciar la simulación, el usuario no tiene contexto sobre la premisa del mundo que se está ejecutando. Mostrar el subtexto de la premisa al inicio mejora la experiencia y deja claro qué escenario se está simulando.

## What Changes

- Al arrancar la simulación, se imprime el bloque de concepto del mundo (`Premise`, `Flavor`, `Rules`) como subtexto visible en stdout, antes de que comiencen los ticks.

## Capabilities

### New Capabilities

- `simulation-premise-display`: Imprime en stdout el concepto del mundo (premise, flavor, rules) al inicio de la simulación, si el campo `Premise` no está vacío.

### Modified Capabilities

- `simulation-engine`: Se agrega la impresión del concepto del mundo al inicio de `Engine.Run`, usando el `Scenario.World.Concept` ya disponible en el `Config`.

## Impact

- `cmd/simulator/main.go` o `internal/simulation/engine.go`: se añade la lógica de impresión al inicio del run.
- No se requieren cambios en estructuras de datos; `WorldConcept` ya tiene los campos necesarios.
