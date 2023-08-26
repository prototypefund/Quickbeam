package marionette

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
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
		"path":      extension.Name(),
		"temporary": true,
	}
	if _, err := ff.transport.Send("Addon:Install", arguments); err != nil {
		return fmt.Errorf("Error installing quickbeam extension: %w", err)
	}
	return nil
}

// MessageHandler is an example for a callback that handles the message
// received from the websocket connection.
// The message is expected to be a JSON object, and can be deserialized
// for further processing.
func ChatMessageHandler(msg []byte) {
	buf := bytes.NewBuffer(msg)

	message := struct {
		Scope     string `json:"scope"`
		User      string `json:"user"`
		Text      string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	}{}

	err := json.NewDecoder(buf).Decode(&message)
	if err != nil {
		log.Println(err)
	}
	log.Println(time.Unix(0, message.Timestamp*1000000).Format("2. Jan 2006 15:04"))
	log.Println(message)
}
