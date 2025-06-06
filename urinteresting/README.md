# urinteresting

Accept URLs on stdin, output the ones that look 'interesting' based on a set of predefined heuristics.

## Installation

To install the `urinteresting` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/urinteresting@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `urinteresting` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `urinteresting` project directory:**
    ```bash
    cd path/to/your/new-hacks/urinteresting
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `urinteresting` in the current directory.

## Usage

Pipe a list of URLs (one URL per line) to `urinteresting` via standard input. The tool will print only the URLs that match its internal "interesting" criteria.

```bash
cat list_of_urls.txt | urinteresting
```

Or using `echo`:
```bash
echo -e "http://example.com/boring.html\nhttp://example.com/admin.php?debug=true" | urinteresting
```
Example output (if the second URL is deemed interesting):
```
http://example.com/admin.php?debug=true
```

## "Interesting" Criteria

The tool considers a URL interesting if it matches one or more of the following conditions:

*   **Query String Parameters:**
    *   Contains parameters with keys like `redirect`, `debug`, `password`, `file`, `template`, `include`, `url`, `src`, `href`, `callback`, etc.
    *   Contains parameter values that look like URLs (start with `http`), JSON/arrays (`{`, `[`), paths (`/`, `\`), HTML (`<`), or Base64-like JWTs (`eyj`).
    *   Excludes common `utm_*` tracking parameters.
*   **File Extensions:**
    *   Ends with extensions like `.php`, `.asp`, `.aspx`, `.json`, `.xml`, `.cgi`, `.pl`, `.py`, `.sh`, `.yaml`, `.ini`, `.md`, `.do`, `.jsp`, etc.
    *   (Excludes common static file extensions like `.js`, `.html`, `.css`, images - see `isBoringStaticFile` function in code).
*   **Path Components:**
    *   Contains keywords like `ajax`, `jsonp`, `admin`, `include`, `src`, `redirect`, `proxy`, `test`, `tmp`, `temp` in the path.
*   **Non-Standard Ports:**
    *   Uses a port other than 80, 443, or no explicit port.

The tool also ensures that each unique URL (based on hostname, path, and sorted query parameter names) is output only once, even if it matches multiple "interesting" criteria or appears multiple times in the input.
