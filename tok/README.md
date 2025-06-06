# tok

`tok` is a command-line tool for extracting tokens (words, numbers, or alphanumeric strings) from standard input. It reads input character by character, groups sequences of letters and numbers, and outputs them according to customizable rules.

## Features

- Extracts tokens containing letters and/or numbers from input.
- Supports minimum and maximum token length filtering.
- Optionally only outputs tokens containing both letters and numbers.
- Allows specifying delimiter exceptions (characters that should not be treated as token boundaries).
- Detects and decodes URL-encoded tokens.

## Usage

```
cat file.txt | tok [options]
```

### Options

- `-min int`  
  Minimum length of string to be output (default: 1)

- `-max int`  
  Maximum length of string to be output (default: 25)

- `-alpha-num-only`  
  Only output strings containing at least one letter and one number

- `-delim-exceptions string`  
  Characters that should not be treated as delimiters

## Example

Extract tokens of at least 5 characters, containing both letters and numbers:

```
cat input.txt | tok -min 5 -alpha-num-only
```

## How it works

`tok` reads input one rune at a time. It groups consecutive letters and numbers into tokens. When a non-letter/number character is encountered (unless it's in the delimiter exceptions), it checks if the current token meets the length and content requirements, then outputs it. If a token contains a `%` character, it attempts to URL-decode it before outputting.

---
