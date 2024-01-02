package anthropic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/derwiki/askgpt/common"
	"github.com/rs/zerolog/log"
)

// https://docs.anthropic.com/claude/reference/complete_post

type CompletionRequest struct {
	Model             string `json:"model"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`
	Prompt            string `json:"prompt"`
}

type CompletionResponse struct {
	Completion string `json:"completion"`
}

func GetChatCompletions(prompt string, config common.Config, model string) (string, error) {
	if config.AnthropicApiKey == "" {
		log.Error().Msg("Anthropic API key is not set in the config")
		return "", errors.New("Anthropic API key is not set in the config")
	}

	client := &http.Client{}
	requestData := &CompletionRequest{
		Model:             model,
		MaxTokensToSample: 300, // This can also be part of the config if it varies
		Prompt:            fmt.Sprintf("Human: %s\n\nAssistant:", prompt),
	}

	requestBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error marshaling request data")
		return "", fmt.Errorf("error marshaling request data: %v", err)
	}

	BaseURL := "https://api.anthropic.com/v1/complete"
	req, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		log.Error().Str("error creating request", err.Error())
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.AnthropicApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error sending request")
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error reading response body")
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Error().Int("response code", resp.StatusCode).Str("response body", string(body)).Msg("Claude non-200 response")
		return "", fmt.Errorf("%s", string(body))
	}

	var completionResp CompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		log.Error().Str("error", err.Error()).Msg("error unmarshaling response")
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	log.Info().Str("completionsResp", completionResp.Completion).Msg("completionsResp")

	return completionResp.Completion, nil
}
