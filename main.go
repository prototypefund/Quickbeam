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
	case "activate":
		meeting.Activate()
	case "leave":
		meeting.Leave()
		os.Exit(0)
	case "attendees":
		for _, a := range meeting.GetAttendees() {
			fmt.Println(a)
		}
		fmt.Println()
	case "subscribe":
		meeting.SubscribeAttendanceChange(attendanceChange)
	}
	return nil
}

func attendanceChange(_ *bbb.Meeting) {
	fmt.Println("Change in attendance detected.")
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "> ",
		HistoryFile: "./.quickbeam_history",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				meeting.Leave()
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
}
