package director

import (
	"fmt"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
)

// setWeatherAction sets state.Weather and appends a public event.
type setWeatherAction struct{}

func (setWeatherAction) Name() string { return "set_weather" }

func (setWeatherAction) Summary(args map[string]any) string {
	if typ, ok := stringArg(args, "type"); ok {
		return "set_weather: " + typ
	}
	return "set_weather"
}

func (setWeatherAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	typ, ok := stringArg(args, "type")
	if !ok {
		return fmt.Errorf("set_weather: missing required arg 'type'")
	}
	state.Weather = typ
	state.AppendEvent(world.Event{
		Type:        "weather",
		Description: fmt.Sprintf("The weather changes to %s.", typ),
		Visibility:  "public",
	})
	return nil
}

// setTimeOfDayAction sets state.TimeOfDay and appends a public event.
type setTimeOfDayAction struct{}

func (setTimeOfDayAction) Name() string { return "set_time_of_day" }

func (setTimeOfDayAction) Summary(args map[string]any) string {
	if moment, ok := stringArg(args, "moment"); ok {
		return "set_time_of_day: " + moment
	}
	return "set_time_of_day"
}

func (setTimeOfDayAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	moment, ok := stringArg(args, "moment")
	if !ok {
		return fmt.Errorf("set_time_of_day: missing required arg 'moment'")
	}
	state.TimeOfDay = moment
	state.AppendEvent(world.Event{
		Type:        "time",
		Description: fmt.Sprintf("The time shifts to %s.", moment),
		Visibility:  "public",
	})
	return nil
}

// setAtmosphereAction sets state.Atmosphere and appends a public event.
type setAtmosphereAction struct{}

func (setAtmosphereAction) Name() string { return "set_atmosphere" }

func (setAtmosphereAction) Summary(args map[string]any) string {
	if descriptor, ok := stringArg(args, "descriptor"); ok {
		return "set_atmosphere: " + descriptor
	}
	return "set_atmosphere"
}

func (setAtmosphereAction) Execute(args map[string]any, state *world.State, _ *[]*character.Character) error {
	descriptor, ok := stringArg(args, "descriptor")
	if !ok {
		return fmt.Errorf("set_atmosphere: missing required arg 'descriptor'")
	}
	state.Atmosphere = descriptor
	state.AppendEvent(world.Event{
		Type:        "atmosphere",
		Description: fmt.Sprintf("The atmosphere becomes %s.", descriptor),
		Visibility:  "public",
	})
	return nil
}

// stringArg extracts a string value from args by key. Returns ("", false) if missing or wrong type.
func stringArg(args map[string]any, key string) (string, bool) {
	v, ok := args[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
