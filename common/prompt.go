package common

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func GetPrompt(config Config) string {
	var prompt string

	args := flag.Args()
	log.Info().Msg(fmt.Sprintf("flag.Args(): %s", args))
	if len(args) > 0 {
		prompt = args[0]
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
	// TODO(derwiki) make this model specific
	PromptModelMax := 4097
	prompt = config.PromptPrefix + prompt
	promptTokenCount, err := GetTokenCount(prompt)
	if promptTokenCount > PromptModelMax {
		log.Panic().Msg(fmt.Sprintf("Prompt token count exceeds model max: %d/%d", promptTokenCount, PromptModelMax))
		panic("token count too long")
	}

	if !config.SkipHistory {
		lines := HistoryLastNRecords(config.HistoryLineCount)
		context := ""
		runningTokenCount := promptTokenCount
		for i, record := range lines {
			log.Debug().Msg(fmt.Sprintf("i: %d, record: %s", i, record.Line))

			if record.TokenCount+runningTokenCount >= PromptModelMax {
				// nothing
			} else {
				context += record.Line + "\n"
				runningTokenCount += record.TokenCount
			}
		}
		log.Info().Msg(fmt.Sprintf("runningTokenCount: %d", runningTokenCount))

		prompt = context + "\n" + prompt
		err = WriteHistory(config, fmt.Sprintf("Q: %s", prompt))
		if err != nil {
			panic("bar")
		}
	} else {
		log.Info().Msg("SkipHistory set, not building history context")
	}
	log.Debug().Msg("prompt: " + prompt)
	return prompt
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
