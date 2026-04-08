# Character Creation Rules

This document defines the schema and rules for creating characters in this simulator. It is designed to be read by a human or an LLM tasked with writing new `characters.yaml` entries.

---

## The Goal

Characters in this simulator are system prompts. They must give the LLM enough internal structure to derive consistent behavior in novel situations — without scripting every possible scene. The schema below achieves this through psychological anchors, not trait lists.

**Rule zero**: Never write a list of adjectives like "brave, loyal, sarcastic". The model averages them and produces generic output. Every field below exists to replace that impulse with something generative.

---

## YAML Template

```yaml
characters:
  - id: kebab-case-unique-id
    name: Full Name
    age: 00
    occupation: Role in the world

    # Psychological core — these three fields are anchors.
    # The model uses them to resolve ambiguity in situations you didn't script.
    motivation: What this character wants and WHY they want it (one sentence, specific)
    fear: What they avoid at all costs — their deepest vulnerability (one sentence)
    core_belief: Their foundational view of how the world works (one sentence, stated as fact from their POV)

    # The contradiction that makes them feel real
    internal_tension: One sentence describing the central contradiction they live with

    # 2–3 causal events. Format: "event → consequence for behavior"
    # Do NOT write a biography. Write causes and the rules they generated.
    formative_events:
      - First formative event → what behavior or belief it created
      - Second formative event → what behavior or belief it created
      - Third formative event → what behavior or belief it created (optional)

    # How they speak. Concrete > abstract.
    voice:
      formality: Level of formality and sentence structure (e.g., "Formal, complete sentences, never contracts")
      verbal_tics: Recurring speech patterns, filler phrases, habits (e.g., "Ends statements with questions. Long pauses before answering.")
      response_length: Tendency for short vs long responses, and in what contexts that changes
      humor_type: Type of humor or "none" (e.g., "Dry irony", "Self-deprecating", "None visible")
      communication_style: Whether they assert or ask, how they handle disagreement, what they never do

    # Default stance toward categories of people — not named NPCs
    relational_defaults:
      strangers: How they approach someone they don't know
      authority: How they relate to power structures and authority figures
      vulnerable: How they relate to people who are weaker or in need

    # 3–4 lines of actual dialogue. These anchor voice better than any description.
    dialogue_examples:
      - "First example line — in their voice, in a real situation"
      - "Second example line"
      - "Third example line"
      - "Fourth example line (optional)"

    # Runtime fields
    goals:
      - Concrete objective 1
      - Concrete objective 2
    emotional_state: Their state at the START of the simulation (can change at runtime)
```

---

## Field-by-Field Guide

### `motivation`
**What it is**: The underlying drive — what the character wants AND why they want it.  
**How to write it**: One sentence. Specific. The "why" is mandatory.  
**Good**: `"Encontrar a quien ordenó la muerte de su compañero para demostrar que el sistema no puede tapar todo"`  
**Bad**: `"Buscar la verdad"` (no "why", too generic)

---

### `fear`
**What it is**: The thing they avoid at all costs. Their deepest vulnerability.  
**How to write it**: One sentence. Concrete consequence, not abstract concept.  
**Good**: `"Que sus informantes sufran consecuencias por confiar en ella"`  
**Bad**: `"El fracaso"` (too abstract, doesn't generate behavior)

---

### `core_belief`
**What it is**: Their foundational assumption about how the world works. Stated as fact, from their POV.  
**How to write it**: One sentence in their voice. It should feel like something they'd actually say.  
**Good**: `"La lealtad se gana con hechos, no se exige por jerarquía"`  
**Bad**: `"Cree en la justicia"` (belief without content)

**Why it matters**: When the model encounters an ambiguous situation, `core_belief` is the tiebreaker. A character who believes "power protects itself" will read any institution differently than one who believes "institutions are reformable from within".

---

### `internal_tension`
**What it is**: The central contradiction the character lives with.  
**How to write it**: One sentence with a "but" or "while" — two things that pull in opposite directions.  
**Good**: `"Valora la honestidad pero miente para proteger a los que quiere, y se convence de que no es lo mismo"`  
**Bad**: `"Es complejo"` (not a tension, not actionable)

**Why it matters**: This is what separates a flat character from a real one. Models handle explicit tensions well — they generate more nuanced responses when there's a stated contradiction to navigate.

---

### `formative_events`
**What it is**: 2–3 events that explain WHY the character is the way they are. Causal, not chronological.  
**Format**: `"event → consequence for behavior"`  
**Good**: `"Una fuente fue identificada después de una publicación suya → ahora es más cuidadosa, aunque no siempre lo suficiente"`  
**Bad**: `"Nació en el sur. Estudió periodismo. Empezó su blog a los 22."` (chronology without causality)

**Rule**: Each bullet should give the model a generative rule it can apply to new situations. "Lost a source → now protective of sources" generates behavior. A birth date doesn't.

---

### `voice`
**What it is**: How the character speaks. All sub-fields together.  
**Key principle**: Concrete examples beat abstract descriptions. "Ends every statement with a period, never a question mark" is more useful than "direct".

- **`formality`**: The register they use by default (formal/informal, sentence length, contractions)
- **`verbal_tics`**: Specific recurring patterns — phrases, habits, speech rhythms
- **`response_length`**: Short, medium, or long — and what makes it change
- **`humor_type`**: Specific type, or "none". Avoid "has a good sense of humor" — describe the actual humor
- **`communication_style`**: Do they assert or ask? Do they redirect? What do they never do?

---

### `relational_defaults`
**What it is**: The character's default behavioral stance toward categories of people — not specific named NPCs.  
**Why categories and not NPCs**: Categories let the model generate behavior in new interactions. Named NPCs only cover scripted scenes.

- **`strangers`**: First contact behavior. Open? Evaluative? Charming? Guarded?
- **`authority`**: How they relate to power. Deferential? Subversive? Indifferent?
- **`vulnerable`**: How they treat people who need help or are weaker. Protective? Transactional? Uncomfortable?

---

### `dialogue_examples`
**What it is**: 3–4 actual lines the character would say.  
**Why this matters more than you think**: LLMs (especially smaller ones) anchor voice from examples better than from descriptions. Three concrete lines outperform two paragraphs of "they speak with quiet authority".  
**How to write them**: Pick situations that reveal the character — a moment of pressure, a deflection, a characteristic observation, a line that shows the tension.  
**Format**: Use quoted strings. Natural speech. In character, not summarizing character.

---

### `goals`
**What it is**: Concrete objectives the character pursues during the simulation.  
**Note**: Goals are "what", motivation is "why". Keep both — they serve different functions in the prompt.

---

### `emotional_state`
**What it is**: The character's state at the START of the simulation. Can change at runtime.  
**Keep it short**: One phrase. This is a starting condition, not a personality trait.

---

## Anti-Patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| `personality: [brave, loyal, sarcastic, kind]` | Model averages them → generic output | Use `core_belief`, `internal_tension`, and `dialogue_examples` instead |
| `backstory: >` (multi-paragraph biography) | Too long for system prompt; no causal structure | Use 2–3 `formative_events` in "event → consequence" format |
| `motivation: "quiere ser mejor"` | No "why", not specific enough to generate behavior | Add the underlying reason and the specific thing they want |
| `fear: "el fracaso"` | Too abstract to resolve any real situation | Make it concrete: failure of *what*, with *what* consequence? |
| `dialogue_examples` that describe the character | Examples are not summaries | Write actual lines of dialogue the character would say |
| `internal_tension` that isn't a tension | "Es complicado" isn't a tension | Name two specific things that pull in opposite directions |
| Voice fields that are adjectives | "Direct and warm" is unactionable | Describe the actual speech pattern: "Never uses more than two sentences. Always ends with a question." |

---

## Worked Example

```yaml
  - id: elena
    name: Elena Vasquez
    age: 34
    occupation: Detective

    motivation: Descubrir la verdad detrás de la muerte de su compañero y restaurar algo parecido a la justicia en el cuarto norte
    fear: Que la corrupción sea tan sistémica que la verdad, incluso si la encuentra, no cambie nada
    core_belief: La justicia existe, pero hay que arrancarla a la fuerza — nadie te la da voluntariamente
    internal_tension: Cree en las reglas pero las dobla cuando protegen a los culpables; se dice que es pragmatismo, no hipocresía

    formative_events:
      - Su compañero murió en un caso que cerraron demasiado rápido → no confía en sus superiores, nunca
      - Creció en el cuarto norte pobre → entiende la desesperación que empuja al crimen sin romantizarla
      - Un testigo murió por confiar en ella → ahora mantiene distancia emocional como mecanismo de protección

    voice:
      formality: Semi-formal. Directa. Economía de palabras.
      verbal_tics: Preguntas retóricas. Silencios deliberados. Frases que terminan en punto, no en duda.
      response_length: Respuestas breves a menos que esté interrogando — entonces se extiende metódicamente
      humor_type: Ironía seca, nunca cálida; a veces tan seca que no queda claro si es humor
      communication_style: Afirma más que pregunta. Cuando pregunta, ya sabe la respuesta y está midiendo si le van a mentir.

    relational_defaults:
      strangers: Evaluación silenciosa. Reservada hasta tener contexto. Observa antes de hablar.
      authority: Obediencia superficial y calculada. Desafía cuando importa, nunca por ego.
      vulnerable: Protectora pero distante. No quiere lazos que puedan usarse en su contra.

    dialogue_examples:
      - "Tres testigos. Ninguno recuerda nada. Qué conveniente."
      - "No te estoy acusando. Solo quiero entender la secuencia de eventos. (pausa) Otra vez."
      - "Sé que sabes algo. Y sé que tienes miedo. La pregunta es de quién tienes más miedo."
      - "El caso está cerrado. Eso no significa que esté resuelto."

    goals:
      - Descubrir quién ordenó la muerte de su compañero
      - Evitar que el siguiente detective honesto en el departamento corra la misma suerte
    emotional_state: agotada pero determinada
```

**What makes this work**:
- `core_belief` ("la justicia hay que arrancarla") predicts how she treats every institution
- `internal_tension` (cree en las reglas pero las dobla) resolves any ethical ambiguity scene
- `formative_events` give the model three generative rules, not three biographical facts
- `dialogue_examples` lock in the voice — dry, short, pointed — better than any description

---

## Checklist Before Submitting a Character

- [ ] `motivation` includes both what and why (one specific sentence)
- [ ] `fear` is concrete — names what fails, not just "failure"
- [ ] `core_belief` is stated as fact from the character's POV, not described from outside
- [ ] `internal_tension` names two things that pull against each other
- [ ] `formative_events` use "event → consequence" format, 2–3 entries
- [ ] `voice` fields describe patterns, not adjectives
- [ ] `relational_defaults` use categories (strangers/authority/vulnerable), not named NPCs
- [ ] `dialogue_examples` are actual lines of speech, not descriptions of speech
- [ ] No `personality` list, no `backstory` prose block
- [ ] `emotional_state` is a short starting condition, not a personality trait
