package topic

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func processLine(line string, topic *Topic) error {
	splits := strings.SplitN(line, " ", 2)
	if len(splits) > 2 {
		return errors.Errorf("Invalid line with %d splits", len(splits))
	}
	timestamp, err := time.Parse(time.RFC3339, splits[0])
	if err != nil {
		return err
	}
	topic.Start = timestamp
	topic.End = time.Now()
	if len(splits) == 2 {
		topic.Data = splits[1]
	}
	return nil
}

func WriteTopics(ctx context.Context, w io.Writer, ch <-chan Topic, done func()) {
	defer done()
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()
	for {
		select {
		case topic, ok := <-ch:
			if !ok {
				return
			}
			if topic.Empty() {
				continue
			}
			err := csvWriter.Write([]string{
				topic.Start.Format(time.RFC3339),
				topic.Duration().Round(time.Second).String(),
				topic.Data,
			})

			if err != nil {
				panic(err)
			}
		case <-ctx.Done():
			break
		}
	}
}

// Generates a report calculating the difference between topic timestamps
func generateReport(ctx context.Context, r io.Reader, w io.Writer) error {
	var err error
	var currentTopic, previousTopic *Topic

	readerSize := MaxTopicLength * 4
	src := bufio.NewReaderSize(r, readerSize)
	var lineBuilder strings.Builder
	emit := make(chan Topic, 10)
	// Print data
	ctx, cancel := context.WithCancel(ctx)
	go WriteTopics(ctx, w, emit, cancel)

	var line string
	for err == nil {
		// Get line
		buf, isPrefix, err := src.ReadLine()
		if err != nil {
			break
		}
		lineBuilder.Write(buf)
		if isPrefix {
			if err == io.EOF {
				return errors.Wrap(err, "Read buffer overflow and EOF condition is invalid")
			}
			continue
		}
		line = strings.TrimSpace(lineBuilder.String())
		lineBuilder.Reset()
		if line == "" {
			continue
		}
		currentTopic = new(Topic)
		err = processLine(line, currentTopic)
		if err != nil {
			return errors.Wrapf(err, "Unable to process: %s", line)
		}
		if previousTopic != nil {
			previousTopic.End = currentTopic.Start
			emit <- *previousTopic
		}
		previousTopic = currentTopic
	}
	emit <- *previousTopic
	close(emit)
	<-ctx.Done()
	if err != nil && err != io.EOF {
		errors.Wrapf(err, "Could not read lines %s", line)
		return err
	}
	return err
}

// Prints report to writer
func Report(ctx context.Context, w io.Writer) error {
	filepath := GetDefaultHistoryFilepath()
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err, "Could not open file: %s", filepath)
	}
	if err := generateReport(ctx, f, w); err != nil {
		return errors.Wrap(err, "Error processing report")
	}
	return nil
}
