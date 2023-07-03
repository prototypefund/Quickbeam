package web

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	for _, value := range(a) {
		counts[value] += 1
	}
	for _, value := range(b) {
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

func NavigateTester(webpage Page, t *testing.T) {
	url, closer := setupTest()
	defer closer()
	webpage.Start()
	webpage.Navigate(url)
	webpage.Close()
	expectedPaths := []string{"/", "/index.js"}
	if !subset(requestedPaths, expectedPaths) {
		t.Errorf("Simple invocation is expected to request %v, but %v were requested", expectedPaths, requestedPaths)
	}
}

func SubnodeTester(webpage Page, t *testing.T) {
	url, closer := setupTest()
	defer closer()
	webpage.Start()
	webpage.Navigate(url)
	root := webpage.Root()
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

func MaybeSubnodeTester(webpage Page, t *testing.T) {
	url, closer := setupTest()
	defer closer()
	webpage.Start()
	webpage.Navigate(url)

	root := webpage.Root()
	ghost, ok, _ := root.MaybeSubNode("div ol ul p", "")
	if ok {
		t.Errorf("MaybeSubNode: found something we did not want: %v", ghost)
	}
}

func InteractionTester(webpage Page, t *testing.T) {
	root := webpage.Root()

	button, _ := root.SubNode("button", "")
	changingList, _ := root.SubNode("#changingList", "")
	c, _ := changingList.SubscribeSubtree()
	quit := make(chan bool)
	go func() {
		time.Sleep(time.Second * 1)
		quit<- true
	}()

	for i := 0; i < 5; i++ {
		_ = button.Click()
	}

	changes := make([]SubtreeChange, 0)
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
	if _, ok := changes[0].(*NodeAdded); !ok {
		t.Errorf("First change not an addition: %T", changes[0])
	}
	if _, ok := changes[2].(*NodeRemoved); !ok {
		t.Errorf("Third change not a removal: %T", changes[2])
	}

}
