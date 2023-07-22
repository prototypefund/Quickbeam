package marionette

import (
	"regexp"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
)

type Node struct {
	element *marionette_client.WebElement
	subscription chan web.SubtreeChange
}

func (n *Node) Element() *marionette_client.WebElement {
	return n.element
}

func NewNode(element *marionette_client.WebElement) *Node {
	return &Node{element, nil}
}

type elementsFinder interface {
	FindElements(marionette_client.By, string) ([]*marionette_client.WebElement, error)
}

func findElements(parent elementsFinder, selector string, re string, xpath string) (res []web.Noder, err error) {
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
	for _, e := range(elements) {
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
		res = append(res, NewNode(e))
	}
	return res, nil
}

func (n *Node) SubNode(selector string, regexp string) (node web.Noder, err error) {
	nodes, err := findElements(n.element, selector, regexp, "")
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		err = web.ErrNotFound{}
		return nil, err
	}
	return nodes[0], nil
}

func (n *Node) SubNodes(selector string) (nodes []web.Noder, err error) {
	nodes, err = findElements(n.element, selector, "", "")
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (n *Node) MaybeSubNode(selector string, regexp string) (web.Noder, bool, error) {
	nodes, err := findElements(n.element, selector, regexp, "")
	if err != nil {
		return nil, false, err
	}
	if len(nodes) == 0 {
		return nil, false, nil
	}
	return nodes[0], true, nil
}

func waitForElements(element elementsFinder, selector string, regexp string, xpath string, timeout time.Duration) (found bool, res []web.Noder, err error) {
	nodes := make(chan web.Noder)
	failure := make(chan error)

	go func(){
		for {
			ns, err := findElements(element, selector, regexp, xpath)
			if err != nil {
				failure <- err
				return
			}
			for _, n := range(ns) {
				nodes <- n
			}
		}
	}()
	select {
	case node := <- nodes:
		res = append(res, node)
	case err = <- failure:
		return false, []web.Noder{}, err
	case <- time.After(timeout):
		return false, []web.Noder{}, nil
	}
	return true, res, nil
}

func (n *Node) WaitSubNode(selector string, regexp string) (web.Noder, error) {
	found, nodes, err := waitForElements(n.element, selector, regexp, "", time.Duration(10) * time.Second)
	if err != nil {
		return nil, err
	}
	if found && len(nodes) > 0 {
		return nodes[0], nil
	}
	return nil, web.ErrNotFound{}
}

func (n *Node) SubscribeSubtree() (c <-chan web.SubtreeChange, err error) {
	if n.subscription != nil {
		return n.subscription, nil
	}
	return make(chan web.SubtreeChange), nil
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


var _ web.Noder = &Node{}
