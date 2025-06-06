# URLTeam Downloader (`urlteamdl`)

This directory contains a Go program (`main.go`) and two helper shell scripts (`extract.sh`, `search.sh`) for working with URLTeam data from archive.org.

## Components

### 1. `main.go` (Go Program)

**Purpose**: Fetches download links for ZIP files from the "UrlteamWebCrawls" collection on archive.org, starting from a specified date.

**How it Works**:
1.  Takes a single command-line argument: a "since" date in `YYYY-MM-DD` format.
2.  Constructs a search query for `archive.org/advancedsearch.php` to find identifiers in the `UrlteamWebCrawls` collection added from the "since" date to the current date.
3.  Retrieves up to 500 identifiers, sorted by most recent. (Note: Does not currently handle pagination for more than 500 results in the date range).
4.  For each identifier, it queries `archive.org/metadata/{identifier}` to get a list of files.
5.  Filters these files to find those with `Format: "ZIP"`.
6.  Prints the direct download URL for each ZIP file found (e.g., `https://archive.org/download/{identifier}/{filename.zip}`).
7.  Uses an HTTP client with a 30-second timeout for requests.
8.  Includes error handling for network requests, HTTP status codes, and JSON decoding, printing errors to standard error.

**Installation (`main.go`)**:
Ensure Go (version 1.24.3 or later, as specified in `go.mod`) is installed.
```bash
# Clone the repository (if not already done)
# git clone https://github.com/0x1Jar/new-hacks.git
# cd new-hacks/urlteamdl

go build
# or
go install github.com/0x1Jar/new-hacks/urlteamdl
```

**Usage (`main.go`)**:
```bash
./urlteamdl YYYY-MM-DD > download_links.txt
# or if installed to PATH
urlteamdl YYYY-MM-DD > download_links.txt
```
Example:
```bash
urlteamdl 2023-01-01 > urlteam_zips_since_2023.txt
```
This will output a list of URLs, one per line, suitable for use with download managers like `wget` or `aria2c`.

### 2. `extract.sh` (Shell Script)

**Purpose**: Extracts all `.zip` files in the current directory into their own subdirectories.

**How it Works**:
-   Lists all `*.zip` files in the current directory.
-   For each zip file, it creates a subdirectory named after the zip file (without the `.zip` extension) and extracts the contents of the zip file into that subdirectory.
-   Uses `xargs` and `unzip`.

**Usage (`extract.sh`)**:
1.  Ensure `unzip` and `xargs` are installed.
2.  Make the script executable: `chmod +x extract.sh`
3.  Run in a directory containing downloaded `.zip` files (e.g., those from `urlteamdl`):
    ```bash
    ./extract.sh
    ```

### 3. `search.sh` (Shell Script)

**Purpose**: Searches for a given pattern within `.txt.xz` files (presumably found within the extracted ZIPs).

**How it Works**:
-   Recursively finds all files named `*.txt.xz` in the current directory and its subdirectories.
-   For each found file:
    -   Decompresses it using `xzcat`.
    -   Filters out lines starting with `#` (comments).
    -   Extracts the second field, assuming fields are delimited by `|`.
    -   Searches for the provided pattern (first argument to `search.sh`) in this second field.
-   Uses `xargs`, `xzcat`, `grep`, and `cut`.

**Usage (`search.sh`)**:
1.  Ensure `xz`, `grep`, `cut`, and `xargs` are installed.
2.  Make the script executable: `chmod +x search.sh`
3.  Run in a directory structure containing `.txt.xz` files (e.g., after running `extract.sh`):
    ```bash
    ./search.sh "your_search_pattern"
    ```
    Example: To search for lines containing "example.com" in the relevant fields:
    ```bash
    ./search.sh "example.com" > search_results.txt
    ```

## Workflow Example

1.  Fetch download links for new URLTeam ZIP archives:
    ```bash
    urlteamdl 2023-05-01 > links_to_download.txt
    ```
2.  Download the ZIP files (e.g., using `wget`):
    ```bash
    wget -i links_to_download.txt
    ```
3.  Extract all downloaded ZIP files:
    ```bash
    ./extract.sh
    ```
4.  Search within the extracted text files:
    ```bash
    ./search.sh "targetdomain.com"
