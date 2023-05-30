package main

import (
	"fmt"
	"github.com/derwiki/askgpt/clients/google"
	openaiclient "github.com/derwiki/askgpt/clients/openai"
	"github.com/derwiki/askgpt/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	config, err := common.LoadConfig()
	if err != nil {
		fmt.Println("error: Fatal occurred in loadConfig")
		common.UsageAndQuit()
	}

	prompt := common.GetPrompt(config)

	var wg sync.WaitGroup

	var llmFuncMap = map[string]func(string, common.Config, string) (string, error){
		openai.GPT3Dot5Turbo: openaiclient.GetChatCompletions,
		openai.GPT4:          openaiclient.GetChatCompletions,
		"text-davinci-003":   openaiclient.GetTextCompletion,
		"bard":               google.GetBardCompletion,
	}

	var llmRequests []LLMRequest

	log.Info().Msg(fmt.Sprintf("Loading LLM_MODELS: %s", config.LLMModels))
	for _, model := range config.LLMModels {
		if fn, ok := llmFuncMap[model]; ok {
			llmRequests = append(llmRequests, LLMRequest{
				Prompt: prompt,
				Config: config,
				Model:  model,
				Fn:     fn,
			})
		} else {
			log.Error().Msg(fmt.Sprintf("Unknown LLM model: %s", model))
		}
	}

	results := make(chan LLMResponse, len(llmRequests))

	for _, llmRequest := range llmRequests {
		wg.Add(1)
		go func(j LLMRequest) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- LLMResponse{Output: output, Err: err, Model: j.Model}
			err = common.WriteHistory(fmt.Sprintf("A: %s", output))
		}(llmRequest)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if len(results) > 1 {
			fmt.Println("----------")
			fmt.Println(result.Model)
		}
		fmt.Println(result.Output)
	}
}
