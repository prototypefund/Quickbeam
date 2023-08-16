package bbb

import (
	"strings"

	"git.sr.ht/~michl/quickbeam/deprecated/browser"
	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

type Meeting struct {
	url string
	b *browser.Browser
}

func NewMeeting(url string) *Meeting {
	b := browser.New(url)
	return &Meeting{
		url: url,
		b: b,
	}
}

func (m *Meeting) userList() *rod.Element {
	page := m.b.Page
	userList := page.MustElement("[aria-label='Users list']")
	return userList
}

func (m *Meeting) GetAttendees() []string {
	a := []string{}
	userList := m.userList()
	users := userList.MustElements("[aria-label]")
	// first element is users list itself
	for _, u := range users[1:] {
		text := u.MustText()
		text = strings.ReplaceAll(text, "\n", " ")
		a = append(a, text)
	}
	return a
}

func (m *Meeting) Join(microphone bool) {
	// m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > div > span > button:nth-child(1)").MustClick()
	page := m.b.Page
	button := page.MustElement("[aria-label='Microphone']")
	button.MustWaitInteractable()
	button.MustClick()
}

func (m *Meeting) Yes() {
	selector := "[data-test='echoYesBtn']"
	page := m.b.Page
	page.MustWaitElementsMoreThan(selector, 0)
	modal := page.MustElement(".ReactModalPortal")
	button := modal.MustElement("[data-test='echoYesBtn']")
	button.MustWaitInteractable()
	button.MustClick();
}

type MeetingCallback func(m *Meeting)

func (m *Meeting) SubscribeAttendanceChange(cb MeetingCallback) {
	page := m.b.Page
	c := make(chan bool)
	observerCallback := func(v gson.JSON) (interface {}, error) {
		c <- true
		return nil, nil
	}
	page.MustExpose("userListCallback", observerCallback)
	page.MustEval(`() => {
	  userList = document.querySelector("[aria-label='Users list']")
	  obsCallback = (mutations, observer) => {
	    for (const mutation of mutations) {
	      console.log(mutation)
	      userListCallback({ a: "Hello", b: "bua", c: mutation.type})
	    }
	  }
	  obs = new MutationObserver(obsCallback)
	  obs.observe(userList, { attributes: false, childList: true, subtree: true })
	}`)
	go func() {
		for {
			<- c
			cb(m)
		}
	}()
}

func (m *Meeting) ToggleMute() {
	p := m.b.Page
	selMute := "[aria-label=\"Mute\"]"
	selUnmute := "[aria-label=\"Unmute\"]"
	if ok, b, _ := p.Has(selMute); ok {
		b.MustClick()
	} else if ok, b, _ := p.Has(selUnmute); ok {
		b.MustClick()
	}
}

func (m *Meeting) ToggleRaiseHand() {
	p := m.b.Page
	selRaise := "[aria-label=\"Raise hand\"]"
	selLower := "[aria-label=\"Lower hand\"]"
	if ok, b, _ := p.Has(selRaise); ok {
		b.MustClick()
	} else if ok, b, _ := p.Has(selLower); ok {
		b.MustClick()
	}
}

func (m *Meeting) Leave() {
	p := m.b.Page
	p.MustElement("header [aria-label=\"Options\"]").MustClick()
	p.MustElementR("[role=\"menuitem\"]", "Leave meeting").MustClick()
	m.b.Page.Close()
}

func (m *Meeting) Activate() {
	m.b.Page.MustActivate()
}
