package common

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/rs/zerolog/log"
)

func GetPrompt(config Config) (string, string) {
	var prompt string

	args := flag.Args()
	log.Info().Msg(fmt.Sprintf("flag.Args(): %s", args))
	if len(args) > 0 {
		prompt = strings.TrimSpace(args[0])
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

	origPrompt := prompt
	// TODO(derwiki) make this model specific
	PromptModelMax := 4097
	prompt = config.PromptPrefix + prompt
	promptTokenCount, err := GetTokenCount(prompt)
	if promptTokenCount > PromptModelMax {
		log.Panic().Int("promptTokenCount", promptTokenCount).Int("PromptModelMax", PromptModelMax).Msg("Prompt token count exceeds model max")
		panic("Token count too long")
	}

	if !config.SkipHistory {
		lines := HistoryLastNRecords(config.HistoryLineCount)
		context := ""
		runningTokenCount := promptTokenCount
		for i, record := range lines {
			log.Debug().Str("line", record.Line).Int("index", i).Msg("History record")

			if record.TokenCount+runningTokenCount >= PromptModelMax {
				// nothing
			} else {
				context += record.Line + "\n"
				runningTokenCount += record.TokenCount
			}
		}
		log.Info().Int("count", runningTokenCount).Msg("runningTokenCount")

		prompt = context + "\n" + prompt
		err = WriteHistory(config, fmt.Sprintf("Q: %s", prompt))
		if err != nil {
			log.Panic().Str("err", err.Error()).Msg("WriteHistory")
			panic("WriteHistory panic")
		}
	} else {
		log.Info().Msg("SkipHistory set, not building history context")
	}
	log.Debug().Str("prompt", prompt).Msg("Prompt")
	return origPrompt, prompt
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
