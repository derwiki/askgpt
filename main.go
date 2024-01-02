package main

import (
	"fmt"
	"os"
	"reflect"
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
		log.Error().Msg("Fatal occurred in loadConfig")
		common.UsageAndQuit()
	}

	if config.UseInfo {
		log.Info().Msg("in main(), setting global log level to InfoLevel")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	origPrompt, fullPrompt := common.GetPrompt(config)
	fmt.Printf("Q: %s\n", origPrompt)

	var wg sync.WaitGroup

	var llmFuncMap sync.Map

	llmFuncMap.Store(
		openai.GPT3Dot5Turbo, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		openai.GPT4, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		openai.GPT4TurboPreview, openaiclient.GetChatCompletions)
	llmFuncMap.Store(
		"bard", google.GetBardCompletion)
	llmFuncMap.Store(
		"claude-2.1", anthropic.GetChatCompletions)

	var llmRequests []LLMRequest

	log.Info().Strs("Loading LLM_MODELS", config.LLMModels).Msg("")
	for _, model := range config.LLMModels {
		fn, ok := llmFuncMap.Load(model)
		if ok {
			fnTyped, ok := fn.(func(string, common.Config, string) (string, error))
			if !ok {
				log.Error().Str("type", reflect.TypeOf(fn).String()).Msg("Stored function not of expected type")
				continue
			}
			// filter out models we don't have an API key for
			if model == "bard" && len(config.BardApiKey) == 0 {
				log.Info().Msg("excluding bard, missing api key")
			} else if model == "claude" && len(config.AnthropicApiKey) == 0 {
				log.Info().Msg("excluding claude, missing api key")
			} else if model != "bard" && len(config.OpenAIApiKey) == 0 {
				log.Info().Str("excluding model, missing api key", model).Msg("")
			} else {
				llmRequest := LLMRequest{
					Prompt: fullPrompt,
					Config: config,
					Model:  model,
					Fn:     fnTyped,
				}
				log.Info().Str("model", llmRequest.Model).Msg("Loaded model")
				llmRequests = append(llmRequests, llmRequest)
			}
		} else {
			log.Error().Str("model", model).Msg("Unknown LLM model")
		}
	}

	if len(llmRequests) == 0 {
		log.Error().Strs("models", config.LLMModels).Msg("No LLMModels selected that have an API key set")
	}

	results := make(chan LLMResponse, len(llmRequests)*2)

	for _, llmRequest := range llmRequests {
		wg.Add(1)
		go func(j LLMRequest) {
			defer wg.Done()
			log.Info().Str("model", j.Model).Msg("Calling j.Fn")
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
			fmt.Printf("A(%s): %s\n", result.Model, result.Output)
		} else {
			fmt.Printf("A(%s): Error: %s\n", result.Model, result.Err)
		}
		os.Stdout.Sync()
	}
}
