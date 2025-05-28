package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var flip bool
	flag.BoolVar(&flip, "f", false, "")
	flag.BoolVar(&flip, "flip", false, "")

	var separator string
	flag.StringVar(&separator, "s", "", "")
	flag.StringVar(&separator, "separator", "", "")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	prefixFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	suffixFile, err := os.Open(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// use 'a' and 'b' because which is the prefix
	// and which is the suffix depends on if we're in
	// flip mode or not.
	fileA := prefixFile
	fileB := suffixFile

	if flip {
		fileA, fileB = fileB, fileA
	}

	// Read the second file (fileB) into memory to avoid repeated seeks/scans
	var linesB []string
	scannerB := bufio.NewScanner(fileB)
	for scannerB.Scan() {
		linesB = append(linesB, scannerB.Text())
	}
	if err := scannerB.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", flag.Arg(1), err) // Use actual filename based on flip or not
		os.Exit(1)
	}
	fileB.Close() // Close fileB as it's now in memory

	scannerA := bufio.NewScanner(fileA)
	for scannerA.Scan() {
		lineA := scannerA.Text()
		for _, lineB := range linesB {
			if flip {
				// When flipped, fileA was original suffixFile, fileB was original prefixFile
				// So, lineA is a suffix, lineB is a prefix
				fmt.Printf("%s%s%s\n", lineB, separator, lineA)
			} else {
				// Normal mode: lineA is a prefix, lineB is a suffix
				fmt.Printf("%s%s%s\n", lineA, separator, lineB)
			}
		}
	}
	if err := scannerA.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", flag.Arg(0), err) // Use actual filename
		os.Exit(1)
	}
	fileA.Close() // Close fileA
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Combine the lines from two files in every combination\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  comb [OPTIONS] <prefixfile> <suffixfile>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -f, --flip             Flip mode (order by suffix)\n")
		fmt.Fprintf(os.Stderr, "  -s, --separator <str>  String to place between prefix and suffix\n")
	}
}
