## Context

The city simulator already supports multi-character scenarios loaded from `simulations/<name>/` directories. The honey-heist scenario proved the full pipeline: world.yaml, characters.yaml, scenario.yaml, and dynamic director-spawned characters. The DOOM hell crusade scenario reuses this same pipeline but inverts the character ratio: one pre-loaded hero, all antagonists and allies spawned dynamically by the director.

## Goals / Non-Goals

**Goals:**
- Create a fully playable DOOM scenario with a single pre-loaded protagonist (Doom Guy)
- Design the world as a 5-location Hell campaign with clear dramatic arc
- Configure the director to spawn characters from a typed pool (demon renegade, trapped soul, prince's herald)
- Demonstrate that the simulator supports single-protagonist + director-as-ensemble pattern

**Non-Goals:**
- No changes to Go source code or existing specs
- No combat system or health points — narrative tension only
- Not a direct retelling of any specific DOOM game plot

## Decisions

**Decision: Single pre-loaded character**
The scenario sets up only Doom Guy in characters.yaml. All other characters (demon renegades, trapped souls, heralds of the Prince) are spawned by the director as the narrative demands. This tests the `spawn_character` director action in a high-stakes solo context.

Alternative considered: pre-load 2-3 demon allies. Rejected because it dilutes the "lone marine in Hell" tension central to the DOOM identity.

**Decision: 5 locations, 3 narrative acts**
- Act 1 (Arrival): The Flesh Gate, The Lava Wastes
- Act 2 (Ascent): The Cathedral of Bone, The Necropolis Vault
- Act 3 (Confrontation): The Throne Sanctum

This mirrors classic DOOM level structure (entry → exploration → boss arena) and gives the director clear environmental cues for pacing.

**Decision: Hellstone Convergence as the central threat**
Rather than a vague "destroy Earth", the Demon Prince (Malphas, Prince of the Seventh Circle) is activating the Hellstone Convergence — an ancient ritual that will collapse the dimensional membrane between Hell and Earth. This gives the simulation a concrete ticking-clock goal and a specific object (the Hellstone) that Doom Guy must destroy.

**Decision: Director persona as tactical intelligence**
The director character is framed as a spectral Watcher — an ancient Hell entity that observes but doesn't interfere directly. Its role is to introduce complications (spawning reinforcements, revealing intel, triggering environmental events) rather than narrate. This keeps the tone dark and diegetic.

## Risks / Trade-offs

- [Risk] Director spawning too many demons creates chaos without narrative → Mitigation: `max_spawned_characters: 4` and `character_spawn_rule` explicitly restricts spawn types to three archetypes
- [Risk] Single-character scenario feels lonely without NPC interaction → Mitigation: the director is instructed to spawn a trapped human soul early (turn 3–5) to give Doom Guy an information source and moral anchor

## Open Questions

- Should Doom Guy have any dialogue examples that break the "silent protagonist" convention? Current decision: yes, sparse internal monologue style ("Rip and tear. But first — information.")
