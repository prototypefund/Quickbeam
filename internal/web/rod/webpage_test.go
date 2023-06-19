package rod

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/web"
)

func TestRod(t *testing.T) {
	page := New()
	web.Tester(page, t)
}
