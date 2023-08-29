package web

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed test/index.html
var indexHtml []byte

//go:embed test/index.js
var indexJs []byte

var requestedPaths []string

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	requestedPaths = append(requestedPaths, path)
	if path == "/" {
		w.Write(indexHtml)
	}
	if path == "/index.js" {
		w.Write(indexJs)
	}
}

func setupTest() (url string, destruct func()) {
	h := http.HandlerFunc(handler)
	ts := httptest.NewServer(h)
	requestedPaths = make([]string, 0)
	return ts.URL, ts.Close
}

func subset(a []string, b []string) bool {
	var counts = make(map[string]int)
	for _, value := range a {
		counts[value] += 1
	}
	for _, value := range b {
		count, found := counts[value]
		if !found {
			return false
		}
		if count < 1 {
			return false
		}
		counts[value] = count - 1
	}
	return true
}

func NavigateTester(t *testing.T, webpage Page) {
	url, closer := setupTest()
	defer closer()
	webpage.Navigate(url)
	webpage.Close()
	expectedPaths := []string{"/", "/index.js"}
	if !subset(requestedPaths, expectedPaths) {
		t.Errorf("Simple invocation is expected to request %v, but %v were requested", expectedPaths, requestedPaths)
	}
}

func ExecuteTester(t *testing.T, page Page) {
	testCases := []struct{
		js string
		expected string
		error bool
	}{
		{`return "hello"`, "hello", false},
		{`const a = 1 + 2;`, "", false},
	}
	for i, test := range testCases {
		got, err := page.Execute(test.js)
		if (err != nil) && !test.error {
			t.Errorf("Expected no error but got: %s for case %d", err, i)
		}
		if (err == nil) && test.error {
			t.Errorf("Expected error, got nil for case %d", i)
		}
		assert.Equal(t, test.expected, got)
	}
}

func SubnodeTester(t *testing.T, webpage Page) {
	url, closer := setupTest()
	defer closer()
	webpage.Navigate(url)
	root, _ := webpage.Root()
	heading, _ := root.SubNode("h1", "")
	actualText, _ := heading.Text()
	expectedText := "Hello, world!"
	if actualText != expectedText {
		t.Errorf("SubNode/Text: want %q, got %q", expectedText, actualText)
	}
	divs, _ := root.SubNodes("div")
	if len(divs) != 4 {
		t.Errorf("SubNodes: expected 2 difs, got %d", len(divs))
	}
	div2 := divs[1]
	actualText, _ = div2.Text()
	expectedText = "bar"
	if actualText != expectedText {
		t.Errorf("SubNodes/Text: want %q, got %q", expectedText, actualText)
	}

	div1, ok, _ := root.MaybeSubNode("div p", "f.o")
	if !ok {
		t.Errorf("MaybeSubNode: did not find div1")
	} else {
		expectedText = "foo"
		actualText, _ = div1.Text()
		if actualText != expectedText {
			t.Errorf("MaybeSubNode/Text: want %q, got %q", expectedText, actualText)
		}
	}
}

func MaybeSubnodeTester(t *testing.T, webpage Page) {
	url, closer := setupTest()
	defer closer()
	webpage.Navigate(url)

	root, _ := webpage.Root()
	ghost, ok, _ := root.MaybeSubNode("div ol ul p", "")
	if ok {
		t.Errorf("MaybeSubNode: found something we did not want: %v", ghost)
	}
}

func InteractionTester(t *testing.T, webpage Page) {
	url, closer := setupTest()
	defer closer()
	webpage.Navigate(url)

	root, _ := webpage.Root()
	button, _ := root.SubNode("button", "")
	changingList, _ := root.SubNode("#changingList", "")
	c, _ := changingList.SubscribeSubtree()
	quit := make(chan bool)
	go func() {
		time.Sleep(time.Second * 1)
		//quit<- true
	}()

	for i := 0; i < 5; i++ {
		_ = button.Click()
	}
	go func() {
		quit <- true
	}()

	changes := make([]SubtreeChange, 0)
	// for i := 0; i < 5; i++ {
	//	change := <- c
	//	changes = append(changes, change)
	// }
collect:
	for {
		select {
		case change := <-c:
			changes = append(changes, change)
		case <-quit:
			break collect
		}
	}

	if len(changes) != 5 {
		t.Errorf("SubscribeSubtree: Wrong number of changes: %d, wanted 5", len(changes))
		return
	}
	if c, ok := changes[0].(*ChildlistChange); !ok {
		t.Errorf("First change not a ChildlistChange: %T", changes[0])
	} else {
		if c.Removals >= c.Additions {
			t.Errorf("First change is not an addition: %v", changes[0])
		}
	}
	if c, ok := changes[2].(*ChildlistChange); !ok {
		t.Errorf("Third change not a ChildlistChange: %T", changes[2])
	} else {
		if c.Additions > 0 {
			t.Errorf("Third change is not a removal: %v", changes[2])
		}
	}

}
