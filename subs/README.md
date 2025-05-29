# Subdomain Enumerator (`subs`)

`subs` is a Go tool that attempts to discover valid subdomains for a given list of apex domains using a provided list of common subdomain names. It performs DNS lookups to verify the existence of constructed `subdomain.apexdomain` combinations.

## How it Works

1.  **Input Files**:
    *   **Domains File**: A list of apex domains (e.g., `example.com`, `example.org`), one per line.
    *   **Subdomains File**: A list of potential subdomain prefixes (e.g., `www`, `dev`, `api`), one per line.

2.  **Wildcard Check**:
    *   For each apex domain from the domains file, the script first performs a wildcard check. It does this by trying to resolve `lkj23lk52lkn23kjh23mnbzzxckjhasdqwe.apexdomain`.
    *   If this random-looking hostname resolves, the apex domain is assumed to have a wildcard DNS record, and it's skipped (no subdomains will be checked for it). This is to avoid generating many false positives from wildcard domains.

3.  **Subdomain Probing**:
    *   For each apex domain that does *not* appear to have a wildcard:
        *   The script iterates through each prefix in the subdomains file.
        *   It constructs a candidate hostname: `prefix.apexdomain`.
        *   It performs a DNS lookup for this candidate.
        *   If the candidate hostname resolves, it's printed to standard output.

4.  **Concurrency**:
    *   The script processes each non-wildcard apex domain concurrently. For each such domain, it spawns two worker goroutines.
    *   These workers read subdomain prefixes from a channel associated with their apex domain and perform the DNS lookups.

## Files Provided

-   `main.go`: The Go source code for the tool.
-   `domains`: An example input file containing apex domains.
-   `subdomains`: An example input file containing common subdomain prefixes.

## Installation

Ensure you have Go (version 1.24.3 or later, as specified in `go.mod`) installed.

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/0x1Jar/new-hacks.git
    cd new-hacks/subs
    ```

2.  **Build the tool:**
    ```bash
    go build
    ```
    This will create a `subs` executable in the current directory.

Alternatively, you can install it directly if your Go environment is set up:
```bash
go install github.com/0x1Jar/new-hacks/subs
```

## Usage

The tool takes two optional command-line arguments:
1.  Path to the domains file.
2.  Path to the subdomains file.

```bash
./subs [path_to_domains_file] [path_to_subdomains_file]
```

**Defaults:**
-   If `path_to_domains_file` is not provided, it defaults to a file named `apexes` in the current directory.
-   If `path_to_subdomains_file` is not provided, it defaults to a file named `subdomains` in the current directory.

**Example using provided files:**
Ensure `domains` and `subdomains` files are in the same directory as the `subs` executable (or provide paths to them).
```bash
./subs domains subdomains > found_subdomains.txt
```
This command will:
-   Read apex domains from the `domains` file.
-   Read subdomain prefixes from the `subdomains` file.
-   Attempt to find valid subdomains.
-   Print any discovered subdomains to standard output, which are then redirected to `found_subdomains.txt`.

**Example using default file names:**
If you rename your input files to `apexes` and `subdomains` (or create symlinks), you can run:
```bash
./subs > found_subdomains.txt
```

## Output

The script prints discovered, resolvable subdomains to standard output, one per line.
Example:
```
www.tomnomnom.com
blog.tomnomnom.com
```
Error messages (e.g., file not found) are printed to standard error.
