package bbb

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web/marionette"
)

func stateActive(states []string, state string) bool {
	for _, s := range states {
		if state == s {
			return true
		}
	}
	return false
}

func TestState(t *testing.T) {
	if testing.Short() {
		t.Skip("do not launch firefox in short mode.")
	}
	firefox := marionette.NewFirefox()
	firefox.Headless = true
	firefox.Start()
	defer firefox.Quit()
	page, _ := firefox.NewPage()
	page.Navigate("https://bbb.cyber4edu.org/b/dan-stu-3a2-0bi")
	res, err := State(page)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !stateActive(res, "greenlight") {
		//t.Errorf("Expected state greenlight, got: %v", res)
	}
}
