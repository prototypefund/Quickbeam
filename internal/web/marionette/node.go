package marionette

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
)

type Node struct {
	client       *marionette_client.Client
	element      *marionette_client.WebElement
	subscriber   *nodeSubscriptions
	subscription chan web.SubtreeChange
}

func (n *Node) Element() *marionette_client.WebElement {
	return n.element
}

func NewNode(client *marionette_client.Client, element *marionette_client.WebElement, subscriber *nodeSubscriptions) *Node {
	return &Node{client, element, subscriber, nil}
}
func (n *Node) SubNode(selector string, regexp string) (node web.Noder, err error) {
	nodes, err := findElements(n.element, copySpawner(n), selector, regexp, "")
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		err = protocol.WebpageError(
			fmt.Sprintf("Element with '%s' and regex '%s' not found",
				selector, regexp))
		return nil, err
	}
	return nodes[0], nil
}

func (n *Node) SubNodes(selector string) (nodes []web.Noder, err error) {
	nodes, err = findElements(n.element, copySpawner(n), selector, "", "")
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (n *Node) MaybeSubNode(selector string, regexp string) (web.Noder, bool, error) {
	nodes, err := findElements(n.element, copySpawner(n), selector, regexp, "")
	if err != nil {
		return nil, false, err
	}
	if len(nodes) == 0 {
		return nil, false, nil
	}
	return nodes[0], true, nil
}
func (n *Node) WaitSubNode(selector string, regexp string) (web.Noder, error) {
	found, nodes, err := waitForElements(copySpawner(n), n.element, selector, regexp, "", time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}
	if found && len(nodes) > 0 {
		return nodes[0], nil
	}
	return nil, web.ErrNotFound{}
}

func (n *Node) SubscribeSubtree() (c <-chan web.SubtreeChange, err error) {
	// if n.subscription != nil {
	//	return n.subscription, nil
	// }
	_, c = n.subscriber.new(n)
	return
}

func (n *Node) Text() (string, error) {
	return n.element.Text(), nil
}

func (n *Node) Click() error {
	n.element.Click()
	return nil
}

func (n *Node) SendKeys(sequence string) error {
	return n.element.SendKeys(sequence)
}

func (n *Node) Attribute(name string) (value string, found bool, err error) {
	value = n.element.Attribute(name)
	if value != "" {
		found = true
	}
	return value, found, nil
}

func (n Node) MarshalJSON() ([]byte, error) {
	e := n.element
	value := map[string]string{
		"element-6066-11e4-a52e-4f735466cecf": e.Id(),
	}
	return json.Marshal(value)
}

var _ web.Noder = &Node{}

// nodeSpawner is a helper struct for the creation of a new Node.
// It encapsulates the context needed by a Node and has methods that
// return new Nodes utilizing that context.
type nodeSpawner struct {
	client     *marionette_client.Client
	subscriber *nodeSubscriptions
}

func newSpawner(client *marionette_client.Client, subscriber *nodeSubscriptions) nodeSpawner {
	return nodeSpawner{
		client:     client,
		subscriber: subscriber,
	}
}

func copySpawner(n *Node) nodeSpawner {
	return nodeSpawner{
		client:     n.client,
		subscriber: n.subscriber,
	}
}

func (s nodeSpawner) fromElement(e *marionette_client.WebElement) *Node {
	return NewNode(s.client, e, s.subscriber)
}

// elementsFinder is a helper interface in order to facilate testing.
// Its only method has the same interface as marionette_client.WebElement.FindElements.
// In testing, a mock implementation of that method can be provided such that we do
// no need to start a marionette server like firefox.
type elementsFinder interface {
	FindElements(marionette_client.By, string) ([]*marionette_client.WebElement, error)
}

func findElements(parent elementsFinder, spawner nodeSpawner, selector string, re string, xpath string) (res []web.Noder, err error) {
	var elements []*marionette_client.WebElement
	if len(xpath) > 0 {
		elements, err = parent.FindElements(
			marionette_client.By(marionette_client.XPATH),
			xpath)
	} else {
		elements, err = parent.FindElements(
			marionette_client.By(marionette_client.CSS_SELECTOR),
			selector)
	}
	if err != nil {
		return res, err
	}
	res = make([]web.Noder, 0)
	for _, e := range elements {
		if re != "" {
			text := []byte(e.Text())
			matched, err := regexp.Match(re, text)
			if err != nil {
				return nil, err
			}
			if !matched {
				continue
			}
		}
		res = append(res, spawner.fromElement(e))
	}
	return res, nil
}

func waitForElements(spawner nodeSpawner, element elementsFinder, selector string, regexp string, xpath string, timeout time.Duration) (found bool, res []web.Noder, err error) {
	nodes := make(chan web.Noder)
	failure := make(chan error)

	go func() {
		for {
			ns, err := findElements(element, spawner, selector, regexp, xpath)
			if err != nil {
				failure <- err
				return
			}
			for _, n := range ns {
				nodes <- n
			}
		}
	}()
	select {
	case node := <-nodes:
		res = append(res, node)
	case err = <-failure:
		return false, []web.Noder{}, err
	case <-time.After(timeout):
		return false, []web.Noder{}, nil
	}
	return true, res, nil
}
