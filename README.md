# wslint

A file formatter that eliminates trailing whitespaces and ensures there is exactly one blank line at the end of the file.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Command Line Flags](#command-line-flags)
- [Default Exclusion Patterns](#default-exclusion-patterns)
- [Disclaimer](#disclaimer)

## Overview

wslint is designed to help keep your codebase clean and consistent by removing unnecessary whitespaces and enforcing a single blank line at the end of each file.

## Installation

    go install github.com/idelchi/wslint@latest

## Usage

    wslint [flags] [path ...]

You can specify one or more glob patterns (or simply paths) as path arguments.

## Examples

1. Lint all text files in the current directory and its subdirectories:

   `wslint **/*.txt`

2. Format all `.js` and `.css` files in the `src` directory, including hidden files:

   `wslint -w -a src/*.js src/*.css`

3. Lint all `.py` files in the `app` directory, excluding `__init__.py` files and the `tests` folder:

   `wslint -e app/**/__init__.py,app/tests/* app/**/*.py`

4. Run wslint on the `my_project` directory with four parallel jobs:

   `wslint -j 4 my_project/**/*`

## Command Line Flags

| Flag | Description                                                        |
| ---- | ------------------------------------------------------------------ |
| `-w` | Format the files. Without this flag, wslint only performs linting. |
| `-a` | Include hidden files and folders.                                  |
| `-e` | Exclude patterns, separated by commas (e.g., `*.log,*.tmp`).       |
| `-j` | Set the number of parallel jobs.                                   |
| `-h` | Print help information.                                            |
| `-v` | Print the version number.                                          |
| `-d` | Show debug output.                                                 |

## Default Exclusion Patterns

By default, wslint excludes the following patterns. These patterns represent common files or folders that are either binary, temporary, or irrelevant to code formatting:

- `**/*.exe`
- `**/.git/**`
- `**/.vscode-server/**`
- `**/node_modules/**`
- `**/vendor/**`
- `**/.task/**`
- `**/.cache/**`
- The executable itself
- Files identified as binary
- Folders and files starting with '`.`', unless the `-a` flag is used

## Disclaimer

> **Warning**
> This project serves as a learning exercise for Go and its surrounding ecosystem and tooling.
> Users are advised against using this tool on files that are not under version control.
