# wslint

A file formatter that eliminates trailing whitespaces and ensures there is exactly one blank line at the end of the file.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Command Line Flags](#command-line-flags)
- [Default Exclusion Patterns](#default-exclusion-patterns)
- [Disclaimer](#disclaimer)

## Overview

wslint is designed to help keep your codebase clean and consistent by removing unnecessary whitespaces
and enforcing a single blank line at the end of each file.

## Installation

### From source

```sh
go install github.com/idelchi/wslint/cmd/wslint@latest
```

### From installation script

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/wslint/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
wslint [flags] [path ...]
```

Paths can be specified as one or more glob patterns or simple file paths.

Enclose path arguments in quotes to prevent shell expansion.

## Examples

Lint all files in the current directory and its subdirectories:

```sh
    wslint "**"
```

Lint all text files in the current directory and its subdirectories:

```sh
    wslint "**/*.txt"
```

Format all `.js` and `.css` files in the `src` directory, including hidden files:

```sh
    wslint -w -a "src/*.js" "src/*.css"
```

Lint all `.py` files in the `app` directory, excluding `__init__.py` files and the `tests` folder:

```sh
    wslint -e "app/**/__init__.py,app/tests/**" "app/**/*.py"
```

Run wslint on the `my_project` directory with four parallel jobs:

```sh
    wslint -j 4 "my_project/**"
```

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
| `-q` | Suppress messages.                                                 |
| `-x` | Enable experimental features.                                      |

## Default Exclusion Patterns

By default, wslint excludes the following patterns. These patterns represent common files or folders that
are either binary, temporary, or irrelevant to code formatting:

- `**/*.exe`
- `**/.git/**`
- `**/.vscode-server/**`
- `**/node_modules/**`
- `**/vendor/**`
- `**/.task/**`
- `**/.cache/**`
- The executable itself
- Folders and files starting with '`.`', unless the `-a` flag is used

To include files that are normally excluded, either:

- Write a full path to the file
- Use glob patterns only in the filename portion of the path (excluding the extension)

Example:

```sh
wslint ".git/config"
wslint ".task/*.env"
wslint "*.exe"
```

## Disclaimer

> **Warning**
> This project serves as a learning exercise for Go and its surrounding ecosystem and tooling.
> Users are advised against using this tool on files that are not under version control.

## Experimental Features

Experimental features are disabled by default. To enable them, use the `-x` flag.

Beaware that applying formatting to files that are not under version control may result in data loss,
as experimental features are not yet fully tested.

The current experimental features are:

<details>
  <summary><strong>stutters</strong></summary>

- **Description**: Remove stuttering words (e.g., `the the` -> `the`).

- **Limitations**: If exceptions are desired, they must be placed in a file [config/stutters](./config/stutters)
  relative to the current working directory.

- **Issues**: Will not respect case, punctuation, as it will always select the second occurrence when fixing.

</details>
