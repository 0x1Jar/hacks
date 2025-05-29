# htmlattribs - Extract HTML Attribute Values

`htmlattribs` is a command-line tool that reads HTML content from standard input and extracts the values of specified HTML attributes. If no attribute names are provided as arguments, it extracts all attribute values.

## Features

*   Reads HTML from stdin.
*   Uses Go's standard `golang.org/x/net/html` tokenizer for parsing.
*   Allows specifying one or more attribute names as command-line arguments to extract only their values.
*   If no attribute names are specified, it extracts the values of all attributes found.
*   Prints each attribute value on a new line.

## Installation

To install the `htmlattribs` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/htmlattribs@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `htmlattribs` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `htmlattribs` project directory:**
    ```bash
    cd path/to/your/new-hacks/htmlattribs
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `htmlattribs` in the current directory.

## Usage

Pipe HTML content to `htmlattribs` via standard input. Optionally, provide attribute names as arguments.

```bash
cat page.html | htmlattribs [attribute_name_1] [attribute_name_2] ...
```

### Examples

**1. Extract specific attribute values (`href`, `src`) from a local HTML file:**
```bash
cat page.html | htmlattribs href src
```

**2. Extract all attribute values from a URL fetched with `curl`:**
```bash
curl -s https://example.com | htmlattribs
```

**3. Extract `type`, `name`, and `href` attributes from `example.com`:**
```bash
curl -s https://example.com | htmlattribs type name href
```
Example Output (order may vary):
```
viewport
text/css
http://www.iana.org/domains/example
```

### Output Format
The tool prints each extracted attribute value on a new line. Only non-empty attribute values are printed.

## How it Works
The tool uses the `html.NewTokenizer` from the `golang.org/x/net/html` package to parse the HTML input stream. It iterates through the tokens. For each token, it inspects its attributes. If specific attribute keys are provided as command-line arguments, it only prints values for those keys. If no arguments are provided, it prints the value of every non-empty attribute it encounters.
