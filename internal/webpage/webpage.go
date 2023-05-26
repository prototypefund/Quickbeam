package webpage

type Webpage interface {
	Close()
	Navigate(url string) error
	Back()
	Forward()
	Start() error
	Running() bool
}
