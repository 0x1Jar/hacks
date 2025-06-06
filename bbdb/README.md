# bbdb - Bug Bounty Database

`bbdb` is a command-line tool for managing a simple database of domains, typically for bug bounty purposes. It uses a SQLite database (`bbdb.db`) to store domain names.

## Features

*   Add domains to the database.
*   Delete domains from the database.
*   List all domains in the database.

## Installation

To install the `bbdb` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using the following command:
```bash
go install github.com/0x1Jar/new-hacks/bbdb@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `bbdb` binary to your Go binary directory (usually `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH` to run the `bbdb` command directly from any location.

The `bbdb.db` SQLite database file will be created in the directory where you run `bbdb init` or any other `bbdb` command for the first time if it doesn't already exist.

**For local development or building from source:**

1.  **Navigate to the `bbdb` project directory:**
    ```bash
    cd path/to/your/new-hacks/bbdb
    ```

2.  **Build the executable:**
    ```bash
    go build
    ```
    This will create an executable file named `bbdb` (or `bbdb.exe` on Windows) in the current directory. You would then run it as `./bbdb`.

## Initialization

Before using `bbdb` for the first time, you need to initialize the database. This creates the `bbdb.db` file and sets up the necessary tables.

```bash
./bbdb init
```

## Usage

`bbdb` can accept commands either as command-line arguments or from standard input (stdin).

**1. Using Command-Line Arguments:**

This is suitable for single operations.

The command format is: `./bbdb [action] [type] [argument]`
*   **`action`**: `add`, `delete`, `all`
*   **`type`**: `domain` or `domains` (case-insensitive)
*   **`argument`**: The domain name (required for `add` and `delete`, omitted for `all`)

**Examples:**
*   Add a domain: `./bbdb add domain example.com`
*   List all domains: `./bbdb all domains`
*   Delete a domain: `./bbdb delete domain example.com`

**2. Using Standard Input (stdin):**

This is suitable for interactive use or batch processing multiple commands. If no command-line arguments (other than `init`) are provided, `bbdb` will read commands from stdin. Each command should be on a new line.

The command format is: `[action] [type] [argument]`

*   **`action`**: `add`, `delete`, `all`
*   **`type`**: `domain` or `domains` (case-insensitive)
*   **`argument`**: The domain name (required for `add` and `delete`)

### Stdin Examples

**1. Add a domain interactively:**
```bash
./bbdb
add domain example.com
add domain another.example.org
# Press Ctrl+D to finish
```
Output:
```
Added: example.com
Added: another.example.org
```

**2. List all domains interactively:**
```bash
./bbdb
all domains
# Press Ctrl+D to finish
```
Output (if domains exist):
```
example.com
another.example.org
```
Or if no domains:
```
No domains found.
```

**3. Delete a domain interactively:**
```bash
./bbdb
delete domain example.com
# Press Ctrl+D to finish
```
Output:
```
Deleted: example.com
```

**4. Batch Processing via Stdin:**

You can pipe a list of commands from a file or another command:
```bash
# commands.txt:
# add domain test1.com
# add domain test2.com
# all domains

cat commands.txt | ./bbdb
```
Or using `echo`:
```bash
echo -e "add domain test3.com\nall domains" | ./bbdb
```
Output:
```
Added: test3.com
test1.com
test2.com
test3.com
```
(Assuming test1.com and test2.com were already added or added by previous lines in commands.txt)
