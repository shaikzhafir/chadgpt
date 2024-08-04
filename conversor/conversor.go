package conversor

import (
	"context"
	"errors"
	"log"

	gogpt "github.com/sashabaranov/go-openai"
)

type Conversor struct {
	client         *gogpt.Client
	messageHistory []gogpt.ChatCompletionMessage
}

func NewConversor(authKey string) *Conversor {
	return &Conversor{
		client: gogpt.NewClient(authKey),
	}
}

func (c *Conversor) Ask(ctx context.Context, userMessage string) (string, error) {
	if c.client == nil {
		return "", errors.New("client is nil")
	}
	if userMessage == "" {
		return "", errors.New("userMessage is empty")
	}
	newMsg := gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: userMessage,
	}

	// dont add new message to history until we know it was successful
	messages := append(c.messageHistory, newMsg)
	log.Printf("messages: %+v", messages)
	req := gogpt.ChatCompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		Messages:  messages,
		MaxTokens: 1000,
		User:      "user",
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Choices[0].Message.Content == "" {
		return "", errors.New("response is nil")
	}
	// add the ask message and the response to the history
	c.messageHistory = append(c.messageHistory, newMsg, resp.Choices[0].Message)
	return resp.Choices[0].Message.Content, nil
}
