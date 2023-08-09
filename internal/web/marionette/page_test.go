package marionette

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web"
)

func TestPage(t *testing.T) {
	if testing.Short() {
		t.Skip("do not launch firefox in short mode.")
	}
	t.Log("TestPage")
	firefox := NewFirefox()
	firefox.Headless = false
	firefox.Start()
	defer firefox.Quit()
	page, _ := firefox.NewPage()
	web.NavigateTester(page, t)
}
