# perms - Permutation Generator

`perms` is a command-line tool that generates permutations of lines provided on standard input. It allows for control over the depth of permutations, separators, prefixes, suffixes, and whether input lines can be repeated in a single permutation.

**WARNING**: If you provide too many lines of input, specify too high of a depth, or too many separators (or any combination of these), you can easily generate an enormous number of combinations. This may consume more RAM, disk space, or time than available. Use with caution and an understanding of the potential output size.

## Features

*   Reads an "alphabet" of lines from stdin.
*   Generates permutations based on specified depths.
*   Allows custom separators between elements in a permutation.
*   Can add a prefix and/or suffix to each generated permutation.
*   Option to prevent repeated use of input lines within a single permutation (`-no-repeats`).

## Installation

To install the `perms` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended).

You can install it using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/perms@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `perms` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `perms` project directory:**
    ```bash
    cd path/to/your/new-hacks/perms
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `perms` in the current directory.

## Usage

Pipe a list of lines (one line per element of your "alphabet") to `perms` via standard input.

```bash
cat input_lines.txt | perms [options]
```

### Options

*   `-depth <depth_spec>`: Specifies the recursion depth(s) for permutations. This determines how many elements from the input alphabet are combined.
    *   Can be a single number (e.g., `-depth 3`).
    *   Can be a range (e.g., `-depth 2-4` for depths 2, 3, and 4).
    *   Can be a comma-separated list (e.g., `-depth 1,3` for depths 1 and 3).
    *   Can be specified multiple times (e.g., `-depth 1 -depth 3-4`).
    *   Default: `1,2` (generates permutations of depth 1 and depth 2).
*   `-prefix <string>`: A string to prepend to every generated permutation.
*   `-suffix <string>`: A string to append to every generated permutation.
*   `-sep <string>`: Separator string to use between elements in a permutation.
    *   Can be specified multiple times to use different separators. Each permutation will be generated for each separator.
    *   Default: `""` (empty string, i.e., direct concatenation).
*   `-no-repeats`: If set, each line from the input alphabet can only be used once within a single generated permutation. By default, lines can be repeated.

### Examples

**1. Basic permutations of depth 1 and 2 (default):**
```bash
echo -e "a\nb\nc" | perms
```
Output:
```
a
b
c
aa
ab
ac
ba
bb
bc
ca
cb
cc
```

**2. Permutations of depth 2 with a hyphen separator:**
```bash
echo -e "one\ntwo" | perms -depth 2 -sep "-"
```
Output:
```
one-one
one-two
two-one
two-two
```

**3. Permutations of depth 3, no repeats, with prefix and suffix:**
```bash
echo -e "apple\nbanana\ncherry" | perms -depth 3 -no-repeats -prefix "fruit: " -suffix "!" -sep ", "
```
Example Output (order may vary for no-repeats due to internal iteration):
```
fruit: apple, banana, cherry!
fruit: apple, cherry, banana!
... (all 6 permutations)
```

**4. Permutations of depth 1 and 2, using multiple separators (`.` and `-`):**
```bash
echo -e "x\ny" | perms -depth 1,2 -sep "." -sep "-"
```
Output:
```
x
y
x.x
x.y
y.x
y.y
x-x
x-y
y-x
y-y
```

## How it Works
The tool reads all lines from stdin to form an "alphabet". It then uses a recursive function (`list`) to generate permutations. The `depths` flag controls which levels of recursion produce output. The `seps` flag allows trying different separators for each combination. The `no-repeats` flag modifies the alphabet available at each step of the recursion to prevent re-using elements. Prefixes and suffixes are added to the final output strings.
