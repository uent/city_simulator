package director

import (
	"fmt"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// Registry maps action names to their implementations.
type Registry struct {
	actions map[string]Action
}

// NewRegistry creates a Registry pre-loaded with all built-in director actions.
func NewRegistry() *Registry {
	r := &Registry{actions: make(map[string]Action)}

	// Environment actions
	r.register(setWeatherAction{})
	r.register(setTimeOfDayAction{})
	r.register(setAtmosphereAction{})

	// NPC actions
	r.register(moveNPCAction{})
	r.register(introduceNPCAction{})
	r.register(addNPCConditionAction{})
	r.register(removeNPCConditionAction{})

	// World / event actions
	r.register(modifyLocationAction{})
	r.register(triggerEncounterAction{})
	r.register(triggerEventAction{})
	r.register(revealInformationAction{})
	r.register(escalateTensionAction{})
	r.register(narrateAction{})

	return r
}

func (r *Registry) register(a Action) {
	r.actions[a.Name()] = a
}

// Names returns the list of all registered action names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.actions))
	for n := range r.actions {
		names = append(names, n)
	}
	return names
}

// Summarize returns the Summary of the named action, or name alone if the action
// is not registered.
func (r *Registry) Summarize(name string, args map[string]any) string {
	if a, ok := r.actions[name]; ok {
		return a.Summary(args)
	}
	return name
}

// Dispatch finds the action by name and calls Execute. Returns an error if the
// action is unknown or Execute fails. State is not mutated on error.
func (r *Registry) Dispatch(name string, args map[string]any, state *world.State, chars *[]*character.Character) error {
	a, ok := r.actions[name]
	if !ok {
		return fmt.Errorf("unknown director action %q", name)
	}
	return a.Execute(args, state, chars)
}
