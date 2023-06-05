package web

type Page interface {
	Close()
	Navigate(url string) error
	Back()
	Forward()
	Start() error
	Running() bool
	Root() Node
}

type SubtreeChange interface {
	isSubtreeChange()
	Node() Node
}

func (_ *NodeAdded) isSubtreeChange() {}
func (c *NodeAdded) Node() Node {return c.node}
type NodeAdded struct {
	node Node
}

func (_ *NodeRemoved) isSubtreeChange() {}
func (c *NodeRemoved) Node() Node {return c.node}
type NodeRemoved struct {
	node Node
}

type Node interface {
	SubNode(selector string, regexp string) Node
	SubNodes(selector string) []Node
	MaybeSubNode(selector string, regexp string) (Node, bool)
	SubscribeSubtree() <-chan SubtreeChange
	Text() string
	Click()
}
