package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

func getBardCompletion(content string) (string, error) {
	// Create the request body
	prompt := Prompt{Text: content}
	data := map[string]Prompt{"prompt": prompt}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request to the API endpoint
	BardApiKey := os.Getenv("BARDAI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta2/models/text-bison-001:generateText?key=%s", BardApiKey)
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

func hasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func main() {
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
		return
	}

	output, err := getBardCompletion(prompt)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)
}
