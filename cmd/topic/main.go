package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tripledogdare/go-topic-timer/topic"
	// "github.com/pkg/errors"
)

func main() {
	ctx := context.Background()
	os.Exit(Main(ctx, os.Stdin, os.Stdout, os.Stderr, os.Args))
}

// Accepts data streams to use as std in/out/err channels
// returns exit code
func Main(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	commands := new(CommandSet)
	commands.Add(
		Command{"report", Report, "print report"},
		Command{"prompt", Prompt, "prompt user to set a topic"},
		Command{"latest", Latest, "report most recent topic"},
		Command{"new", NewTopic, "sets a new topic"},
	)
	if f := commands.GetCommandFunc(args); f != nil {
		return f(ctx, stdin, stdout, stderr, args[2:])
	}
	commands.Usage(stderr)
	return 1
}

func Usage(flagset *flag.FlagSet) {
	fmt.Fprintf(flagset.Output(), "Usage of %s:\n", flagset.Arg(0))
	flagset.PrintDefaults()
}

func Report(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprintln(stderr, "running report")
	if err := topic.Report(ctx, stdout); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func NewTopic(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	data := strings.Join(args, " ")
	fmt.Fprintf(stderr, "adding topic \"%s\"\n", data)
	if err := topic.Append(data); err != nil {
		fmt.Fprintln(stderr, err)
	}
	return 0
}

func Prompt(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	flags := flag.NewFlagSet("prompt", flag.ContinueOnError)
	duration := flags.Duration("duration", 30*time.Minute, "How long to check")
	flags.Parse(args)

	fmt.Println(duration)
	panic("not implemented")
	return 0
}

func Latest(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	// initialize flagset
	flagset := flag.NewFlagSet("prompt", flag.ContinueOnError)
	flagset.SetOutput(stderr)
	// set flags
	timestamp := flagset.Bool("timestamp", false, "print timestamp")
	duration := flagset.Bool("duration", false, "print duration")
	printTopic := flagset.Bool("topic", false, "print topic")
	// handle flag errors
	if err := flagset.Parse(args); err != nil {
		fmt.Fprintf(stderr, "Invalid args")
		flagset.PrintDefaults()
		return 1
	}

	result := new(topic.Topic)
	err := topic.Recent(ctx, result)
	if err != nil {
		fmt.Fprintf(stderr, "Error\n%s", err.Error())
		return 1
	}
	if *duration {
		fmt.Println(result.Duration())
	} else if *timestamp {
		fmt.Println(result.Start.Format(time.RFC3339))
	} else if *printTopic {
		fmt.Println(result.Data)
	} else {
		fmt.Fprintln(stdout, result)
	}
	return 0
}
