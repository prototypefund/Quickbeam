package web

type ErrNotFound struct {}
func (_ ErrNotFound) Error() string {
	return "Element was not found on webpage"
}

type Page interface {
	Start() error
	Close()
	Running() bool
	Navigate(url string) error
	Back()
	Forward()
	Root() Noder
	Execute(js string) (string, error)
}

type SubtreeChange interface {
	isSubtreeChange()
	//Node() Node
}

func (_ *NodeAdded) isSubtreeChange() {}
//func (c *NodeAdded) Node() Node {return c.node}
type NodeAdded struct {
	Data interface{}
}

func (_ *NodeRemoved) isSubtreeChange() {}
//func (c *NodeRemoved) Node() Node {return c.node}
type NodeRemoved struct {
	Data interface{}
}

func (_ *UnknownChange) isSubtreeChange() {}
//func (c *UnknownChange) Node() Node {return c.node}
type UnknownChange struct {
	Data interface{}
}

type Noder interface {
	SubNode(selector string, regexp string) (Noder, error)
	SubNodes(selector string) ([]Noder, error)
	MaybeSubNode(selector string, regexp string) (Noder, bool, error)
	SubscribeSubtree() (<-chan SubtreeChange, error)
	Text() (string, error)
	Click() error
}
