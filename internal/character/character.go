package character

import (
	"fmt"
	"os"

	"github.com/jnn-z/city_simulator/internal/world"
	"gopkg.in/yaml.v3"
)

// VoiceProfile describes how a character speaks and communicates.
type VoiceProfile struct {
	Formality          string `yaml:"formality"`
	VerbalTics         string `yaml:"verbal_tics"`
	ResponseLength     string `yaml:"response_length"`
	HumorType          string `yaml:"humor_type"`
	CommunicationStyle string `yaml:"communication_style"`
}

// RelationalProfile describes a character's default stance toward categories of people.
type RelationalProfile struct {
	Strangers  string `yaml:"strangers"`
	Authority  string `yaml:"authority"`
	Vulnerable string `yaml:"vulnerable"`
}

// Character represents an autonomous agent in the simulation.
type Character struct {
	ID         string `yaml:"id"`
	Type       string `yaml:"type"` // "" or "character" = regular; "game_director" = Game Director
	Name       string `yaml:"name"`
	Age        int    `yaml:"age"`
	Occupation string `yaml:"occupation"`

	// Psychological core — static anchors the LLM uses to resolve ambiguity
	Motivation      string `yaml:"motivation"`
	Fear            string `yaml:"fear"`
	CoreBelief      string `yaml:"core_belief"`
	InternalTension string `yaml:"internal_tension"`

	// History as causal events, not biography (2-3 "event → consequence" bullets)
	FormativeEvents []string `yaml:"formative_events"`

	// Voice and relational stance
	Voice              VoiceProfile      `yaml:"voice"`
	RelationalDefaults RelationalProfile `yaml:"relational_defaults"`

	// Concrete dialogue anchors (3-4 representative lines)
	DialogueExamples []string `yaml:"dialogue_examples"`

	// Runtime state
	Location       string   `yaml:"location"`       // current location name, matches world.Location.Name
	Goals          []string `yaml:"goals"`
	EmotionalState string   `yaml:"emotional_state"`

	Memory    []MemoryEntry `yaml:"-"`
	MaxMemory int           `yaml:"-"`
	Inbox     []world.Event `yaml:"-"` // private events addressed to this character; flushed after prompt build
}

type charactersFile struct {
	Characters []Character `yaml:"characters"`
}

// LoadCharacters reads a YAML file and returns all defined characters.
func LoadCharacters(path string) ([]Character, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read characters file %q: %w", path, err)
	}
	var cf charactersFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("cannot parse characters file %q: %w", path, err)
	}
	for i := range cf.Characters {
		if cf.Characters[i].EmotionalState == "" {
			cf.Characters[i].EmotionalState = "neutral"
		}
		if cf.Characters[i].MaxMemory == 0 {
			cf.Characters[i].MaxMemory = 20
		}
		cf.Characters[i].Inbox = []world.Event{}
	}
	return cf.Characters, nil
}
