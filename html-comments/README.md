# html-comments - Extract HTML Comments

`html-comments` is a simple command-line tool that reads HTML content from standard input and extracts all HTML comments.

## Features

*   Reads HTML from stdin.
*   Uses Go's standard `golang.org/x/net/html` tokenizer for parsing.
*   Prints each non-empty HTML comment found, with newlines within comments replaced by spaces.

## Installation

To install the `html-comments` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/html-comments@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `html-comments` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `html-comments` project directory:**
    ```bash
    cd path/to/your/new-hacks/html-comments
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `html-comments` in the current directory.

## Usage

Pipe HTML content to `html-comments` via standard input.

**Example 1: Using a local HTML file**
```bash
cat page.html | html-comments
```

**Example 2: Fetching a URL with `curl` and piping its content**
```bash
curl -s https://example.com | html-comments
```

**Example 3: Piping `echo`'d HTML**
```bash
echo "<html><body><!-- This is a comment --><h1>Title</h1><!-- Another comment --></body></html>" | html-comments
```
Output for Example 3:
```
This is a comment
Another comment
```

### Output Format
The tool prints each extracted HTML comment on a new line. Newlines within a single comment are replaced with spaces, and leading/trailing whitespace from the comment data is trimmed. Empty comments are ignored.

## How it Works
The tool uses the `html.NewTokenizer` from the `golang.org/x/net/html` package to parse the HTML input stream. It iterates through the tokens, and when it encounters a `html.CommentToken`, it extracts the comment data, cleans it up (replaces newlines, trims whitespace), and prints it if it's not empty.
