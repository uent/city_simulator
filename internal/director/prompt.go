package director

import (
	"fmt"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// BuildDirectorPrompt constructs the system prompt for the Game Director.
// The prompt describes the current world state and provides a <tools> block
// listing all 13 available actions. The director is expected to respond with
// a <tool_calls> JSON array followed by optional free-form narration.
func BuildDirectorPrompt(state *world.State, chars []*character.Character, tick int) string {
	var sb strings.Builder

	sb.WriteString("You are the Game Director for a city simulation. Your role is to generate autonomous world events that make the simulation feel alive and reactive.\n\n")
	sb.WriteString(fmt.Sprintf("== CURRENT WORLD STATE (tick %d) ==\n", tick))
	sb.WriteString(fmt.Sprintf("Time: %s\n", state.TimeOfDay))
	if state.Weather != "" {
		sb.WriteString(fmt.Sprintf("Weather: %s\n", state.Weather))
	}
	if state.Atmosphere != "" {
		sb.WriteString(fmt.Sprintf("Atmosphere: %s\n", state.Atmosphere))
	}
	sb.WriteString(fmt.Sprintf("Tension: %d/10\n\n", state.Tension))

	sb.WriteString("Locations:\n")
	for _, loc := range state.Locations {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", loc.Name, loc.Description))
	}
	sb.WriteString("\n")

	if len(chars) > 0 {
		sb.WriteString("Characters:\n")
		for _, c := range chars {
			loc := c.Location
			if loc == "" {
				loc = "unknown"
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s @ %s\n", c.ID, c.Name, loc))
		}
		sb.WriteString("\n")
	}

	events := state.EventLog
	if len(events) > 10 {
		events = events[len(events)-10:]
	}
	if len(events) > 0 {
		sb.WriteString("Recent events:\n")
		for _, e := range events {
			if e.Visibility != "private" {
				sb.WriteString(fmt.Sprintf("- [tick %d][%s] %s\n", e.Tick, e.Type, e.Description))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(toolsBlock)
	sb.WriteString("\n\n")
	sb.WriteString("Respond with a <tool_calls> block containing a JSON array of action calls, followed by optional narration.\n")
	sb.WriteString("Call 0 to 5 actions. If nothing should happen, use an empty array.\n\n")
	sb.WriteString("Example response:\n")
	sb.WriteString("<tool_calls>\n")
	sb.WriteString(`[{"name":"set_weather","args":{"type":"heavy rain"}},{"name":"narrate","args":{"text":"Thunder rolls across the city."}}]`)
	sb.WriteString("\n</tool_calls>\n")
	sb.WriteString("The city braces for the storm.\n")

	return sb.String()
}

// toolsBlock is the static <tools> schema embedded in the director prompt.
const toolsBlock = `<tools>
Available actions — call them by name with the listed args:

ENVIRONMENT
  set_weather       args: type(string)          — change the weather (e.g. "rain","fog","clear")
  set_time_of_day   args: moment(string)         — shift time (e.g. "dawn","noon","midnight")
  set_atmosphere    args: descriptor(string)     — set mood (e.g. "tense","calm","oppressive")

NPC
  move_npc          args: id(string), destination(string), reason?(string)
  introduce_npc     args: id(string), name(string), role?(string), attitude?(string), motivation?(string), location?(string)
  add_npc_condition    args: id(string), condition(string)
  remove_npc_condition args: id(string), condition(string)

WORLD EVENTS
  modify_location   args: name(string), description?(string), details?(string)
  trigger_encounter args: participants([]string), context(string)
  trigger_event     args: type(string), description(string), location?(string), severity?(int 1-10)
  reveal_information args: recipient(string — character id), content(string)
  escalate_tension  args: delta(int, positive=raise, negative=lower)
  narrate           args: text(string)           — public broadcast, no state change
</tools>`
