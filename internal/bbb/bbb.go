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

func userList(w web.Page) web.Node {
	return w.Root().SubNode("[aria-label='Users list']", "")
}

type Attendee struct {
	Name string `json:"name"`
	Muted bool `json:"muted"`
	Status string `json:"status"`
}

type AttendeeResult struct{
	Attendees []Attendee `json:"attendees"`
}
func GetAttendees(_ EmptyArgs, w web.Page) (res AttendeeResult, err error) {
	attendees := []Attendee{}
	userList := userList(w)
	users := userList.SubNodes("[aria-label]")
	for _, u := range users[1:] {
		name := u.Text()
		name = strings.ReplaceAll(name, "\n", " ")
		attendees = append(attendees, Attendee{Name: name,})
	}
	return AttendeeResult{attendees}, nil
}

type ChangeResult struct {
	Change string `json:"change"`
}
func WaitAttendanceChange(_ EmptyArgs, w web.Page) (res ChangeResult, err error) {
	userList := userList(w)
	c := userList.SubscribeSubtree()
	for {
		change := <-c
		switch change.(type) {
		case *web.NodeAdded:
			return ChangeResult{Change: "added"}, nil
		case *web.NodeRemoved:
			return ChangeResult{Change: "removed"}, nil
			//case *web.UnknownChange:
			//log.Warn(fmt.Sprintf("Unknown subscription change: %v", c))
		}
	}
	return ChangeResult{}, nil
}
