// Package wslint provides the main functionality of the linter.
package wslint

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/wslint/internal/checkers"
	"github.com/idelchi/wslint/internal/linter"
	"github.com/idelchi/wslint/internal/worker"
	"github.com/idelchi/wslint/pkg/matcher"
)

// Wslint acts as a wrapper for the main functionality.
type Wslint struct {
	Options Options
	Files   []linter.Linter
	Usage   func()
	Version string
}

// Match stores the files that match the patterns.
func (w *Wslint) Match() {
	verboseLog := w.Options.Logger
	patterns := w.Options.Patterns
	hidden := w.Options.Hidden
	exclude := w.Options.Exclude

	// Create a matcher
	matcher := matcher.New(hidden, exclude, verboseLog)

	// Collect the files to inspect, ranging over the patterns
	for _, arg := range patterns {
		if err := matcher.Match(arg); err != nil {
			log.Printf("Error: %v", err)

			return
		}
	}

	// Fill the slice with files
	for _, file := range matcher.ListFiles() {
		// Get the relative path to the execution directory
		if fileRel, err := filepath.Rel(".", file); err == nil {
			file = fileRel
		}

		// TODO(Idelchi) Set up a factory function for this
		lint := linter.New(file)

		if w.Options.Experimental {
			stutter := checkers.Stutter{}
			// Load file in config/stutters and read into a slice of strings
			// Pass the slice to the checker
			// TODO(Idelchi): The configuration file is hardcoded
			configurationFile := "settings/stutters"
			if _, err := os.Stat(configurationFile); !os.IsNotExist(err) {
				if content, err := os.ReadFile(configurationFile); err == nil {
					stutter.Exceptions = strings.Split(string(content), "\n")
				}
			}

			lint.InsertChecker("stutter", stutter)
		}

		// Append the linter to the slice
		w.Files = append(w.Files, *lint)

		verboseLog.Printf("<included> %q", file)
	}

	if len(w.Files) == 0 {
		log.Println("No files found")
	}
}

// Process processes the files, prints out the results and returns the exit code.
func (w *Wslint) Process() int {
	numberOfFiles := len(w.Files)
	w.Options.Logger.Printf("Processing %d files", numberOfFiles)

	w.Options.NumberOfWorkers = min(w.Options.NumberOfWorkers, numberOfFiles)

	workerPool := worker.Pool{
		NumberOfWorkers: w.Options.NumberOfWorkers,
		NumberOfJobs:    numberOfFiles,
		Fix:             w.Options.Fix,
		Files:           w.Files,
		Logger:          w.Options.Logger,
	}

	// Create channels for sending and receiving jobs and results
	jobs := make(chan linter.Linter, numberOfFiles)
	results := make(chan linter.Linter, numberOfFiles)

	workerPool.Start(jobs, results)

	exitCode := 0

	// Collect the results
	for range w.Files {
		result := <-results

		if ok := result.Summary(); !ok {
			exitCode = 1
		}
	}

	workerPool.Stats()

	if exitCode == 0 {
		log.Println("No issues found")
	}

	return exitCode
}
