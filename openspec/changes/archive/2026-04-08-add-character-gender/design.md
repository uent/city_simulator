## Context

El struct `Character` en `internal/character/character.go` define todos los atributos del personaje. Actualmente no incluye género, por lo que el motor no puede referirse a un personaje con pronombres correctos ni incluirlo en la línea de identidad del system prompt. Los cuatro escenarios de simulación existentes tienen archivos `characters.yaml` que necesitarán el nuevo campo.

## Goals / Non-Goals

**Goals:**
- Añadir `Gender string` al struct `Character` (YAML key: `gender`)
- Actualizar `BuildSystemPrompt` para incluir el género en la línea de identidad
- Poblar el campo en los cuatro `characters.yaml` existentes

**Non-Goals:**
- Validación de valores de género (campo libre, no enum)
- Pronombres o lógica de concordancia automática en el motor
- Cambios en el esquema de la base de datos o persistencia

## Decisions

**Campo string libre (no enum)**
El género se representa como string libre (e.g. `"femenino"`, `"masculino"`, `"no binario"`). Un enum restringiría valores culturalmente y complicaría escenarios de simulación creativos. La coherencia la garantiza el autor del YAML.

**Campo opcional con comportamiento silent-omit**
Si `gender` se omite, `BuildSystemPrompt` no lo incluye en la línea de identidad (mantiene el comportamiento actual). Esto garantiza compatibilidad hacia atrás sin errores.

**Integración en línea de identidad**
La línea de identidad pasa de `"You are {Name}, a {Age}-year-old {Occupation}."` a `"You are {Name}, a {Age}-year-old {Gender} {Occupation}."` cuando `gender` está presente. No se añade sección separada para no fragmentar el prompt.

## Risks / Trade-offs

- **Riesgo: YAMLs de terceros sin el campo** → Mitigación: campo optional, carga sin error si se omite.
- **Trade-off: string libre vs enum** → Mayor flexibilidad a costa de inconsistencias tipográficas entre archivos. Aceptable para un simulador narrativo.

## Migration Plan

1. Añadir campo al struct en Go
2. Actualizar `BuildSystemPrompt`
3. Actualizar los cuatro `characters.yaml`
4. Sin rollback necesario — cambio aditivo y backward-compatible
