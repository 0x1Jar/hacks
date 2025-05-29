# Tokenizer (`tok`)

`tok` is a Go tool that tokenizes input from standard input. It reads text rune by rune and splits it into tokens based on character types (letters, numbers) and specified delimiters. It offers options to filter tokens by length, content (alphanumeric), and to perform URL decoding.

## How it Works

1.  **Input Reading**: Reads text from standard input rune by rune.
2.  **Token Building**:
    *   A token is built by accumulating consecutive characters that are either letters, numbers, or characters specified in the `-delim-exceptions` flag.
    *   Any character not meeting these criteria acts as a delimiter, signaling the end of the current token.
3.  **State Tracking**: For each token being built, the script tracks:
    *   If it contains any letters (`includesLetters`).
    *   If it contains any numbers (`includesNumbers`).
    *   If it might be URL-encoded (i.e., if a `%` character has been encountered: `maybeURLEncoded`).
4.  **Token Processing (when a delimiter is found or EOF is reached)**:
    *   The accumulated string (token) is processed.
    *   **Length Check**: If the token's length is less than `-min` or greater than `-max`, it's discarded.
    *   **Alpha-Numeric Only Check** (if `-alpha-num-only` is true):
        *   The token is discarded if it does not contain at least one letter AND at least one number.
        *   The token is discarded if it's identical to the previously printed token (simple deduplication for this mode, based on pre-decoded string).
    *   **URL Decoding**: If `maybeURLEncoded` is true for the token, the script attempts to URL-unescape it (using `url.QueryUnescape`, which decodes `+` to space). If decoding is successful, the decoded string becomes the token.
    *   **Output**: If the token passes all checks, it's printed to standard output on a new line.
    *   The state (flags for letters/numbers/URL-encoding, and the token buffer) is then reset for the next token.

## Installation

Ensure you have Go (version 1.24.3 or later, as specified in `go.mod`) installed.

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/0x1Jar/new-hacks.git
    cd new-hacks/tok
    ```

2.  **Build the tool:**
    ```bash
    go build
    ```
    This will create a `tok` executable in the current directory.

Alternatively, you can install it directly if your Go environment is set up:
```bash
go install github.com/0x1Jar/new-hacks/tok
```

## Usage

Pipe text to the tool via standard input. If you have installed the tool using `go install` (and your Go bin directory is in your system's PATH), you can run:

```bash
cat input.txt | tok [flags]
```

If you have built the tool from source (e.g., using `go build`) and are running it from its directory, you would use:
```bash
cat input.txt | ./tok [flags]
```

### Flags

-   `-min int`: Minimum length of a token to be output (default: `1`).
-   `-max int`: Maximum length of a token to be output (default: `25`).
-   `-alpha-num-only`: If true, only output tokens that contain at least one letter AND at least one number. Also enables simple deduplication of consecutively identical (pre-decoded) tokens in this mode. (default: `false`).
-   `-delim-exceptions string`: A string of characters that should NOT be treated as delimiters, even if they are not letters or numbers (e.g., providing `.-_` would allow tokens like `file-name_v1.0`).

## Examples

1.  **Basic Tokenization:**
    ```bash
    echo "word1 word2_v3 word4!" | ./tok
    ```
    Output:
    ```
    word1
    word2
    v3
    word4
    ```
    (Here, `_` and `!` are delimiters by default)

2.  **With Delimiter Exceptions:**
    ```bash
    echo "word1 word2_v3 word4!" | ./tok -delim-exceptions "_."
    ```
    Output:
    ```
    word1
    word2_v3
    word4
    ```
    (`_` is no longer a delimiter, `.` would also be allowed if present)

3.  **Length Filtering:**
    ```bash
    echo "a ab abc abcd abcde" | ./tok -min 3 -max 4
    ```
    Output:
    ```
    abc
    abcd
    ```

4.  **Alpha-Numeric Only:**
    ```bash
    echo "word only_letters num123 mixed_v1 mixed_v1" | ./tok -alpha-num-only -delim-exceptions "_"
    ```
    Output:
    ```
    num123
    mixed_v1 
    ```
    (`word` and `only_letters` are skipped. The second `mixed_v1` is skipped due to deduplication in this mode.)

5.  **URL Decoding:**
    ```bash
    echo "test%20value path%2Fname" | ./tok
    ```
    Output:
    ```
    test value
    path/name
    ```
    (Note: `url.QueryUnescape` decodes `%20` to space and `+` to space.)

## .gitignore

The `.gitignore` file includes:
-   `*.sw*`: Ignores Vim swap files.
-   `tok`: Ignores the compiled binary.
-   `testfile.html`: Likely a file used for local testing.
