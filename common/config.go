package common

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	BardApiKey       string
	OpenAIApiKey     string
	AnthropicApiKey  string
	MaxTokens        int
	PromptPrefix     string
	LLMModels        []string
	SkipHistory      bool
	HistoryLineCount int
}

func LoadConfig() (Config, error) {
	// parse info first so other info logs in this file are outputted
	var useInfo bool
	var skipHistory bool
	var useGpt4 bool
	var useClaude bool
	var useBard bool
	flag.BoolVar(&useInfo, "info", false, "If set, show info and above logs")
	flag.BoolVar(&skipHistory, "skip-history", false, "If set, history will not be written to or read from.")
	flag.BoolVar(&useGpt4, "gpt4", false, "If set, shortcut to LLM_MODELS=gpt-4-1106-preview")
	flag.BoolVar(&useBard, "bard", false, "If set, shortcut to LLM_MODELS=bard")
	flag.BoolVar(&useClaude, "claude", false, "If set, shortcut to LLM_MODELS=claude-2.1")
	flag.Parse()
	log.Info().Msg(fmt.Sprintf("config useInfo: %b", useInfo))
	log.Info().Msg(fmt.Sprintf("config skipHistory: %b", skipHistory))
	log.Info().Msg(fmt.Sprintf("config useGpt4: %b", useGpt4))
	log.Info().Msg(fmt.Sprintf("config useBard: %b", useBard))
	log.Info().Msg(fmt.Sprintf("config useClaude: %b", useClaude))

	if useInfo {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg(fmt.Sprintf("config LoadConfig() enter"))
	config := Config{}

	config.PromptPrefix = os.Getenv("PROMPT_PREFIX")
	log.Info().Msg(fmt.Sprintf("config.PromptPrefix: %s", config.PromptPrefix))
	num, err := strconv.Atoi(os.Getenv("HISTORY_LINE_COUNT"))
	if err != nil {
		config.HistoryLineCount = 10 // default
	} else {
		config.HistoryLineCount = num
	}

	config.OpenAIApiKey = os.Getenv("OPENAI_API_KEY")
	log.Info().Msg(fmt.Sprintf("config.OpenAIApiKey length: %d", len(config.OpenAIApiKey)))
	config.BardApiKey = os.Getenv("BARDAI_API_KEY")
	log.Info().Msg(fmt.Sprintf("config.BardApiKey length: %d", len(config.BardApiKey)))
	config.AnthropicApiKey = os.Getenv("ANTHROPIC_API_KEY")
	log.Info().Msg(fmt.Sprintf("config.AnthropicApiKey length: %d", len(config.BardApiKey)))

	// read LLM models as an array
	log.Info().Msg(fmt.Sprintf("config LLM_MODELS: %s", os.Getenv("LLM_MODELS")))
	models := strings.Split(os.Getenv("LLM_MODELS"), ",")
	if models[0] != "" {
		config.LLMModels = models
	} else {
		config.LLMModels = []string{openai.GPT4TurboPreview, "bard", "claude-2.1"}
	}

	maxTokensStr := os.Getenv("MAX_TOKENS")
	log.Info().Msg(fmt.Sprintf("config maxTokenStr: %s", maxTokensStr))
	if maxTokensStr == "" {
		config.MaxTokens = 200
	} else {
		maxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			return config, err
		}
		config.MaxTokens = maxTokens
	}

	config.SkipHistory = skipHistory
	if useGpt4 {
		config.LLMModels = []string{openai.GPT4TurboPreview}
	} else if useBard {
		config.LLMModels = []string{"bard"}
	} else if useClaude {
		config.LLMModels = []string{"claude-2.1"}
	}
	return config, nil
}
