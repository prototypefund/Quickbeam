package bbb

import (
	"strings"

	"git.sr.ht/~michl/quickbeam/internal/web"
)


type JoinResult struct {success string}
func Join(_ EmptyArgs, w web.Page) (res JoinResult, err error) {
	// m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > div > span > button:nth-child(1)").MustClick()
	button := w.Root().SubNode("[aria-label='Microphone']", "")
	button.Click()
	return JoinResult{"success"}, nil
}

func Yes(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	selector := "[data-test='echoYesBtn']"
	button := w.Root().SubNode(selector, "")
	button.Click()
	return EmptyResult{}, nil
}

type EmptyArgs struct{}
type EmptyResult struct{}
func ToggleMute(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	selMute := "[aria-label=\"Mute\"]"
	selUnmute := "[aria-label=\"Unmute\"]"
	root := w.Root()
	if b, ok := root.MaybeSubNode(selMute, ""); ok {
		b.Click()
	} else if b, ok := root.MaybeSubNode(selUnmute, ""); ok {
		b.Click()
	}
	return EmptyResult{}, nil
}

func Leave(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	w.Root().SubNode("header [aria-label=\"Options\"]", "").Click()
	w.Root().SubNode("[role=\"menuitem\"]", "Leave meeting").Click()
	return EmptyResult{}, nil
}

type Attendee struct {
	Name string `json:"name"`
	Muted bool `json:"muted"`
	Status string `json:"status"`
}

type AttendeeResult struct{
	Attendees []Attendee `json:"attendees"`
}
func GetAttendess(_ EmptyArgs, w web.Page) (res AttendeeResult, err error) {
	attendees := []Attendee{}
	userList := w.Root().SubNode("[aria-label='Users list']", "")
	users := userList.SubNodes("[aria-label]")
	for _, u := range users[1:] {
		name := u.Text()
		name = strings.ReplaceAll(name, "\n", " ")
		attendees = append(attendees, Attendee{Name: name,})
	}
	return AttendeeResult{attendees}, nil
}
