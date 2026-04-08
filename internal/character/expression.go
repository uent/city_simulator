package character

import "strings"

// Expression holds the two components of a character's LLM response:
// what they say aloud (Speech) and what they physically do (Action).
// Action is empty when the LLM produces no *...* markers.
type Expression struct {
	Speech string
	Action string
}

// ParseExpression extracts the first *...* block from raw as Action and
// treats the remaining text (before + after the block) as Speech.
// If no markers are found, the full trimmed text becomes Speech.
func ParseExpression(raw string) Expression {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Expression{}
	}

	start := strings.Index(raw, "*")
	if start == -1 {
		return Expression{Speech: raw}
	}

	end := strings.Index(raw[start+1:], "*")
	if end == -1 {
		// Unclosed marker — treat everything as speech.
		return Expression{Speech: raw}
	}
	end = start + 1 + end // absolute index of closing *

	action := strings.TrimSpace(raw[start+1 : end])
	before := strings.TrimSpace(raw[:start])
	after := strings.TrimSpace(raw[end+1:])

	var speech string
	switch {
	case before != "" && after != "":
		speech = before + " " + after
	case before != "":
		speech = before
	default:
		speech = after
	}

	return Expression{Action: action, Speech: speech}
}

// FormatExpression recombines an Expression into the *action* speech wire
// format used when injecting one character's turn into another's context.
func FormatExpression(e Expression) string {
	if e.Action == "" {
		return e.Speech
	}
	if e.Speech == "" {
		return "*" + e.Action + "*"
	}
	return "*" + e.Action + "* " + e.Speech
}
