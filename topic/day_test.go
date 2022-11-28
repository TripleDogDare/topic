package topic

import (
	"testing"
	"time"
)

func TestIsOnSameDay(t *testing.T) {
	for _, tt := range []struct {
		name   string
		a      time.Time
		b      time.Time
		result bool
	}{
		{name: "now", a: time.Now(), b: time.Now(), result: true},
		{name: "now+25h", a: time.Now(), b: time.Now().Add(25 * time.Hour), result: false},
		{name: "zero-utc", a: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), b: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), result: true},
		{name: "zero-one-utc", a: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), b: time.Date(0, 0, 1, 0, 0, 0, 0, time.UTC), result: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := isSameDate(tt.a, tt.b)
			if r != tt.result {
				t.Logf("%s == %s", tt.a, tt.b)
				t.Fatalf("expected %v but got %v", tt.result, r)
			}
		})
	}
}

func TestIsTopicOnSameDay(t *testing.T) {
	for _, tt := range []struct {
		name   string
		topic  *Topic
		time   time.Time
		result bool
	}{
		{name: "nil", topic: nil, time: time.Now(), result: false},
		{name: "start-now", topic: &Topic{Start: time.Now()}, time: time.Now(), result: true},
		{name: "end-now", topic: &Topic{End: time.Now()}, time: time.Now(), result: true},
		{name: "start-now+25h", topic: &Topic{Start: time.Now()}, time: time.Now().Add(25 * time.Hour), result: false},
		{name: "end-now+25h", topic: &Topic{End: time.Now()}, time: time.Now().Add(25 * time.Hour), result: false},
		{name: "start-zero-utc", topic: &Topic{Start: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)}, time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), result: true},
		{name: "end-zero-utc", topic: &Topic{End: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)}, time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), result: true},
		{name: "start-zero-one-utc", topic: &Topic{Start: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)}, time: time.Date(0, 0, 1, 0, 0, 0, 0, time.UTC), result: false},
		{name: "end-zero-one-utc", topic: &Topic{End: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)}, time: time.Date(0, 0, 1, 0, 0, 0, 0, time.UTC), result: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := isTopicOnDay(tt.topic, tt.time)
			if r != tt.result {
				t.Logf("is %v on %s?", tt.topic, tt.time)
				t.Fatalf("expected %v but got %v", tt.result, r)
			}
		})
	}
}
