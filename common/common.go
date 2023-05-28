package common

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

func LoadConfig() (Config, error) {
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

func refineAnswers() {
	//refine := fmt.Sprintf("Which of the following answers is best? \n\n%s\n\n%s\n\n%s\n\n%s", gpt3TurboRes, gpt3Davinci003Res, gpt3Davinci002Res, textDavinci002Res)
	// refined := libopenai.getChatCompletions(refine, config, openai.GPT4)
	fmt.Println("\n> Which of those answers is best?")
	// fmt.Println(refined)
}

func printUsage() {
	fmt.Println(`
Usage:
./chatgpt [PROMPT]
echo "PROMPT" | ./chatgpt
cat chatgpt.go | PROMPT_PREFIX="Improve this program" ./chatgpt

Description:
A Go command-line interface to communicate with OpenAI's ChatGPT API.
This program sends a prompt or question to the ChatGPT API for several models,
prints the generated response for each, and then sends all the responses to
gpt-4 to ask which is best.

Required Options:
PROMPT              The question or prompt to send to the ChatGPT API.

Environment Variables:
OPENAI_API_KEY      Your OpenAI API key.
MAX_TOKENS          The maximum number of tokens to generate in the response. (default: 100)
PROMPT_PREFIX       A prefix to add to each prompt.
GPT_MODEL           The model to use. If not specified, all models will be used.

Example:
./chatgpt "What is the capital of Ohio?"

> Chat Completion (gpt-3.5-turbo):
The capital of Ohio is Columbus.

> Chat Completion (text-davinci-003):
The capital of Ohio is Columbus.

> Chat Completion (text-davinci-002):
The capital of Ohio is Columbus.

> Text Completion (da-vinci-002):
The capital of Ohio is Columbus.

> Chat Completion (gpt-4):
The capital of Ohio is Columbus.

> Which of those answers is best?
All of the answers are the same and correct.`)
}

func GetPrompt(config Config) string {
	var prompt string

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
