/*
Wslint formats files, removing trailing whitespaces and enforcing exactly one blank line at the end of the file.

Usage:

	wslint [flags] [path ...]

The flags are:

	-w		formats the files in-place. Running without this flag only lints
	-a		include hidden files and folders (starting with a dot) in the search
	-e		exclude pattern, comma separated
	-j		number of parallel jobs (defaults to number of CPUs)
	-d		debug output
	-h		prints help
	-v		prints version

The path arguments are one or more glob patterns (or simply paths to files).
*/
package main
