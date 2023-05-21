package main

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"git.sr.ht/~michl/quickbeam/internal/api"
	"github.com/sourcegraph/jsonrpc2"
)

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
	var params map[string]interface{}
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
	} else {
		params = nil
	}
	return api.Dispatch(req.Method, params)
}

func main() {
	stdInOutCloser := &myReadWriteCloser{
		in: os.Stdin,
		out: os.Stdout,
	}
	stdInOutStream := jsonrpc2.NewPlainObjectStream(stdInOutCloser)
	handler := jsonrpc2.HandlerWithError(handlerFunc)
	conn := jsonrpc2.NewConn(context.TODO(), stdInOutStream, handler, func (_ *jsonrpc2.Conn) {})
	<-conn.DisconnectNotify()
}
