# lsinteresting - List Files with "Interesting" Sizes

`lsinteresting` is a command-line tool that lists files in a specified directory. It attempts to identify "interesting" files by only showing those whose sizes are sufficiently different from the sizes of files it has already listed from that directory.

## Features

*   Scans a target directory.
*   Compares file sizes to identify those that are unique enough based on a percentage threshold.
*   Prints the path of files deemed "interesting" (i.e., not too similar in size to previously seen files).
*   Ignores subdirectories.
*   Allows customization of the similarity threshold.

## Installation

To install the `lsinteresting` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/lsinteresting@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `lsinteresting` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `lsinteresting` project directory:**
    ```bash
    cd path/to/your/new-hacks/lsinteresting
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `lsinteresting` in the current directory.

## Usage

```bash
lsinteresting <directory> [threshold_percentage]
```

*   `<directory>`: (Required) The path to the directory you want to scan.
*   `[threshold_percentage]`: (Optional) A floating-point number representing the percentage difference. If a file's size is within this percentage of a previously seen file size, it's considered not "interesting" enough and will be skipped. Defaults to `1.0` (i.e., 1%).

### Examples

**1. List interesting files in `/var/log` using the default threshold (1%):**
```bash
lsinteresting /var/log
```

**2. List interesting files in `my_data_folder` using a 5% threshold:**
```bash
lsinteresting my_data_folder 5.0
```

**3. List interesting files in the current directory with a 0.1% threshold:**
```bash
lsinteresting . 0.1
```

### Output Format
The tool prints the full path to each "interesting" file, one per line.

## How it Works
The tool reads the contents of the specified directory. For each file (it ignores directories), it calculates its size. It maintains a list of sizes of files it has already deemed "interesting". For the current file, it compares its size against all previously recorded "interesting" sizes. If the current file's size is different by more than the specified `threshold_percentage` from all previously recorded sizes, it's considered "interesting", its path is printed, and its size is added to the list for future comparisons. Zero-byte files are handled such that only the first encountered zero-byte file is listed (unless the threshold logic makes others appear different relative to non-zero files).
