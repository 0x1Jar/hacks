# Strip Wildcards

`strip-wildcards` is a Go tool that reads a list of domain names from standard input, attempts to identify wildcard DNS records affecting these domains or their parent domains, and prints only those domain names that do not appear to be covered by such wildcards.

## How it Works

For each input domain (e.g., `sub.example.com`):
1.  The tool generates a unique random string (e.g., `randomprefix`) once at the start of its execution.
2.  It then checks its parent domains for wildcard behavior. For an input like `one.two.target.com`, it will perform DNS lookups for:
    -   `randomprefix.target.com`
    -   `randomprefix.two.target.com`
3.  A domain (e.g., `target.com`) is considered to have a wildcard if a DNS lookup for `randomprefix.target.com` successfully resolves.
4.  If any of these checks indicate a wildcard, the original input domain (e.g., `one.two.target.com`) is considered covered by a wildcard and is **not** printed.
5.  If none of the checks indicate a wildcard, the original input domain is printed to standard output.
6.  The tool caches the results of wildcard checks for base domains (like `target.com`, `two.target.com`) to avoid redundant DNS lookups.

This method helps filter out domains that are likely part of a wildcard DNS setup (e.g., `*.example.com` resolving all non-existent subdomains to an IP).

## Installation

Ensure you have Go (version 1.20 or later recommended) installed on your system. You can install `strip-wildcards` using:

```bash
go install github.com/0x1Jar/new-hacks/strip-wildcards
```
Alternatively, you can build from source:
```bash
git clone https://github.com/0x1Jar/new-hacks.git
cd new-hacks/strip-wildcards
go build
```
This will create a `strip-wildcards` executable in the current directory.

## Usage

Pipe a list of domain names (one per line) to the tool via standard input:

```bash
cat list_of_domains.txt | strip-wildcards
```

Or, after building:
```bash
cat list_of_domains.txt | ./strip-wildcards
```

### Flags

-   `-v`: Verbose mode. If specified, the tool will print the total count of DNS lookups performed to standard error when it finishes.
    ```bash
    cat list_of_domains.txt | strip-wildcards -v
    ```

## Example

Suppose `example.com` has a wildcard DNS record `*.example.com` that resolves.
And `specific.target.com` is a normal A record, but `*.target.com` does not resolve (i.e., no wildcard for `target.com`).

Input (`domains.txt`):
```
test1.example.com
foo.specific.target.com
another.example.com
specific.target.com
```

Command:
```bash
cat domains.txt | strip-wildcards
```

Expected Output:
```
foo.specific.target.com
specific.target.com
```

**Explanation:**
-   `test1.example.com`: The tool checks `randomprefix.example.com`. Since `*.example.com` is a wildcard, this resolves. So, `test1.example.com` is stripped.
-   `foo.specific.target.com`: The tool checks `randomprefix.target.com` (no wildcard) and `randomprefix.specific.target.com` (no wildcard, assuming `specific.target.com` is not itself a wildcard). So, `foo.specific.target.com` is printed.
-   `another.example.com`: Stripped for the same reason as `test1.example.com`.
-   `specific.target.com`: The tool checks `randomprefix.target.com` (no wildcard). So, `specific.target.com` is printed.

If run with `-v`, it would also output something like:
```
DNS lookup count: X
```
to stderr, where X is the number of unique DNS lookups performed (due to caching).
