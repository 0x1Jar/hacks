# inscope - Scope Filtering Tool

`inscope` is a command-line tool that filters a list of URLs and/or domain names supplied on standard input, outputting only those that match a set of regular expressions defined in a `.scope` file. It's designed to help filter the output of other security tools to ensure results are within the defined scope for a bug bounty program or penetration test.

## Features

*   Reads URLs or domain names from stdin.
*   Filters input based on regular expressions in a `.scope` file.
*   Automatically searches for the `.scope` file in the current directory and then recursively up through parent directories.
*   Supports positive match patterns and negative match patterns (lines starting with `!`).
*   If a URL is provided, only its hostname is checked against the scope.

## Installation

To install the `inscope` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/inscope@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `inscope` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `inscope` project directory:**
    ```bash
    cd path/to/your/new-hacks/inscope
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `inscope` in the current directory.

## Usage

Pipe URLs and/or domain names (one per line) to `inscope` via standard input.

```bash
cat list_of_urls_or_domains.txt | inscope
```

### Example

Given a `testinput` file:
```
https://example.com/footle
https://inscope.example.com/some/path?foo=bar
https://outofscope.example.net/bar
example.com
example.net
http://sub.example.com
```

And a `.scope` file in the current directory or a parent directory (see example below).

Running the command:
```bash
cat testinput | inscope
```
Might produce the following output (depending on the `.scope` file content):
```
https://example.com/footle
https://inscope.example.com/some/path?foo=bar
example.com
http://sub.example.com
```

## Scope File (`.scope`)

The tool reads regexes from a file called `.scope` in the current working directory.
If it doen't find one it recursively checks the parent directory until it hits the root.

Here's an example `.scope` file:

```
.*\.example\.com$
^example\.com$
.*\.example\.net$
!.*outofscope\.example\.net$
```

Each line is a regular expression to match domain names. When URLs are provided as input they
are parsed and only the hostname/domain portion is checked against the regex.

Line starting with `!` are treated as negative matches - i.e. any domain matching that regex will
be considered out of scope even if it matches one of the other regexes.
