# unisub - Unicode Substitution Finder

`unisub` is a Go tool that helps find Unicode characters that might be confusable with or act as substitutions for a given input character. It checks a predefined list of "fallback" translations and also searches the Unicode range for characters that match the input character when case-folded (toLower or toUpper).

## How it Works

1.  **Input**: Takes a single character as a command-line argument. If multiple characters are provided, only the first one is used.
2.  **Fallback Translations**:
    *   It first consults an internal map (`translations.go`) for predefined "fallback" characters related to the input character.
    *   If found, these are printed with their Unicode code point and URL-escaped form.
3.  **Case-Folding Search**:
    *   The tool then iterates through Unicode code points from `U+0001` to `U+10FFFF`.
    *   For each character in this range (excluding the input character itself):
        *   It checks if `strings.ToLower(string(char_from_range)) == input_char_to_match`.
        *   It checks if `strings.ToUpper(string(char_from_range)) == input_char_to_match`.
    *   Matching characters found via case-folding are printed with their details.
4.  **Output**: Results are printed to standard output, indicating the type of match ("fallback", "toLower", "toUpper"), the character itself, its Unicode code point, and its URL-escaped representation.

**Note on Performance**: The search for `toLower`/`toUpper` matches involves iterating over a very large range of Unicode code points, which can be slow.

## Installation

Ensure you have Go (version 1.24.3 or later, as specified in `go.mod`) installed.

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/0x1Jar/new-hacks.git
    cd new-hacks/unisub
    ```

2.  **Build the tool:**
    ```bash
    go build
    ```
    This will create a `unisub` executable in the current directory.

Alternatively, you can install it directly if your Go environment is set up:
```bash
go install github.com/0x1Jar/new-hacks/unisub
```

## Usage

```
▶ unisub '@'
＠ U+FF20 %EF%BC%A0
﹫ U+FE6B %EF%B9%AB
```

```
▶ unisub '<'
＜ U+FF1C %EF%BC%9C
﹤ U+FE64 %EF%B9%A4
```

```
▶ unisub 's'
ｓ U+FF53 %EF%BD%93
𝐬 U+1D42C %F0%9D%90%AC
𝑠 U+1D460 %F0%9D%91%A0
𝒔 U+1D494 %F0%9D%92%94
𝓈 U+1D4C8 %F0%9D%93%88
𝓼 U+1D4FC %F0%9D%93%BC
𝔰 U+1D530 %F0%9D%94%B0
𝕤 U+1D564 %F0%9D%95%A4
𝖘 U+1D598 %F0%9D%96%98
𝗌 U+1D5CC %F0%9D%97%8C
𝘀 U+1D600 %F0%9D%98%80
𝘴 U+1D634 %F0%9D%98%B4
𝙨 U+1D668 %F0%9D%99%A8
𝚜 U+1D69C %F0%9D%9A%9C
ⓢ U+24E2 %E2%93%A2
ˢ U+02E2 %CB%A2
ₛ U+209B %E2%82%9B
ſ U+017F %C5%BF
```
