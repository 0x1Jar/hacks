package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"unicode"
)

func main() {
	var minlength int
	flag.IntVar(&minlength, "min", 1, "min length of string to be output")
	var maxlength int
	flag.IntVar(&maxlength, "max", 25, "max length of string to be output")
	var alphaNumOnly bool
	flag.BoolVar(&alphaNumOnly, "alpha-num-only", false, "return only strings containing at least one letter and one number")

	var delimExceptions string
	flag.StringVar(&delimExceptions, "delim-exceptions", "", "don't use the characters provided as delimiters")

	flag.Parse()

	r := bufio.NewReader(os.Stdin)
	var out strings.Builder

	maybeURLEncoded := false
	includesLetters := false
	includesNumbers := false
	last := "" // Used for deduplication with alphaNumOnly

	// processToken is a helper function to handle the logic for processing a collected token.
	// It returns the new 'last' printed string.
	processToken := func(tokenContent string, isMaybeURLEncoded, hasLetters, hasNumbers bool, currentLast string) string {
		if len(tokenContent) == 0 {
			return currentLast
		}

		if len(tokenContent) < minlength {
			return currentLast
		}
		if len(tokenContent) > maxlength {
			return currentLast
		}

		if alphaNumOnly && (!hasLetters || !hasNumbers || tokenContent == currentLast) {
			return currentLast
		}

		finalToken := tokenContent
		if isMaybeURLEncoded {
			dec, err := url.QueryUnescape(finalToken)
			if err == nil {
				finalToken = dec
			}
		}

		// Second check for alphaNumOnly if decoding changed the string or for the initial check if not deduplicating yet
		// This handles cases where decoding might make it non-alpha-numeric or identical to last.
		// However, the original logic only deduplicates based on the pre-decoded string if alphaNumOnly is true.
		// To keep behavior similar, we might only re-check length constraints if str changed.
		// For simplicity and robustness, let's assume the primary check on 'tokenContent' is sufficient for alphaNumOnly's includesLetters/Numbers.
		// The 'str == last' check for alphaNumOnly should ideally be on 'finalToken' if 'last' stores finalTokens.
		// The original code compared 'str' (pre-decode) to 'last' (post-decode of previous). This is a bit inconsistent.
		// Let's make 'last' store the pre-decoded string for consistency with the original 'str == last' check.
		// So, the alphaNumOnly check `tokenContent == currentLast` is against pre-decoded.

		fmt.Println(finalToken)
		return tokenContent // Return the pre-decoded string to be set as 'last' if alphaNumOnly
	}

	resetState := func() {
		maybeURLEncoded = false
		includesLetters = false
		includesNumbers = false
		out.Reset()
	}

	for {
		runeChar, _, err := r.ReadRune()
		if err != nil {
			// End of input or error, process any remaining token
			if out.Len() > 0 {
				last = processToken(out.String(), maybeURLEncoded, includesLetters, includesNumbers, last)
			}
			break // Exit loop
		}

		l := unicode.In(runeChar, unicode.L)
		if l {
			includesLetters = true
		}

		n := unicode.In(runeChar, unicode.N)
		if n {
			includesNumbers = true
		}

		// Check if the rune is a delimiter
		isDelimiter := !l && !n && !isDelimException(runeChar, delimExceptions)

		if isDelimiter {
			if out.Len() > 0 { // Process token if buffer is not empty
				last = processToken(out.String(), maybeURLEncoded, includesLetters, includesNumbers, last)
			}
			resetState() // Reset for the next token
			continue     // Skip adding the delimiter to the current token
		}

		// If not a delimiter, add to token
		if runeChar == '%' {
			maybeURLEncoded = true
		}
		out.WriteRune(runeChar)
	}
}

func isDelimException(r rune, delims string) bool {
	for _, comp := range delims {
		if r == comp {
			return true
		}
	}

	return false
}
