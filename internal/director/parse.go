package director

import (
	"encoding/json"
	"strings"
)

// ParseToolCalls extracts tool calls from a raw LLM response.
// It scans for the first <tool_calls>...</tool_calls> block, parses the JSON
// array inside it, and returns the valid entries. Entries missing a "name"
// field are skipped. Returns an empty slice (not an error) when no block is
// found or the JSON is malformed.
func ParseToolCalls(raw string) ([]ToolCall, error) {
	const open = "<tool_calls>"
	const close = "</tool_calls>"

	start := strings.Index(raw, open)
	end := strings.Index(raw, close)
	if start == -1 || end == -1 || end <= start {
		return []ToolCall{}, nil
	}

	content := strings.TrimSpace(raw[start+len(open) : end])

	var calls []ToolCall
	if err := json.Unmarshal([]byte(content), &calls); err != nil {
		return []ToolCall{}, nil
	}

	// Filter out entries with no name.
	valid := calls[:0]
	for _, c := range calls {
		if strings.TrimSpace(c.Name) != "" {
			valid = append(valid, c)
		}
	}
	return valid, nil
}
