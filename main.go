package main

import (
	"fmt"
	"github.com/derwiki/askgpt/clients/google"
	openai_client "github.com/derwiki/askgpt/clients/openai"
	"github.com/derwiki/askgpt/common"
	"github.com/pkoukk/tiktoken-go"
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

	encoding := "r50k_base"
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return
	}
	numTokens := len(tke.Encode(prompt, nil, nil))
	fmt.Println(fmt.Sprintf("Number of tokens in prompt: %d", numTokens))
	var wg sync.WaitGroup

	var llmFuncMap = map[string]func(string, common.Config, string) (string, error){
		openai.GPT3Dot5Turbo: openai_client.GetChatCompletions,
		openai.GPT4:          openai_client.GetChatCompletions,
		"text-davinci-003":   openai_client.GetTextCompletion,
		"bard":               google.GetBardCompletion,
	}

	var llmRequests []LLMRequest

	fmt.Println("Loading LLM_MODELS: ", config.LLMModels)
	for _, model := range config.LLMModels {
		if fn, ok := llmFuncMap[model]; ok {
			llmRequests = append(llmRequests, LLMRequest{
				Prompt: prompt,
				Config: config,
				Model:  model,
				Fn:     fn,
			})
		} else {
			fmt.Printf("Unknown LLM model: %s\n", model)
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
