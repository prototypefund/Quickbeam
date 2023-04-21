package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Browser struct {
	url string
	Page *rod.Page
}

func New(url string) *Browser {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(false).Devtools(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(url)

	perms := proto.BrowserGrantPermissions{
		Permissions: []proto.BrowserPermissionType{proto.BrowserPermissionTypeAudioCapture},
		Origin: "",
		BrowserContextID: "",
	}
	perms.Call(page)

	return &Browser{url, page,}
}
