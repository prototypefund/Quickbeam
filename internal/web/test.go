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

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		w.Write(indexHtml)
	}
	if path == "/index.js" {
		w.Write(indexJs)
	}
}

func Tester(webpage Page, t *testing.T) {
	handler := http.HandlerFunc(handler)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	webpage.Start()
	webpage.Navigate(ts.URL)
	root := webpage.Root()
	heading := root.SubNode("h1", "")
	text := heading.Text()
	want := "Hello, world!"
	if text != want {
		t.Errorf("SubNode/Text: want %q, got %q", want, text)
	}

	divs := root.SubNodes("div")
	if len(divs) != 4 {
		t.Errorf("SubNodes: expected 2 difs, got %d", len(divs))
	}
	div2 := divs[1]
	text = div2.Text()
	want = "bar"
	if text != want {
		t.Errorf("SubNodes/Text: want %q, got %q", want, text)
	}

	div1, ok := root.MaybeSubNode("div p", "f.o")
	if !ok {
		t.Errorf("MaybeSubNode: did not find div1")
	} else {
		want = "foo"
		text = div1.Text()
		if text != want {
			t.Errorf("MaybeSubNode/Text: want %q, got %q", want, text)
		}
	}

	ghost, ok := root.MaybeSubNode("div ol ul p", "")
	if ok {
		t.Errorf("MaybeSubNode: found something we did not want: %v", ghost)
	}

	button := root.SubNode("button", "")
	changingList := root.SubNode("#changingList", "")
	c := changingList.SubscribeSubtree()
	quit := make(chan bool)
	go func() {
		time.Sleep(time.Second * 1)
		quit<- true
	}()

	for i := 0; i < 5; i++ {
		button.Click()
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
		t.Errorf("SubscribeSubtree: Got changes, did not want to")
	}
	if _, ok = changes[0].(*NodeAdded); !ok {
		t.Errorf("First change not an addition: %T", changes[0])
	}
	if _, ok = changes[2].(*NodeRemoved); !ok {
		t.Errorf("Third change not a removal: %T", changes[2])
	}

}
