package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Browser struct {
	url  string
	Page *rod.Page
}

func New(url string) *Browser {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(false).Devtools(false).Set("audio").Delete("mute-audio").Delete("disable-audio-input").Delete("disable-audio-output").MustLaunch()
	//u := "ws://127.0.0.1:9222/devtools/browser/1b5a56d2-ca33-4e59-aaac-4a786bc8674c"
	page := rod.New().ControlURL(u).MustConnect().MustPage(url)

	perms := proto.BrowserGrantPermissions{
		Permissions: []proto.BrowserPermissionType{
			proto.BrowserPermissionTypeAudioCapture,
		},
		Origin:           "",
		BrowserContextID: "",
	}
	perms.Call(page)

	return &Browser{url, page}
}
