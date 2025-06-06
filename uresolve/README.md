# uresolve - Concurrent Host Resolver

`uresolve` is a simple Go tool that reads a list of domain names (or hostnames) from standard input, attempts to resolve each one using DNS lookups, and prints the ones that successfully resolve to standard output.

## How it Works

1.  **Input**: Reads domain names line by line from `stdin`.
2.  **Concurrent Resolution**: For each domain name read:
    *   A new goroutine is launched to handle the DNS lookup for that specific domain.
    *   Inside the goroutine, `net.LookupHost(domain)` is called.
3.  **Output**:
    *   If `net.LookupHost(domain)` returns no error (meaning the domain resolved to at least one IP address), the original domain name is printed to `stdout`.
    *   Domains that fail to resolve are silently ignored (no output for them).
4.  **Synchronization**: The tool uses a `sync.WaitGroup` to ensure that all launched goroutines complete their DNS lookups before the main program exits.
5.  **Error Handling**: If an error occurs while reading from standard input, the program will terminate and print an error message to standard error.

**Note on Output Order**: Due to the concurrent nature of the DNS lookups, the order of resolved domains printed to output is not guaranteed to match the order in the input list.

## Installation

Ensure you have Go (version 1.24.3 or later, as specified in `go.mod`) installed.

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/0x1Jar/new-hacks.git
    cd new-hacks/uresolve
    ```

2.  **Build the tool:**
    ```bash
    go build
    ```
    This will create a `uresolve` executable in the current directory.

Alternatively, you can install it directly if your Go environment is set up:
```bash
go install github.com/0x1Jar/new-hacks/uresolve
```

## Usage

Pipe a list of domain names (one per line) to the tool via standard input:

```bash
cat list_of_domains.txt | ./uresolve
```

Or, if installed to your PATH:
```bash
cat list_of_domains.txt | uresolve
```

**Example:**

Input (`domains.txt`):
```
google.com
thisdomainprobablydoesnotexist123abc.com
github.com
anothernonexistentdomainzyxw.org
```

Command:
```bash
cat domains.txt | uresolve
```

Expected Output (order may vary):
```
google.com
github.com
```

Only `google.com` and `github.com` are printed because they are resolvable. The non-existent domains are silently ignored.
