package wslint

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

// exit prints the message and exits with the specified exit code.
func (w *Wslint) exit(code int, msg string) {
	log.Println(msg)

	if code != 0 {
		w.Usage()
	}

	//nolint:forbidigo // This is the only place where os.Exit() is used.
	os.Exit(code)
}

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
	Experimental    bool
	Interactive     bool
}

// Parse collects the commandline arguments and returns them as a CLIOptions struct.
//
//nolint:funlen // This function is long, but has one dedicated function.
func (w *Wslint) Parse() {
	// Flags for the CLI
	var (
		fix          = flag.Bool("w", false, "format file in-place")
		verbose      = flag.Bool("d", false, "debug output")
		exclude      = flag.String("e", "", "exclude pattern, comma separated")
		hidden       = flag.Bool("a", false, "show hidden files & folders")
		parallel     = flag.Int("j", runtime.NumCPU(), "number of parallel jobs, defaults to number of CPUs")
		version      = flag.Bool("v", false, "print version")
		quiet        = flag.Bool("q", false, "suppress messages")
		experimental = flag.Bool("x", false, "enable experimental features")
		interactive  = flag.Bool("i", false, "interactive mode")
	)

	// No time stamp in the log output
	log.SetFlags(0)
	// Set the usage message & parse the flags
	flag.Usage = w.Usage
	flag.Parse()

	switch {
	// If the version flag is set, print the version and exit
	case *version:
		if *verbose {
			if info, available := debug.ReadBuildInfo(); available {
				w.Version += fmt.Sprintf("\nruntime version information: %v", info.Main.Version)
			}
		}

		w.exit(0, w.Version)
	// If no arguments are given, raise an error message
	case flag.NArg() == 0:
		w.exit(1, "Error: Need to provide at least one path element")
	// If the number of parallel jobs is less than 1, raise an error message
	case *parallel <= 0:
		w.exit(1, "Error: Number of parallel jobs must be greater than 0")
	// Interactive is not implemented yet
	case *interactive:
		w.exit(1, "Error: Interactive mode is not implemented yet")
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

	w.Options = Options{
		Exclude:         excludes,
		NumberOfWorkers: *parallel,
		Fix:             *fix,
		Logger:          verboseLog,
		Patterns:        flag.Args(),
		Hidden:          *hidden,
		Quiet:           *quiet,
		Verbose:         *verbose,
		Experimental:    *experimental,
		Interactive:     *interactive,
	}
}
