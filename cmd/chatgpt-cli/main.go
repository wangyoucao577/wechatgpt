package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"
	openai "github.com/sashabaranov/go-openai"
)

var flags struct {
	apiKey string
}

func init() {
	flag.StringVar(&flags.apiKey, "apiKey", "", "Your api_key of OpenAI platform.")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	client := openai.NewClient(flags.apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		glog.Errorf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
