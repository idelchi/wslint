package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

// Options contains the options for the linter.
type Options struct {
	// Exclude patterns, separated by commas (e.g., *.log,*.tmp).
	Exclude []string
	// Number of parallel jobs.
	NumberOfWorkers int
	Fix             bool
	Logger          *log.Logger
	Patterns        []string
	Hidden          bool
	Quiet           bool
	Verbose         bool
	Exp             bool
}

// Parse collects the commandline arguments and returns them as a CLIOptions struct.
func (w *Wslint) Parse() {
	// Flags for the CLI
	var (
		fix      = flag.Bool("w", false, "format file in-place")
		verbose  = flag.Bool("d", false, "debug output")
		exclude  = flag.String("e", "", "exclude pattern, comma separated")
		hidden   = flag.Bool("a", false, "show hidden files & folders")
		parallel = flag.Int("j", runtime.NumCPU(), "number of parallel jobs, defaults to number of CPUs")
		version  = flag.Bool("v", false, "print version")
		quiet    = flag.Bool("q", false, "suppress messages")
		exp      = flag.Bool("exp", false, "enable experimental features")
	)

	// No time stamp in the log output
	log.SetFlags(0)
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

	// Create a logger for debug messages
	verboseLog := log.New(os.Stdout, "", 0)
	if !*verbose {
		// Disable debug messages if the verbose flag is not set,
		verboseLog.SetOutput(io.Discard)
	}

	// Disable the logger if the quiet flag is set
	if *quiet {
		log.SetOutput(io.Discard)
		verboseLog.SetOutput(io.Discard)
	}

	// Split the exclude patterns into a slice
	excludes := strings.Split(*exclude, ",")

	for i, exclude := range excludes {
		// Remove any leading and trailing whitespace
		exclude = strings.TrimSpace(exclude)
		// Remove "./" from the beginning of the pattern, if it exists
		exclude = strings.TrimPrefix(exclude, "./")
		excludes[i] = exclude
	}

	w.options = Options{
		Exclude:         excludes,
		NumberOfWorkers: *parallel,
		Fix:             *fix,
		Logger:          verboseLog,
		Patterns:        flag.Args(),
		Hidden:          *hidden,
		Quiet:           *quiet,
		Verbose:         *verbose,
		Exp:             *exp,
	}
}
