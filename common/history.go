package common

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type HistoryRecord struct {
	TimestampSec int
	TokenCount   int
	Line         string
}

func historyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Err(err)
		return ""
	}
	return filepath.Join(home, ".askgpt_history")
}

func HistoryLastNRecords(n int) []HistoryRecord {
	// Open the CSV file for reading
	f, err := os.Open(historyPath())
	if err != nil {
		if os.IsNotExist(err) {
			// Handle file not found error here
			log.Info().Msg("History file not found, creating new one.")
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
		timestampSec, err := strconv.Atoi(records[i][0])
		if err != nil {
			panic(err)
		}
		tokenCount, err := strconv.Atoi(records[i][1])
		if err != nil {
			panic(err)
		}
		str := records[i][2]

		record := HistoryRecord{TimestampSec: timestampSec, TokenCount: tokenCount, Line: str}
		log.Info().Int("TimestampSec", record.TimestampSec).Int("TokenCount", record.TokenCount).Int("Line length", len(record.Line)).Msg("HistoryRecord")
		buffer = append(buffer, record)
	}
	return buffer
}

func WriteHistory(config Config, line string) error {
	if config.SkipHistory {
		log.Info().Msg("SkipHistory set, not executing WriteHistory")
		return nil
	}

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
