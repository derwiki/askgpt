package common

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func HasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal().Err(err)
		UsageAndQuit()
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func UsageAndQuit() {
	fmt.Println(`Usage: askgpt [OPTIONS] [PROMPT]
    OPTIONS:
        --info            Show info and above logs.
        --skip-history    Skip writing to the history.
        --gpt4            Use GPT-4 model.
        --bard            Use Bard model.
        --claude          Use Claude model.
    PROMPT               A string prompt to send to the GPT models, surrounded by quotes if it has spaces.

    Environment variables:
      PROMPT_PREFIX       A prefix to add to the prompt read from STDIN.
      OPENAI_API_KEY      API key for OpenAI
      ANTHROPIC_API_KEY   API key for Anthropic
      BARDAI_API_KEY      API key for Bard AI
      LLM_MODELS          Comma-separated list of LLM models
      MAX_TOKENS          Maximum number of tokens for a prompt
      HISTORY_LINE_COUNT  Number of history records to keep

    Examples:
      askgpt "Generate go code to iterate over a list"
      askgpt "Refactor that generated code to be in a function names Scan()"
      cat main.go | PROMPT_PREFIX="Generate a code review: " askgpt
      askgpt --skip-history "Generate go code to iterate over a list"
      askgpt --gpt4 "What is the meaning of life?"
      askgpt --bard "Tell me a story about a robot."
      askgpt --claude "Explain quantum computing in simple terms."`)
	os.Exit(-1)
}
