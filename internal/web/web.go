package web

type WebPage interface {
	Close()
	Navigate(url string) error
	Back()
	Forward()
	Start() error
	Running() bool
	Root() WebNode
}

type WebNode interface {
	SubNode(selector string) WebNode
	Click()
}
