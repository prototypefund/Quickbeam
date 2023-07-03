package marionette

import (
	"testing"
	"git.sr.ht/~michl/quickbeam/internal/web"
)

func TestNode(t *testing.T) {
	page := NewPage()
	page.Headless = false
	page.StartBrowser()
	defer page.KillBrowser()
	web.SubnodeTester(page, t)
	web.MaybeSubnodeTester(page, t)
	web.InteractionTester(page, t)
}
