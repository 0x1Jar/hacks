# ettu - Recursive DNS Brute-forcer

`ettu` is a recursive DNS brute-forcer designed to discover subdomains. It attempts to avoid "dead-ends" by checking DNS responses to infer if further recursion on a subdomain might be fruitful, even if the subdomain itself doesn't resolve directly.

*Disclaimer: This tool is based on an original concept and may have rough edges.*

## Installation

To install the `ettu` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/ettu@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `ettu` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**
1.  **Navigate to the `ettu` project directory:**
    ```bash
    cd path/to/your/new-hacks/ettu
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `ettu` in the current directory.

## Usage

```bash
ettu [options] <domain> [<wordfile>|-]
```

### Arguments
*   `<domain>`: The target domain to brute-force (e.g., `example.com`).
*   `[<wordfile>|-]`: (Optional) Path to a file containing words to use for subdomain generation (one word per line). If omitted or specified as `-`, `ettu` reads words from standard input.

### Options
*   `-d, --depth <int>`: Maximum recursion depth (default: 4).
*   `-w <int>`: Number of concurrent workers for DNS lookups (default: 10).

### Examples

**1. Using a wordlist file with specified depth and workers:**
```bash
ettu -d 2 -w 50 example.com wordlist.txt
```

**2. Reading words from stdin:**
```bash
cat wordlist.txt | ettu example.com -
```
Or simply, if no wordlist file is provided as an argument:
```bash
cat wordlist.txt | ettu example.com
```

**3. Example from original README:**
```bash
echo -e "www\none\ntwo\nthree" | ettu tomnomnom.uk
```
Expected output (if `one.two.three.tomnomnom.uk` resolves):
```
one.two.three.tomnomnom.uk
```

## Dead-end Avoidance

`ettu` attempts to be smarter than a simple brute-forcer. Ordinarily, if a DNS name doesn't exist, you might get an `NXDOMAIN` error:
```
$ host nonexistentsub.example.com
Host nonexistentsub.example.com not found: 3(NXDOMAIN)
```
However, sometimes you might get an empty response (no error, but no records) for a name like `sub1.example.com`. This can happen if another name, such as `deeper.sub1.example.com`, *does* exist. `ettu` uses this distinction:
*   If a lookup for `word.domain.com` results in `NXDOMAIN` or a timeout, `ettu` stops recursing down that path.
*   If the lookup resolves, or if it results in an error that is *not* `NXDOMAIN` or timeout (e.g., an empty success, or a server failure), `ettu` will attempt to recurse further (e.g., `anotherword.word.domain.com`), up to the specified depth.

This helps focus efforts on potentially fruitful branches of the domain tree.
