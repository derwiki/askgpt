package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/derwiki/askgpt/common"
	"io/ioutil"
	"net/http"
)

type Prompt struct {
	Text string `json:"text"`
}

type Candidate struct {
	Output        string          `json:"output"`
	SafetyRatings []SafetyRatings `json:"safetyRatings"`
}

type SafetyRatings struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

func GetBardCompletion(prompt string, config common.Config, model string) (string, error) {
	// Create the request body
	data := map[string]Prompt{"prompt": {Text: prompt}}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request to the API endpoint
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta2/models/text-bison-001:generateText?key=%s", config.BardApiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	if len(response.Candidates) > 0 {
		return response.Candidates[0].Output, nil
	}

	return "", fmt.Errorf("no candidate output found in response")
}
