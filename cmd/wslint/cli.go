package main

import (
	"flag"
	"log"
	"runtime"
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

// CLIOptions contains the configuration parsed from the CLI.
type CLIOptions struct {
	fix      bool
	verbose  bool
	exclude  string
	hidden   bool
	parallel int
	version  bool
	quiet    bool
	patterns []string
}

func parse() CLIOptions {
	// Flags for the CLI
	var (
		fix      = flag.Bool("w", false, "format file in-place")
		verbose  = flag.Bool("d", false, "debug output")
		exclude  = flag.String("e", "", "exclude pattern, comma separated")
		hidden   = flag.Bool("a", false, "show hidden files & folders")
		parallel = flag.Int("j", runtime.NumCPU(), "number of parallel jobs, defaults to number of CPUs")
		version  = flag.Bool("v", false, "print version")
		quiet    = flag.Bool("q", false, "suppress messages")
	)

	// Set the usage message & parse the flags
	flag.Usage = usage
	flag.Parse()

	switch {
	// If the version flag is set, print the version and exit
	case *version:
		exit(0, versionStamp)
	// If no arguments are given, raise an error message
	case flag.NArg() == 0:
		exit(1, "Error: Need to provide at least one path element")
	// If the number of parallel jobs is less than 1, raise an error message
	case *parallel <= 0:
		exit(1, "Error: Number of parallel jobs must be greater than 0")
	}

	return CLIOptions{
		fix:      *fix,
		verbose:  *verbose,
		exclude:  *exclude,
		hidden:   *hidden,
		parallel: *parallel,
		version:  *version,
		quiet:    *quiet,
		patterns: flag.Args(),
	}
}
