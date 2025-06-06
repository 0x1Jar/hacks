# tok

`tok` is a command-line tool for extracting tokens (words, numbers, or alphanumeric strings) from standard input. It reads input character by character, groups sequences of letters and numbers, and outputs them according to customizable rules.

## Features

- Extracts tokens containing letters and/or numbers from input.
- Supports minimum and maximum token length filtering.
- Optionally only outputs tokens containing both letters and numbers.
- Allows specifying delimiter exceptions (characters that should not be treated as token boundaries).
- Detects and decodes URL-encoded tokens.

## Installation

Make sure you have Go installed (version 1.18 or newer recommended).

Clone the repository and build the tool:

```sh
git clone https://github.com/0x1Jar/new-hacks.git
cd new-hacks/tok
go build
```

This will produce a `tok` binary in the current directory. You can move it to a directory in your `$PATH` for easier use:

```sh
mv tok ~/go/bin/
```

Or install directly with:

```sh
go install github.com/0x1Jar/new-hacks/tok@latest
```

Now you can use `tok` from anywhere in your terminal.

## Usage

Pipe text to the tool via standard input:

```
cat file.txt | tok [options]
```
Or run interactively and type/paste input, then press Ctrl+D (EOF) to process.

### Options

- `-min int`  
  Minimum length of string to be output (default: 1)
- `-max int`  
  Maximum length of string to be output (default: 25)
- `-alpha-num-only`  
  Only output strings containing at least one letter and one number
- `-delim-exceptions string`  
  Characters that should not be treated as delimiters

## Examples

### Basic Tokenization
Input:
```
hello world123 test42 foo_bar
```
Command:
```
echo "hello world123 test42 foo_bar" | tok
```
Output:
```
hello
world123
test42
foo
bar
```

### Minimum Length and Alpha-Numeric Only
Input:
```
abc 123 ab12 1a2b3c
```
Command:
```
echo "abc 123 ab12 1a2b3c" | tok -min 4 -alpha-num-only
```
Output:
```
ab12
1a2b3c
```

### Delimiter Exceptions
Input:
```
file-name_v1.0 test.data
```
Command:
```
echo "file-name_v1.0 test.data" | tok -delim-exceptions "-."
```
Output:
```
file-name_v1.0
test.data
```

### URL Decoding
Input:
```
hello%20world test%2Fdata
```
Command:
```
echo "hello%20world test%2Fdata" | tok
```
Output:
```
hello world
test/data
```

### Generating a Wordlist from Waybackurls
Suppose you use [waybackurls](https://github.com/tomnomnom/waybackurls) to gather URLs for a domain, and want to extract unique words or tokens for use in wordlists or further analysis.

Command:
```
echo "https://example.com/login.php?id=123\nhttps://example.com/assets/js/app.js\nhttps://example.com/profile/username" | tok -min 3
```
Output:
```
https
://
example
com
login
php
id
123
assets
js
app
profile
username
```

Or, in a real workflow:
```
waybackurls example.com | tok -min 4 -alpha-num-only
```
This will extract tokens of at least 4 characters, containing both letters and numbers, from all URLs found by waybackurls.

## How it works

`tok` reads input one rune at a time. It groups consecutive letters and numbers into tokens. When a non-letter/number character is encountered (unless it's in the delimiter exceptions), it checks if the current token meets the length and content requirements, then outputs it. If a token contains a `%` character, it attempts to URL-decode it before outputting.

---
