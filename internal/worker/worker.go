// Package worker provides a concurrent mechanism to process a set of jobs using a pool of workers.
// The primary component of this package is the Pool, which manages a set of goroutines (workers)
// to process jobs in parallel. Each job represents a file that needs linting.
//
// A typical use case involves:
// 1. Initializing a Pool with a specified number of workers.
// 2. Sending files (jobs) to the Pool for processing.
// 3. Starting the Pool, which dispatches the jobs to the workers.
// 4. Each worker processes its assigned jobs, performing linting and optionally fixing issues.
// 5. Once all jobs are processed, the Pool can provide statistics about the processing duration.
//
// The worker package ensures efficient and safe concurrent processing of jobs,
// allowing for faster linting of large sets of files.
package worker

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/atomic"

	"github.com/idelchi/wslint/internal/linter"
	"golang.org/x/exp/slices"
)

// Pool represents a pool of workers.
type Pool struct {
	// The number of workers in the pool
	NumberOfWorkers int
	// The number of jobs to process
	NumberOfJobs int
	// Logger
	Logger *log.Logger
	// Fix
	Fix bool
	// Files
	Files []linter.Linter
	// Time spent processing the files
	ProcessingTime time.Duration
}

// Start the worker pool.
func (p *Pool) Start(jobs, results chan linter.Linter) {
	// Create a wait group to ensure all workers have finished
	var waitGroup sync.WaitGroup

	// Start the workers
	for i := 0; i < p.NumberOfWorkers; i++ {
		waitGroup.Add(1)

		i := i

		go func() {
			defer waitGroup.Done()
			worker(i+1, p.Logger, p.Fix, jobs, results)
		}()
	}

	// Measure the time it takes to process all the files
	start := time.Now()

	// Send the jobs to the workers
	for _, file := range p.Files {
		jobs <- file
	}

	close(jobs)

	// Wait for all the workers to finish
	waitGroup.Wait()

	// Measure the time it takes to process all the files
	p.ProcessingTime = time.Since(start)
}

// Stats prints the stats of the worker pool run.
func (p *Pool) Stats() {
	p.Logger.Printf("<processed> %d files in %s", len(p.Files), p.ProcessingTime)
}

// worker processes jobs.
// https://twin.sh/articles/39/go-concurrency-goroutines-worker-pools-and-throttling-made-simple
func worker(
	identifier int,
	logger *log.Logger,
	fix bool,
	files <-chan linter.Linter,
	results chan<- linter.Linter,
) {
	jobsProcessed := 0

	for file := range files {
		logger.Printf("<processing> %q", file.Name)

		func() {
			// TODO(Idelchi): Work with []byte instead of string
			content, err := os.ReadFile(file.Name)
			if err != nil {
				panic("failed to read file")
			}

			src := strings.Split(string(content), "\n")
			res := make([]string, len(src))
			copy(res, src)

			res = file.Format(res)

			if !slices.Equal(src, res) && fix {
				info, _ := os.Lstat(file.Name)
				if err := atomic.WriteFile(file.Name, strings.NewReader(strings.Join(res, "\n"))); err != nil {
					panic("failed to write file")
				}

				if err = os.Chmod(file.Name, info.Mode()); err != nil {
					panic("failed to change file permissions")
				}
			}
		}()

		results <- file

		jobsProcessed++
	}

	logger.Printf("<worker %d> processed %d jobs", identifier, jobsProcessed)
}
