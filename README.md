# wslint

A file formatter, removing trailing whitespaces and enforcing exactly one blank line at the end of the file.

Usage:

    wslint [flags] [path ...]

The flags are:

    -w formats the files, running without this flag only lints
    -a include hidden files and folders
    -e exclude pattern, comma separated
    -j number of parallel jobs
    -h prints help
    -v prints version
    -d debug output

The path arguments are one or more glob patterns (or simply paths).

The following exclude patterns are used by default:

- `**/*.exe`
- `**/.git/**`
- `**/.vscode-server/**`
- `**/node_modules/**`
- `**/vendor/**`
- `**/.task/**`
- `**/.cache/**`
- the executable itself
- files identified as binary
- folders and files starting with `.`, unless the `-a` flag is used

> **Warning**
> This is a toy project to learn Go and its surrounding ecosystem and tooling.
> The user is warned against using the tool on files which are not under version control.

## TODOs

- Race detected, rerun with race flag to see the problem. Might have to run several times to see the problem.
