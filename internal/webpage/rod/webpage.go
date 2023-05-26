package rod

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type RodWebpage struct {
	browser *rod.Browser
	page *rod.Page
}

func New() *RodWebpage {
	return &RodWebpage{}
}

func (b RodWebpage) Running() bool {
	return b.browser != nil
}

func (b RodWebpage) Start() error {
	path, _ := launcher.LookPath()
	control := launcher.New().Bin(path).Headless(false).Devtools(false).Set("audio").Delete("mute-audio").Delete("disable-audio-input").Delete("disable-audio-output").MustLaunch()
	b.browser = rod.New().ControlURL(control).MustConnect()
	b.page = nil
	return nil
}

func (b RodWebpage) Close() {
	if b.browser != nil {
		b.browser.MustClose()
	}
	b.browser = nil
	b.page = nil
}

func (b RodWebpage) Navigate(url string) error {
	if b.browser != nil {
		b.page = b.browser.MustPage(url)
	}
	return nil
}

func (b RodWebpage) Back() {
	if b.page != nil {
		b.page.NavigateBack()
	}
}

func (b RodWebpage) Forward() {
	if b.page != nil {
		b.page.NavigateForward()
	}
}
