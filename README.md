# Topic - Simple Time Management

A simple command line tool to keep track of time. Topic uses simple timestamps and free form text. Allowing you to manage tasks with minimal intrusion while you work.

Set the topic by executing `topic new [task]`. For example:
```
topic new Some meeting that should have been an email
```

## Why?

If you're the kind of person who wants text based time management software then this might be for you.

## Use cases
1. Put the current task and duration someplace visible such as command prompt or status bar.
2. Generate a report of starting timestamps and durations

## Examples
1. [i3](https://i3wm.org/)/[i3blocks](https://github.com/vivien/i3blocks): [example](doc/i3blocks-example.sh)

## Build
Clone the repository and execute `./build.sh` will create a binary at `./bin/topic`.
