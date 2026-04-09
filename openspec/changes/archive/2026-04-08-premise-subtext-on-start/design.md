## Context

El simulador ya carga el concepto del mundo (`WorldConcept`) desde `world.yaml` via `scenario.Load`. El struct `WorldConcept` expone `Premise`, `Flavor` y `Rules`. Al arrancar, `main.go` imprime una línea de estado pero no muestra nada sobre la premisa del escenario. El engine recibe el `Scenario` completo en su `Config`.

## Goals / Non-Goals

**Goals:**
- Imprimir el concepto del mundo (premise, flavor, rules) al inicio de `Engine.Run`, antes del primer tick.
- No requerir cambios en estructuras de datos ni en la carga de escenarios.

**Non-Goals:**
- Formatear la salida como rich-text, colores o markdown renderizado.
- Mostrar el subtexto en el output JSONL.
- Internacionalizar el encabezado del bloque de concepto.

## Decisions

**Dónde imprimir**: En `Engine.Run` al inicio, no en `main.go`.
- Razón: el engine ya tiene acceso al `Scenario` y es la unidad responsable de la ejecución de la simulación. Poner la lógica ahí mantiene `main.go` liviano y la responsabilidad cohesionada.
- Alternativa descartada: imprimir en `main.go` después de `engine.Run` — no, tiene que ser antes de los ticks.

**Cuándo omitir**: Solo imprimir si `Concept.Premise != ""`. Si el campo está vacío, no se imprime nada (compatibilidad con escenarios sin concepto definido).

**Formato de salida**:
```
=== World Concept ===
Premise: <premise>
Flavor:  <flavor>        (omitir si vacío)
Rules:
  - <rule1>              (omitir sección si Rules es vacío)
  - <rule2>
=====================
```

## Risks / Trade-offs

- [Riesgo] Escenarios legacy sin `concept:` en `world.yaml` → Mitigación: el bloque se omite completamente si `Premise == ""`.
- [Trade-off] El output va a stdout, no al JSONL; quien parsee el log no verá la premisa. Esto es intencional (es información para el operador, no para el análisis de la simulación).
