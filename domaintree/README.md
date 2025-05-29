# domaintree

`domaintree` is a command-line tool that reads domain names from standard input, parses them, and then prints a hierarchical tree structure of the domains based on their components (TLD, domain, subdomains).

## Features

*   Reads domain names from stdin.
*   Uses `github.com/Cgboal/DomainParser` to accurately identify TLDs.
*   Outputs a text-based tree representing the domain hierarchy.

## Installation

To install the `domaintree` command-line tool, ensure you have Go installed (version 1.21 or newer is recommended, as specified in `go.mod`).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/domaintree@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `domaintree` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `domaintree` project directory:**
    ```bash
    cd path/to/your/new-hacks/domaintree
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `domaintree` in the current directory.

## Usage

Pipe a list of domain names (one domain per line) to `domaintree` via standard input.

```bash
cat list_of_domains.txt | domaintree
```

Or using `echo`:
```bash
echo -e "sub1.example.com\nsub2.example.com\nanother.example.co.uk\nexample.com" | domaintree
```

### Input Format

The tool expects one domain name per line.

### Output Format

The tool prints a tree structure to standard output. Each level of the domain is indented.

**Example Input (`domains.txt`):**
```
www.example.com
api.example.com
example.com
test.staging.example.org
staging.example.org
example.org
beta.test.com
test.com
```

**Example Command:**
```bash
cat domains.txt | domaintree
```

**Example Output:**
```
com
  example
    api
    www
  test
    beta
org
  example
    staging
      test
```
(Note: The exact order of sibling nodes at the same level, e.g., `com` and `org`, or `example` and `test` under `com`, might vary as map iteration order is not guaranteed in Go. However, the hierarchical structure will be correct.)

## How it Works

The tool reads each domain, uses a domain parser to identify the TLD, and then splits the remaining part of the domain into its components. It builds an internal tree data structure where nodes represent parts of the domain (e.g., "com", "example", "www"). Finally, it traverses this tree to print the hierarchical representation.
