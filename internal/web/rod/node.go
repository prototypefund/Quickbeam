package rod

import (
	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/go-rod/rod"
)

type RodNode struct {
	element *rod.Element
}

func (n *RodNode) SubNode(selector string) web.WebNode {
	e := n.element.MustElement(selector)
	return &RodNode{element: e}
}

func (n *RodNode) Click() {
	n.element.MustClick()
}
