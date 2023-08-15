package main

import (
	"log"
	"path/filepath"

	"github.com/idelchi/wslint/internal/checkers"
	"github.com/idelchi/wslint/internal/linter"
	"github.com/idelchi/wslint/internal/worker"
	"github.com/idelchi/wslint/matcher"
)

// Wslint acts as a wrapper for the main functionality.
type Wslint struct {
	options Options
	files   []linter.Linter
}

// Match stores the files that match the patterns.
func (w *Wslint) Match() {
	verboseLog := w.options.Logger
	patterns := w.options.Patterns
	hidden := w.options.Hidden
	exclude := w.options.Exclude

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
		if w.options.Exp {
			lint.InsertChecker(&checkers.Stutters{})
		}

		// Append the linter to the slice
		w.files = append(w.files, *lint)

		verboseLog.Printf("<included> %q", file)
	}

	if len(w.files) == 0 {
		log.Println("No files found")
	}
}

// Process processes the files, prints out the results and returns the exit code.
func (w *Wslint) Process() int {
	numberOfFiles := len(w.files)
	w.options.Logger.Printf("Processing %d files", numberOfFiles)

	w.options.NumberOfWorkers = min(w.options.NumberOfWorkers, numberOfFiles)

	workerPool := worker.Pool{
		NumberOfWorkers: w.options.NumberOfWorkers,
		NumberOfJobs:    numberOfFiles,
		Fix:             w.options.Fix,
		Files:           w.files,
		Logger:          w.options.Logger,
	}

	// Create channels for sending and receiving jobs and results
	jobs := make(chan linter.Linter, numberOfFiles)
	results := make(chan linter.Linter, numberOfFiles)

	workerPool.Start(jobs, results)

	exitCode := 0

	// Collect the results
	for range w.files {
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
