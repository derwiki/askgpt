package openai

import (
	"context"

	"github.com/derwiki/askgpt/common"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type ChatGPTCompletionsResponseChoice struct {
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	LogProbs     string `json:"logprobs"`
	Text         string `json:"text"`
}
type ChatGPTCompletionsRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

func GetChatCompletions(content string, config common.Config, model string) (string, error) {
	if model == "" {
		model = openai.GPT4TurboPreview
	}
	// TODO(derwiki) assert model exists in openai package
	client := openai.NewClient(config.OpenAIApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: config.PromptPrefix + content,
				},
			},
			MaxTokens: config.MaxTokens,
		},
	)

	if err != nil {
		log.Error().Str("ChatCompletion error", err.Error())
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
