# mirror - Detect Reflected Query String Values

`mirror` is a command-line tool that takes URLs (from stdin or command-line arguments) and checks if their query string parameter values are reflected in the HTTP response body.

## Features

*   Reads URLs from stdin or as command-line arguments.
*   For each URL, iterates through its query parameters.
*   Fetches the URL and checks if the value of each parameter is present in the response body.
*   Prints information about reflected parameters, including a small snippet of context.
*   Configurable minimum length for parameter values to check (to reduce false positives).
*   Configurable User-Agent.
*   Option to skip TLS certificate verification.
*   Configurable request timeout.

## Installation

To install the `mirror` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/mirror@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `mirror` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `mirror` project directory:**
    ```bash
    cd path/to/your/new-hacks/mirror
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `mirror` in the current directory.

## Usage

```bash
mirror [options] [url...]
```
If URLs are provided as arguments, they will be processed. Otherwise, URLs are read from standard input (one URL per line).

### Options

*   `-min-len <number>`: Minimum length of a parameter value to check for reflection (default: 4).
*   `-ua <string>`: User-Agent string for requests (default: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36").
*   `-k`: Skip TLS certificate verification (insecure).
*   `-t <seconds>`: Request timeout in seconds (default: 10).

### Examples

**1. Check a single URL provided as an argument:**
```bash
mirror -min-len 3 "http://example.com/search?query=test&id=12345"
```

**2. Check URLs from a file, piped via stdin:**
```bash
cat urls.txt | mirror -ua "MyMirrorBot/1.0"
```

### Output Format
If a parameter's value is found in the response body (and meets the minimum length criteria), the tool prints a line in the format:
`<URL>: '<param_key>=<param_value>' reflected in response body (...<context_snippet>...)`

Example:
```
http://test.com/page?name=Alice&debug=true: 'name=Alice' reflected in response body (...lo, Alice. Wel...)
```

Error messages (e.g., URL parsing errors, request errors) are printed to stderr.

## Future Enhancements / TODO (from original README)

*   Check for URL-encoded versions of values.
*   Check for reflection in HTTP headers.
*   Option for checking path reflection (not just query parameters).
*   A way to send and check for reflection in POST data.
*   More sophisticated context extraction from the response body.
*   Concurrency for processing multiple URLs faster.

## How it Works
The tool parses input URLs. For each URL, it iterates through its query parameters. If a parameter's value meets the specified minimum length, the tool makes an HTTP GET request to the URL. It then reads the response body and uses a regular expression to search for the parameter's value within the body, capturing a small amount of surrounding context if found.
