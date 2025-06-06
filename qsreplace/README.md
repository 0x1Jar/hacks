# qsreplace - Query String Parameter Replacer (Archived Version)

**❗ This Tool Has Moved!**
The author, Tom Hudson (tomnomnom), has moved `qsreplace` to its own dedicated repository. For the latest version and ongoing development, please visit:
[https://github.com/tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace)

This version in the `new-hacks` repository is an older snapshot. The information below pertains to this archived version.

---

`qsreplace` is a command-line tool that reads URLs from standard input and replaces the values of their query string parameters with a user-provided value. It can also append the value instead of replacing it.

## Features (This Archived Version)

*   Reads URLs from stdin.
*   Replaces all query parameter values with a specified string.
*   Optionally appends the string to existing parameter values.
*   Ensures unique output for combinations of hostname, path, and parameter names (sorted).

## Installation (This Archived Version)

To install this specific version from the `new-hacks` repository, ensure you have Go installed (version 1.16 or newer is recommended).

```bash
go install github.com/0x1Jar/new-hacks/qsreplace@latest
```
This will install the `qsreplace` binary from this repository to your Go binary directory.

## Usage (This Archived Version)

Pipe URLs (one per line) to `qsreplace` via standard input, followed by the replacement value as a command-line argument.

```bash
cat list_of_urls.txt | qsreplace [options] <replacement_value>
```

### Options

*   `-a`: Append the `<replacement_value>` to existing parameter values instead of replacing them.

### Arguments

*   `<replacement_value>`: (Required) The string to use for replacing or appending to query parameter values.

### Examples

**1. Replace all parameter values with "TEST":**
```bash
echo "http://example.com/search?q=old&page=1" | qsreplace TEST
```
Output:
```
http://example.com/search?page=TEST&q=TEST
```
*(Note: Parameter order in output may vary due to map iteration and sorting for uniqueness.)*

**2. Append "XYZ" to all parameter values:**
```bash
echo "http://example.com/api?token=abc&id=123" | qsreplace -a XYZ
```
Output:
```
http://example.com/api?id=123XYZ&token=abcXYZ
```

**3. Processing multiple unique URLs:**
```bash
echo -e "http://example.com/path?a=1&b=2\nhttp://example.com/path?b=3&a=4" | qsreplace INJECTED
```
Output (only one line because hostname, path, and sorted parameter names 'a,b' are the same):
```
http://example.com/path?a=INJECTED&b=INJECTED
```

## How it Works (This Archived Version)
The tool reads each URL from stdin, parses it, and iterates through its query parameters. For each parameter, it either replaces its value with the user-supplied `<replacement_value>` or appends the `<replacement_value>` if the `-a` flag is used. To avoid duplicate output for URLs that only differ in parameter values (but have the same set of parameter names), it creates a unique key based on the hostname, path, and sorted parameter names.
