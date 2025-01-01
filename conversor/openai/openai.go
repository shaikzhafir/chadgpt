package openai

import (
	"chat-backend/conversor"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	gogpt "github.com/sashabaranov/go-openai"
)

type openAIConversor struct {
	client         *gogpt.Client
	sessionID      string                        // sessionID is used to keep track of the conversation and retain context
	messageHistory []gogpt.ChatCompletionMessage // messageHistory is used to keep track of the conversation and retain context
	model          string
}

func NewConversor(authKey string, baseUrl string, model string) conversor.Conversor {
	config := gogpt.DefaultConfig(authKey)
	if baseUrl != "" {
		config.BaseURL = baseUrl
	}
	if model == "" {
		model = "gpt-4o"
	}
	client := gogpt.NewClientWithConfig(config)
	return &openAIConversor{
		client:    client,
		sessionID: uuid.New().String(),
		model:     model,
	}
}

func (c *openAIConversor) Ask(ctx context.Context, sessionID, userMessage string) (string, string, error) {
	if c.client == nil {
		return "", "", errors.New("client is nil")
	}
	if userMessage == "" {
		return "", "", errors.New("userMessage is empty")
	}
	newMsg := gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: userMessage,
	}

	// dont add new message to history until we know it was successful
	messages := append(c.messageHistory, newMsg)
	log.Printf("messages: %+v", messages)
	req := gogpt.ChatCompletionRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: 1000,
		User:      "user",
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp.Choices[0].Message.Content == "" {
		return "", "", errors.New("response is nil")
	}
	log.Printf("response: %+v", resp.Choices[0].Message.Content)
	// add the ask message and the response to the history
	c.messageHistory = append(c.messageHistory, newMsg, resp.Choices[0].Message)
	return resp.Choices[0].Message.Content, sessionID, nil
}
