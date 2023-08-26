package marionette

import (
	"errors"
	"log"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
)

type Page struct {
	pageName string
	client   *marionette_client.Client
	firefox  *Firefox
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
	//_, err = p.Execute(string(chatJs))
	err = p.firefox.initJavascript()
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

func (p *Page) Root() (web.Noder, error) {
	spawner := newSpawner(p.client, &p.firefox.nodeSubscriptions)
	found, roots, err := waitForElements(spawner, p.client, "body", "", "", time.Second*time.Duration(10))
	if err != nil || !found {
		return nil, errors.New("Could not find root node")
	}
	return roots[0], nil
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
