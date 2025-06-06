# get-title

`get-title` is a command-line tool that fetches HTML pages from URLs provided on standard input and extracts their `<title>` tags.

## Features

*   Reads URLs from stdin.
*   Concurrently fetches and processes URLs.
*   Extracts and prints the content of the `<title>` tag for each page.
*   Configurable concurrency level.
*   Option to skip TLS certificate verification.

## Installation

To install the `get-title` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/get-title@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `get-title` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `get-title` project directory:**
    ```bash
    cd path/to/your/new-hacks/get-title
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `get-title` in the current directory.

## Usage

Pipe a list of URLs (one URL per line) to `get-title` via standard input.

```bash
cat list_of_urls.txt | get-title [options]
```

Or using `echo`:
```bash
echo -e "https://example.com\nhttps://google.com" | get-title -c 10 -k
```

### Options

*   `-c <number>`: Set the concurrency level for fetching URLs (default: 20).
*   `-k`, `--skip-verify`: Skip TLS certificate verification (insecure). By default, TLS certificates are verified.

### Output Format

For each URL that is successfully fetched and has a title, the tool prints a line in the format:
`Page Title (URL)`

**Example Output:**
```
Example Domain (https://example.com)
Google (https://google.com)
```
Error messages for failed requests are not explicitly printed by `extractTitle` but `gahttp` might log them or they might be silent.

## How it Works

The tool reads URLs from stdin and uses the `gahttp` library to concurrently fetch them. For each successful response, it parses the HTML using `golang.org/x/net/html` to find the first `<title>` tag and prints its text content along with the original URL. The concurrency of HTTP requests and TLS verification behavior can be controlled via command-line flags.
