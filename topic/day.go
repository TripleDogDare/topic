package topic

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func isSameDate(a time.Time, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func isTopicOnDay(t *Topic, d time.Time) bool {
	if t == nil {
		return false
	}
	return isSameDate(d, t.Start) || isSameDate(d, t.End)
}

func isTopicStartOnDay(t *Topic, d time.Time) bool {
	if t == nil {
		return false
	}
	return isSameDate(d, t.Start)
}

func generateDay(ctx context.Context, date time.Time, r io.Reader, w io.Writer) error {
	var err error
	var currentTopic, previousTopic *Topic
	var totalDuration time.Duration
	lineCounter := 0
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	src := bufio.NewScanner(r)
	for src.Scan() {
		lineCounter++
		line := strings.TrimSpace(src.Text())
		if line == "" {
			// ignore empty lines
			continue
		}
		currentTopic = new(Topic)
		err = processLine(line, currentTopic)
		if err != nil {
			return fmt.Errorf("Unable to process: %w", err)
		}
		if previousTopic == nil {
			previousTopic = currentTopic
			continue
		}
		previousTopic.End = currentTopic.Start
		if !previousTopic.Empty() && isTopicStartOnDay(previousTopic, date) {
			err := previousTopic.CsvWrite(csvWriter)
			if err != nil {
				return err
			}
			totalDuration = totalDuration + previousTopic.Duration()
		}
		previousTopic = currentTopic
	}
	if src.Err() != nil {
		return fmt.Errorf("Could not read line %d: %w", lineCounter, src.Err())
	}

	if previousTopic != nil && !previousTopic.Empty() && isTopicOnDay(previousTopic, date) {
		// emit last topic
		err = previousTopic.CsvWrite(csvWriter)
		if err != nil {
			return err
		}
		totalDuration = totalDuration + previousTopic.Duration()
	}
	csvWriter.Flush()

	// emit totals
	fmt.Fprintln(w, "\nTotals")
	err = csvWriter.Write([]string{
		fmt.Sprintf("%.02f", totalDuration.Hours()),
		totalDuration.Round(time.Second).String(),
	})
	if err != nil {
		return err
	}
	return nil
}

func DayReport(ctx context.Context, w io.Writer, day time.Time) error {
	filepath := GetDefaultHistoryFilepath()
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Could not open file: %s: %w", filepath, err)
	}
	if err := generateDay(ctx, day, f, w); err != nil {
		return fmt.Errorf("Error processing report: %w", err)
	}
	return nil
}
