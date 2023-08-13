package main

import (
	"log"

	"github.com/idelchi/wslint/internal/linter"
)

func process(options Options, files []linter.Linter) int {
	numberOfFiles := len(files)
	options.Logger.Printf("Processing %d files", numberOfFiles)

	if numberOfFiles < options.NumberOfWorkers {
		options.NumberOfWorkers = numberOfFiles
	}

	workerPool := WorkerPool{
		NumberOfWorkers: options.NumberOfWorkers,
		NumberOfJobs:    numberOfFiles,
		Fix:             options.Fix,
		Files:           files,
		Logger:          options.Logger,
	}

	// Create channels for sending and receiving jobs and results
	jobs := make(chan linter.Linter, numberOfFiles)
	results := make(chan linter.Linter, numberOfFiles)

	workerPool.Start(jobs, results)

	exitCode := 0

	// Collect the results
	for range files {
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
