# manyreqs - Send Many Requests with Dynamic Parameters and Headers

`manyreqs` is a command-line tool that reads URLs from standard input and sends multiple GET requests to each. It dynamically constructs query parameters and headers based on patterns defined in local files named `params` and `headers`. Parameters are chunked to avoid overly long URLs.

## Features

*   Reads base URLs from stdin.
*   Loads parameter patterns from a local `params` file.
*   Loads header patterns from a local `headers` file.
*   Supports `%s` in pattern values, which will be replaced by the hostname of the target URL.
*   Chunks parameters into multiple requests if the total number of parameters exceeds a configurable limit.
*   Configurable concurrency level for sending requests.
*   Configurable number of parameters per request chunk.
*   Skips TLS certificate verification by default and does not follow redirects.

## Prerequisites

*   **Go**: Version 1.16 or newer is recommended.
*   **`params` file (optional but recommended for functionality):**
    *   Create a file named `params` in the same directory where `manyreqs` is run.
    *   Each line should be a key-value pair for a query parameter, formatted as `key=value`.
    *   Example `params` file:
        ```
        q=searchterm
        debug=true
        target=%s
        ```
*   **`headers` file (optional but recommended for functionality):**
    *   Create a file named `headers` in the same directory.
    *   Each line should be an HTTP header, formatted as `Header-Name: Value`.
    *   Example `headers` file:
        ```
        User-Agent: MyCustomAgent/1.0
        X-Target-Host: %s
        Referer: https://example.com/
        ```

If `params` or `headers` files are not found, the tool will proceed with empty parameter/header sets for those, effectively just making requests to the base URLs provided.

## Installation

To install the `manyreqs` command-line tool, ensure you have Go installed.

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/manyreqs@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `manyreqs` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `manyreqs` project directory:**
    ```bash
    cd path/to/your/new-hacks/manyreqs
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `manyreqs` in the current directory.

## Usage

Pipe a list of base URLs (one URL per line) to `manyreqs` via standard input.

```bash
cat list_of_urls.txt | manyreqs [options]
```

### Options

*   `-p <number>`: Number of parameters to include per request chunk (default: 40).
*   `-c <number>`: Concurrency level (number of worker goroutines) (default: 20).

### Example

Assuming you have `params` and `headers` files in your current directory:

**`params` file:**
```
param1=value1
param2=test
site=%s
```

**`headers` file:**
```
X-Custom-Header: MyValue
X-Forwarded-For: 127.0.0.1
```

**Command:**
```bash
echo "https://api.example.com/data" | manyreqs -p 2 -c 5
```

**Behavior:**
`manyreqs` will read `https://api.example.com/data`.
It will load `param1=value1`, `param2=test`, and `site=api.example.com` from the `params` file.
It will load `X-Custom-Header: MyValue` and `X-Forwarded-For: 127.0.0.1` from the `headers` file.

Since `-p 2` is set, it will make requests in chunks of 2 parameters:
1.  Request to `https://api.example.com/data?param1=value1&param2=test&` (with the defined headers).
2.  Request to `https://api.example.com/data?site=api.example.com&` (with the defined headers).

The tool primarily discards the response bodies but will print errors to stderr if requests fail.

## How it Works
The tool reads URLs from stdin. For each URL, it reads `params` and `headers` files to get patterns. It substitutes `%s` in these patterns with the hostname of the current URL. It then groups the generated parameters into chunks based on the `-p` flag. For each chunk, it constructs a full URL with those parameters and sends a GET request with the generated headers. This is done concurrently using a pool of workers defined by the `-c` flag.
