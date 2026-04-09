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

// CoverIdentity describes how a character presents themselves within this world.
// Nil means the character has no cover — they are who they appear to be.
type CoverIdentity struct {
	Alias     string   `yaml:"alias"`     // name used in this world
	Role      string   `yaml:"role"`      // claimed occupation or social position
	Backstory string   `yaml:"backstory"` // one sentence of invented personal history
	Weaknesses []string `yaml:"weaknesses"` // behaviours/topics that could expose their true nature
}

// CharacterJudgment holds one character's subjective opinion of another.
// It is formed from observable information only and colored by the judging
// character's own psychology. Never persisted to YAML.
type CharacterJudgment struct {
	About       string // ID of the character being judged
	Name        string // name snapshot at judgment time (cover alias if applicable)
	Impression  string // first-person internal narrative opinion (2–3 sentences)
	Trust       string // "high" | "medium" | "low" | "none"
	Interest    string // "high" | "medium" | "low"
	Threat      string // "high" | "medium" | "low" | "none"
	FormedTick  int    // tick at which judgment was formed (0 = pre-simulation)
	UpdatedTick int    // tick at which judgment was last updated (0 = never updated)
}

// Character represents an autonomous agent in the simulation.
type Character struct {
	ID         string `yaml:"id"`
	Type       string `yaml:"type"` // "" or "character" = regular; "game_director" = Game Director
	Name       string `yaml:"name"`
	Age        int    `yaml:"age"`
	Gender     string `yaml:"gender"`
	Occupation string `yaml:"occupation"`

	// Appearance describes how this character presents to others on first encounter.
	// It is visible to other characters when forming judgments but not included in
	// the character's own system prompt.
	Appearance string `yaml:"appearance"`

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

	// Cover identity — nil if character has no cover in this scenario
	CoverIdentity *CoverIdentity `yaml:"cover_identity"`

	// Runtime state
	Location       string   `yaml:"location"`       // current location name, matches world.Location.Name
	Goals          []string `yaml:"goals"`
	EmotionalState string   `yaml:"emotional_state"`

	Memory    []MemoryEntry              `yaml:"-"`
	MaxMemory int                        `yaml:"-"`
	Inbox     []world.Event              `yaml:"-"` // private events addressed to this character; flushed after prompt build
	Judgments map[string]CharacterJudgment `yaml:"-"` // keyed by character ID; subjective opinions of other characters
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
		cf.Characters[i].Judgments = make(map[string]CharacterJudgment)
	}
	return cf.Characters, nil
}
