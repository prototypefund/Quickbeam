package api

import (
	"fmt"

	"git.sr.ht/~michl/quickbeam/internal/bbb"
	"git.sr.ht/~michl/quickbeam/internal/web"
)

var (
	actions map[string]interface{} = map[string]interface{}{
		"greet": greet,
		"bbb/join": bbb.Join,
		"bbb/yes": bbb.Yes,
		"bbb/toggle_mute": bbb.ToggleMute,
		"bbb/toggle_raised_hand": bbb.ToggleRaisedHand,
		"bbb/leave": bbb.Leave,
		"bbb/attendees": bbb.GetAttendees,
		"bbb/wait_attendance_change": bbb.WaitAttendanceChange,
	}
)

type greetParams struct {
	Name string `json:"name"`
}

type greetReturn struct {
	Greeting string
}

func greet(p greetParams, w web.Page, a *Api) (greetReturn, error) {
	name := p.Name
	greeting := fmt.Sprintf("Hello, %s!", name)
	return greetReturn{greeting}, nil
}

func (a *Api) dispatchAction(action string, params map[string]interface{}) (result interface{}, err error) {
	act, ok := actions[action]
	if !ok {
		return nil, ActionNotAvailableError{action}
	}
	return a.dispatchFunc(act, params)
}
