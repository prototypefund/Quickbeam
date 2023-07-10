package bbb

import (
	"strings"

	"git.sr.ht/~michl/quickbeam/internal/web"
)


type JoinResult struct {success string}
func Join(_ EmptyArgs, w web.Page) (res JoinResult, err error) {
	// m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > div > span > button:nth-child(1)").MustClick()
	button, err := w.Root().WaitSubNode("[aria-label='Microphone']", "")
	if err != nil {
		return res, err
	}
	button.Click()
	return JoinResult{"success"}, nil
}

func Yes(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	selector := "[data-test='echoYesBtn']"
	button, err := w.Root().WaitSubNode(selector, "")
	if err != nil {
		return res, err
	}
	button.Click()
	return EmptyResult{}, nil
}

type EmptyArgs struct{}
type EmptyResult struct{}
func ToggleMute(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	selMute := "[aria-label=\"Mute\"]"
	selUnmute := "[aria-label=\"Unmute\"]"
	root := w.Root()
	b, ok, err := root.MaybeSubNode(selMute, "")
	if err != nil {
		return res, err
	}
	if ok {
		b.Click()
	} else {
		b, ok, err := root.MaybeSubNode(selUnmute, "")
		if err != nil {
			return res, err
		}
		if ok {
			b.Click()
		}
	}
	return EmptyResult{}, nil
}

func Leave(_ EmptyArgs, w web.Page) (res EmptyResult, err error) {
	hamburger, err := w.Root().SubNode("header [aria-label=\"Options\"]", "")
	if err != nil {
		return
	}
	err = hamburger.Click()
	if err != nil {
		return
	}
	leave, err := w.Root().SubNode("[role=\"menuitem\"]", "Leave meeting")
	if err != nil {
		return
	}
	err = leave.Click()
	if err != nil {
		return
	}
	return
}

func userList(w web.Page) (web.Noder, error) {
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
	userList, err := userList(w)
	if err != nil {
		return
	}
	users, err := userList.SubNodes("[aria-label]")
	if err != nil {
		return
	}
	for _, u := range users[1:] {
		name, _ := u.Text()
		name = strings.ReplaceAll(name, "\n", " ")
		attendees = append(attendees, Attendee{Name: name,})
	}
	return AttendeeResult{attendees}, nil
}

type ChangeResult struct {
	Change string `json:"change"`
}
func WaitAttendanceChange(_ EmptyArgs, w web.Page) (res ChangeResult, err error) {
	userList, err := userList(w)
	if err != nil {
		return
	}
	c, err := userList.SubscribeSubtree()
	if err != nil {
		return
	}
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
