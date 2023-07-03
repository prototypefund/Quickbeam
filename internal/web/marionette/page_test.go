package marionette

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web"
)

func TestPage(t *testing.T) {
	page := NewPage()
	page.Headless = false
	page.StartBrowser()
	defer page.KillBrowser()
	web.NavigateTester(page, t)
}
