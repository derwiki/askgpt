package common

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BardApiKey       string
	OpenAIApiKey     string
	MaxTokens        int
	PromptPrefix     string
	LLMModels        []string
	SkipHistory      bool
	HistoryLineCount int
}

func LoadConfig() (Config, error) {
	config := Config{}

	config.PromptPrefix = os.Getenv("PROMPT_PREFIX")
	num, err := strconv.Atoi(os.Getenv("HISTORY_LINE_COUNT"))
	if err != nil {
		config.HistoryLineCount = 10 // default
	} else {
		config.HistoryLineCount = num
	}

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
		config.MaxTokens = 200
	} else {
		maxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			return config, err
		}
		config.MaxTokens = maxTokens
	}

	var skipHistory bool
	var useGpt4 bool
	var useBard bool
	var useInfo bool
	flag.BoolVar(&skipHistory, "skip-history", false, "If set, history will not be written to or read from.")
	flag.BoolVar(&useGpt4, "gpt4", false, "If set, shortcut to LLM_MODELS=gpt-4")
	flag.BoolVar(&useBard, "bard", false, "If set, shortcut to LLM_MODELS=bard")
	flag.BoolVar(&useInfo, "info", false, "If set, show info and above logs")
	flag.Parse()
	log.Info().Msg(fmt.Sprintf("skipHistory: %b", skipHistory))
	log.Info().Msg(fmt.Sprintf("useGpt4: %b", useGpt4))
	log.Info().Msg(fmt.Sprintf("useBard: %b", useBard))
	log.Info().Msg(fmt.Sprintf("useInfo: %b", useInfo))
	config.SkipHistory = skipHistory
	if useGpt4 {
		config.LLMModels = []string{openai.GPT4}
	} else if useBard {
		config.LLMModels = []string{"bard"}
	}

	if useInfo {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return config, nil
}
