package topic

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

func appendTopic(w io.Writer, topic string, date time.Time) error {
	s := fmt.Sprintf("%s %s\n", date.Format(time.RFC3339), topic)
	_, err := io.WriteString(w, s)
	return err
}
func Append(topic string) error {
	filepath := GetDefaultHistoryFilepath()
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, PermissionBitsFile)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err, "Could not open file: %s", filepath)
	}
	return appendTopic(f, topic, time.Now())
}
