# gronval - Extract Values from gron Output

`gronval` is a command-line tool that takes `gron` (a tool that transforms JSON into discrete assignments) output from standard input and extracts only the values from these assignments.

## Features

*   Reads `gron`-style assignments from stdin.
*   Parses these assignments.
*   Outputs only the value part of each assignment.
    *   String values are unquoted.
    *   Numbers, booleans, and null are output as their literal values.
    *   Lines representing empty arrays (`[]`) or empty objects (`{}`) are skipped (produce no output).

## Installation

To install the `gronval` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/gronval@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `gronval` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `gronval` project directory:**
    ```bash
    cd path/to/your/new-hacks/gronval
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `gronval` in the current directory.

## Usage

Pipe the output of `gron` (or any `gron`-formatted text) to `gronval` via standard input.

```bash
cat gron_output.txt | gronval
# or
gron some_file.json | gronval
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

The tool outputs each extracted value on a new line.
*   String values are unquoted.
*   Numeric, boolean, and null values are output as literals.
*   Empty arrays or objects in the input produce no output line.

**Example (based on the input above):**
```
example
1.0
fast
small
2023-10-26
true
null
```

## How it Works
The tool shares much of its parsing logic with `gron2shell`. It includes a lexer (`token.go`, `ungron.go`) to parse the `gron` statement strings into a sequence of tokens. The `main.go` then iterates through these tokens. Its `formatStatement` function is specifically designed to extract and print only the value part of each assignment, unquoting strings and handling other literals appropriately.
