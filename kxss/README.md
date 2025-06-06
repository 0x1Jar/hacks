# kxss - XSS Parameter Reflection Checker

`kxss` is a tool designed to identify parameters in URLs that are reflected in the HTTP response. This is often a first step in finding potential Cross-Site Scripting (XSS) vulnerabilities.

*Original Author's Note: "I don't know what this name is... Like all tools in my hacks repo it's alpha-quality (at best) so you're likely to find rough edges."*

## Idea / How it Works

The tool processes URLs from standard input through a multi-stage pipeline:

1.  **Initial Reflection Check (`checkReflected`):**
    *   Takes a URL.
    *   Makes a GET request to the URL.
    *   Checks if any of the URL's query parameter *values* are found in the response body.
    *   This stage can have false positives.

2.  **Append Check (`checkAppend` with random string):**
    *   For each parameter identified in the first stage, a unique, random-like string (`kXssRand0mStr1ng`) is appended to its value.
    *   A new request is made with this modified URL.
    *   The `checkReflected` logic is run again. If the *original parameter name* is still identified as "reflected" (meaning the server likely processed and reflected the parameter name itself, or the modified value containing the random string was reflected in a way that still includes the original parameter's value context), it's considered a more confident reflection.

3.  **Special Character Check (`checkAppend` with probe characters):**
    *   For parameters that passed the append check, this stage tests if common XSS-related special characters can be reflected.
    *   It appends a payload like `kXssT3st<char>P4yL0ad` (where `<char>` is one of `"',<>()\`;{}`) to the parameter's value.
    *   If the parameter is still identified as "reflected" by `checkReflected` (meaning the special character didn't break the reflection and the parameter name/value context is still found), it prints that the parameter allows that specific character.

## Installation

To install the `kxss` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/kxss@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `kxss` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**
1.  **Navigate to the `kxss` project directory:**
    ```bash
    cd path/to/your/new-hacks/kxss
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `kxss` in the current directory.

## Usage

Pipe URLs (with query parameters) to `kxss` via standard input.

```bash
cat list_of_urls.txt | kxss [options]
```

### Options
*   `-c <number>`: Number of concurrent workers per stage (default: 20).
*   `-ua <string>`: User-Agent string for requests (default: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36").

### Example

```bash
echo "http://testsite.com/search?query=test&page=1" | kxss -c 10
```
Example Output (if 'query' parameter reflects special characters):
```
param query is reflected and allows " on http://testsite.com/search?query=test&page=1
param query is reflected and allows ' on http://testsite.com/search?query=test&page=1
param query is reflected and allows < on http://testsite.com/search?query=test&page=1
param query is reflected and allows > on http://testsite.com/search?query=test&page=1
...
```
Error messages are printed to stderr.

*(Note: The original README mentioned a test server in `cmd/testserver`. This directory was not provided in the current context, so specific examples using it have been omitted. You can create a simple local server that reflects parameters to test `kxss`.)*

## Further Development Ideas (from original README)

*   **Support POST parameters.**
*   **Rate-limiting:** The tool can generate many requests; rate-limiting per host would be beneficial.
*   **Contextual Payloads:** Determine the context of reflection (HTML, attribute, script) to prioritize and tailor special character/payload testing.
*   **Full XSS Payload Testing & Validation:** Beyond individual characters, try full XSS payloads and potentially use headless Chrome (e.g., with `chromedp`) for validation, though this is resource-intensive.
