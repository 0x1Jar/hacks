# b64d

`b64d` is a command-line tool that finds and decodes base64-encoded strings within a given file or from standard input. It's designed to quickly extract potentially interesting data from files containing mixed content.

## Usage

The tool can be used in two main ways:

1.  **Reading from a file:**

    ```bash
    ▶ b64d <filename>
    ```
    Replace `<filename>` with the path to the file you want to process.

    *Example:*
    If you have a file named `data.txt` containing `SGVsbG8gV29ybGQh`, running:
    ```bash
    ▶ b64d data.txt
    ```
    Would output:
    ```
    Hello World!
    ```

2.  **Reading from standard input (stdin):**

    You can pipe the output of another command directly to `b64d`.

    ```bash
    ▶ cat <filename> | b64d
    ```
    Or, for example, using `echo`:
    ```bash
    ▶ echo "Found this: SGVsbG8gV29ybGQh in the logs" | b64d
    ```
    Would output:
    ```
    Hello World!
    ```

## How it Works

`b64d` scans the input for patterns that look like base64 encoded strings. It specifically looks for:
- Strings composed of valid base64 characters (A-Z, a-z, 0-9, +, /).
- Optional padding characters (= or ==) at the end.
- A non-base64 character preceding the potential base64 string (this helps to avoid false positives on long random-looking strings that are not actually base64).

Once a potential base64 string is identified, `b64d` attempts to decode it. If the decoded output consists only of printable ASCII characters, it is printed to standard output. Each decoded string is printed on a new line.

## Install

To install `b64d`, ensure you have Go installed and configured on your system. Then, run the following command:

```bash
▶ go install github.com/0x1Jar/new-hacks/b64d@latest
```
This command will download the source code, compile it, and install the `b64d` binary into your Go bin directory (usually `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH` environment variable to run `b64d` from any location.
