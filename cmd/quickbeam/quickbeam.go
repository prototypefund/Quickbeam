package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~michl/quickbeam/internal/api"
	"git.sr.ht/~michl/quickbeam/internal/bbb"
	"git.sr.ht/~michl/quickbeam/internal/web/marionette"
	"github.com/sourcegraph/jsonrpc2"
)

var bbbChatMessages = api.Collection{
	Identifier:    "bbb/chat_message",
	GetAllMembers: bbb.ChatAllMessages,
	Subscribe:     bbb.ChatSubscribeMessages,
}

var a = api.New()

type myReadWriteCloser struct {
	in  io.ReadCloser
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
	defer func() {
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
	var headless bool = true
	for _, arg := range os.Args {
		switch arg {
		case "--version":
			fmt.Println("quickbeam v0.4")
			os.Exit(0)
		case "--headless":
			headless = true
		case "--no-headless":
			headless = false
		}
	}

	firefox := marionette.NewFirefox()
	firefox.Headless = headless
	cleanup := func() {
		firefox.Quit()
	}
	defer cleanup()
	err := firefox.Start()
	if err != nil {
		cleanup()
		log.Fatal(err)
	}
	a.WebPage, err = firefox.NewPage()
	if err != nil {
		cleanup()
		log.Fatal(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		firefox.Wait()
		os.Exit(1)
	}()
	go func() {
		for {
			<-sigs
			cleanup()
			os.Exit(0)
		}
	}()
	a.RegisterCollection(bbbChatMessages)
	a.RegisterState("bbb", bbb.State)
	a.RegisterAction("bbb/join", bbb.Join)
	a.RegisterAction("bbb/yes", bbb.Yes)
	a.RegisterAction("bbb/toggle_mute", bbb.ToggleMute)
	a.RegisterAction("bbb/toggle_raise_hand", bbb.ToggleRaisedHand)
	a.RegisterAction("bbb/leave", bbb.Leave)
	a.RegisterAction("bbb/attendees", bbb.GetAttendees)
	a.RegisterAction("bbb/wait_attendance_change", bbb.WaitAttendanceChange)
	a.RegisterAction("bbb/send_chat_message", bbb.SendChatMessage)
	a.RegisterAction("bbb/log_user_list", bbb.LogUserList)

	stdInOutCloser := &myReadWriteCloser{
		in:  os.Stdin,
		out: os.Stdout,
	}
	stdInOutStream := jsonrpc2.NewBufferedStream(stdInOutCloser, jsonrpc2.VSCodeObjectCodec{})
	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(handlerFunc))
	conn := jsonrpc2.NewConn(context.TODO(), stdInOutStream, handler, func(_ *jsonrpc2.Conn) {})
	a.CallBack = func(method string, params interface{}) {
		err := conn.Notify(context.TODO(), method, params)
		if err != nil {
			log.Println(err)
		}
	}
	<-conn.DisconnectNotify()
}
