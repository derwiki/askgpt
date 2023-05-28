package main

import (
	"fmt"
	"github.com/derwiki/askgpt/clients/google"
	openai_client "github.com/derwiki/askgpt/clients/openai"
	"github.com/derwiki/askgpt/common"
	"github.com/sashabaranov/go-openai"
	"os"
	"sync"
)

// Job Define a struct to hold parameters for function
type Job struct {
	Prompt string
	Config common.Config
	Model  string
	Fn     func(string, common.Config, string) (string, error)
}

// Result Define a struct to hold results from function
type Result struct {
	Output string
	Err    error
	Model  string
	FnName string
}

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		fmt.Println("error: Fatal occurred in loadConfig")
		os.Exit(-1)
	}

	prompt := common.GetPrompt(config)

	var wg sync.WaitGroup

	// TODO(derwiki) if a model is specified, only call that model and exit

	// Define your jobs
	jobs := []Job{
		{
			Prompt: prompt,
			Config: config,
			Model:  openai.GPT3Dot5Turbo,
			Fn:     openai_client.GetChatCompletions,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  openai.GPT4,
			Fn:     openai_client.GetChatCompletions,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  "text-davinci-003",
			Fn:     openai_client.GetTextCompletion,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  "bard",
			Fn:     google.GetBardCompletion,
		},
	}

	results := make(chan Result, len(jobs))

	for _, job := range jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- Result{Output: output, Err: err, Model: j.Model, FnName: "foo"}
		}(job)
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
