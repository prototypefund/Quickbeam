package marionette

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web"
)

func TestNode(t *testing.T) {
	if testing.Short() {
		t.Skip("do not launch firefox in short mode.")
	}
	t.Log("TestNode")
	firefox := NewFirefox()
	firefox.Headless = true
	firefox.Start()
	defer firefox.Quit()
	page, _ := firefox.NewPage()
	web.SubnodeTester(page, t)
	web.MaybeSubnodeTester(page, t)
	web.InteractionTester(page, t)
}
