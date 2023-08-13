package main

import (
	"log"
)

// LinterOptions contains the options for the linter.
type LinterOptions struct {
	Exclude         []string
	NumberOfWorkers int
	Fix             bool
	Logger          *log.Logger
	Patterns        []string
	Hidden          bool
	Quiet           bool
	Verbose         bool
}

// Options contains the options for the linter.
type Options struct {
	Exclude         []string
	NumberOfWorkers int
	Fix             bool
	Logger          *log.Logger
	Patterns        []string
	Hidden          bool
	Quiet           bool
	Verbose         bool
}

// CLIOptions contains the options for the CLI.
type CLIOptions struct{}
