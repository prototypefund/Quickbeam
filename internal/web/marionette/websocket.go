package marionette

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (ff *Firefox) startWebsocketServer() {
	errChannel := make(chan error)
	ff.websocketErrors = errChannel
	wsEndpoint := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			errChannel <- err
		}
		read(ws, &ff.nodeSubscriptions)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", wsEndpoint)
	server := http.Server{
		Handler: mux,
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		errChannel <- err
	}
	go func() {
		err := server.Serve(listener)
		if err != nil {
			errChannel <- err
		}

	}()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		errChannel <- err
	}
	ff.websocketPort = port
}

//go:embed init.js
var initJs []byte

func (ff *Firefox) initJavascript() error {
	if ff.websocketPort == "" {
		return protocol.InternalError(
			"Cannot install javascript backchannel: websocket server was not startet correctly")
	}
	args := []interface{}{
		ff.websocketPort,
	}
	_, err := ff.client.ExecuteScript(string(initJs), args, 10000, false)
	return err
}

func subscriptionHandler(subscriptions *nodeSubscriptions, msg []byte) {
	message := struct {
		Type      string `json:"type"`
		Id        string `json:"id"`
		Additions int    `json:"additions"`
		Removals  int    `json:"removals"`
	}{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("error handling websocket message `%s`: %s\n", msg, err)
		return
	}
	id, err := strconv.Atoi(message.Id)
	if err != nil {
		log.Printf("error parsing Id field of node subscription message: %s", err)
		return
	}
	c, err := subscriptions.get(id)
	if err != nil {
		log.Printf("could not find subscription id %d\n", id)
		return
	}
	c <- &web.ChildlistChange{
		Removals:  message.Removals,
		Additions: message.Additions,
	}
}

func read(conn *websocket.Conn, subscriptions *nodeSubscriptions) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if len(msg) > 0 {
			buf := bytes.NewBuffer(msg)
			message := struct {
				Type string `json:"type"`
			}{}
			err := json.NewDecoder(buf).Decode(&message)
			if err != nil {
				log.Println(err)
			}
			switch message.Type {
			case "chat":
				go ChatMessageHandler(msg)
			case "subscription":
				go subscriptionHandler(subscriptions, msg)
			default:
				log.Printf("Unknown message type: %s in %s\n", message.Type, msg)
			}
		}
	}
}
