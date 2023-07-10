package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/web/marionette"
)

//go:embed internal/web/test/index.html
var indexHtml []byte
//go:embed internal/web/test/index.js
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

func exec(page *marionette.Page, js string) {
	r, err := page.Execute(js)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func main() {
	h := http.HandlerFunc(handler)
	ts := httptest.NewServer(h)
	defer ts.Close()

	page := marionette.NewPage()
	page.Headless = false
	page.StartBrowser()
	defer page.KillBrowser()

	page.Start()
	page.Navigate(ts.URL)

	b, _ := page.Root().SubNode("button", "")
	button, ok := b.(*marionette.Node)
	if !ok {
		fmt.Println("Could no convert")
	}
	args := []interface{}{button.Element(), button.Element().Id()}
	ret, err := page.Exec(`
window.document.getElementById("foo").click()
console.log(arguments)
return arguments[0]
`, args)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)

	time.Sleep(time.Second * 2)

	js := `
window.mycallback = (mutations, observer) => {
  for (const mutation of mutations) {

    console.log(mutation)
  }
}
obs = new MutationObserver(window.mycallback)
obs.observe(document.querySelector("body"), { attributes: false, childList: true, subtree: true })

window.mybutton = document.querySelector("button")
window.mybutton.click()
`
	exec(page, js)

	exec(page, `
window.mybutton.click()
`)
	for {}
}
