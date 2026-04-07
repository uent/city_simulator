## Context

The `simulation-scenarios` change introduces the `simulations/<name>/` directory convention and the `internal/scenario` package. `honey-heist` is a pure data deliverable — three YAML files dropped into `simulations/honey-heist/` — that exercises that infrastructure without modifying any Go code.

The scenario premise: a crew of bears with criminal specialties converges on HoneyCon, the world's largest honey convention, to steal the legendary Golden Comb from the secured vault.

## Goals / Non-Goals

**Goals:**
- Ship a fully self-contained, launchable scenario at `simulations/honey-heist/`
- Demonstrate a non-trivial cast (5–7 characters with distinct roles) and a multi-room world
- Use `scenario.yaml` to set tone (system prompt prefix) and a short default run length so the heist resolves quickly

**Non-Goals:**
- Any Go code changes — the scenario loader handles all reading
- Balancing or tuning LLM outputs — the YAML just primes the characters and world
- Adding new YAML fields not already supported by the scenario-loader spec

## Decisions

### 1. Cast design — 6 bears with distinct archetypes

```
characters.yaml
```

| Name | Role | Trait |
|------|------|-------|
| Grizwald | Mastermind | calculating, terse |
| Honeydrop | Safe-cracker | nervous, precise |
| Claws McGee | Muscle | jovial, impulsive |
| Lady Marmalade | Grifter | charming, duplicitous |
| Patches | Wheelman | taciturn, loyal |
| Dr. Snuffles | Tech Bear | verbose, anxious |

Six characters gives enough conversational variety without overwhelming the LLM context budget per tick. Each character's `personality` field in `characters.yaml` encodes their archetype in 1–2 sentences so the LLM can stay in character.

**Alternative considered:** 10+ characters — rejected because it would balloon token usage per tick and dilute individual character voice.

### 2. World layout — 6 locations inside and around HoneyCon

```
world.yaml
```

| Location | Description |
|----------|-------------|
| Convention Lobby | Busy entrance, security desk, bear crowds |
| Vendor Hall | Chaotic honey stalls, cover for movement |
| Security Office | Guard station with camera feeds |
| Vault Antechamber | Locked corridor leading to the vault |
| Vault | The Golden Comb rests here |
| Alley (Exit) | Dark alley behind the building — the getaway point |

Linear escalation (lobby → vendor hall → security → antechamber → vault → alley) creates natural narrative tension. Initial events in `world.yaml` seed the scene: guards have just changed shift, the convention is at peak attendance.

### 3. `scenario.yaml` for tone and run length

```yaml
turns: 20
system_prompt_prefix: |
  You are narrating a comedic bear heist. All characters are bears with criminal
  specialties. Keep the tone playful but tense. Characters care deeply about honey.
```

20 turns is short enough for the heist to feel urgent. The system prompt prefix is prepended to the engine's base system prompt so every LLM call carries the heist framing.

**Rationale:** Without a tone anchor the LLM may narrate a mundane city simulation rather than a heist comedy. The prefix is the cheapest way to enforce genre.

## Risks / Trade-offs

- **Depends on `simulation-scenarios` being implemented first** → This change should not be applied until `simulation-scenarios` tasks are complete. The `honey-heist` tasks.md will call this out explicitly.
- **`system_prompt_prefix` field may not be supported yet** → If the scenario-loader spec does not include this field, we omit it from `scenario.yaml` and rely on character personality strings alone for tone. Tasks will verify this before writing the field.
- **LLM creative drift** → Characters may ignore their assigned roles over long runs. Mitigated by the short 20-turn default and strong personality strings in `characters.yaml`.

## Migration Plan

No migration required — this is additive. Running `--scenario honey-heist` before the `simulations/honey-heist/` directory exists will produce a clear error from the scenario loader ("scenario not found").

## Open Questions

- Does `scenario.yaml` support a `system_prompt_prefix` field in the current loader spec? Check `openspec/changes/simulation-scenarios/specs/scenario-loader/spec.md` before writing `scenario.yaml`.
- Should `world.yaml` include `initial_events` to seed the scene? Recommended yes — at least one event: "The annual HoneyCon convention is at peak attendance. Guards changed shift 10 minutes ago."
