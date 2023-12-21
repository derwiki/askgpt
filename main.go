package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/derwiki/askgpt/clients/anthropic"
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

	var llmFuncMap sync.Map

	llmFuncMap.Store(
		openai.GPT3Dot5Turbo, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		openai.GPT4, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		openai.GPT4TurboPreview, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		"text-davinci-003", openaiclient.GetTextCompletion)
	llmFuncMap.Store(
		"bard", google.GetBardCompletion)
	llmFuncMap.Store(
		"claude-2.1", anthropic.GetChatCompletions)

	var llmRequests []LLMRequest

	log.Info().Msg(fmt.Sprintf("Loading LLM_MODELS: %s", config.LLMModels))
	for _, model := range config.LLMModels {
		fn, ok := llmFuncMap.Load(model)
		if ok {
			fnTyped, ok := fn.(func(string, common.Config, string) (string, error))
			if !ok {
				log.Error().Msg(fmt.Sprintf("Stored function not of expected type: %T", fn))
				continue
			}
			// filter out models we don't have an API key for
			if model == "bard" && len(config.BardApiKey) == 0 {
				log.Info().Msg("excluding bard, missing api key")
			} else if model != "bard" && len(config.OpenAIApiKey) == 0 {
				log.Info().Msg(fmt.Sprintf("excluding %s, missing api key", model))
			} else {
				llmRequest := LLMRequest{
					Prompt: prompt,
					Config: config,
					Model:  model,
					Fn:     fnTyped,
				}
				log.Info().Msg(fmt.Sprintf("Loaded model: %s", llmRequest.Model))
				llmRequests = append(llmRequests, llmRequest)
			}
		} else {
			log.Error().Msg(fmt.Sprintf("Unknown LLM model: %s", model))
		}
	}

	if len(llmRequests) == 0 {
		log.Error().Msg(fmt.Sprintf("No LLMModels selected that have an API key set: %s", config.LLMModels))
	}

	results := make(chan LLMResponse, len(llmRequests)*2)

	for _, llmRequest := range llmRequests {
		wg.Add(1)
		go func(j *LLMRequest) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- LLMResponse{Output: output, Err: err, Model: j.Model}
			err = common.WriteHistory(config, fmt.Sprintf("A(%s): %s", j.Model, output))
		}(&llmRequest)
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
