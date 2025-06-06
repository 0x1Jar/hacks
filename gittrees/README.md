# gittrees

`gittrees` is a command-line tool that lists all unique blob object hashes and their corresponding filenames from all trees in a specified Git repository.

## Features

*   Traverses all tree objects in a Git repository.
*   For each tree, iterates through all files.
*   Outputs the SHA1 hash of each blob object and its filename.
*   Ensures each unique blob hash is printed only once, even if the same file content appears under different names or in different trees.
*   Accepts a path to a Git repository as an argument, or defaults to the current directory.

## Installation

To install the `gittrees` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/gittrees@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `gittrees` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `gittrees` project directory:**
    ```bash
    cd path/to/your/new-hacks/gittrees
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `gittrees` in the current directory.

## Usage

```bash
gittrees [path_to_repo]
```
*   `[path_to_repo]`: (Optional) The path to the Git repository. If not provided, it defaults to the current directory (`.`).

### Output Format
The tool prints one line per unique file blob, in the format:
`<blob_sha1_hash> <filename>`

### Examples

**1. List blobs for a repository at a specific path:**
```bash
gittrees ~/projects/my-git-repo
```
Example Output:
```
79e0815a8c82d0f1f8f467b01d97143a8fed7048 .gitignore
b7214221c4b62b4712f783c2b7c3b8fc4f076fbc .travis.yml
2fa1a8731c67462f4eef8ae6b92623a2cc726a5d ADVANCED.md
...
```

**2. List blobs for the Git repository in the current directory:**
```bash
gittrees
```

**3. Grep contents of specific file types for a keyword:**
This example shows how you might use `gittrees` in a pipeline to search for 'console.log' within all JavaScript files in the history of the current repository:
```bash
gittrees | grep '\.js$' | cut -d' ' -f1 | xargs -I{} git cat-file blob {} | grep -a 'console.log'
```
*   `gittrees`: Lists all blob hashes and filenames.
*   `grep '\.js$'`: Filters for lines ending with `.js`.
*   `cut -d' ' -f1`: Extracts just the blob hash (the first field).
*   `xargs -I{} git cat-file blob {}`: For each hash, shows its content.
*   `grep -a 'console.log'`: Searches the content for 'console.log'.

## How it Works
The tool uses the `gopkg.in/src-d/go-git.v4` library to open and inspect a Git repository. It iterates over all tree objects, and for each tree, it iterates over its files. It keeps track of seen blob hashes to ensure each unique file content is listed only once with one of its associated filenames.
