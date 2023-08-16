package bbb

import (
	"encoding/json"
	"log"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
)

type ChatScope int

const (
	ChatScopePrivate ChatScope = iota
	ChatScopePublic
)

type ChatMessage struct {
	Scope ChatScope `json:"scope"`
	Time time.Time `json:"time"`
	User string `json:"user"`
	Text string `json:"text"`
}

func (msg *ChatMessage) MarshalJSON() ([]byte, error) {
	jsonMsg := jsonChatMessage{}
	jsonMsg.encode(*msg)
	return json.Marshal(jsonMsg)
}

type GetAllMessagesReturn struct {
	Messages []ChatMessage `json:"messages"`
	Count int `json:"count"`
}

func GetAllMessages(_ EmptyArgs, page web.Page) (resp GetAllMessagesReturn, err error) {
	root, err := page.Root()
	if err != nil {
		return
	}

	// chatButton, ok, err := root.MaybeSubNode("div#chat-toggle-button", "")
	// if err != nil {
	//	return
	// }
	// if !ok {
	//	return resp, protocol.WebpageError("Could not find public chat button")
	// }
	// err = chatButton.Click()
	// if err != nil {
	//	return
	// }

	chatPanel, ok, err := root.MaybeSubNode(
		`div[data-test="publicChat"]>div[data-test="chatMessages"]`, "")
	if err != nil {
		return
	}
	if !ok {
		return resp, protocol.WebpageError("Could not find public chat panel")
	}

	messageCandidates, err := chatPanel.SubNodes(`span`)
	if err != nil {
		return
	}

	msgNodes := []web.Noder{}
	for _, c := range messageCandidates {
		_, found, err := c.MaybeSubNode(`p[data-test="chatUserMessageText"]`, "")
		if err != nil {
			return resp, err
		}
		if found {
			msgNodes = append(msgNodes, c)
		}
	}

	messages := []ChatMessage{}
	for _, n := range msgNodes {
		scope := ChatScopePublic

		userNode, ok, err := n.MaybeSubNode(
			`div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)`, "")
		if err != nil {
			return resp, err
		}
		var user string
		if ok {
			user, _ = userNode.Text()
		}

		textNode, ok, err := n.MaybeSubNode(
			`[data-test="chatUserMessageText"]`, "")
		if err != nil {
			return resp, err
		}
		var text string
		if ok {
			text, _ = textNode.Text()
		}

		timeNode, ok, err := n.MaybeSubNode(
			`div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > time`, "")
		if err != nil {
			return resp, err
		}
		var timelit string
		if ok {
			timelit, _, _ = timeNode.Attribute("datetime")
		}
		timestamp, err := time.Parse("Mon Jan 2 2006 15:04:05 GMT-0700", timelit[0:33])

		messages = append(messages, ChatMessage{
			Scope: scope,
			User: user,
			Text: text,
			Time: timestamp,
		})
	}

	time1, _ := time.Parse("2.1.2006 15:04:05", "11.8.2023 08:20:00")
	time2, _ := time.Parse("2.1.2006 15:04:05", "10.8.2023 19:12:00")
	time2 = time.Now()
	return GetAllMessagesReturn{
		Messages: append(messages, []ChatMessage{
			{
				Scope: ChatScopePrivate,
				User: "Alice",
				Text: "Good morning!",
				Time: time1,
			},
			{
				Scope: ChatScopePublic,
				User: "Bob",
				Text: "Good evening!",
				Time: time2,
			},
		}...),
		Count: 2,
	}, nil
}

type jsonChatMessage struct {
	Scope string `json:"scope"`
	User  string `json:"user"`
	Text  string `json:"text"`
	Timestamp int64 `json:"timestamp"`
}

func (jm *jsonChatMessage) encode(in ChatMessage) {
	switch in.Scope {
	case ChatScopePrivate:
		jm.Scope = "private"
	case ChatScopePublic:
		jm.Scope = "public"
	}
	jm.Timestamp = in.Time.Unix()
	jm.User = in.User
	jm.Text = in.Text
}

func (jm *jsonChatMessage) decode() ChatMessage {
	timestamp := time.Unix(jm.Timestamp, 0)
	var scope ChatScope
	switch jm.Scope {
	case "public":
		scope = ChatScopePublic
	case "private":
		scope = ChatScopePrivate
	}

	return ChatMessage{
		Scope: scope,
		Time: timestamp,
		User: jm.User,
		Text: jm.Text,
	}
}

func ChatMessageHandler(msg []byte) {
	parseMsg := jsonChatMessage{}

	err := json.Unmarshal(msg, &parseMsg)
	if err != nil {
		log.Println(err)
	}

	message := parseMsg.decode()
	log.Println(message)
}
