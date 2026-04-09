## 1. Character Struct Extensions

- [x] 1.1 Add `Appearance string` field (YAML key `appearance`) to `Character` struct in `internal/character/character.go`
- [x] 1.2 Add `CharacterJudgment` struct to `internal/character/character.go` with fields: `About`, `Name`, `Impression`, `Trust`, `Interest`, `Threat`, `FormedTick`, `UpdatedTick`
- [x] 1.3 Add `Judgments map[string]CharacterJudgment` field (YAML tag `"-"`) to `Character` struct
- [x] 1.4 Initialize `Judgments` map in `LoadCharacters` for each loaded character (avoid nil map panics)

## 2. Observable Profile and Snapshot Logic

- [x] 2.1 Define `ObservableProfile` struct in `internal/character/judgment.go` with fields: `ID`, `Name`, `Age`, `Occupation`, `EmotionalState`, `Appearance`, `Location`
- [x] 2.2 Implement `ObservableSnapshot(c Character) ObservableProfile` — use `CoverIdentity.Alias`/`Role` when non-nil, else `Character.Name`/`Occupation`; always use real `Age`, `EmotionalState`, `Appearance`, `Location`

## 3. Judgment Prompt Builders and Parser

- [x] 3.1 Implement `BuildJudgmentPrompt(judge Character, target ObservableProfile, language string) string` — includes judge's full profile + target's observable fields + JSON response instruction
- [x] 3.2 Implement `BuildUpdatePrompt(judge Character, target ObservableProfile, prior CharacterJudgment, recentHistory []string, language string) string` — same as above plus prior judgment and conversation history
- [x] 3.3 Implement `ParseJudgmentResponse(raw string, about string, name string) CharacterJudgment` — extracts JSON from raw LLM response, validates enum values (`trust`, `interest`, `threat`), falls back to `"medium"`/`"medium"`/`"none"` + `"No strong impression yet."` on parse failure
- [x] 3.4 Implement `FormatJudgmentForPrompt(j CharacterJudgment) string` — renders the judgment block injected into system prompts

## 4. Judgment Formation Functions

- [x] 4.1 Implement `FormInitialJudgments(ctx context.Context, chars []*Character, llmClient llmClient, language string) error` — fires N×(N-1) goroutines concurrently, writes results to `char.Judgments[targetID]`, collects errors (non-fatal: failed pairs get neutral fallback)
- [x] 4.2 Implement `FormJudgmentsForNew(ctx, newChar *Character, existing []*Character, llmClient, language)` — new character judges all existing; parallel calls
- [x] 4.3 Implement `FormJudgmentsOfNew(ctx, existing []*Character, newChar *Character, llmClient, language)` — all existing characters judge the new character; parallel calls

## 5. Engine Integration

- [x] 5.1 Add `pairConversations map[string]int` field to `Engine` struct in `internal/simulation/engine.go`; initialize in `NewEngine`
- [x] 5.2 Add `pairKey(idA, idB string) string` helper that returns `sorted(idA) + ":" + sorted(idB)` for symmetric pair keys
- [x] 5.3 Call `character.FormInitialJudgments` in `engine.Run()` before the tick loop begins; log completion
- [x] 5.4 In `registerSpawnedChars`, after starting the new character actor, call `FormJudgmentsForNew` and `FormJudgmentsOfNew` for the new character
- [x] 5.5 In the conversation prompt build section of `engine.Run()`, inject `FormatJudgmentForPrompt(pair.Initiator.Judgments[pair.Responder.ID])` into `initiatorSystem` (if judgment exists)
- [x] 5.6 In the conversation prompt build section of `engine.Run()`, inject `FormatJudgmentForPrompt(pair.Responder.Judgments[pair.Initiator.ID])` into `responderSystem` (if judgment exists)
- [x] 5.7 After each conversation, increment `pairConversations[pairKey(A.ID, B.ID)]`
- [x] 5.8 After incrementing, check `count % 10 == 0 && count > 0`; if true, call `UpdateJudgment` for both A→B and B→A concurrently using last 5 public conversation events from `world.EventLog` as history

## 6. Scenario YAML Updates — Appearance Field

- [x] 6.1 Add `appearance` to all 3 characters in `simulations/default/characters.yaml` (Elena Vasquez, Marcus Thorne, Nadia Osei)
- [x] 6.2 Add `appearance` to all 6 non-director characters in `simulations/honey-heist/characters.yaml` (Grizwald, Honeydrop, Claws McGee, Lady Marmalade, Patches, Dr. Snuffles) — write appearance based on their cover identity presentation
- [x] 6.3 Add `appearance` to all 2 non-director characters in `simulations/doom-hell-crusade/characters.yaml` (Doom Slayer, Vael)
- [x] 6.4 Add `appearance` to both characters in `simulations/test-scenario/characters.yaml` (Alice, Bob)

## 7. Verification

- [x] 7.1 Run `go build ./...` and confirm no compile errors
- [x] 7.2 Run a short simulation (`--turns 5`) on `default` scenario and confirm pre-sim judgment log output appears before tick 1
- [x] 7.3 Run a short simulation on `honey-heist` and confirm cover identity names appear in judgment blocks within prompts (check via log)
- [x] 7.4 Confirm simulation completes without errors when `appearance` is absent from a character (backward compat)
