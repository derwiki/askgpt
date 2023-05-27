package main

import (
	"fmt"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

// Define a Job struct to hold parameters for function
type Job struct {
	Prompt string
	Config Config
	Model  string
	Fn     func(string, Config, string) (string, error)
}

// Define a Result struct to hold results from function
type Result struct {
	Output string
	Err    error
}

func main() {
	prompt := getPrompt()
	var wg sync.WaitGroup

	// Define your config
	config := Config{}

	// Define your jobs
	jobs := []Job{
		{
			Prompt: prompt,
			Config: config,
			Model:  openai.GPT3Dot5Turbo,
			Fn:     getChatCompletions,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  openai.GPT4,
			Fn:     getChatCompletions,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  "text-davinci-003",
			Fn:     getTextCompletion,
		},
		{
			Prompt: prompt,
			Config: config,
			Model:  "text-davinci-003",
			Fn:     getBardCompletion,
		},
	}

	results := make(chan Result, len(jobs))

	for _, job := range jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			output, err := j.Fn(j.Prompt, j.Config, j.Model)
			results <- Result{Output: output, Err: err}
		}(job)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
