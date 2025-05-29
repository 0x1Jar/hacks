# geteventlisteners

`geteventlisteners` is a command-line tool that uses a headless Chrome browser (via ChromeDP) to navigate to URLs and extract JavaScript event listeners attached to the `window` object. The extracted listeners are saved as beautified JavaScript files.

## Features

*   Reads URLs from stdin or a command-line argument.
*   Uses ChromeDP to fetch pages and execute JavaScript.
*   Extracts event listeners attached to the `window` object.
*   Filters listeners by event type (e.g., `click`, `mouseover`).
*   Beautifies the extracted JavaScript code.
*   Saves listeners to `.js` files, organized into directories by domain.
*   Configurable concurrency for processing multiple URLs.
*   Configurable timeout for page loading.
*   Verbose mode for more detailed output.

## Prerequisites

*   **Go**: Version 1.16 or newer is recommended.
*   **Chrome/Chromium**: `chromedp` requires a Chrome or Chromium browser executable to be installed and accessible in your system's PATH or at a standard location.

## Installation

To install the `geteventlisteners` command-line tool, ensure you have Go and Chrome/Chromium installed.

You can install `geteventlisteners` using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/geteventlisteners@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `geteventlisteners` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `geteventlisteners` project directory:**
    ```bash
    cd path/to/your/new-hacks/geteventlisteners
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `geteventlisteners` in the current directory.

## Usage

```bash
geteventlisteners [options] [url]
```
If `[url]` is provided as an argument, it will process that single URL. Otherwise, it reads URLs from standard input (one URL per line).

### Options

*   `-f, --filter <event_type>`: Event type to filter for (e.g., `click`, `mouseover`). Can be specified multiple times to include multiple event types. If not used, all event listeners are extracted.
*   `-v`: Verbose mode (prints the current URL being processed and other info).
*   `-c <number>`: Number of concurrent browser contexts/tabs to use (default: 5).
*   `-t <seconds>`: Timeout in seconds for each URL processing (page navigation and script evaluation) (default: 20).
*   `-o <directory>`: Output directory to save listener files (default: `eventlisteners_out`). Files will be organized into subdirectories based on the domain name.

### Examples

**1. Process a single URL:**
```bash
geteventlisteners -v https://example.com
```

**2. Process URLs from a file, filtering for 'click' listeners:**
```bash
cat urls.txt | geteventlisteners -f click
```

**3. Process URLs with higher concurrency and a custom output directory:**
```bash
cat urls.txt | geteventlisteners -c 10 -o my_listeners_output
```

## Output

*   For each processed URL where matching event listeners are found, a `.js` file is created.
*   Files are saved in the structure: `[outputDir]/[domain]/[sanitized_path_and_query].js`.
    *   Example: For `https://example.com/path/page?id=1`, the file might be `eventlisteners_out/example.com/path-page-id-1.js`.
*   The script prints messages to stdout indicating which files are saved (e.g., `Saved listeners for https://example.com to eventlisteners_out/example.com/index.js`).
*   Error messages are printed to stderr.

The content of the saved `.js` files will be beautified JavaScript, with each listener assigned to a variable (e.g., `let onclick1 = function(){...};`).
