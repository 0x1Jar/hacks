# check-cnames - CNAME Resolution Checker

`check-cnames` is a command-line tool that reads domain names from standard input, resolves their CNAME records, and then checks if those CNAME targets actually resolve to an IP address. It's useful for identifying potentially misconfigured DNS records or dangling CNAMEs.

## Features

*   Reads domain names from stdin.
*   Resolves CNAME records for each domain using a list of public DNS servers.
*   Checks if the CNAME target resolves to an IP address.
*   Prints domains whose CNAME targets do not resolve.
*   Uses concurrent workers for faster processing.

## Prerequisites

*   **Go**: Version 1.16 or newer.

## Installation

You can install `check-cnames` using `go install`:

```bash
go install github.com/0x1Jar/new-hacks/check-cnames@latest
```
This will install the `check-cnames` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Ensure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `check-cnames` project directory:**
    ```bash
    cd path/to/your/new-hacks/check-cnames
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `check-cnames` in the current directory. You would then run it as `./check-cnames`.

## Usage

Pipe a list of domain names to `check-cnames` via standard input.

```bash
cat list_of_domains.txt | check-cnames
```

Or using `echo`:
```bash
echo -e "subdomain1.example.com\nsubdomain2.example.org" | check-cnames
```

### Input Format

The tool expects one domain name per line.

### Output Format

If a domain's CNAME target does not resolve, the tool will print a message in the following format:
`<cname_target> does not resolve (pointed at by <original_domain>)`

**Example Output:**
```
dangling.cname.target.com does not resolve (pointed at by subdomain1.example.com)
```

## DNS Servers Used

The tool uses a hardcoded list of public DNS servers for CNAME lookups. These include servers from Google, Cloudflare, and Quad9. A random server from this list is chosen for each lookup.

## Concurrency

The tool uses a default of 20 concurrent workers to perform DNS lookups. This can be configured using the `-c` flag.

### Options

*   `-c int`: Set the concurrency level (default: 20).

**Example with custom concurrency:**
```bash
cat list_of_domains.txt | check-cnames -c 50
```
