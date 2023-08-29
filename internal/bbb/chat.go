package bbb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/api"
	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
)

type ChatScope int

const (
	ChatScopePrivate ChatScope = iota
	ChatScopePublic
)

func (s *ChatScope) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	switch value {
	case "public":
		*s = ChatScopePublic
	case "private":
		*s = ChatScopePrivate
	default:
		return protocol.UserError(
			fmt.Sprintf("Unknown scope ('public' or 'private') value:%s", value))
	}

	return nil
}

func (s ChatScope) MarshalJSON() ([]byte, error) {
	var value string
	switch s {
	case ChatScopePrivate:
		value = "private"
	case ChatScopePublic:
		value = "public"
	default:
		return nil, protocol.InternalError(
			fmt.Sprintf("Invalid value for ChatScope: %v", s))
	}

	return json.Marshal(value)
}

type ChatMessage struct {
	Scope ChatScope `json:"scope"`
	Time  time.Time `json:"time"`
	User  string    `json:"user"`
	Text  string    `json:"text"`
}

func (msg *ChatMessage) MarshalJSON() ([]byte, error) {
	jsonMsg := jsonChatMessage{}
	jsonMsg.encode(*msg)
	return json.Marshal(jsonMsg)
}

type GetAllMessagesReturn struct {
	Messages []ChatMessage `json:"messages"`
	Count    int           `json:"count"`
}

func getChatPanel(root web.Noder, button web.Noder, public bool) (web.Noder, error) {
	log.Printf("getChatPanel(root, %v, %t)\n", button, public)
	err := button.Click()
	if err != nil {
		return nil, err
	}

	chatPanel, ok, err := findChatPanel(root, public)
	log.Printf("findChatPanel: %s, %t, %s\n", chatPanel, ok, err)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = button.Click()
		if err != nil {
			return nil, err
		}
		chatPanel, ok, err = findChatPanel(root, public)
		log.Printf("findChatPanel: %s, %t, %s\n", chatPanel, ok, err)
		if err != nil {
			return nil, err
		}
		if !ok {
			log.Println("throwing error")
			return nil, protocol.WebpageError("Could not find chat panel")
		}
	}
	return chatPanel, nil
}

func findChatPanel(root web.Noder, public bool) (web.Noder, bool, error) {
	log.Printf("findChatPanel(root, %t)\n", public)
	var mainDivTestData string
	if public {
		mainDivTestData = "publicChat"
	} else {
		mainDivTestData = "privateChat"
	}
	return root.MaybeSubNode(
		fmt.Sprintf(`div[data-test="%s"]>div[data-test="chatMessages"]`,
			mainDivTestData),
		"")
}

func parseMessageNode(node web.Noder, scope ChatScope) (msg ChatMessage, err error) {
	userNode, ok, err := node.MaybeSubNode(
		`div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)`, "")
	if err != nil {
		return
	}
	var user string
	if ok {
		user, _ = userNode.Text()
	}

	textNodes, err := node.SubNodes(
		`[data-test="chatUserMessageText"]`)
	if err != nil {
		return
	}
	var texts []string
	for _, tn := range textNodes {
		t, err := tn.Text()
		if err == nil {
			texts = append(texts, t)
		}
	}
	text := strings.Join(texts, "\n")

	timeNode, ok, err := node.MaybeSubNode(
		`div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > time`, "")
	if err != nil {
		return
	}
	var timelit string
	if ok {
		timelit, _, _ = timeNode.Attribute("datetime")
	}
	timestamp, err := time.Parse("Mon Jan 2 2006 15:04:05 GMT-0700", timelit[0:33])

	return ChatMessage{
		Scope: scope,
		User:  user,
		Text:  text,
		Time:  timestamp,
	}, nil
}

func getScopeMessages(root web.Noder, button web.Noder) (messages []interface{}, err error) {
	t, err := button.Text()
	if err != nil {
		return
	}
	var scope ChatScope
	if t == "Public Chat" {
		scope = ChatScopePublic
	} else {
		scope = ChatScopePrivate
	}

	chatPanel, err := getChatPanel(root, button, scope == ChatScopePublic)
	if err != nil {
		return
	}

	time.Sleep(time.Millisecond * 200)
	messageCandidates, err := chatPanel.SubNodes(`span`)
	if err != nil {
		return
	}

	msgNodes := []web.Noder{}
	for _, c := range messageCandidates {
		_, found, err := c.MaybeSubNode(`p[data-test="chatUserMessageText"]`, "")
		if err != nil {
			return messages, err
		}
		if found {
			msgNodes = append(msgNodes, c)
		}
	}

	for _, n := range msgNodes {
		msg, err := parseMessageNode(n, scope)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}

	button.Click()
	return
}

func collectChatMessages(root web.Noder) (messages []interface{}, err error) {
	chatButtons, err := root.SubNodes("div#chat-toggle-button")
	if err != nil {
		return
	}
	log.Printf("number of buttons: %d\n", len(chatButtons))
	for _, button := range chatButtons {
		button.LogToConsole("chatbutton")
		t, err := button.Text()
		log.Printf("button.Text: %s, %s\n", t, err)
		scopeMessages, err := getScopeMessages(root, button)
		if err != nil {
			return messages, err
		}
		messages = append(messages, scopeMessages...)
	}
	return
}

func ChatAllMessages(_ EmptyArgs, page web.Page) (res api.CollectionGetFunctionResult, err error) {
	root, err := page.Root()
	if err != nil {
		return
	}

	messages, err := collectChatMessages(root)
	if err != nil {
		return
	}

	return api.CollectionGetFunctionResult{
		Members: messages,
	}, nil
}

func ChatSubscribeMessages(_ EmptyArgs, page web.Page) (res api.CollectionSubsribeFunctionResult, err error) {
	log.Println("ChatSubscribeMessage")
	root, err := page.Root()
	if err != nil {
		return
	}
	log.Println("Got root")
	userList, err := getUserList(root)
	if err != nil {
		return
	}
	log.Println("Got userList")
	subtreeChanges, err := userList.SubscribeSubtree()
	if err != nil {
		return
	}
	log.Println("Subscribed")
	returnedTicks := make(chan interface{})
	res.Channel = returnedTicks
	go func() {
		for {
			log.Println("Waiting for change")
			<-subtreeChanges
			log.Println("Received change")
			returnedTicks <- struct{}{}
			log.Println("Sent struct")
		}
	}()
	log.Println("Launched go routine")
	return
}

type jsonChatMessage struct {
	Scope     string `json:"scope"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
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
		Time:  timestamp,
		User:  jm.User,
		Text:  jm.Text,
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

type SendChatMessageArgs struct {
	Text      string `json:"text"`
	Recipient string `json:"recipient"`
	Scope     string `json:"scope"`
}

func LogUserList(_ EmptyArgs, page web.Page) (res EmptyResult, err error) {
	root, err := page.Root()
	if err != nil {
		return
	}
	userList, err := getUserList(root)
	if err != nil {
		return
	}
	args := []interface{}{
		userList,
	}
	_, err = root.SubscribeSubtree()
	if err != nil {
		return
	}
	_, err = page.Exec(`console.log("hello eins");console.log(arguments[0]);`, args)
	if err != nil {
		return
	}
	_, err = page.Exec(`console.log("hello zwei!");`, args)
	if err != nil {
		return
	}
	return
}

func getUserList(root web.Noder) (userList web.Noder, err error) {
	userList, err = root.SubNode(`div[aria-label="Users list"]>div>div[aria-label="Users list"]`, "")
	return
}

func getUserButton(root web.Noder, user string) (userButton web.Noder, err error) {
	userList, err := getUserList(root)
	if err != nil {
		return
	}
	userRegex := fmt.Sprintf("^%s\\s+(\\(You\\))?$", user)
	userButton, err = userList.SubNode("span div[role=\"button\"]", userRegex)
	return
}

func getPrivateChatButton(root web.Noder) (chatButton web.Noder, err error) {
	chatButton, err = root.SubNode("ul[role=\"menu\"] li[role=\"menuitem\"]", "Start a private chat")
	return
}

func SendChatMessage(args SendChatMessageArgs, page web.Page) (res EmptyResult, err error) {
	root, err := page.Root()
	if err != nil {
		return
	}

	if args.Scope == "private" {
		if args.Recipient == "" {
			return res, protocol.UserError("Recipient argument empty")
		}
		userButton, err := getUserButton(root, args.Recipient)
		if err != nil {
			return res, err
		}
		err = userButton.Click()
		if err != nil {
			return res, err
		}
		chatButton, err := getPrivateChatButton(root)
		if err != nil {
			return res, err
		}
		err = chatButton.Click()
		if err != nil {
			return res, err
		}
	} else {
		_, found, err := findChatPanel(root, true)
		if err != nil {
			return res, err
		}
		if !found {
			button, err := root.SubNode("div#chat-toggle-button", "^Public Chat$")
			if err != nil {
				return res, err
			}
			err = button.Click()
			if err != nil {
				return res, err
			}
		}
	}
	messageArea, err := root.SubNode("textarea#message-input", "")
	if err != nil {
		return
	}
	messageArea.SendKeys(args.Text)
	sendButton, err := root.SubNode("button[aria-label=\"Send message\"]", "")
	if err != nil {
		return
	}
	err = sendButton.Click()
	return
}
