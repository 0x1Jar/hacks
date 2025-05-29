package main

import (
	"bufio"
	"fmt"
	// "io/ioutil" // Replaced with os.WriteFile
	"os"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		path := sc.Text()
		if path == "" {
			continue
		}

		options := jsbeautifier.DefaultOptions()
		// BeautifyFile returns (*string)
		beautified := jsbeautifier.BeautifyFile(path, options)
		
		// Check if beautified is nil, which might indicate an error or empty input.
		// The library itself doesn't seem to return an error from BeautifyFile.
		// It might panic or return nil on error, so a nil check is important.
		if beautified == nil {
			fmt.Fprintf(os.Stderr, "Error beautifying file %s or file is empty/invalid.\n", path)
			continue // Skip to the next file
		}

		// Write the beautified content back to the file
		// os.WriteFile truncates the file if it exists.
		err := os.WriteFile(path, []byte(*beautified), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
	}
}
