# cors-blimey - CORS Misconfiguration Scanner

`cors-blimey` is a command-line tool that reads URLs from standard input and tests them for common Cross-Origin Resource Sharing (CORS) misconfigurations by sending requests with various crafted `Origin` headers.

## Features

*   Reads target URLs from stdin.
*   Tests a range of `Origin` header permutations against each target URL.
*   Checks `Access-Control-Allow-Origin` (ACAO) and `Access-Control-Allow-Credentials` (ACAC) headers in responses.
*   Highlights potentially vulnerable configurations (e.g., reflected origin with credentials, wildcard with credentials).
*   Configurable concurrency and HTTP client timeout.
*   Uses a custom HTTP client that does not follow redirects and can skip TLS verification (currently hardcoded to skip).

## Installation

To install the `cors-blimey` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/cors-blimey@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `cors-blimey` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `cors-blimey` project directory:**
    ```bash
    cd path/to/your/new-hacks/cors-blimey
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `cors-blimey` in the current directory.

## Usage

Pipe a list of target URLs (one URL per line) to `cors-blimey` via standard input.

```bash
cat list_of_target_urls.txt | cors-blimey [options]
```

Or using `echo`:
```bash
echo "https://api.example.com/user" | cors-blimey -c 10 -t 15
```

### Options

*   `-c <number>`: Set the concurrency level for making requests (default: 20).
*   `-t <seconds>`: Set the HTTP client timeout in seconds (default: 10).
*   *(Future)* `-origins <filepath>`: Path to a custom file of origins to test.
*   *(Future)* `-patterns <filepath>`: Path to a custom file of origin patterns.
*   *(Future)* `-skip-verify <true|false>`: Skip TLS certificate verification (currently hardcoded to true).

### Output Format

The tool prints findings to standard output. Each line represents a potentially interesting CORS configuration found for a target URL with a specific tested origin.

**Example Output:**
```
Target: https://api.example.com/data | Origin: https://evil.com | ACAO: https://evil.com | ACAC: true (VULNERABLE!)
Target: https://vulnerable.site/api | Origin: null | ACAO: null | ACAC: true (VULNERABLE!)
Target: https://test.com/resource | Origin: https://sub.test.com | ACAO: * | ACAC: true (VULNERABLE!)
Target: https://another.api/info | Origin: https://sub.another.api.attacker.com | ACAO: https://sub.another.api.attacker.com | ACAC: false
Target: https://service.example.org/endpoint | Origin: https://example.org | ACAO: https://example.org (Potentially interesting parent domain reflection) | ACAC: true
```
Error messages for invalid URLs or request failures are printed to standard error.

## Tested Origin Permutations

For each target URL, `cors-blimey` generates and tests `Origin` headers including:
*   `null`
*   `https://evil.com`, `http://evil.com`
*   The target URL itself (scheme + hostname)
*   A generic subdomain of the target (e.g., `https://sub.targethostname`)
*   Variations using `attacker.com` (e.g., `https://targethostname.attacker.com`, `https://sub.targethostname.attacker.com`)
*   The TLD+1 of the target (e.g., if target is `api.dept.example.com`, it tests `https://example.com`)
*   And more (see `getPermutations` function in `main.go`).

## TODO (from original README & code)
*   Implement the idea: "i have also seen stuff like they don't allow target.com.evil.com or anything after the first . but if you put target.anothertld it works. Then you just purcahse that domain and you are set to go" - This involves more sophisticated TLD manipulation in `getPermutations`.
*   Make HTTP client options (MaxIdleConns, IdleConnTimeout, InsecureSkipVerify, KeepAlive for Dialer) configurable via flags.
*   Allow custom lists of origins and origin patterns to be provided via files.
*   More granular error reporting for HTTP requests.
