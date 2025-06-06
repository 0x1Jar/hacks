package main

import (
	"flag"
	"fmt"
	// "io/ioutil" // Replaced by os.ReadDir
	// "io/fs" // Not strictly needed as os.DirEntry implements fs.DirEntry
	"log"
	"math"
	"os" // For os.ReadDir
	"path"
	"strconv"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Lists files in a directory whose sizes are unique enough.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <directory> [threshold_percentage]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  <directory>: The directory to scan.\n")
		fmt.Fprintf(os.Stderr, "  [threshold_percentage]: (Optional) Percentage difference to consider sizes similar (default: 1.0).\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: No directory specified.")
		flag.Usage()
		os.Exit(1)
	}
	dir := flag.Arg(0)

	threshold := 1.0
	if flag.NArg() > 1 {
		t, err := strconv.ParseFloat(flag.Arg(1), 64)
		if err == nil {
			threshold = t
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Invalid threshold '%s', using default 1.0%%.\n", flag.Arg(1))
		}
	}

	contents, err := os.ReadDir(dir) // Changed from ioutil.ReadDir
	if err != nil {
		log.Fatalf("Failed to read dir %s: %s", dir, err)
	}

	sizes := make([]int64, 0)

	for _, entry := range contents { // entry is fs.DirEntry
		if entry.IsDir() {
			continue
		}

		fileInfo, err := entry.Info() // Get fs.FileInfo from fs.DirEntry
		if err != nil {
			log.Printf("Failed to get info for %s: %v. Skipping.", entry.Name(), err)
			continue
		}
		
		fileSize := fileInfo.Size()

		isDifferent := true
		if fileSize == 0 { // Handle zero-byte files separately or consider them different
			// Depending on desired behavior, zero-byte files might always be "different"
			// or grouped. For now, let's assume they are distinct unless another zero-byte file is seen.
			for _, s := range sizes {
				if s == 0 {
					isDifferent = false
					break
				}
			}
		} else {
			for _, s := range sizes {
				if s == 0 { continue } // Avoid division by zero if a previous size was 0
				diff := math.Abs((float64(s-fileSize) / float64(s)) * 100)
				if diff < threshold {
					isDifferent = false
					break // Found a similar enough size
				}
			}
		}

		if isDifferent {
			sizes = append(sizes, fileSize)
			fmt.Println(path.Join(dir, entry.Name()))
		}
	}
}
