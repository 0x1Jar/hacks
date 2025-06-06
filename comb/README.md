# comb

Combine the lines from two files in every combination.

## Installation

To install the `comb` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using the following command:
```bash
go install github.com/0x1Jar/new-hacks/comb@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `comb` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH` to run the `comb` command directly from any location.

**For local development or building from source:**

1.  **Navigate to the `comb` project directory:**
    ```bash
    cd path/to/your/new-hacks/comb
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `comb` in the current directory. You would then run it as `./comb`.


## Usage

The basic command structure is:
```bash
comb [OPTIONS] <prefixfile> <suffixfile>
```

To see help information:
```bash
comb -h
```
Output:
```
Combine the lines from two files in every combination

Usage:
  comb [OPTIONS] <prefixfile> <suffixfile>

Options:
  -f, --flip             Flip mode (order by suffix)
  -s, --separator <str>  String to place between prefix and suffix
```

## Examples

Assume you have two files:

**prefixes.txt:**
```
1
2
```

**suffixes.txt:**
```
A
B
C
```

**Normal mode:**
```bash
comb prefixes.txt suffixes.txt
```
Output:
```
1A
1B
1C
2A
2B
2C
```

**Flip mode (order by suffix):**
```bash
comb -f prefixes.txt suffixes.txt
# or comb --flip prefixes.txt suffixes.txt
```
Output:
```
1A
2A
1B
2B
1C
2C
```

**Separator:**
```bash
comb -s "-" prefixes.txt suffixes.txt
# or comb --separator="-" prefixes.txt suffixes.txt
```
Output:
```
1-A
1-B
1-C
2-A
2-B
2-C
```

## Can't you just do this with a couple of nested bash loops?

Yes, but it's a PITA to type.
