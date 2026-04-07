package world

import "fmt"

// Location represents a named place in the city.
type Location struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Event records something that happened during the simulation.
type Event struct {
	Tick         int      `yaml:"-"          json:"tick"`
	Type         string   `yaml:"event_type" json:"type"`
	Description  string   `yaml:"description" json:"description"`
	Participants []string `yaml:"participants,omitempty" json:"participants,omitempty"`
}

// WorldConfig is loaded from a scenario's world.yaml.
type WorldConfig struct {
	Locations     []Location `yaml:"locations"`
	InitialEvents []Event    `yaml:"initial_events"`
}

// State holds the current world state for a simulation run.
type State struct {
	Tick      int
	TimeOfDay string
	Locations []Location
	EventLog  []Event
}

var timeOfDayLabels = []string{"morning", "afternoon", "evening", "night"}

// NewState creates a fresh world from a WorldConfig with tick 0 and time "morning".
// InitialEvents from cfg are pre-populated into the event log.
func NewState(cfg WorldConfig) *State {
	events := make([]Event, len(cfg.InitialEvents))
	copy(events, cfg.InitialEvents)
	return &State{
		Tick:      0,
		TimeOfDay: "morning",
		Locations: cfg.Locations,
		EventLog:  events,
	}
}

// AdvanceTick increments the tick and updates the time-of-day label.
// Cycle: morning (0-5), afternoon (6-11), evening (12-17), night (18-23), repeating.
func (s *State) AdvanceTick() {
	s.Tick++
	s.TimeOfDay = timeOfDayLabels[(s.Tick/6)%4]
}

// AppendEvent adds an event to the world log.
func (s *State) AppendEvent(e Event) {
	s.EventLog = append(s.EventLog, e)
}

// Summary returns a concise description of the world for LLM context.
func (s *State) Summary() string {
	recent := s.EventLog
	if len(recent) > 5 {
		recent = recent[len(recent)-5:]
	}
	if len(recent) == 0 {
		return fmt.Sprintf("It is currently %s. No events have occurred yet.", s.TimeOfDay)
	}
	summary := fmt.Sprintf("It is currently %s. Recent events:\n", s.TimeOfDay)
	for _, e := range recent {
		summary += fmt.Sprintf("- [tick %d] %s\n", e.Tick, e.Description)
	}
	return summary
}
