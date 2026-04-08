package character

import (
	"fmt"
	"strings"
)

// BuildSystemPrompt constructs a system-role prompt from a character's persona.
// Sections with empty or nil fields are silently omitted.
func BuildSystemPrompt(c Character) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s, a %d-year-old %s.\n\n", c.Name, c.Age, c.Occupation))

	if c.Motivation != "" {
		sb.WriteString(fmt.Sprintf("Motivación: %s\n\n", c.Motivation))
	}
	if c.Fear != "" {
		sb.WriteString(fmt.Sprintf("Miedo: %s\n\n", c.Fear))
	}
	if c.CoreBelief != "" {
		sb.WriteString(fmt.Sprintf("Creencia central: %s\n\n", c.CoreBelief))
	}
	if c.InternalTension != "" {
		sb.WriteString(fmt.Sprintf("Tensión interna: %s\n\n", c.InternalTension))
	}

	if len(c.FormativeEvents) > 0 {
		sb.WriteString("Eventos formativos:\n")
		for _, e := range c.FormativeEvents {
			sb.WriteString(fmt.Sprintf("- %s\n", e))
		}
		sb.WriteString("\n")
	}

	v := c.Voice
	if v.Formality != "" || v.VerbalTics != "" || v.ResponseLength != "" || v.HumorType != "" || v.CommunicationStyle != "" {
		sb.WriteString("Voz:\n")
		if v.Formality != "" {
			sb.WriteString(fmt.Sprintf("- Formalidad: %s\n", v.Formality))
		}
		if v.VerbalTics != "" {
			sb.WriteString(fmt.Sprintf("- Muletillas: %s\n", v.VerbalTics))
		}
		if v.ResponseLength != "" {
			sb.WriteString(fmt.Sprintf("- Extensión: %s\n", v.ResponseLength))
		}
		if v.HumorType != "" {
			sb.WriteString(fmt.Sprintf("- Humor: %s\n", v.HumorType))
		}
		if v.CommunicationStyle != "" {
			sb.WriteString(fmt.Sprintf("- Estilo: %s\n", v.CommunicationStyle))
		}
		sb.WriteString("\n")
	}

	r := c.RelationalDefaults
	if r.Strangers != "" || r.Authority != "" || r.Vulnerable != "" {
		sb.WriteString("Relaciones default:\n")
		if r.Strangers != "" {
			sb.WriteString(fmt.Sprintf("- Extraños: %s\n", r.Strangers))
		}
		if r.Authority != "" {
			sb.WriteString(fmt.Sprintf("- Autoridad: %s\n", r.Authority))
		}
		if r.Vulnerable != "" {
			sb.WriteString(fmt.Sprintf("- Vulnerables: %s\n", r.Vulnerable))
		}
		sb.WriteString("\n")
	}

	if len(c.Goals) > 0 {
		sb.WriteString("Objetivos:\n")
		for _, g := range c.Goals {
			sb.WriteString(fmt.Sprintf("- %s\n", g))
		}
		sb.WriteString("\n")
	}

	if c.EmotionalState != "" {
		sb.WriteString(fmt.Sprintf("Estado emocional actual: %s\n\n", c.EmotionalState))
	}

	if len(c.DialogueExamples) > 0 {
		sb.WriteString("Ejemplos de diálogo:\n")
		for _, d := range c.DialogueExamples {
			sb.WriteString(fmt.Sprintf("— %s\n", d))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Stay in character at all times. Respond as this person would. Keep responses concise.")
	return sb.String()
}

// FlushInbox returns a "Private information" section built from c.Inbox and clears it.
// Returns an empty string if the inbox is empty.
func FlushInbox(c *Character) string {
	if len(c.Inbox) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\nPrivate information you recently learned:\n")
	for _, ev := range c.Inbox {
		sb.WriteString(fmt.Sprintf("- %s\n", ev.Description))
	}
	c.Inbox = c.Inbox[:0]
	return sb.String()
}

// BuildZoneRoster returns a map from location name to the list of character names
// currently at that location. Characters with an empty Location are omitted.
func BuildZoneRoster(chars []*Character) map[string][]string {
	roster := make(map[string][]string)
	for _, c := range chars {
		if c.Location == "" {
			continue
		}
		roster[c.Location] = append(roster[c.Location], c.Name)
	}
	return roster
}

// BuildZoneContext renders a zone roster as a human-readable block suitable for
// appending to LLM system prompts. Returns an empty string for an empty roster.
func BuildZoneContext(roster map[string][]string) string {
	if len(roster) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\nWho is where:\n")
	for loc, names := range roster {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", loc, strings.Join(names, ", ")))
	}
	return sb.String()
}

// BuildMovementPrompt constructs a prompt asking a character where to move next.
// The expected response is an exact location name from the provided list, or "stay".
// zoneRoster (location → character names) is appended as a "Who is where" section
// when non-empty, so the character can reason about where others are.
func BuildMovementPrompt(c Character, locations []string, zoneRoster map[string][]string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s.\n", c.Name))
	if c.Motivation != "" {
		sb.WriteString(fmt.Sprintf("Your motivation: %s\n", c.Motivation))
	}
	if c.Fear != "" {
		sb.WriteString(fmt.Sprintf("Your fear: %s\n", c.Fear))
	}

	sb.WriteString(fmt.Sprintf("\nYou are currently at: %s\n", c.Location))
	sb.WriteString("Available locations you can move to:\n")
	for _, loc := range locations {
		sb.WriteString(fmt.Sprintf("- %s\n", loc))
	}

	if zoneCtx := BuildZoneContext(zoneRoster); zoneCtx != "" {
		sb.WriteString(zoneCtx)
	}

	sb.WriteString("\nBased on your goals and what just happened, decide where to go next.")
	sb.WriteString("\nRespond with ONLY the exact location name from the list above.")
	sb.WriteString("\nIf you want to stay where you are, respond with exactly: stay")
	return sb.String()
}
