# unfurl

`unfurl` is a Go tool that parses URLs provided on standard input and extracts specified components or formats them according to a given template.

## How it Works

1.  **Input**: Reads URLs line by line from `stdin`.
2.  **Parsing**: Each line is parsed as a URL. If parsing fails, an error can be shown with the `-v` (verbose) flag.
3.  **Mode of Operation**: The tool operates in one of several modes (e.g., `keys`, `values`, `domains`, `paths`, `format`):
    *   `keys`: Extracts all unique keys from the URL's query string.
    *   `values`: Extracts all values from the URL's query string.
    *   `domains`: Extracts the hostname.
    *   `paths`: Extracts the path component.
    *   `format`: Formats the URL components according to a user-supplied format string using specific directives (e.g., `%s` for scheme, `%d` for domain).
4.  **Output**: The extracted or formatted strings are printed to `stdout`, one per line.
5.  **Uniqueness**: The `-u` or `--unique` flag can be used to ensure that only unique values are printed.

## Help

```
▶ unfurl -h
Format URLs provided on stdin

Usage:
  unfurl [OPTIONS] [MODE] [FORMATSTRING]

Options:
  -u, --unique   Only output unique values
  -v, --verbose  Verbose mode (output URL parse errors)

Modes:
  keys     Keys from the query string (one per line)
  values   Values from the query string (one per line)
  domains  The hostname (e.g. sub.example.com)
  paths    The request path (e.g. /users)
  format   Specify a custom format (see below)

Format Directives:
  %%  A literal percent character
  %s  The request scheme (e.g. https)
  %d  The domain (e.g. sub.example.com)
  %P  The port (e.g. 8080)
  %p  The path (e.g. /users)
  %q  The raw query string (e.g. a=1&b=2)
  %f  The page fragment (e.g. page-section)

Examples:
  cat urls.txt | unfurl keys
  cat urls.txt | unfurl format %s://%d%p?%q

```

## Installation

Ensure you have Go installed on your system. You can install `unfurl` using:

```bash
go install github.com/0x1Jar/new-hacks/unfurl
```
This will place the compiled binary in your Go bin directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`).

Alternatively, to build from source:

```bash
git clone https://github.com/0x1Jar/new-hacks.git
cd new-hacks/unfurl
go build
```
This will create an `unfurl` executable in the current directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
