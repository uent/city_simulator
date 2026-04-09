## Context

Characters currently interact with zero prior context about each other. Each conversation starts cold — no instinct, no preconception, no bias. The only shared information is the zone roster (who is where), which gives names and locations but nothing psychological.

The engine builds system prompts for each character per tick in `engine.go`. Prompt content is constructed from the character's own profile, world state, inbox events, and zone context. There is currently no mechanism to inject one character's read of another.

The `CharacterActor` in `actor.go` handles conversation generation and maintains per-pair conversation history (`history map[string][]Turn`). The engine maintains no per-pair state beyond scheduling.

## Goals / Non-Goals

**Goals:**
- Each character forms a subjective first impression of every other character before tick 1
- Judgments are stored on the character and injected into conversation prompts
- Judgments refresh every 10 conversations between the same pair
- New characters (director-spawned) are immediately judged by and judge all existing characters
- Cover identities are respected — judgments target the observable persona, not the true identity

**Non-Goals:**
- Judgments do not affect movement decisions
- Judgments are not persisted to disk or included in JSONL output
- Judgment history (prior versions) is not tracked — only the current impression
- Characters do not share their judgments with each other (private, not spoken aloud)
- No UI or log output for individual judgments (internal state only)

## Decisions

### D1: New `judgment.go` file in `internal/character/`

All judgment logic lives in a single new file: `internal/character/judgment.go`. This keeps the `character` package self-contained and avoids scattering judgment concerns across the engine.

The file exposes:
- `ObservableSnapshot(c Character) ObservableProfile` — builds the filtered public view
- `BuildJudgmentPrompt(judge Character, target ObservableProfile, language string) string`
- `BuildUpdatePrompt(judge Character, target ObservableProfile, prior CharacterJudgment, history []string, language string) string`
- `ParseJudgmentResponse(raw string) (CharacterJudgment, error)` — parses LLM JSON output
- `FormatJudgmentForPrompt(j CharacterJudgment) string` — renders judgment as prompt-injectable text

**Alternative considered:** Placing judgment formation in the engine directly. Rejected because it would bloat `engine.go` and mix character-domain logic with simulation orchestration.

### D2: Observable snapshot uses CoverIdentity alias/role when present

When forming a judgment about character B, the judging character A only sees:
- Name: `CoverIdentity.Alias` if set, else `Character.Name`
- Occupation: `CoverIdentity.Role` if set, else `Character.Occupation`
- Age, EmotionalState, Appearance, Location always from the real Character struct

Private fields — Motivation, Fear, CoreBelief, InternalTension, FormativeEvents, Goals — are never included.

This is the central design choice for narrative realism: in honey-heist, every character judges the covers, not the bears underneath.

### D3: Judgment formation runs in parallel goroutines before tick 1

`FormInitialJudgments(ctx, chars, llmClient, language string)` iterates all pairs and fires one goroutine per judgment. For N characters, this is N×(N-1) concurrent LLM calls. Results are collected and written to `char.Judgments[targetID]` under a mutex.

For typical scenario sizes (2–8 characters) this is 2–56 calls — acceptable parallelism. A `sync.WaitGroup` with error collection handles completion.

**Alternative considered:** Sequential judgment formation. Rejected — adds unnecessary latency before tick 1 with no benefit.

### D4: Engine maintains `pairConversations map[string]int` for update tracking

The engine adds a new field tracking how many conversations each pair has had. The key is `sortedID(A) + ":" + sortedID(B)` so A↔B and B↔A map to the same counter.

After each conversation, the engine increments the counter. When `count % 10 == 0 && count > 0`, it calls `UpdateJudgment` for both A→B and B→A in parallel.

The update call receives the last 5 conversation turns from `world.EventLog` as context (filtering by participants), plus the prior judgment text.

**Alternative considered:** Tracking conversation count inside `CharacterActor`. Rejected — the actor doesn't own the counter for both directions, and the engine already coordinates post-conversation work (movement decisions, log output).

### D5: Judgment injected into system prompt at conversation build time

In `engine.go`, when building `initiatorSystem` and `responderSystem`, a judgment block is appended if it exists:

```
initiatorSystem += character.FormatJudgmentForPrompt(pair.Initiator.Judgments[pair.Responder.ID])
responderSystem  += character.FormatJudgmentForPrompt(pair.Responder.Judgments[pair.Initiator.ID])
```

The formatted block reads like:
```
Your prior impression of [Name] ([occupation]):
"[impression text]"
Trust: [level] | Interest: [level] | Perceived threat: [level]
```

Only the judgment for the current conversation partner is injected — not the full map. This keeps prompts lean.

### D6: LLM returns judgment as JSON

The judgment prompt instructs the LLM to respond with a JSON object:
```json
{
  "impression": "...",
  "trust": "high|medium|low|none",
  "interest": "high|medium|low",
  "threat": "high|medium|low|none"
}
```

`ParseJudgmentResponse` extracts the JSON from the raw response (tolerating surrounding text) and validates enum values, defaulting to "medium"/"low" on parse failure.

**Alternative considered:** Free-text response with regex extraction. Rejected — JSON parsing is more reliable and the structured fields (trust, interest, threat) are needed for the formatted prompt injection.

### D7: Spawned characters trigger judgment formation in `registerSpawnedChars`

After a character is registered in the bus and started, the engine calls:
1. `FormJudgmentsForNew(ctx, newChar, existingChars, llmClient, language)` — new char judges all existing (parallel)
2. `FormJudgmentsOfNew(ctx, existingChars, newChar, llmClient, language)` — all existing judge new char (parallel)

Both complete before the simulation tick continues.

## Risks / Trade-offs

**Latency at simulation start** → Pre-sim judgment adds N×(N-1) LLM calls before tick 1. For 8 characters, that's 56 calls. Mitigated by full parallelism; in practice this is bounded by LLM throughput, not sequential time. Acceptable trade-off for richer simulation.

**LLM JSON parse failures** → The judgment prompt may occasionally produce malformed JSON or invalid enum values. Mitigation: `ParseJudgmentResponse` has a lenient fallback that produces a neutral judgment (`impression: "No strong impression yet."`, all levels `"medium"`) rather than failing the simulation.

**Stale observable snapshot in updates** → The observable snapshot passed to update calls reflects the character's *current* state at update time, which may differ from when the judgment was formed (e.g., emotional state changed). This is actually desirable — the updater sees the character as they are now. No mitigation needed.

**Cover identity break scenarios** → If a character's cover is blown mid-simulation (via director action modifying character state), the judgment was formed against the cover. The next update (at the 10-interaction threshold) will reflect any observable changes, but the initial impression persists until then. This is a realistic limitation and creates interesting narrative tension.

**Memory pressure for large casts** → For very large character counts, the judgment map grows as O(N²) strings. Not a concern for current scenario sizes (≤10 characters).

## Migration Plan

1. Add `appearance` and `Judgments` fields — fully backward compatible; `appearance` is optional in YAML, `Judgments` is runtime-only
2. Add `appearance` to existing `characters.yaml` files — additive, no schema breaks
3. `FormInitialJudgments` is called unconditionally in `engine.Run()` before the tick loop — no flag needed
4. No rollback complexity: removing the feature means removing the pre-sim call and prompt injection lines; all other code is additive
