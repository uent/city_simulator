## 1. Update Character Struct

- [x] 1.1 Add `VoiceProfile` struct to `internal/character/character.go` with fields: `Formality`, `VerbalTics`, `ResponseLength`, `HumorType`, `CommunicationStyle` (all string, yaml-tagged)
- [x] 1.2 Add `RelationalProfile` struct with fields: `Strangers`, `Authority`, `Vulnerable` (all string, yaml-tagged)
- [x] 1.3 Remove `Personality []string` and `Backstory string` from `Character` struct
- [x] 1.4 Add new fields to `Character`: `Motivation`, `Fear`, `CoreBelief`, `InternalTension` (string), `FormativeEvents []string`, `Voice VoiceProfile`, `RelationalDefaults RelationalProfile`, `DialogueExamples []string`

## 2. Update System Prompt Builder

- [x] 2.1 Rewrite `BuildSystemPrompt` in `internal/llm/prompt.go` to use the new structured template (identity → motivation → fear → core_belief → internal_tension → formative_events → voice → relational_defaults → goals → emotional_state → dialogue_examples → closing instruction)
- [x] 2.2 Ensure sections with empty/nil fields are silently omitted from the output

## 3. Rewrite Default Scenario Characters

- [x] 3.1 Rewrite `simulations/default/characters.yaml` for Elena Vasquez using the new schema (remove personality/backstory, add all new fields with 3 formative events, voice block, relational_defaults, 3–4 dialogue examples)
- [x] 3.2 Rewrite `simulations/default/characters.yaml` for Marcus Thorne using the new schema
- [x] 3.3 Rewrite `simulations/default/characters.yaml` for Nadia Osei using the new schema

## 4. Rewrite Honey Heist Characters

- [x] 4.1 Rewrite `simulations/honey-heist/characters.yaml` for Grizwald using the new schema
- [x] 4.2 Rewrite for Honeydrop using the new schema
- [x] 4.3 Rewrite for Claws McGee using the new schema
- [x] 4.4 Rewrite for Lady Marmalade using the new schema
- [x] 4.5 Rewrite for Patches using the new schema
- [x] 4.6 Rewrite for Dr. Snuffles using the new schema

## 5. Add Character Generation Rulebook

- [x] 5.1 Create `simulations/CHARACTER_RULES.md` containing: the complete YAML template with all fields annotated, field-by-field descriptions, anti-pattern guidance, and one fully worked example character in the new schema

## 6. Verify

- [x] 6.1 Run `go build ./...` and confirm no compilation errors
- [x] 6.2 Run the simulator with `--scenario default` and inspect the generated system prompt to confirm new sections appear
- [x] 6.3 Run the simulator with `--scenario honey-heist` and confirm all 6 bear characters load without errors
