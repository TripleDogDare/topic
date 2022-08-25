package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/tripledogdare/go-topic-timer/topic"
	// "github.com/pkg/errors"
)

var MainVar string

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
		Command{"version", Version, "print version info"},
		Command{"buildinfo", BuildInfo, "print build info"},
	)
	if len(args) >= 2 {
		if f := commands.GetCommandFunc(args[1]); f != nil {
			return f(ctx, stdin, stdout, stderr, args[2:])
		}
	}
	commands.Usage(stderr)
	return 1
}

func BuildInfo(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(stdout, "Build info not available")
	} else {
		fmt.Fprintf(stdout, "Build path: %s\n", info.Path)
		fmt.Fprintf(stdout, "Build Module: %s\n", info.Main.Path)
		fmt.Fprintf(stdout, "Build version: %s\n", info.Main.Version)
	}
	return 0
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
	printTimestamp := flagset.Bool("timestamp", false, "print timestamp")
	printDuration := flagset.Bool("duration", false, "print duration")
	printTopic := flagset.Bool("topic", false, "print topic")
	truncate := flagset.Duration("truncate", 0, "truncate duration")

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

	printAll := !(*printTimestamp || *printDuration || *printTopic)

	var sb strings.Builder
	if *printTimestamp || printAll {
		sb.WriteString(result.Start.Format(time.RFC3339))
	}
	if *printDuration || printAll {
		if sb.Len() != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(result.Duration().Truncate(*truncate).String())
	}
	if *printTopic || printAll {
		if sb.Len() != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(result.Data)
	}
	fmt.Println(sb.String())
	return 0
}

func Version(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprintf(stdout, "Commit: %s\n", GitCommit)
	fmt.Fprintf(stdout, "GitDescribe: %s\n", GitDescribe)
	fmt.Fprintf(stdout, "GitTag: %s\n", GitTag)
	goVersion := GoVersion
	if strings.TrimSpace(goVersion) == "" {
		goVersion = runtime.Version()
	}
	fmt.Fprintf(stdout, "GoVersion: %s\n", goVersion)
	fmt.Fprintf(stdout, "CommitsSinceTag: %s\n", CommitsSinceTag)
	return 0
}
