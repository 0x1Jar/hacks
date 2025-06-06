# compres - HTTP Response Comparator

`compres` (Compare Responses) is a command-line tool that fetches two URLs (a "base" and a "candidate") and then compares their HTTP responses, highlighting various differences.

## Features

*   Fetches two specified URLs.
*   Allows customization of HTTP method and headers for both requests.
*   Supports using an HTTP proxy.
*   Configurable client timeout and keep-alive settings.
*   Compares:
    *   Status code and HTTP protocol version.
    *   Headers (missing, added, or different values).
    *   Body length and SHA256 hash.
    *   Sniffed MIME types.
    *   JSON content differences (if both responses are JSON).
    *   Counts of specific keywords (e.g., "error", "warn") in the body.
*   Outputs a list of detected differences.

## Prerequisites

*   **Go**: Version 1.18 or newer (due to generics used in the original `slicesDiffer` and `inSlice` which are now simplified or could be adapted). The current `go.mod` specifies 1.18.

## Installation

You can install `compres` using `go install`:

```bash
go install github.com/0x1Jar/new-hacks/compres@latest
```
This will install the `compres` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Ensure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `compres` project directory:**
    ```bash
    cd path/to/your/new-hacks/compres
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `compres` in the current directory.

## Usage

```bash
compres -burl <base_url> -curl <candidate_url> [options]
```

### Required Arguments

*   `-burl <url>`: The base URL to fetch.
*   `-curl <url>`: The candidate URL to fetch and compare against the base.

### Options

*   `-bmethod <method>`: HTTP method for the base URL (default: "GET").
*   `-cmethod <method>`: HTTP method for the candidate URL (default: "GET").
*   `-bH "Header: Value"`: Add a header to the base URL request. Can be used multiple times (e.g., `-bH "User-Agent: MyClient" -bH "X-Custom: Base"`).
*   `-cH "Header: Value"`: Add a header to the candidate URL request. Can be used multiple times.
*   `-proxy <proxy_url>`: HTTP proxy to use for requests (e.g., `http://127.0.0.1:8080`).
*   `-timeout <seconds>`: HTTP client timeout in seconds (default: 10).
*   `-keepalive`: Enable HTTP keep-alives (default: false, i.e., disabled).
*   `-h, --help`: Show help information.

### Example

```bash
compres \
  -burl "https://httpbin.org/get?param1=base" \
  -bH "X-Base-Header: BaseValue" \
  -curl "https://httpbin.org/get?param1=candidate" \
  -cH "X-Cand-Header: CandValue" \
  -cH "User-Agent: CompresTool/1.0"
```

This command will fetch `https://httpbin.org/get?param1=base` with an `X-Base-Header` and compare its response to `https://httpbin.org/get?param1=candidate` fetched with `X-Cand-Header` and a custom `User-Agent`.

## Output Format

The tool first prints the URLs it's fetching. If differences are found, it prints a section titled "Differences found:", followed by a list of differences. Each difference is formatted as:
`[kind] key: Base="base_value", Candidate="candidate_value"`
or for JSON diffs:
`[body-json-diff] content: <description_of_difference>`

**Example Difference Output:**
```
[status-code] code: Base="200 OK", Candidate="404 Not Found"
[header-missing] X-Powered-By: Base="Express", Candidate=""
[header-added] X-New-Header: Base="", Candidate="NewValue"
[body] length: Base="1024", Candidate="512"
[body] hash: Base="...", Candidate="..."
[keyword-count] error: Base="0", Candidate="1"
[body-json-diff] content: value at '/args/param1' is different: "base" vs "candidate"
```
If no differences are found, it will print "No differences found."
