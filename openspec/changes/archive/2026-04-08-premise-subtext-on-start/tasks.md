## 1. Implementar la función de impresión del concepto

- [x] 1.1 En `internal/simulation/engine.go`, añadir una función `printWorldConcept(w io.Writer, concept world.Concept)` que imprima el bloque con el formato definido en el spec
- [x] 1.2 Omitir la línea `Flavor:` si `Concept.Flavor` es vacío
- [x] 1.3 Omitir la sección `Rules:` si `Concept.Rules` está vacío o es nil

## 2. Integrar en Engine.Run

- [x] 2.1 Llamar a `printWorldConcept` al inicio de `Engine.Run`, antes del primer tick, usando el `OutputWriter` del `Config` (o stdout directo si es más adecuado)
- [x] 2.2 Verificar que el bloque no se imprime cuando `Concept.Premise` es vacío

## 3. Verificación manual

- [x] 3.1 Ejecutar la simulación con un escenario que tenga `concept.premise` en `world.yaml` y confirmar que el bloque aparece antes del `[Tick 1]`
- [x] 3.2 Ejecutar con un escenario sin `concept:` y confirmar que no se imprime nada extra

