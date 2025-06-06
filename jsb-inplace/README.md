# jsb-inplace - In-Place JavaScript Beautifier

`jsb-inplace` is a command-line tool that reads JavaScript file paths from standard input and beautifies these files in place using the `github.com/ditashi/jsbeautifier-go` library.

## Features

*   Reads file paths from stdin (one path per line).
*   Beautifies JavaScript files using default options from `jsbeautifier-go`.
*   Overwrites the original files with their beautified versions.
*   Basic error reporting to stderr if beautification or file writing fails.

## Installation

To install the `jsb-inplace` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/jsb-inplace@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `jsb-inplace` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `jsb-inplace` project directory:**
    ```bash
    cd path/to/your/new-hacks/jsb-inplace
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `jsb-inplace` in the current directory.

## Usage

Pipe a list of JavaScript file paths (one path per line) to `jsb-inplace` via standard input.

**Example 1: Beautify a single file**
```bash
echo "/path/to/your/script.js" | jsb-inplace
```

**Example 2: Beautify multiple files listed in a file**
```bash
cat list_of_js_files.txt | jsb-inplace
```

**Example 3: Beautify all `.js` files in a directory (use with caution)**
```bash
find . -type f -name "*.js" | jsb-inplace
```
**Warning:** This command will modify all found JavaScript files in place. Ensure you have backups or are using version control.

### Input Format
The tool expects one valid file path per line from stdin, pointing to a JavaScript file.

### Output
The tool does not produce output to stdout upon successful beautification. It modifies the files directly.
Error messages are printed to stderr if:
*   Beautification fails for a file (e.g., file is empty, invalid, or `jsbeautifier` encounters an issue).
*   Writing the beautified content back to the file fails.
*   Reading from stdin fails.

## How it Works
The tool reads each line from stdin, interpreting it as a file path. For each path, it calls `jsbeautifier.BeautifyFile` with default options. If successful, it overwrites the original file with the beautified content using `os.WriteFile`.
