package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

//nolint:gochecknoglobals // Global variable for CI stamping.
var versionStamp = "unknown - unofficial & generated by unknown"

func usage() {
	log.Println("wslint checks or fixes files with trailing whitespaces and enforces final newlines")
	log.Println("Usage: wslint [flags] [path ...]")
	flag.PrintDefaults()
}

// exit prints the message and exits with the specified exit code.
func exit(code int, msg string) {
	log.Println(msg)

	if code != 0 {
		usage()
	}

	os.Exit(code)
}

func main() {
	// No time stamp in the log output
	log.SetFlags(0)

	cli := parse()

	// Create a logger for debug messages
	verboseLog := log.New(os.Stdout, "", 0)
	if !cli.verbose {
		// Disable debug messages if the verbose flag is not set,
		verboseLog.SetOutput(io.Discard)
	}

	// Disable the logger if the quiet flag is set
	if cli.quiet {
		log.SetOutput(io.Discard)
		verboseLog.SetOutput(io.Discard)
	}

	// Split the exclude patterns into a slice
	excludes := strings.Split(cli.exclude, ",")

	for i, exclude := range excludes {
		// Remove any leading and trailing whitespace
		exclude = strings.TrimSpace(exclude)
		// Remove "./" from the beginning of the pattern, if it exists
		exclude = strings.TrimPrefix(exclude, "./")
		excludes[i] = exclude
	}

	// Create the options
	options := LinterOptions{
		Exclude:         excludes,
		NumberOfWorkers: cli.parallel,
		Fix:             cli.fix,
		Logger:          verboseLog,
		Patterns:        cli.patterns,
		Hidden:          cli.hidden,
	}

	os.Exit(match(options))
}
