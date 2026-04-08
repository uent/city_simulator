package director

import (
	"fmt"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// modifyLocationAction updates a location's Description and/or Details.
type modifyLocationAction struct{}

func (modifyLocationAction) Name() string { return "modify_location" }

func (modifyLocationAction) Summary(args map[string]any) string {
	if name, ok := stringArg(args, "name"); ok {
		return "modify_location: " + name
	}
	return "modify_location"
}

func (modifyLocationAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	name, ok := stringArg(args, "name")
	if !ok {
		return fmt.Errorf("modify_location: missing required arg 'name'")
	}
	for i := range state.Locations {
		if state.Locations[i].Name == name {
			if desc, ok := stringArg(args, "description"); ok {
				state.Locations[i].Description = desc
			}
			if details, ok := stringArg(args, "details"); ok {
				state.Locations[i].Details = details
			}
			return nil
		}
	}
	return fmt.Errorf("modify_location: location %q not found", name)
}

// triggerEncounterAction appends a public encounter event.
type triggerEncounterAction struct{}

func (triggerEncounterAction) Name() string { return "trigger_encounter" }

func (triggerEncounterAction) Summary(args map[string]any) string {
	if ctx, ok := stringArg(args, "context"); ok {
		return "trigger_encounter: " + ctx
	}
	return "trigger_encounter"
}

func (triggerEncounterAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	ctx, ok := stringArg(args, "context")
	if !ok {
		return fmt.Errorf("trigger_encounter: missing required arg 'context'")
	}
	var participants []string
	if raw, ok := args["participants"]; ok {
		switch v := raw.(type) {
		case []any:
			for _, item := range v {
				if s, ok := item.(string); ok {
					participants = append(participants, s)
				}
			}
		case []string:
			participants = v
		}
	}
	desc := ctx
	if len(participants) > 0 {
		desc = fmt.Sprintf("%s encounter: %s", strings.Join(participants, " and "), ctx)
	}
	state.AppendEvent(world.Event{
		Type:         "encounter",
		Description:  desc,
		Participants: participants,
		Visibility:   "public",
	})
	return nil
}

// triggerEventAction appends a public event; high severity raises tension.
type triggerEventAction struct{}

func (triggerEventAction) Name() string { return "trigger_event" }

func (triggerEventAction) Summary(args map[string]any) string {
	typ, typOk := stringArg(args, "type")
	desc, descOk := stringArg(args, "description")
	if typOk && descOk {
		return "trigger_event: " + typ + " — " + desc
	}
	if typOk {
		return "trigger_event: " + typ
	}
	return "trigger_event"
}

func (triggerEventAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	typ, ok := stringArg(args, "type")
	if !ok {
		return fmt.Errorf("trigger_event: missing required arg 'type'")
	}
	desc, ok := stringArg(args, "description")
	if !ok {
		return fmt.Errorf("trigger_event: missing required arg 'description'")
	}
	location, _ := stringArg(args, "location")

	severity := 0
	if raw, ok := args["severity"]; ok {
		switch v := raw.(type) {
		case float64:
			severity = int(v)
		case int:
			severity = v
		}
	}

	state.AppendEvent(world.Event{
		Type:        typ,
		Description: desc,
		Visibility:  "public",
		Location:    location,
	})

	if severity > 5 {
		state.Tension++
		if state.Tension > 10 {
			state.Tension = 10
		}
	}
	return nil
}

// revealInformationAction sends a private event to a specific character's inbox.
type revealInformationAction struct{}

func (revealInformationAction) Name() string { return "reveal_information" }

func (revealInformationAction) Summary(args map[string]any) string {
	recipient, recOk := stringArg(args, "recipient")
	content, contOk := stringArg(args, "content")
	if recOk && contOk {
		return "reveal_information: " + recipient + " — " + content
	}
	if recOk {
		return "reveal_information: " + recipient
	}
	return "reveal_information"
}

func (revealInformationAction) Execute(args map[string]any, state *world.State, chars *[]*character.Character) error {
	recipient, ok := stringArg(args, "recipient")
	if !ok {
		return fmt.Errorf("reveal_information: missing required arg 'recipient'")
	}
	content, ok := stringArg(args, "content")
	if !ok {
		return fmt.Errorf("reveal_information: missing required arg 'content'")
	}
	c := findChar(*chars, recipient)
	if c == nil {
		return fmt.Errorf("reveal_information: recipient %q not found", recipient)
	}
	ev := world.Event{
		Type:             "revelation",
		Description:      content,
		Visibility:       "private",
		PrivateRecipient: recipient,
	}
	state.AppendEvent(ev)
	c.Inbox = append(c.Inbox, ev)
	return nil
}

// escalateTensionAction adjusts state.Tension by delta, clamped to [0, 10].
type escalateTensionAction struct{}

func (escalateTensionAction) Name() string { return "escalate_tension" }

func (escalateTensionAction) Summary(args map[string]any) string {
	if raw, ok := args["delta"]; ok {
		switch v := raw.(type) {
		case float64:
			return fmt.Sprintf("escalate_tension: %+d", int(v))
		case int:
			return fmt.Sprintf("escalate_tension: %+d", v)
		}
	}
	return "escalate_tension"
}

func (escalateTensionAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	delta := 0
	if raw, ok := args["delta"]; ok {
		switch v := raw.(type) {
		case float64:
			delta = int(v)
		case int:
			delta = v
		}
	} else {
		return fmt.Errorf("escalate_tension: missing required arg 'delta'")
	}
	state.Tension += delta
	if state.Tension > 10 {
		state.Tension = 10
	}
	if state.Tension < 0 {
		state.Tension = 0
	}
	return nil
}

// narrateAction appends a public narration event without mutating other state.
type narrateAction struct{}

func (narrateAction) Name() string { return "narrate" }

func (narrateAction) Summary(args map[string]any) string {
	if text, ok := stringArg(args, "text"); ok {
		return "narrate: " + text
	}
	return "narrate"
}

func (narrateAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	text, ok := stringArg(args, "text")
	if !ok {
		return fmt.Errorf("narrate: missing required arg 'text'")
	}
	state.AppendEvent(world.Event{
		Type:        "narration",
		Description: text,
		Visibility:  "public",
	})
	return nil
}
