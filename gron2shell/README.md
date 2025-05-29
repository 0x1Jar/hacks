# gron2shell - Convert gron Output to Shell Variables

`gron2shell` is a command-line tool that takes `gron` (a tool that transforms JSON into discrete assignments) output from standard input and converts it into a format suitable for shell variable assignments.

## Features

*   Reads `gron`-style assignments from stdin.
*   Parses these assignments (which include type information).
*   Formats the path and value into shell-compatible variable assignments.
    *   Strips the leading `json` base key.
    *   Converts array indexing (e.g., `[0]`) and object key access (e.g., `.foo` or `["foo bar"]`) into underscore-separated variable names.
    *   Outputs string values without their surrounding quotes.
    *   Outputs numbers, booleans, and null as their literal values.
    *   Ignores lines representing empty arrays (`[]`) or empty objects (`{}`).

## Installation

To install the `gron2shell` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/gron2shell@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `gron2shell` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `gron2shell` project directory:**
    ```bash
    cd path/to/your/new-hacks/gron2shell
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `gron2shell` in the current directory.

## Usage

Pipe the output of `gron` (or any `gron`-formatted text) to `gron2shell` via standard input.

```bash
cat gron_output.txt | gron2shell
# or
gron some_file.json | gron2shell
```

### Input Format

The tool expects `gron`-style assignments, one per line. For example:
```
json.name = "example";
json.version = 1.0;
json.features[0] = "fast";
json.features[1] = "small";
json.metadata["release-date"] = "2023-10-26";
json.enabled = true;
json.empty_array = [];
json.empty_object = {};
json.description = null;
```

### Output Format

The tool converts each `gron` assignment into a shell variable assignment.
*   Path components are joined by underscores.
*   String values are unquoted.
*   Numeric, boolean, and null values are output as literals.
*   Assignments for empty arrays or objects are skipped.

**Example (based on the input above):**
```bash
name=example
version=1.0
features_0=fast
features_1=small
metadata_release-date=2023-10-26
enabled=true
description=null
```

You can then use these in a shell script, for example, by `eval`-ing the output (use with caution if the input is untrusted):
```bash
eval $(gron some_file.json | gron2shell)
echo $name
echo $features_0
```

## How it Works
The tool includes a lexer (`token.go`, `ungron.go`) to parse the `gron` statement strings into a sequence of tokens. The `main.go` then iterates through these tokens, stripping the initial `json` key and formatting the path by replacing dots and bracketed access with underscores, and then outputs the assignment with the unquoted or literal value. The `identifier.go` and `statements.go` files provide helpers for tokenizing and managing statements, though not all of their functionality (like full ungronning to an interface) is used by the current `main.go` for `gron2shell`.
