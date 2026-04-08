package director

import (
	"fmt"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// moveNPCAction updates a character's Location and appends a public event.
type moveNPCAction struct{}

func (moveNPCAction) Name() string { return "move_npc" }

func (moveNPCAction) Summary(args map[string]any) string {
	id, idOk := stringArg(args, "id")
	dest, destOk := stringArg(args, "destination")
	if idOk && destOk {
		return "move_npc: " + id + " → " + dest
	}
	return "move_npc"
}

func (moveNPCAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	id, ok := stringArg(args, "id")
	if !ok {
		return fmt.Errorf("move_npc: missing required arg 'id'")
	}
	dest, ok := stringArg(args, "destination")
	if !ok {
		return fmt.Errorf("move_npc: missing required arg 'destination'")
	}
	reason, _ := stringArg(args, "reason")

	c := findChar(*chars, id)
	if c == nil {
		return fmt.Errorf("move_npc: character %q not found", id)
	}
	c.Location = dest

	desc := fmt.Sprintf("%s moves to %s.", c.Name, dest)
	if reason != "" {
		desc = fmt.Sprintf("%s moves to %s (%s).", c.Name, dest, reason)
	}
	state.AppendEvent(world.Event{
		Type:        "movement",
		Description: desc,
		Visibility:  "public",
	})
	return nil
}

// introduceNPCAction appends a new character to the engine's character slice.
type introduceNPCAction struct{}

func (introduceNPCAction) Name() string { return "introduce_npc" }

func (introduceNPCAction) Summary(args map[string]any) string {
	name, nameOk := stringArg(args, "name")
	role, _ := stringArg(args, "role")
	if nameOk {
		if role != "" {
			return "introduce_npc: " + name + " (" + role + ")"
		}
		return "introduce_npc: " + name
	}
	return "introduce_npc"
}

func (introduceNPCAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	id, ok := stringArg(args, "id")
	if !ok {
		return fmt.Errorf("introduce_npc: missing required arg 'id'")
	}
	name, ok := stringArg(args, "name")
	if !ok {
		return fmt.Errorf("introduce_npc: missing required arg 'name'")
	}
	role, _ := stringArg(args, "role")
	attitude, _ := stringArg(args, "attitude")
	motivation, _ := stringArg(args, "motivation")
	location, _ := stringArg(args, "location")

	if attitude == "" {
		attitude = "neutral"
	}

	newChar := &character.Character{
		ID:             id,
		Name:           name,
		Occupation:     role,
		EmotionalState: attitude,
		Motivation:     motivation,
		Location:       location,
		MaxMemory:      20,
		Inbox:          []world.Event{},
	}
	*chars = append(*chars, newChar)

	state.AppendEvent(world.Event{
		Type:        "arrival",
		Description: fmt.Sprintf("%s (%s) arrives in the city.", name, role),
		Visibility:  "public",
	})
	return nil
}

// addNPCConditionAction appends a condition string to a character's EmotionalState.
type addNPCConditionAction struct{}

func (addNPCConditionAction) Name() string { return "add_npc_condition" }

func (addNPCConditionAction) Summary(args map[string]any) string {
	id, idOk := stringArg(args, "id")
	condition, condOk := stringArg(args, "condition")
	if idOk && condOk {
		return "add_npc_condition: " + id + " += " + condition
	}
	return "add_npc_condition"
}

func (addNPCConditionAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	id, ok := stringArg(args, "id")
	if !ok {
		return fmt.Errorf("add_npc_condition: missing required arg 'id'")
	}
	condition, ok := stringArg(args, "condition")
	if !ok {
		return fmt.Errorf("add_npc_condition: missing required arg 'condition'")
	}
	c := findChar(*chars, id)
	if c == nil {
		return fmt.Errorf("add_npc_condition: character %q not found", id)
	}
	if c.EmotionalState == "" || c.EmotionalState == "neutral" {
		c.EmotionalState = condition
	} else {
		c.EmotionalState = c.EmotionalState + ", " + condition
	}
	return nil
}

// removeNPCConditionAction removes the first occurrence of a condition from EmotionalState.
type removeNPCConditionAction struct{}

func (removeNPCConditionAction) Name() string { return "remove_npc_condition" }

func (removeNPCConditionAction) Summary(args map[string]any) string {
	id, idOk := stringArg(args, "id")
	condition, condOk := stringArg(args, "condition")
	if idOk && condOk {
		return "remove_npc_condition: " + id + " -= " + condition
	}
	return "remove_npc_condition"
}

func (removeNPCConditionAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	id, ok := stringArg(args, "id")
	if !ok {
		return fmt.Errorf("remove_npc_condition: missing required arg 'id'")
	}
	condition, ok := stringArg(args, "condition")
	if !ok {
		return fmt.Errorf("remove_npc_condition: missing required arg 'condition'")
	}
	c := findChar(*chars, id)
	if c == nil {
		return fmt.Errorf("remove_npc_condition: character %q not found", id)
	}
	// Remove first occurrence of the condition (comma-separated list).
	parts := strings.Split(c.EmotionalState, ", ")
	filtered := parts[:0]
	removed := false
	for _, p := range parts {
		if !removed && strings.TrimSpace(p) == condition {
			removed = true
			continue
		}
		filtered = append(filtered, p)
	}
	c.EmotionalState = strings.Join(filtered, ", ")
	if c.EmotionalState == "" {
		c.EmotionalState = "neutral"
	}
	return nil
}

// spawnCharacterAction creates a new character with a full character sheet.
// The action is only valid when state.Concept.CharacterSpawnRule is non-empty.
type spawnCharacterAction struct{}

func (spawnCharacterAction) Name() string { return "spawn_character" }

func (spawnCharacterAction) Summary(args map[string]any) string {
	name, nameOk := stringArg(args, "name")
	id, _ := stringArg(args, "id")
	if nameOk {
		return "spawn_character: " + name + " (" + id + ")"
	}
	return "spawn_character"
}

func (spawnCharacterAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	if state.Concept.CharacterSpawnRule == "" {
		return fmt.Errorf("spawn_character: no character_spawn_rule defined in this world")
	}

	// Enforce the cap if set (0 = unlimited).
	if max := state.Concept.MaxSpawnedCharacters; max > 0 && state.SpawnedCharacters >= max {
		return fmt.Errorf("spawn_character: max spawned characters (%d) already reached", max)
	}

	id, ok := stringArg(args, "id")
	if !ok || id == "" {
		return fmt.Errorf("spawn_character: missing required arg 'id'")
	}
	name, ok := stringArg(args, "name")
	if !ok || name == "" {
		return fmt.Errorf("spawn_character: missing required arg 'name'")
	}

	// Check for duplicate ID.
	if findChar(*chars, id) != nil {
		return fmt.Errorf("spawn_character: character with id %q already exists", id)
	}

	occupation, _ := stringArg(args, "occupation")
	motivation, _ := stringArg(args, "motivation")
	fear, _ := stringArg(args, "fear")
	coreBelief, _ := stringArg(args, "core_belief")
	internalTension, _ := stringArg(args, "internal_tension")
	location, _ := stringArg(args, "location")
	emotionalState, _ := stringArg(args, "emotional_state")
	if emotionalState == "" {
		emotionalState = "neutral"
	}

	var age int
	if v, ok := args["age"]; ok {
		switch n := v.(type) {
		case float64:
			age = int(n)
		case int:
			age = n
		}
	}

	var formativeEvents []string
	if v, ok := args["formative_events"]; ok {
		if raw, ok := v.([]any); ok {
			for _, item := range raw {
				if s, ok := item.(string); ok {
					formativeEvents = append(formativeEvents, s)
				}
			}
		}
	}

	var goals []string
	if v, ok := args["goals"]; ok {
		if raw, ok := v.([]any); ok {
			for _, item := range raw {
				if s, ok := item.(string); ok {
					goals = append(goals, s)
				}
			}
		}
	}

	newChar := &character.Character{
		ID:              id,
		Name:            name,
		Age:             age,
		Occupation:      occupation,
		Motivation:      motivation,
		Fear:            fear,
		CoreBelief:      coreBelief,
		InternalTension: internalTension,
		FormativeEvents: formativeEvents,
		Location:        location,
		Goals:           goals,
		EmotionalState:  emotionalState,
		MaxMemory:       20,
		Inbox:           []world.Event{},
	}
	*chars = append(*chars, newChar)
	state.SpawnedCharacters++

	state.AppendEvent(world.Event{
		Type:        "spawn",
		Description: fmt.Sprintf("%s (%s) appears in the city.", name, occupation),
		Visibility:  "public",
	})
	return nil
}

// findChar returns the first character in chars with the given ID, or nil.
func findChar(chars []*character.Character, id string) *character.Character {
	for _, c := range chars {
		if c.ID == id {
			return c
		}
	}
	return nil
}
