package marionette

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed extension.xpi
var extensionXpi []byte
//go:embed extension/chat.js
var chatJs []byte

func (ff *Firefox) LoadExtension() error {
	extension, err := os.CreateTemp("", "quickbeam-*.xpi")
	if err != nil {
		return err
	}
	defer os.Remove(extension.Name())

	if _, err = extension.Write(extensionXpi); err != nil {
		return err
	}
	if err = extension.Close(); err != nil {
		return err
	}

	arguments := map[string]interface{}{
		"path": extension.Name(),
		"temporary": true,
	}
	if _, err := ff.transport.Send("Addon:Install", arguments); err != nil {
		return fmt.Errorf("Error installing quickbeam extension: %w", err)
	}
	return nil
}

var upgrade = websocket.Upgrader{
	ReadBufferSize: 8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// MessageHandler is an example for a callback that handles the message
// received from the websocket connection.
// The message is expected to be a JSON object, and can be deserialized
// for further processing.
func ChatMessageHandler(msg []byte) {
	buf := bytes.NewBuffer(msg)

	message := struct {
		Scope string `json:"scope"`
		User  string `json:"user"`
		Text  string `json:"message"`
		Timestamp int64 `json:"timestamp"`
	}{}

	err := json.NewDecoder(buf).Decode(&message)
	if err != nil {
		log.Println(err)
	}
	log.Println(time.Unix(0, message.Timestamp * 1000000).Format("2. Jan 2006 15:04"))
	log.Println(message)
}

func read(conn *websocket.Conn) {
	for {
		_, msg, _ := conn.ReadMessage()
		if len(msg) > 0 {
			buf := bytes.NewBuffer(msg)
			message := struct {
				Type string `json:"type"`
			}{}
			err := json.NewDecoder(buf).Decode(&message)
			if err != nil {
				log.Println(err)
			}
			if message.Type == "chat" {
				ChatMessageHandler(msg)
			} else {
				log.Println(string(msg))
			}
		}
	}
}

func (ff *Firefox) startWebsocketServer() {
	errChannel := make(chan error)
	ff.websocketErrors = errChannel
	wsEndpoint := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			errChannel <- err
		}
		read(ws)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", wsEndpoint)
	go func() {
		err := http.ListenAndServe(":18981", mux)
		if err != nil {
			errChannel <- err
		}
	}()
}

func (ff *Firefox) injectJavascript() error {
	resp, err := ff.client.ExecuteScript(string(chatJs), []interface{}{}, 1000, false)
	log.Println(resp)
	return err
}
