package marionette

import (
	"sync"

	"github.com/njasm/marionette_client"
)

// Transport is a go routine safe marionette_client.Transporter
type Transport struct {
	// mu guards all calls to Transport.Send and Transport.Receive
	mu sync.Mutex
	// All actual work is delegated to the original implementation
	mt marionette_client.MarionetteTransport
}

func newTransport() *Transport {
	return &Transport{}
}

func (t *Transport) MessageID() int {
	return t.mt.MessageID()
}

func (t *Transport) Connect(host string, port int) error {
	return t.mt.Connect(host, port)
}

func (t *Transport) Close() error {
	return t.mt.Close()
}

func (t *Transport) Send(command string, values interface{}) (*marionette_client.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.mt.Send(command, values)
}

func (t *Transport) Receive() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.mt.Receive()
}

var _ marionette_client.Transporter = &Transport{}
