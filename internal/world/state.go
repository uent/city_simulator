package world

import (
	"fmt"
	"log"
	"strings"
)

// Location represents a named place in the city.
type Location struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Details     string `yaml:"details"` // private; shown only to characters present here
}

// Event records something that happened during the simulation.
type Event struct {
	Tick             int      `yaml:"-"                        json:"tick"`
	Type             string   `yaml:"event_type"               json:"type"`
	Description      string   `yaml:"description"              json:"description"`
	Participants     []string `yaml:"participants,omitempty"   json:"participants,omitempty"`
	Visibility       string   `yaml:"visibility"               json:"visibility,omitempty"`        // "public" (default) or "local"
	Location         string   `yaml:"location"                 json:"location,omitempty"`         // location name where event occurred
	Target           string   `yaml:"target,omitempty"         json:"target,omitempty"`           // optional character or location ID the event is about
	PrivateRecipient string   `yaml:"private_recipient,omitempty" json:"private_recipient,omitempty"` // if set, only this character sees the event
}

// WorldConcept describes the fundamental nature of the world and its hidden rules.
// All fields are optional; omitting the concept block in world.yaml leaves this at zero value.
type WorldConcept struct {
	Premise                string   `yaml:"premise"`                  // one sentence: the hidden truth characters must conceal
	Rules                  []string `yaml:"rules"`                    // constraints defining what is "normal" in this world
	Flavor                 string   `yaml:"flavor"`                   // tone/mood descriptor, e.g. "absurdist heist comedy"
	CharacterSpawnRule     string   `yaml:"character_spawn_rule"`     // if set, director may spawn new characters following this rule
	MaxSpawnedCharacters   int      `yaml:"max_spawned_characters"`   // max characters the director may spawn (0 = unlimited)
}

// WorldConfig is loaded from a scenario's world.yaml.
type WorldConfig struct {
	Locations     []Location   `yaml:"locations"`
	InitialEvents []Event      `yaml:"initial_events"`
	Weather       string       `yaml:"weather"`    // optional initial weather
	Atmosphere    string       `yaml:"atmosphere"` // optional initial atmosphere
	Concept       WorldConcept `yaml:"concept"`    // optional world premise and rules
}

// State holds the current world state for a simulation run.
type State struct {
	Tick              int
	TimeOfDay         string
	Weather           string // e.g. "clear", "rain", "fog" — empty means unset
	Atmosphere        string // e.g. "tense", "calm", "oppressive" — empty means unset
	Tension           int    // narrative tension level 0–10
	Concept           WorldConcept
	Locations         []Location
	EventLog          []Event
	SpawnedCharacters int    // count of characters created by the director at runtime
}

var timeOfDayLabels = []string{"morning", "afternoon", "evening", "night"}

// NewState creates a fresh world from a WorldConfig with tick 0 and time "morning".
// InitialEvents from cfg are pre-populated into the event log.
func NewState(cfg WorldConfig) *State {
	// Default missing Visibility to "public" so old YAML files behave correctly.
	events := make([]Event, len(cfg.InitialEvents))
	copy(events, cfg.InitialEvents)
	for i := range events {
		if events[i].Visibility == "" {
			events[i].Visibility = "public"
		}
	}
	return &State{
		Tick:       0,
		TimeOfDay:  "morning",
		Weather:    cfg.Weather,
		Atmosphere: cfg.Atmosphere,
		Concept:    cfg.Concept,
		Locations:  cfg.Locations,
		EventLog:   events,
	}
}

// AdvanceTick increments the tick and updates the time-of-day label.
// Cycle: morning (0-5), afternoon (6-11), evening (12-17), night (18-23), repeating.
func (s *State) AdvanceTick() {
	s.Tick++
	s.TimeOfDay = timeOfDayLabels[(s.Tick/6)%4]
}

// AppendEvent adds an event to the world log, defaulting Visibility to "public" if unset.
func (s *State) AppendEvent(e Event) {
	if e.Visibility == "" {
		e.Visibility = "public"
	}
	s.EventLog = append(s.EventLog, e)
}

// PublicSummary returns the universal world context available to all characters:
// time of day, weather, atmosphere, tension, all location names with their public
// descriptions, and the last 5 public events (visibility == "public" or "").
func (s *State) PublicSummary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("It is currently %s.\n", s.TimeOfDay))
	if s.Weather != "" {
		sb.WriteString(fmt.Sprintf("Weather: %s.\n", s.Weather))
	}
	if s.Atmosphere != "" {
		sb.WriteString(fmt.Sprintf("Atmosphere: %s.\n", s.Atmosphere))
	}
	if s.Tension > 0 {
		sb.WriteString(fmt.Sprintf("Tension: %d/10.\n", s.Tension))
	}
	sb.WriteString("\n")

	sb.WriteString("Known locations:\n")
	for _, loc := range s.Locations {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", loc.Name, loc.Description))
	}

	if s.Concept.Premise != "" {
		sb.WriteString("World Rules:\n")
		sb.WriteString(fmt.Sprintf("  %s\n", s.Concept.Premise))
		if s.Concept.Flavor != "" {
			sb.WriteString(fmt.Sprintf("  Tone: %s\n", s.Concept.Flavor))
		}
		for _, rule := range s.Concept.Rules {
			sb.WriteString(fmt.Sprintf("  - %s\n", rule))
		}
		sb.WriteString("\n")
	}

	var publicEvents []Event
	for _, e := range s.EventLog {
		if e.Visibility == "public" || e.Visibility == "" {
			publicEvents = append(publicEvents, e)
		}
	}
	if len(publicEvents) > 5 {
		publicEvents = publicEvents[len(publicEvents)-5:]
	}
	if len(publicEvents) > 0 {
		sb.WriteString("\nRecent public events:\n")
		for _, e := range publicEvents {
			sb.WriteString(fmt.Sprintf("- [tick %d] %s\n", e.Tick, e.Description))
		}
	}

	return sb.String()
}

// LocalContext returns private context for a specific location: the location's
// Details text, the last 5 local events that occurred there, and optionally the
// names of characters currently present (presentNames).
// Returns an empty string if locationID is empty or does not match any known location.
func (s *State) LocalContext(locationID string, presentNames []string) string {
	if locationID == "" {
		return ""
	}

	var loc *Location
	for i := range s.Locations {
		if s.Locations[i].Name == locationID {
			loc = &s.Locations[i]
			break
		}
	}
	if loc == nil {
		log.Printf("world: LocalContext called with unknown location %q — no local context returned", locationID)
		return ""
	}

	var sb strings.Builder

	if loc.Details != "" {
		sb.WriteString(fmt.Sprintf("Your location — %s:\n%s\n", loc.Name, loc.Details))
	}

	if len(presentNames) > 0 {
		sb.WriteString(fmt.Sprintf("\nCharacters present at %s: %s\n", loc.Name, strings.Join(presentNames, ", ")))
	}

	var localEvents []Event
	for _, e := range s.EventLog {
		if e.Visibility == "local" && e.Location == locationID {
			localEvents = append(localEvents, e)
		}
	}
	if len(localEvents) > 5 {
		localEvents = localEvents[len(localEvents)-5:]
	}
	if len(localEvents) > 0 {
		sb.WriteString("\nLocal events here:\n")
		for _, e := range localEvents {
			sb.WriteString(fmt.Sprintf("- [tick %d] %s\n", e.Tick, e.Description))
		}
	}

	return sb.String()
}
