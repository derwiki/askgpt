package common

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sashabaranov/go-openai"
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
	LLMModels    []string
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

	// read LLM models as an array
	models := strings.Split(os.Getenv("LLM_MODELS"), ",")
	if models[0] != "" {
		config.LLMModels = models
	} else {
		config.LLMModels = []string{openai.GPT3Dot5Turbo, openai.GPT4, "text-davinci-003", "bard"}
	}

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

	return config, nil
}

func hasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
		UsageAndQuit()
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func refineAnswers() {
	//refine := fmt.Sprintf("Which of the following answers is best? \n\n%s\n\n%s\n\n%s\n\n%s", gpt3TurboRes, gpt3Davinci003Res, gpt3Davinci002Res, textDavinci002Res)
	// refined := libopenai.getChatCompletions(refine, config, openai.GPT4)
	fmt.Println("\n> Which of those answers is best?")
	// fmt.Println(refined)
}

func UsageAndQuit() {
	fmt.Println(`UsageAndQuit: askgpt [PROMPT]

    PROMPT        A string prompt to send to the GPT models, surrounded by quotes if it has spaces.

    Environment variables:
      PROMPT_PREFIX    A prefix to add to the prompt read from STDIN.

    Examples:
      askgpt "What is the meaning of life?"
      echo "review this source code" | PROMPT_PREFIX="Generate a code review:" askgpt`)
	os.Exit(-1)
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
		UsageAndQuit()
	}
	// TODO(derwiki) use https://github.com/sugarme/tokenizer to verify token count
	return config.PromptPrefix + prompt
}
