# html-tool

Takes URLs or filenames for HTML documents on stdin and extracts tag contents, attribute values, comments, or matches CSS selectors.

## Installation

To install the `html-tool` command-line tool, ensure you have Go installed (version 1.18 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/html-tool@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `html-tool` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**
1.  **Navigate to the `html-tool` project directory:**
    ```bash
    cd path/to/your/new-hacks/html-tool
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `html-tool` in the current directory.

## Usage

`html-tool` reads locations (URLs or file paths) from standard input.

```bash
html-tool [global_options] <mode> [<mode_args>]
```

### Global Options
*   `-c <number>`: Set the concurrency level for fetching URLs (default: 20).
*   `-s, --source`: Prepend each output line with the source URL or filename.
*   `-skip-verify`: Skip TLS certificate verification for HTTPS URLs (default: false, i.e., verification is enabled).

### Modes and Mode Arguments

*   `tags <tag-name1> [<tag-name2> ...]`
    *   Extracts the text content of the specified HTML tags.
    *   Example: `html-tool tags title h1 p`

*   `attribs <attr-name1> [<attr-name2> ...]`
    *   Extracts the values of the specified HTML attributes.
    *   Example: `html-tool attribs href src data-id`

*   `comments`
    *   Extracts all HTML comments.
    *   No mode arguments needed.
    *   Example: `html-tool comments`

*   `query <css-selector>`
    *   Extracts the text content of elements matching the given CSS selector.
    *   The CSS selector should be a single argument (quote if it contains spaces).
    *   Example: `html-tool query "div.content > p"`

### Examples

**1. Extract titles, H1s, and paragraphs from URLs in a file:**
```bash
cat urls.txt | html-tool tags title h1 p
```

**2. Extract `src` and `href` attributes from local HTML files, showing source:**
```bash
find . -type f -name "*.html" | html-tool -s attribs src href
```
Example output with `-s`:
```
./index.html /path/to/image.jpg
./index.html /another/page.html
http://example.com/page1.html http://example.com/style.css
```

**3. Extract comments from URLs, skipping TLS verification and using higher concurrency:**
```bash
cat urls.txt | html-tool -skip-verify -c 50 comments
```

**4. Extract text from elements matching a CSS selector:**
```bash
echo "http://example.com" | html-tool query "article .summary"
```

## TODO (from original README)
* Support more advanced selectors with https://github.com/ericchiang/css (partially done with `query` mode).
* Option to ignore certificate errors (done with `-skip-verify`).
* Option to output filename/URL with output (done with `-s` / `--source`).
