# filter-resolved

`filter-resolved` is a command-line tool that reads domain names from standard input and outputs only those that successfully resolve to an IP address.

## Features

*   Reads domain names from stdin.
*   Concurrently attempts to resolve each domain.
*   Outputs only the domains that resolve.
*   Configurable concurrency level.

## Installation

To install the `filter-resolved` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/filter-resolved@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `filter-resolved` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `filter-resolved` project directory:**
    ```bash
    cd path/to/your/new-hacks/filter-resolved
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `filter-resolved` in the current directory.

## Usage

Pipe a list of domain names (one domain per line) to `filter-resolved` via standard input.

```bash
cat list_of_domains.txt | filter-resolved [options] > resolved_domains.txt
```

Or using `echo`:
```bash
echo -e "example.com\nnonexistentdomain123abc.com\ngoogle.com" | filter-resolved
```
Example output:
```
example.com
google.com
```

### Options

*   `-c <number>`: Set the concurrency level for DNS lookups (default: 20).

**Example with custom concurrency:**
```bash
cat list_of_domains.txt | filter-resolved -c 50 > resolved_domains.txt
```

## How it Works

The tool reads domain names from stdin and distributes them to a pool of worker goroutines. Each worker attempts to resolve the domain to an IPv4 address using `net.ResolveIPAddr("ip4", domain)`. If the resolution is successful (no error), the domain name is printed to standard output. The number of concurrent workers is controlled by the `-c` flag.
