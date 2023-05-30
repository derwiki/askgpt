package common

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func HasStdinInput() bool {
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
	} else if HasStdinInput() {
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
	PromptModelMax := 8097
	prompt = config.PromptPrefix + prompt
	prompTokenCount, err := GetTokenCount(prompt)
	if prompTokenCount > PromptModelMax {
		panic("token count too long")
	}
	lines := HistoryLastNRecords(4)
	context := ""
	runningTokenCount := prompTokenCount
	for i, record := range lines {
		fmt.Println("record", record, "i", i)

		if record.TokenCount+runningTokenCount >= PromptModelMax {
			// nothing
		} else {
			context += record.Line + "\n"
			runningTokenCount += record.TokenCount
		}
	}
	fmt.Println("runningTokenCount:", runningTokenCount)

	prompt = context + "\n" + prompt
	fmt.Println("prompt", prompt)
	err = WriteHistory(fmt.Sprintf("Q: %s", prompt))
	if err != nil {
		panic("bar")
	}
	return prompt
}

func historyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filepath.Join(home, ".askgpt_history")
}

func WriteHistory(line string) error {

	f, err := os.OpenFile(historyPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	tokenCount, _ := GetTokenCount(line)
	timestamp := time.Now().Unix()
	escapedLine := "\"" + strings.ReplaceAll(line, "\"", "\"\"") + "\""
	buffer := fmt.Sprintf("%d,%d,%s", timestamp, tokenCount, escapedLine)

	w := bufio.NewWriter(f)
	_, err = fmt.Fprintln(w, buffer)
	if err != nil {
		return err
	}
	return w.Flush()
}

func GetTokenCount(line string) (int, error) {
	encoding := "r50k_base"
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return -1, err
	}
	tokens := tke.Encode(line, nil, nil)
	return len(tokens), nil
}

type HistoryRecord struct {
	TimestampSec int
	TokenCount   int
	Line         string
}

func HistoryLastNRecords(n int) []HistoryRecord {
	// Open the CSV file for reading
	f, err := os.Open(historyPath())
	if err != nil {
		if os.IsNotExist(err) {
			// Handle file not found error here
			fmt.Println("History file not found, creating new one.")
			f, err = os.Create(historyPath())
			if err != nil {
				panic(err)
			}
		} else {
			// Handle other errors
			panic(err)
		}
	}

	defer f.Close()

	// Parse the CSV file
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	// Print the last n records (or all if there are less than n)
	if len(records) < n {
		n = len(records)
	}
	var buffer []HistoryRecord
	for i := len(records) - n; i < len(records); i++ {
		int1, err := strconv.Atoi(records[i][0])
		if err != nil {
			panic(err)
		}
		int2, err := strconv.Atoi(records[i][1])
		if err != nil {
			panic(err)
		}
		str := records[i][2]

		record := HistoryRecord{TimestampSec: int1, TokenCount: int2, Line: str}
		fmt.Println(record)
		buffer = append(buffer, record)
	}
	return buffer
}

func AppendToCsv(line string) error {
	f, err := os.OpenFile(historyPath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write([]string{line})
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}
