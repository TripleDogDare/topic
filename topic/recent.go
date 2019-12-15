package topic

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func LastLine(ctx context.Context, r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var result string
	for scanner.Scan() {
		buf := scanner.Text()
		if len(strings.TrimSpace(buf)) > 0 {
			result = buf
		}
	}
	return string(result)
}

func Recent(ctx context.Context, topic *Topic) error {
	filepath := GetDefaultHistoryFilepath()
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err, "Could not open file: %s", filepath)
	}
	line := LastLine(ctx, f)
	if err := processLine(line, topic); err != nil {
		return err
	}
	return nil
}
