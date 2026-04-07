package llm

import (
	"fmt"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
)

// BuildSystemPrompt constructs a system-role prompt from a character's persona.
func BuildSystemPrompt(c character.Character) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s, a %d-year-old %s.\n\n", c.Name, c.Age, c.Occupation))
	sb.WriteString(fmt.Sprintf("Personality: %s.\n\n", strings.Join(c.Personality, ", ")))
	sb.WriteString(fmt.Sprintf("Backstory: %s\n\n", strings.TrimSpace(c.Backstory)))

	if len(c.Goals) > 0 {
		sb.WriteString("Your goals:\n")
		for _, g := range c.Goals {
			sb.WriteString(fmt.Sprintf("- %s\n", g))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Current emotional state: %s.\n\n", c.EmotionalState))
	sb.WriteString("Stay in character at all times. Respond naturally as this person would in conversation. Keep responses concise (2-4 sentences).")

	return sb.String()
}
