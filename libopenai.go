package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const apiBaseURL = "https://api.openai.com/v1/completions"

type TextCompletionResponse struct {
	Choices []ChatGPTCompletionsResponseChoice `json:"choices"`
}
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

func getTextCompletion(prompt string, config Config, model string) (string, error) {
	if model == "" {
		model = "text-davinci-003"
	}
	textCompletionRequest := ChatGPTCompletionsRequest{
		Model:     model,
		Prompt:    config.PromptPrefix + prompt,
		MaxTokens: config.MaxTokens,
	}
	requestBodyBytes, err := json.Marshal(textCompletionRequest)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", apiBaseURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenAIApiKey))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// close the response body at the end of the function
	defer response.Body.Close()

	var responseBody TextCompletionResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		log.Fatal(err)
	}

	if len(responseBody.Choices) == 0 {
		log.Fatal("No choices found in the response body.")
	}

	return strings.TrimSpace(responseBody.Choices[0].Text), nil
}

func getChatCompletions(content string, config Config, model string) (string, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}
	// TODO(derwiki) assert model exists in openai package
	client := openai.NewClient(config.OpenAIApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: config.PromptPrefix + content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
