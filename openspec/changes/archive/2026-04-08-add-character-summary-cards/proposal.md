## Why

After a simulation ends, the narrative summary gives a great story arc but leaves no structured record of who the characters were. Adding a character card section at the end of the summary gives readers an at-a-glance reference for each character's full profile, anchoring the narrative in concrete data.

## What Changes

- The summary output (printed to terminal and saved to file) SHALL append a "Character Cards" section after the narrative text.
- Each card SHALL display the key attributes of a character: name, age, occupation, appearance, motivation, fear, core belief, internal tension, emotional state, location, goals, voice profile, relational defaults, formative events, cover identity (if present), and dialogue examples.
- Fields that are empty or nil SHALL be omitted from the card.
- Game Director characters (`Type == "game_director"`) SHALL be excluded from the cards section.

## Capabilities

### New Capabilities
- `character-summary-cards`: Renders a structured character card section appended after the simulation narrative summary, displaying each non-director character's full attribute set.

### Modified Capabilities
- `simulation-summary`: The summary output now includes a character cards section after the narrative text.

## Impact

- `internal/summary` package: the function that builds and returns the summary string must append the character cards block.
- The saved `.md` file will include the cards section.
- No new dependencies required; all data is already available from `[]*character.Character`.
