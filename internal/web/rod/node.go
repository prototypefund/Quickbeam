package rod

import (
	"fmt"
	"os"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

var nextId uint64 = 1

type RodNode struct {
	element *rod.Element
	id uint64
}

func (n *RodNode) getId() uint64 {
	if n.id == 0 {
		n.id = nextId
		nextId++
	}
	return n.id
}

func (n *RodNode) SubNode(selector string, regexp string) web.Node {
	e := n.element.MustElement(selector)
	return &RodNode{element: e}
}

func (n *RodNode) MaybeSubNode(selector string, regexp string) (web.Node, bool) {
	if ok, e, _ := n.element.Has(selector); ok {
		return &RodNode{element: e}, true
	}
	return nil, false
}

func (n *RodNode) SubNodes(selector string) []web.Node {
	elements := n.element.MustElements(selector)
	nodes := make([]web.Node, 0, len(elements))
	for _, e := range elements {
		n := &RodNode{element: e,}
		nodes = append(nodes, n)
	}
	return nodes
}

func (n *RodNode) SubscribeSubtree() (<-chan web.SubtreeChange) {
	c := make(chan web.SubtreeChange)
	observerCallback := func(v gson.JSON) (interface {}, error) {
		mutations := v.Get("mutations").Arr()
		fmt.Fprintf(os.Stderr, "%s", v.JSON("", ""))
		if len(mutations) == 0 {
			fmt.Fprintln(os.Stderr, "addedNodes not found")
			c <- &web.UnknownChange{Data: mutations,}
		}
		for _, m := range mutations {
			mutation := m.Map()
			t, ok := mutation["type"]
			if !ok {
				c <- &web.UnknownChange{Data: mutation,}
				continue
			}
			mutType := t.String()
			if mutType == "childList" {
				an, ok := mutation["addedNodes"]
				if !ok {
					c <- &web.UnknownChange{Data: mutation,}
					continue
				}
				addedNodes := an.Arr()
				for _, node := range addedNodes {
					c <- &web.NodeAdded{Data: node,}
				}
				rm, ok := mutation["removedNodes"]
				if !ok {
					c <- &web.UnknownChange{Data: mutation,}
				}
				removedNodes := rm.Arr()
				for _, node := range removedNodes {
					c <- &web.NodeRemoved{Data: node,}
				}
			} else {
				c <- &web.UnknownChange{Data: mutation}
			}
		}
		return nil, nil
	}
	callbackName := fmt.Sprintf("subtreeCallback%x", n.getId())
	page := n.element.Page()
	page.MustExpose(callbackName, observerCallback)
	n.element.MustEval(`() => {
	  obs = new MutationObserver(` +
		callbackName +`)
	  obs.observe(this, { attributes: false, childList: true, subtree: true })
	}`)
	return c
}

func (n *RodNode) Text() string {
	text, err := n.element.Text()
	if err != nil {
		text = ""
	}
	return text
}

func (n *RodNode) Click() {
	n.element.MustClick()
}
