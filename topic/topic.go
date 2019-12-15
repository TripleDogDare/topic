package topic

import (
	"fmt"
	"os"
	"path"
	"time"
)

const DefaultDataDirectory = ".topic"
const DefaultTopicHistoryFile = "history"
const MaxTopicLength = 256
const PermissionBitsFile = 0644
const PermissionBitsDir = 0755

func GetDefaultFilepath(args ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	join := []string{
		home,
		DefaultDataDirectory,
	}
	// join = append(join, args...)
	return path.Join(append(join, args...)...)
}

func GetDefaultHistoryFilepath() string {
	return GetDefaultFilepath(DefaultTopicHistoryFile)
}

func Mkdir() error {
	err := os.Mkdir(GetDefaultFilepath(), os.ModeDir|PermissionBitsDir)
	if os.IsExist(err) {
		return nil
	}
	return err
}

type Topic struct {
	Start time.Time
	End   time.Time
	Data  string
}

func (t *Topic) Duration() time.Duration {
	if t.End.IsZero() {
		return time.Since(t.Start)
	}
	return t.End.Sub(t.Start)
}

func (t *Topic) Empty() bool {
	return t.Data == ""
}

func (t *Topic) String() string {
	return fmt.Sprintf("%s,%s,%s", t.Start.Format(time.RFC3339), t.Duration(), t.Data)
}
