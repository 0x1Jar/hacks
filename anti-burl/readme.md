# anti-burl

Takes URLs on stdin or from a file, fetches them concurrently, and prints details to stdout if they return a 200 OK status code.

## Installation

To install the `anti-burl` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

There are two main ways to install it locally:

1.  **From within the `anti-burl` project directory:**
    Navigate to the `anti-burl` directory (e.g., `/Users/0xjar/Documents/github/hacks/anti-burl`) and run:
    ```bash
    go install .
    ```

2.  **From the parent directory (`hacks`):**
    If you are in the parent directory (`/Users/0xjar/Documents/github/hacks`), you can install it using its relative path:
    ```bash
    go install ./anti-burl
    ```

Either of these commands will compile the source code and install the `anti-burl` binary to your Go binary directory (usually `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH` to run the `anti-burl` command directly from any location.

(Note: The module path is `anti-burl`. Using `go install anti-burl@latest` would typically be for modules hosted on a version control system like GitHub, and requires the module to be published there.)

## Usage

```bash
cat list_of_urls.txt | anti-burl [options]
# OR
anti-burl [options] list_of_urls.txt
```

If no file is specified, input is read from stdin.

Output format: `<status_code> <content_length_runes> <word_count> <url>`

## Options

The program accepts the following command-line flags:

-   `-c int`: Set the concurrency level (default: 50)
-   `-t duration`: Set the request timeout (e.g., 5s, 10s, 1m) (default: 5s)
-   `-ms int`: Set the maximum response body size to read in bytes (default: 1024000)
-   `-k bool`: Skip TLS certificate verification (default: true). Set to `-k=false` to enable verification.
-   `-ua string`: Set the User-Agent string (default: "burl/0.1")
-   `-h`: Show help message.

## Examples

**Basic usage with stdin:**

```bash
echo "https://example.com" | anti-burl
```

**Usage with a file and custom concurrency and timeout:**

```bash
anti-burl -c 20 -t 10s myurls.txt
```

**Usage with a different User-Agent and disabled TLS verification skip:**

```bash
cat urls.txt | anti-burl -ua "MyCustomAgent/1.0" -k=false
```

## Output Example

If `https://example.com` returns a 200 OK, the output might look like:

```
200 1256        180 https://example.com
```

This indicates:
- Status Code: 200
- Content Length (in runes/characters): 1256
- Word Count: 180
- URL: https://example.com

Errors during URL parsing or fetching are printed to stderr.
