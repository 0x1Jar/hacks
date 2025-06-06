package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync" // Added for WaitGroup
	"time"
)

const fetchURL = "http://web.archive.org/cdx/search/cdx?url=%s/*&output=json&fl=original&collapse=urlkey"

var (
	concurrency int
	verbose     bool
	httpClient  = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func init() {
	flag.IntVar(&concurrency, "c", 10, "Number of concurrent requests")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
}

func main() {
	flag.Parse()

	var domains []string

	if flag.NArg() > 0 {
		// fetch for a single domain
		domains = []string{flag.Arg(0)}
	} else {
		// fetch for all domains from stdin
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			domains = append(domains, sc.Text())
		}

		if err := sc.Err(); err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
	}

	if len(domains) == 0 {
		if verbose {
			log.Println("No domains provided. Exiting.")
		}
		return
	}

	// Create a buffered channel to act as a worker queue
	domainCh := make(chan string, concurrency)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainCh {
				if verbose {
					log.Printf("Fetching URLs for domain: %s\n", domain)
				}
				urls, err := getWaybackURLs(domain)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to fetch URLs for [%s]: %v\n", domain, err)
					continue
				}

				for _, url := range urls {
					fmt.Println(url)
				}
			}
		}()
	}

	// Send domains to the channel
	for _, domain := range domains {
		domainCh <- domain
	}
	close(domainCh) // Close the channel once all domains are sent

	wg.Wait() // Wait for all goroutines to finish
}

func getWaybackURLs(domain string) ([]string, error) {
	out := make([]string, 0)
	requestURL := fmt.Sprintf(fetchURL, domain)

	res, err := httpClient.Get(requestURL) // Use shared httpClient
	if err != nil {
		return out, fmt.Errorf("requesting %s: %w", requestURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return out, fmt.Errorf("received non-200 status code (%d) from %s", res.StatusCode, requestURL)
	}

	var wrapper [][]string
	// Use json.NewDecoder for potentially large responses
	if err := json.NewDecoder(res.Body).Decode(&wrapper); err != nil {
		return out, fmt.Errorf("decoding JSON from %s: %w", requestURL, err)
	}

	if len(wrapper) > 0 {
		// Check if the first row is the header "original" and skip it.
		// This is more robust than just skipping the first row unconditionally.
		if len(wrapper[0]) == 1 && wrapper[0][0] == "original" {
			wrapper = wrapper[1:]
		}
	}
	
	for _, urls := range wrapper {
		// Each 'urls' here is expected to be a slice containing a single URL string.
		if len(urls) > 0 {
			out = append(out, urls[0])
		}
	}

	return out, nil
}
