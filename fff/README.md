# fff - Fairly Frickin' Fast HTTP Requester

`fff` is a command-line tool that reads URLs from standard input and makes HTTP requests to them concurrently. It allows for customization of request method, headers, body, and provides options for saving responses.

## Installation

To install the `fff` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/fff@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `fff` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**
1.  **Navigate to the `fff` project directory:**
    ```bash
    cd path/to/your/new-hacks/fff
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `fff` in the current directory.

## Usage

`fff` reads URLs from standard input, one URL per line.

```bash
cat list_of_urls.txt | fff [options]
```

### Options

You can view the options by running `fff -h`:
```
Request URLs provided on stdin fairly frickin' fast

Options:
  -c, --concurrency <num>   Number of concurrent requests (default: 20)
  -b, --body <data>         Request body
  -d, --delay <delay>       Delay between issuing requests (ms) (applied per worker, not globally before each request)
  -H, --header <header>     Add a header to the request (can be specified multiple times, e.g., "User-Agent: fff-client")
  -k, --keep-alive          Use HTTP Keep-Alive
  -m, --method              HTTP method to use (default: GET, or POST if body is specified)
  -o, --output <dir>        Directory to save responses in (will be created, default: out)
  -s, --save-status <code>  Save responses with given status code (can be specified multiple times, e.g., -s 200 -s 302)
  -S, --save                Save all responses
```

### Examples

**1. Basic GET requests to URLs from a file:**
```bash
cat urls.txt | fff
```
Output (if not saving responses):
```
http://example.com/page1 200
http://example.com/page2 404
```

**2. POST requests with a JSON body, custom headers, and save all responses:**
```bash
echo "http://api.example.com/submit" | fff -m POST -b '{"key":"value"}' -H "Content-Type: application/json" -H "Authorization: Bearer token" -S
```
Output (if saving responses):
```
out/api.example.com/submit/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.body: http://api.example.com/submit 200
```
This will save the response body to `out/api.example.com/submit/hash.body` and headers/request info to `out/api.example.com/submit/hash.headers`.

**3. GET requests with a specific User-Agent, saving only 200 responses, with higher concurrency:**
```bash
cat urls.txt | fff -H "User-Agent: MyCustomAgent/1.0" -s 200 -c 50
```

**4. Using a delay between requests per worker:**
```bash
cat urls.txt | fff -d 1000 -c 5 
```
This will run 5 concurrent workers, and each worker will wait 1000ms (1 second) before making its next request.

## Output

*   If responses are **not** saved (default, or if `-s` doesn't match and `-S` is not used):
    The tool prints `<URL> <StatusCode>` to stdout for each request.
*   If responses **are** saved (due to `-S` or matching `-s` status code):
    *   The tool prints `<filepath_to_body>: <URL> <StatusCode>` to stdout.
    *   The response body is saved to `[outputDir]/[hostname]/[normalized_path]/[hash].body`.
    *   Request and response headers, along with the request line and original body, are saved to `[outputDir]/[hostname]/[normalized_path]/[hash].headers`.
    *   The `hash` is a SHA1 hash of the method, URL, request body, and custom headers to ensure uniqueness.

Error messages (e.g., failed requests, file creation errors) are printed to stderr.
