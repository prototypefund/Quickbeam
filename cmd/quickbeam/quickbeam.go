package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~michl/quickbeam/internal/api"
	"git.sr.ht/~michl/quickbeam/internal/web/marionette"
	"github.com/sourcegraph/jsonrpc2"
)

var a api.Api = api.Api{
}

type myReadWriteCloser struct {
	in io.ReadCloser
	out io.WriteCloser
}

func (s *myReadWriteCloser) Read(p []byte) (n int, err error) {
	return s.in.Read(p)
}

func (s *myReadWriteCloser) Write(p []byte) (n int, err error) {
	return s.out.Write(p)
}

func (s *myReadWriteCloser) Close() (err error) {
	var errIn error
	errOut := s.out.Close()
	if errIn != nil {
		return errIn
	} else {
		return errOut
	}
}

func handlerFunc(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	defer func () {
		if r := recover(); r != nil {
			result = nil
			err = api.RuntimeError(r)
		}
	}()
	var params map[string]interface{}
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
	} else {
		params = nil
	}
	return a.Dispatch(req.Method, params)
}

func main() {
	firefox := marionette.NewFirefox()
	firefox.Headless = false
	firefox.Start()
	a.WebPage, _ = firefox.NewPage()

	cleanup := func() {
		firefox.Quit()
	}
	defer cleanup()
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			_ = <- sigs
			cleanup()
			os.Exit(0)
		}
	}()

	stdInOutCloser := &myReadWriteCloser{
		in: os.Stdin,
		out: os.Stdout,
	}
	stdInOutStream := jsonrpc2.NewBufferedStream(stdInOutCloser, jsonrpc2.VSCodeObjectCodec{})
	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(handlerFunc))
	conn := jsonrpc2.NewConn(context.TODO(), stdInOutStream, handler, func (_ *jsonrpc2.Conn) {})
	<-conn.DisconnectNotify()
}
