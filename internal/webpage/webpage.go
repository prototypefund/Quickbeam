package webpage

type Webpage interface {
	Close()
	Navigate(url string)
	Back()
	Forward()
}
