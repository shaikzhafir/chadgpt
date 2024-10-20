package anthropic

import (
	"context"
	"errors"
	"log"

	"chat-backend/conversor"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

var sessionMap = make(map[string][]anthropic.MessageParam)

type ClaudeConversor struct {
	client       *anthropic.Client
	systemPrompt string
}

func New(authKey string) conversor.Conversor {
	return &ClaudeConversor{
		client: anthropic.NewClient(
			option.WithAPIKey(authKey),
		),
		systemPrompt: "You are an experienced programmer",
	}
}

func (c *ClaudeConversor) Ask(ctx context.Context, sessionID, userMessage string) (string, string, error) {
	// check if req has session-id. if have, then we use it to reference context

	// lookup sessionID in sessionMap
	// if sessionID is found, then we use it to reference context
	// if sessionID is not found, then we create a new sessionID
	// and add it to the sessionMap
	_, ok := sessionMap[sessionID]
	if !ok {
		sessionMap[sessionID] = []anthropic.MessageParam{}
	}

	if c.client == nil {
		return "", "", errors.New("client is nil")
	}
	if userMessage == "" {
		return "", "", errors.New("userMessage is empty")
	}

	// Add the new user message to the conversation history
	newUserMsg := anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage))
	messagesToSend := append(sessionMap[sessionID], newUserMsg)
	log.Printf("messages: %+v", messagesToSend)
	message, err := c.client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(c.systemPrompt),
		}),
		Messages:      anthropic.F(messagesToSend),
		Model:         anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		StopSequences: anthropic.F([]string{"```\n"}),
	})

	if err != nil {
		return "", "", err
	}

	if message == nil {
		return "", "", errors.New("response is nil")
	}
	// dont add new message to history until we know it was successful
	assistantMsg := message.ToParam()
	sessionMap[sessionID] = append(sessionMap[sessionID], newUserMsg, assistantMsg)

	return message.Content[0].Text, sessionID, nil
}
