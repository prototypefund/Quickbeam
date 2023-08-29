package marionette

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	if testing.Short() {
		t.Skip("do not launch firefox in short mode.")
	}
	firefox := NewFirefox()
	firefox.Headless = true
	firefox.Start()
	defer firefox.Quit()
	page, err := firefox.NewPage()
	assert.Nil(t, err)
	web.NavigateTester(t, page)
	web.ExecuteTester(t, page)
}

