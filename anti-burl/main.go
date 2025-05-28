package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	client    *http.Client
	transport *http.Transport
	wg        sync.WaitGroup

	// Default values
	concurrencyDefault = 50
	timeoutDefault     = 5 * time.Second
	maxSizeDefault     = int64(1024000)
	userAgentDefault   = "burl/0.1"
)

// Global variables to hold flag values
var (
	concurrency  int
	timeout      time.Duration
	maxSize      int64
	insecureSkip bool
	userAgent    string
)

func main() {
	// Define flags
	flag.IntVar(&concurrency, "c", concurrencyDefault, "Set the concurrency level")
	flag.DurationVar(&timeout, "t", timeoutDefault, "Set the request timeout (e.g., 5s, 1m)")
	flag.Int64Var(&maxSize, "ms", maxSizeDefault, "Set the maximum response body size to read (bytes)")
	flag.BoolVar(&insecureSkip, "k", true, "Skip TLS certificate verification")
	flag.StringVar(&userAgent, "ua", userAgentDefault, "Set the User-Agent string")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] [file]:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Reads URLs from stdin or a file, fetches them, and prints details for those returning 200 OK.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf no file is specified, input is read from stdin.\n")
		fmt.Fprintf(os.Stderr, "Output format: <status_code> <content_length_runes> <word_count> <url>\n")
	}

	flag.Parse()

	var input io.Reader
	input = os.Stdin

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
		}
		input = file
	}

	sc := bufio.NewScanner(input)

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        concurrency, // Use the flag value
			MaxIdleConnsPerHost: concurrency, // Use the flag value
			MaxConnsPerHost:     concurrency, // Use the flag value
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkip, // Use the flag value
			},
		},
		Timeout: timeout, // Use the flag value
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	semaphore := make(chan bool, concurrency)

	for sc.Scan() {
		raw := sc.Text()
		wg.Add(1)
		semaphore <- true // Acquire a slot
		go func(raw string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the slot when done

			u, err := url.ParseRequestURI(raw)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing URL %s: %v\n", raw, err)
				return
			}
			resp, ws, err := fetchURL(u)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching URL %s: %v\n", u.String(), err)
				return
			}
			// Only print if status code is 200 OK, as per README
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("%-3d %-9d %-5d %s\n", resp.StatusCode, resp.ContentLength, ws, u.String())
			}
		}(raw)
	}

	wg.Wait()

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}
}

func fetchURL(u *url.URL) (*http.Response, int, error) {
	wordsSize := 0

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", userAgent) // Use the flag value

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	var respbody []byte
	readBodyError := true // Assume error initially

	if resp.ContentLength <= maxSize {
		var readErr error
		respbody, readErr = ioutil.ReadAll(resp.Body)
		if readErr == nil {
			readBodyError = false // Successfully read the body
			resp.ContentLength = int64(utf8.RuneCountInString(string(respbody)))
			wordsSize = len(strings.Split(string(respbody), " "))
		}
		// If readErr is not nil, we'll fall through to io.Copy below
	}

	// If we haven't successfully read the body (either too large or ReadAll failed),
	// ensure the body is consumed to allow connection reuse.
	if readBodyError {
		io.Copy(ioutil.Discard, resp.Body)
	}

	return resp, wordsSize, err
}
