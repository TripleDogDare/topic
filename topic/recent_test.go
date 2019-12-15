package topic

import (
	"context"
	"strings"
	"testing"
)

// test that the last line is retrieved
// extra newlines should be okay
func TestLastLine(t *testing.T) {
	ctx := context.Background()
	r := strings.NewReader("hello\n\nworld")
	result := LastLine(ctx, r)
	if result != "world" {
		t.Fail()
	}
}

// Test that the last non-empty line is used
// space characters should not count
func TestLastLineNewline(t *testing.T) {
	ctx := context.Background()
	r := strings.NewReader("hello\nworld\n\t")
	result := LastLine(ctx, r)
	if result != "world" {
		t.Fail()
	}
}
