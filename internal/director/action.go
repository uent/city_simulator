package director

import (
	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// Action is the interface every director action must implement.
// chars is a pointer to the engine's character slice so actions like introduce_npc
// can append new characters and have them visible to the engine after dispatch.
type Action interface {
	Name() string
	Execute(args map[string]any, state *world.State, chars *[]*character.Character) error
	// Summary returns a concise human-readable description of the action and its
	// key argument values (e.g. "set_weather: storm"). Falls back to Name() if
	// required args are missing.
	Summary(args map[string]any) string
}

// ToolCall represents a single parsed tool call from the director's LLM response.
type ToolCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}
