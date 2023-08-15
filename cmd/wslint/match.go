package main

import (
	"log"
	"path/filepath"

	"github.com/idelchi/wslint/internal/linter"
	"github.com/idelchi/wslint/matcher"
)

func match(options LinterOptions) (files []linter.Linter) {
	verboseLog := options.Logger
	patterns := options.Patterns
	hidden := options.Hidden
	exclude := options.Exclude

	// Create a matcher
	matcher := matcher.New(hidden, exclude, verboseLog)

	// Collect the files to inspect, ranging over the patterns
	for _, arg := range patterns {
		if err := matcher.Match(arg); err != nil {
			log.Printf("Error: %v", err)

			return files
		}
	}

	// Fill the slice with files
	for _, file := range matcher.ListFiles() {
		// Get the relative path to the execution directory
		if fileRel, err := filepath.Rel(".", file); err == nil {
			file = fileRel
		}

		lint := linter.NewLinter(file)

		// Append the linter to the slice
		files = append(files, *lint)

		verboseLog.Printf("<included> %q", file)
	}

	if len(files) == 0 {
		log.Println("No files found")
	}

	return files
}
