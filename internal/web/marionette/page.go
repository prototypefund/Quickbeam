package marionette

import (
	"log"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
)

type Page struct {
	pageName string
	client *marionette_client.Client
}

func (p *Page) activate() (err error) {
	err = p.client.SwitchToWindow(p.pageName)
	if err != nil {
		log.Printf("marionette.Page.activate: %v\n", err)
	}
	return
}

func (p *Page) Close() {
	err := p.activate()
	if err == nil {
		p.client.CloseWindow()
	}
}

func (p *Page) Navigate(url string) (err error) {
	err = p.activate()
	if err != nil {
		return
	}
	_, err = p.client.Navigate(url)
	if err != nil {
		log.Printf("marionette.Page.Navigate: %v", err)
		return
	}
	return
}

func (p *Page) Back() {
	p.activate()
	_ = p.client.Back()
}

func (p *Page) Forward() {
	p.activate()
	_ = p.client.Forward()
}

func (p *Page) Root() web.Noder {
	root, _ := p.client.FindElement(marionette_client.By(marionette_client.CSS_SELECTOR), "body")
	return NewNode(root)
}

func (p *Page) Execute(js string) (string, error) {
	args := []interface{}{}
	r, err := p.client.ExecuteScript(js, args, 10000, false)
	if err != nil {
		return "", err
	}
	return r.Value, nil
}

func (p *Page) Exec(js string, args []interface{}) (string, error) {
	r, err := p.client.ExecuteScript(js, args, 10000, false)
	if err != nil {
		return "", err
	}
	return r.Value, nil
}
