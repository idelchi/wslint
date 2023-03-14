package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/idelchi/wslint/internal/linter"
	"github.com/idelchi/wslint/matcher"
)

//nolint:gochecknoglobals // Global variable for CI stamping.
var versionStamp = "unknown - unofficial & generated by unknown"

func usage() {
	log.Println("wslint checks or fixes files with trailing whitespaces and enforces final newlines")
	log.Println("Usage: wslint [flags] [path ...]")
	flag.PrintDefaults()
}

// Helper to print a message, the usage and exit with error code 1.
func printAndExit(msg string) {
	log.Println(msg)
	flag.Usage()
	os.Exit(1)
}

// Helper to print the version and exit with error code 0.
func displayVersion() {
	log.Println(versionStamp)
	os.Exit(0)
}

// Options contains the options for the linter.
type Options struct {
	Exclude         []string
	NumberOfWorkers int
	Fix             bool
	Logger          *log.Logger
	Patterns        []string
	Hidden          bool
}

func main() {
	var (
		fix      = flag.Bool("w", false, "format file in-place")
		verbose  = flag.Bool("d", false, "debug output")
		exclude  = flag.String("e", "", "exclude pattern, comma separated")
		hidden   = flag.Bool("a", false, "show hidden files & folders")
		parallel = flag.Int("j", runtime.NumCPU(), "number of parallel jobs, defaults to number of CPUs")
		version  = flag.Bool("v", false, "print version")
	)

	// No time stamp in the log output
	log.SetFlags(0)

	// Set the usage message & parse the flags
	flag.Usage = usage
	flag.Parse()

	// If the -v flag is set, print the version and exit
	if *version {
		displayVersion()
	}

	// If no arguments are given, give an error message
	if flag.NArg() == 0 {
		printAndExit("Error: Need to provide at least one path element")
	}

	// If the number of parallel jobs is less than 1, give an error message
	if *parallel <= 0 {
		printAndExit("Error: Number of parallel jobs must be greater than 0")
	}

	// Create a logger for debug messages
	verboseLog := log.New(os.Stdout, "", 0)
	if !*verbose {
		// If the -d flag is not set, disable debug messages
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

	// Create the options
	options := Options{
		Exclude:         excludes,
		NumberOfWorkers: *parallel,
		Fix:             *fix,
		Logger:          verboseLog,
		Patterns:        flag.Args(),
		Hidden:          *hidden,
	}

	os.Exit(match(options))
}

func match(options Options) int {
	verboseLog := options.Logger
	patterns := options.Patterns
	hidden := options.Hidden
	exclude := options.Exclude

	// Create a matcher
	matcher := matcher.New(hidden, exclude, verboseLog)

	// Collect the files to inspect, ranging over the patterns
	for _, arg := range patterns {
		if err := matcher.Match(arg); err != nil {
			log.Fatalln(err)
		}
	}

	// Create a slice of files to inspect
	files := []linter.Linter{}

	// Fill the slice with files
	for _, file := range matcher.ListFiles() {
		// Get the relative path to the execution directory
		if fileRel, err := filepath.Rel(".", file); err == nil {
			file = fileRel
		}

		files = append(files, linter.New(file))

		verboseLog.Printf("<included> %q", file)
	}

	// If no files are found, exit with error code 0
	if len(files) == 0 {
		log.Println("No files found")

		return 1
	}

	return process(options, files)
}

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

	if exitCode == 0 {
		log.Println("No issues found")
	}

	return exitCode
}

// WorkerPool represents a pool of workers.
type WorkerPool struct {
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
}

// Start the worker pool.
func (p *WorkerPool) Start(jobs, results chan linter.Linter) {
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
	p.Logger.Printf("<processed> all files in %s", time.Since(start))
}

// worker processes jobs.
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

		if fix {
			file.Error = file.Fix()
		} else {
			file.Error = file.Lint()
		}

		results <- file

		jobsProcessed++
	}

	logger.Printf("<worker %d> processed %d jobs", identifier, jobsProcessed)
}
