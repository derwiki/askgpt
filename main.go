package main

import (
	"fmt"
	"github.com/derwiki/askgpt/clients/google"
	openai_client "github.com/derwiki/askgpt/clients/openai"
	"github.com/derwiki/askgpt/common"
	"github.com/sashabaranov/go-openai"
	"sync"
)

// LLMRequest Define a struct to hold parameters for function
type LLMRequest struct {
	Prompt string
	Config common.Config
	Model  string
	Fn     func(string, common.Config, string) (string, error)
}

// LLMResponse Define a struct to hold results from function
type LLMResponse struct {
	Output string
	Err    error
	Model  string
}

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		fmt.Println("error: Fatal occurred in loadConfig")
		common.UsageAndQuit()
	}

	prompt := common.GetPrompt(config)

	var wg sync.WaitGroup

	var llmRequests []LLMRequest

	fmt.Println("Loading LLM_MODELS: ", config.LLMModels)
	for _, model := range config.LLMModels {
		if model == openai.GPT3Dot5Turbo {
			llmRequests = append(llmRequests,
				LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  openai.GPT3Dot5Turbo,
					Fn:     openai_client.GetChatCompletions,
				})
		}
		if model == openai.GPT4 {
			llmRequests = append(llmRequests,
				LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  openai.GPT4,
					Fn:     openai_client.GetChatCompletions,
				})
		}
		if model == "text-davinci-003" {
			llmRequests = append(llmRequests,
				LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  "text-davinci-003",
					Fn:     openai_client.GetTextCompletion,
				})
		}
		if model == "bard" {
			llmRequests = append(llmRequests,
				LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  "bard",
					Fn:     google.GetBardCompletion,
				})
		}
	}

	results := make(chan LLMResponse, len(llmRequests))

	for _, llmRequest := range llmRequests {
		wg.Add(1)
		go func(j LLMRequest) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- LLMResponse{Output: output, Err: err, Model: j.Model}
		}(llmRequest)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result.Model)
		fmt.Println(result.Output)
		fmt.Println("----------")
	}
}
