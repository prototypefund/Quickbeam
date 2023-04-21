package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"git.sr.ht/~michl/quickbeam/bbb"
	"github.com/chzyer/readline"
)

func handleExit() {
	os.Exit(0)
}

var (
	meeting *bbb.Meeting
)

func handleCommand(command []string) error {
	switch command[0] {
	case "exit":
		handleExit()
	case "open":
		if len(command) != 2 {
			return errors.New("Invalid number of arguments: need meeting url")
		}
		url := command[1]
		meeting = bbb.NewMeeting(url)
	case "join":
		meeting.Join(true)
		meeting.Yes()
	case "mute":
		meeting.ToggleMute()
	case "raiseHand":
		meeting.ToggleRaiseHand()
	}
	return nil
}

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(rl.Stderr(), "Bye")
				os.Exit(0)
			} else {
				panic(err)
			}
		}
		line = strings.TrimSpace(line)
		command := strings.Fields(line)
		if len(command) > 0 {
			err = handleCommand(command)
			if err != nil {
				fmt.Fprintln(rl.Stderr(), err)
			}
		}
	}
	meeting.Leave()
}
