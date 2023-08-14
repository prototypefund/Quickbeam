package marionette

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

//go:embed extension.xpi
var extensionXpi []byte

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

func read(conn *websocket.Conn) {
	for {
		_, p, _ := conn.ReadMessage()
		log.Println(p)
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
	mux.HandleFunc("/ws", wsEndpoint)
	go func() {
		err := http.ListenAndServe(":18981", mux)
		if err != nil {
			errChannel <- err
		}
	}()
}
