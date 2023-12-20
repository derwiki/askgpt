package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/derwiki/askgpt/clients/google"
	openaiclient "github.com/derwiki/askgpt/clients/openai"
	"github.com/derwiki/askgpt/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
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
		openai.GPT3Dot5Turbo:    openaiclient.GetChatCompletions,
		openai.GPT4:             openaiclient.GetChatCompletions,
		openai.GPT4TurboPreview: openaiclient.GetChatCompletions,
		"text-davinci-003":      openaiclient.GetTextCompletion,
		"bard":                  google.GetBardCompletion,
	}

	var llmRequests []LLMRequest

	log.Info().Msg(fmt.Sprintf("Loading LLM_MODELS: %s", config.LLMModels))
	for _, model := range config.LLMModels {
		if fn, ok := llmFuncMap[model]; ok {
			// filter out models we don't have an API key for
			if model == "bard" && len(config.BardApiKey) == 0 {
				log.Info().Msg("excluding bard, missing api key")
			} else if model != "bard" && len(config.OpenAIApiKey) == 0 {
				log.Info().Msg(fmt.Sprintf("excluding %s, missing api key", model))
			} else {
				llmRequests = append(llmRequests, LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  model,
					Fn:     fn,
				})
			}
		} else {
			log.Error().Msg(fmt.Sprintf("Unknown LLM model: %s", model))
		}
	}

	if len(llmRequests) == 0 {
		log.Error().Msg(fmt.Sprintf("No LLMModels selected that have an API key set: %s", config.LLMModels))
	}

	results := make(chan LLMResponse, len(llmRequests))

	for _, llmRequest := range llmRequests {
		wg.Add(1)
		go func(j LLMRequest) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- LLMResponse{Output: output, Err: err, Model: j.Model}
			err = common.WriteHistory(config, fmt.Sprintf("A(%s): %s", j.Model, output))
		}(llmRequest)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err == nil {
			fmt.Println(fmt.Sprintf("\nA(%s): %s", result.Model, result.Output))
		} else {
			fmt.Println(fmt.Sprintf("\nA(%s): Error: %s", result.Model, result.Err))
		}
		os.Stdout.Sync()
	}
}
