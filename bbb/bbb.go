package bbb

import (
	"git.sr.ht/~michl/quickbeam/browser"
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

func (m *Meeting) Join(microphone bool) {
	m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > div > span > button:nth-child(1)").MustClick()
}

func (m *Meeting) Yes() {
	m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > span > button:nth-child(1)").MustClick()
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
