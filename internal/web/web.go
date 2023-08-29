package web

type ErrNotFound struct{}

func (_ ErrNotFound) Error() string {
	return "Element was not found on webpage"
}

type Browser interface {
	Start() error
	Quit() error
	NewPage() (Page, error)
}

type Page interface {
	Close()
	Navigate(url string) error
	Back()
	Forward()
	Root() (Noder, error)
	// only keep one of them:
	Execute(js string) (string, error)
	Exec(js string, args []interface{}) (string, error)
}

type SubtreeChange interface {
	isSubtreeChange()
	//Node() Node
}

func (_ *ChildlistChange) isSubtreeChange() {}

// func (c *NodeAdded) Node() Node {return c.node}
type ChildlistChange struct {
	Additions int
	Removals  int
}

func (_ *NodeAdded) isSubtreeChange() {}

// func (c *NodeAdded) Node() Node {return c.node}
type NodeAdded struct {
	Data interface{}
}

func (_ *NodeRemoved) isSubtreeChange() {}

// func (c *NodeRemoved) Node() Node {return c.node}
type NodeRemoved struct {
	Data interface{}
}

func (_ *UnknownChange) isSubtreeChange() {}

// func (c *UnknownChange) Node() Node {return c.node}
type UnknownChange struct {
	Data interface{}
}

type Noder interface {
	SubNode(selector string, regexp string) (Noder, error)
	SubNodes(selector string) ([]Noder, error)
	MaybeSubNode(selector string, regexp string) (Noder, bool, error)
	WaitSubNode(selector string, regexp string) (Noder, error)
	SubscribeSubtree() (<-chan SubtreeChange, error)
	Text() (string, error)
	Click() error
	SendKeys(sequence string) error
	Attribute(name string) (string, bool, error)
	// LogToConsole logs the node to the browser console.
	// It is meant for debugging DOM elements and selectors.
	// It ignores potential errors.
	LogToConsole(prefix string)
}
