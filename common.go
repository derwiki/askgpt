package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BardApiKey   string
	OpenAIApiKey string
	MaxTokens    int
	PromptPrefix string
	Model        string
}

func loadConfig() (Config, error) {
	config := Config{}

	config.PromptPrefix = os.Getenv("PROMPT_PREFIX")

	openAiApiKey := os.Getenv("OPENAI_API_KEY")
	if openAiApiKey == "" {
		apiKeyBytes, err := ioutil.ReadFile("./.openai_key")
		if err != nil {
			return config, err
		}
		openAiApiKey = strings.TrimSpace(string(apiKeyBytes))
	}
	config.OpenAIApiKey = openAiApiKey

	bardAiApiKey := os.Getenv("BARDAI_API_KEY")
	if bardAiApiKey == "" {
		apiKeyBytes, err := ioutil.ReadFile("./.bardai_key")
		if err != nil {
			return config, err
		}
		bardAiApiKey = strings.TrimSpace(string(apiKeyBytes))
	}
	config.BardApiKey = bardAiApiKey

	maxTokensStr := os.Getenv("MAX_TOKENS")
	if maxTokensStr == "" {
		config.MaxTokens = 100
	} else {
		maxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			return config, err
		}
		config.MaxTokens = maxTokens
	}

	config.Model = os.Getenv("GPT_MODEL")

	return config, nil
}

func hasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func getPrompt() string {
	config, err := loadConfig()
	var prompt string
	if err != nil {
		fmt.Println("error: Fatal occurred in loadConfig")
		os.Exit(-1)
	}

	if len(os.Args) > 1 {
		prompt = os.Args[1]
	} else if hasStdinInput() {
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Split(bufio.ScanBytes)
		var buffer bytes.Buffer
		for scanner.Scan() {
			buffer.Write(scanner.Bytes())
		}

		prompt = strings.TrimSpace(buffer.String())
	} else {
		fmt.Println("error: No prompt found in args or STDIN")
		printUsage()
		os.Exit(-1)
	}
	return config.PromptPrefix + prompt
}
