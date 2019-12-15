package main

import (
	"context"
	"fmt"
	"io"
)

type CommandFunc func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args []string) int
type Command struct {
	Name        string
	Command     CommandFunc
	Description string
}

type CommandSet struct {
	Commands []Command
}

func (c *CommandSet) Add(cmds ...Command) {
	c.Commands = append(c.Commands, cmds...)
}

func (c CommandSet) GetCommandFunc(args []string) CommandFunc {
	// TODO: use a better data structure
	for _, cmd := range c.Commands {
		if cmd.Name == args[1] {
			return cmd.Command
		}
	}
	return nil
}

func (c CommandSet) Usage(w io.Writer) {
	fmt.Fprintln(w, "Available Commands:")
	for _, cmd := range c.Commands {
		fmt.Fprintf(w, "\t%s\n\t\t%s", cmd.Name, cmd.Description)
	}
}
