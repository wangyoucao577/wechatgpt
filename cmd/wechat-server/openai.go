package main

import (
	"context"
	"flag"
	"time"

	"github.com/golang/glog"
	openai "github.com/sashabaranov/go-openai"
)

var openaiFlags struct {
	apiKey string
}

func init() {
	flag.StringVar(&openaiFlags.apiKey, "api_key", "", "Your api_key of OpenAI platform.")
}

func chatgpt(question string, timeout time.Duration) string {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := openai.NewClient(openaiFlags.apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		glog.Errorf("ChatCompletion error: %v\n", err)
		return err.Error()
	}

	return resp.Choices[0].Message.Content
}
