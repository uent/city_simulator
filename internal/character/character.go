package character

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Character represents an autonomous agent in the simulation.
type Character struct {
	ID             string   `yaml:"id"`
	Name           string   `yaml:"name"`
	Age            int      `yaml:"age"`
	Occupation     string   `yaml:"occupation"`
	Personality    []string `yaml:"personality"`
	Backstory      string   `yaml:"backstory"`
	Goals          []string `yaml:"goals"`
	EmotionalState string   `yaml:"emotional_state"`

	Memory    []MemoryEntry `yaml:"-"`
	MaxMemory int           `yaml:"-"`
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
	}
	return cf.Characters, nil
}
