## ADDED Requirements

<!-- Spec consolidado: absorbe env-config y simulation-language en una sola capacidad. -->

### Requirement: Orden de prioridad de configuración
El binario SHALL aplicar valores de configuración en el siguiente orden de prioridad (de mayor a menor):

1. **CLI flags** — pasados explícitamente en el momento de invocación
2. **Variables de entorno** — del shell o cargadas desde `.env`
3. **scenario.yaml overrides** — defaults por-escenario en el escenario cargado
4. **Defaults hardcodeados** — valores de fallback compilados

#### Scenario: Env var gana sobre scenario.yaml
- **WHEN** `SIM_TURNS=20` está seteada y el `scenario.yaml` define `turns: 5`
- **THEN** el binario SHALL usar `20`

#### Scenario: CLI flag gana sobre env var
- **WHEN** `SIM_TURNS=20` está seteada y se pasa `--turns=15`
- **THEN** el binario SHALL usar `15`

#### Scenario: scenario.yaml gana sobre default hardcodeado
- **WHEN** no hay env var ni CLI flag para `turns` y `scenario.yaml` define `turns: 5`
- **THEN** el binario SHALL usar `5`

---

### Requirement: Variables de entorno soportadas
El binario SHALL leer las siguientes variables de entorno como valores default para sus flags:

| Variable | Flag | Tipo | Default |
|---|---|---|---|
| `OLLAMA_URL` | `--ollama-url` | string | `http://localhost:11434` |
| `OLLAMA_MODEL` | `--model` | string | `llama3` |
| `SIM_TURNS` | `--turns` | int | `10` |
| `SIM_SEED` | `--seed` | int64 | `0` |
| `SIM_OUTPUT` | `--output` | string | `simulation_output.jsonl` |
| `SIM_SCENARIO` | `--scenario` | string | `default` |
| `SIM_LANGUAGE` | `--language` | string | `""` (vacío — sin instrucción de idioma) |

#### Scenario: Env var setea el default cuando el flag no se provee
- **WHEN** `OLLAMA_MODEL=mistral` está en el entorno y `--model` no se pasa
- **THEN** el binario SHALL usar `mistral` como nombre del modelo

#### Scenario: CLI flag sobreescribe env var
- **WHEN** `OLLAMA_MODEL=mistral` está seteada y se pasa también `--model=llama3`
- **THEN** el binario SHALL usar `llama3`

#### Scenario: Env var numérica inválida usa el default hardcodeado
- **WHEN** `SIM_TURNS=abc` está seteada (no-numérica)
- **THEN** el binario SHALL loguear un warning y usar el valor default hardcodeado (`10`)

#### Scenario: SIM_LANGUAGE setea el idioma para los prompts
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada y `--language` no se pasa
- **THEN** el binario SHALL usar `"Spanish"` como valor de idioma para todos los prompts LLM

---

### Requirement: Idioma inyectado en todos los prompts
El valor resuelto de `SIM_LANGUAGE` / `--language` SHALL inyectarse en el system prompt de cada actor de personaje, el Game Director, el prompt de movimiento, y el prompt de generación de resumen, instruyendo al LLM a responder en ese idioma.

Cuando el idioma está seteado, cada prompt SHALL terminar con `"Respond in {language}."`.

Cuando el idioma no está seteado (vacío), ningún prompt SHALL contener una instrucción de idioma.

El valor SHALL aceptarse como string libre — el binario SHALL NOT validar ni normalizar el identificador de idioma.

#### Scenario: Instrucción de idioma aparece en el prompt del personaje
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada
- **THEN** el system prompt de cada personaje SHALL terminar con `"Respond in Spanish."`

#### Scenario: Instrucción de idioma aparece en el prompt del director
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada
- **THEN** el system prompt del Game Director SHALL contener `"Respond in Spanish."`

#### Scenario: Instrucción de idioma aparece en el prompt de movimiento
- **WHEN** `SIM_LANGUAGE=French` está seteada
- **THEN** el prompt de movimiento de cada personaje SHALL terminar con `"Respond in French."`

#### Scenario: Instrucción de idioma aparece en el resumen
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada
- **THEN** el system prompt del resumen SHALL contener `"Respond in Spanish."`

#### Scenario: Sin instrucción cuando idioma no está seteado
- **WHEN** `SIM_LANGUAGE` no está seteada y `--language` no se pasa
- **THEN** ningún prompt SHALL contener una línea de instrucción de idioma

#### Scenario: CLI flag sobreescribe env var de idioma
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada y se pasa `--language=English`
- **THEN** el binario SHALL usar `"English"` como valor de idioma

#### Scenario: Tag BCP-47 aceptado sin error
- **WHEN** `SIM_LANGUAGE=es` está seteada
- **THEN** los prompts SHALL contener `"Respond in es."` sin error

---

### Requirement: Archivo .env.example en la raíz del proyecto
El repositorio SHALL contener un archivo `.env.example` documentando cada variable de entorno soportada con comentarios que describen su propósito y el valor default. Incluye `SIM_LANGUAGE`.

#### Scenario: Developer puede copiar el ejemplo para comenzar
- **WHEN** un developer ejecuta `cp .env.example .env`
- **THEN** el archivo `.env` resultante SHALL contener todas las variables con sus valores default, listo para personalizar

---

### Requirement: .env excluido del control de versiones
El archivo `.gitignore` SHALL contener una entrada para `.env` de modo que los archivos de entorno locales nunca se commiteen.

#### Scenario: .env es ignorado por git
- **WHEN** existe un archivo `.env` en la raíz del proyecto
- **THEN** `git status` SHALL NOT listarlo como archivo tracked o untracked
