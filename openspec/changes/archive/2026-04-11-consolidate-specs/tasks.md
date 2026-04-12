## 1. Crear spec `character` (nuevo — absorbe 3 specs existentes + partes de character-engine)

- [x] 1.1 Crear `openspec/specs/character/spec.md` copiando el contenido consolidado de `specs/character/spec.md` de este change
- [x] 1.2 Verificar que el spec contiene todos los requirements de `character-schema` (psychological fields, formative events, voice, relational defaults, dialogue examples, inventory, initial state, gender, appearance, BuildSystemPrompt, ObservableSnapshot, CHARACTER_RULES.md)
- [x] 1.3 Verificar que el spec contiene todos los requirements de `character-cover-identity` (CoverIdentity struct, parsing YAML, inyección en system prompt)
- [x] 1.4 Verificar que el spec contiene todos los requirements de `character-initial-state` (inventory field, initial_state field)
- [x] 1.5 Verificar que el spec contiene los requirements de datos de `character-engine` (Character struct completo, inbox, type field, LoadCharacters, memory buffer)
- [x] 1.6 Eliminar `openspec/specs/character-schema/spec.md` y su directorio
- [x] 1.7 Eliminar `openspec/specs/character-cover-identity/spec.md` y su directorio
- [x] 1.8 Eliminar `openspec/specs/character-initial-state/spec.md` y su directorio

## 2. Actualizar spec `character-actor` (absorbe character-expression y partes runtime de character-engine)

- [x] 2.1 Agregar al `openspec/specs/character-actor/spec.md` existente los requirements de `character-expression`: Expression type, ParseExpression, instrucción de formato en system prompt, FormatExpression
- [x] 2.2 Agregar los requirements runtime de `character-engine`: CharChatReply payload fields (InitiatorSpeech/Action, ResponderSpeech/Action)
- [x] 2.3 Verificar que el spec unificado de `character-actor` no tiene contenido duplicado con el spec de `character`
- [x] 2.4 Eliminar `openspec/specs/character-expression/spec.md` y su directorio

## 3. Eliminar spec `character-engine` (disuelto)

- [x] 3.1 Verificar que "Character struct with persona fields" está en el nuevo `character`
- [x] 3.2 Verificar que "Character inbox for private events" está en el nuevo `character`
- [x] 3.3 Verificar que "Character type field" está en el nuevo `character`
- [x] 3.4 Verificar que "Character loader from YAML file" está en el nuevo `character`
- [x] 3.5 Verificar que "Per-character memory buffer" está en el nuevo `character`
- [x] 3.6 Verificar que "CharChatReply payload fields" está en el `character-actor` actualizado
- [x] 3.7 Verificar que "Engine renders action and speech separately" está en el `simulation-engine` actualizado
- [x] 3.8 Verificar que "JSONL log entry includes action and speech fields" está en el `simulation-engine` actualizado
- [x] 3.9 Eliminar `openspec/specs/character-engine/spec.md` y su directorio

## 4. Actualizar spec `world` (absorbe world-concept)

- [x] 4.1 Agregar al `openspec/specs/world/spec.md` existente los requirements de `world-concept`: WorldConcept struct, todos sus campos YAML (premise, rules, flavor, character_spawn_rule, max_spawned_characters), initial_location
- [x] 4.2 Verificar que no hay duplicación con el contenido ya existente en `world`
- [x] 4.3 Eliminar `openspec/specs/world-concept/spec.md` y su directorio

## 5. Actualizar spec `director` (absorbe director-character-spawn)

- [x] 5.1 Agregar al `openspec/specs/director/spec.md` existente los requirements de `director-character-spawn`: spawn_character action, gating por character_spawn_rule, cap de max_spawned_characters, registro de actores spawneados, AddCharacter en Scheduler
- [x] 5.2 Verificar que el `spawn_character` tool schema en BuildDirectorPrompt (ya en director) está alineado con la nueva definición de la acción
- [x] 5.3 Eliminar `openspec/specs/director-character-spawn/spec.md` y su directorio

## 6. Actualizar spec `simulation-engine` (absorbe premise-display y rendering de character-engine)

- [x] 6.1 Agregar al `openspec/specs/simulation-engine/spec.md` existente los requirements de `simulation-premise-display`: formato del bloque `=== World Concept ===`, condiciones de omisión de Flavor y Rules
- [x] 6.2 Agregar los requirements de rendering de `character-engine`: formato per-tick con acción/speech en líneas separadas, entrada JSONL con campos de acción y speech
- [x] 6.3 Verificar que el requirement "World concept printed at run start" ya en simulation-engine es consistente con el bloque de formato absorbido (fusionar si están solapados)
- [x] 6.4 Eliminar `openspec/specs/simulation-premise-display/spec.md` y su directorio

## 7. Crear spec `simulation-output` (nuevo — absorbe simulation-summary y character-summary-cards)

- [x] 7.1 Crear `openspec/specs/simulation-output/spec.md` copiando el contenido consolidado de `specs/simulation-output/spec.md` de este change
- [x] 7.2 Verificar que el spec contiene todos los requirements de `simulation-summary` (GenerateSummary, SaveSummary con timestamp)
- [x] 7.3 Verificar que el spec contiene todos los requirements de `character-summary-cards` (renderCharacterCards, todos los campos de la card, exclusión de directores)
- [x] 7.4 Eliminar `openspec/specs/simulation-summary/spec.md` y su directorio
- [x] 7.5 Eliminar `openspec/specs/character-summary-cards/spec.md` y su directorio

## 8. Crear spec `configuration` (nuevo — absorbe env-config y simulation-language)

- [x] 8.1 Crear `openspec/specs/configuration/spec.md` copiando el contenido consolidado de `specs/configuration/spec.md` de este change
- [x] 8.2 Verificar que el spec contiene todos los requirements de `env-config` (priority order, variables table, .env.example, .gitignore)
- [x] 8.3 Verificar que el spec contiene todos los requirements de `simulation-language` (idioma en todos los prompts, valor libre, override CLI)
- [x] 8.4 Eliminar `openspec/specs/env-config/spec.md` y su directorio
- [x] 8.5 Eliminar `openspec/specs/simulation-language/spec.md` y su directorio

## 9. Verificación final

- [x] 9.1 Confirmar que el total de specs en `openspec/specs/` bajó de 24 a 15
- [x] 9.2 Verificar que los specs sin cambios siguen intactos: `character-judgment`, `message-bus`, `llm-client`, `scenario-loader`, `project-makefile`, `honey-heist-scenario`, `everyday-lives-scenario`, `doom-hell-crusade-scenario`
- [x] 9.3 Hacer una búsqueda de referencias cruzadas en specs restantes para asegurarse de que no apuntan a specs eliminados
