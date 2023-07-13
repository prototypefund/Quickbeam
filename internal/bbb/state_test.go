package bbb

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web/marionette"
)

func stateActive(states []string, state string) bool {
	for _, s := range(states) {
		if state == s {
			return true
		}
	}
	return false
}

func TestState(t *testing.T) {
	firefox := marionette.NewPage()
	firefox.Headless = false
	firefox.StartBrowser()
	defer firefox.KillBrowser()
	firefox.Start()
	defer firefox.Close()
	firefox.Navigate("https://bbb.cyber4edu.org/b/dan-stu-3a2-0bi")
	res, err := State(firefox)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !stateActive(res, "greenlight") {
		t.Errorf("Expected state greenlight, got: %v", res)
	}
}
