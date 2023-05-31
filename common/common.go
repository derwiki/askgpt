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
	fmt.Println(`Usage: askgpt [OPTIONS] PROMPT
    OPTIONS:
        --skip-history    Skip reading and writing to the history.
    PROMPT               A string prompt to send to the GPT models, surrounded by quotes if it has spaces.

    Environment variables:
      PROMPT_PREFIX       A prefix to add to the prompt read from STDIN.
      OPENAI_API_KEY      API key for OpenAI
      BARDAI_API_KEY      API key for Bard AI
      LLM_MODELS          Comma-separated list of LLM models
      MAX_TOKENS          Maximum number of tokens for a prompt

    Examples:
      askgpt "Generate go code to iterate over a list"
      askgpt "Refactor that generated code to be in a function names Scan()"
      cat main.go | PROMPT_PREFIX="Generate a code review: " askgpt
      askgpt --skip-history "Generate go code to iterate over a list"`)
	os.Exit(-1)
}
