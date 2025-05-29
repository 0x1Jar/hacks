package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	// "io/ioutil" // Replaced by io
	"net/http"
	"net/url"
	"os"
	"regexp"
	// "regexp" // Removed duplicate
	"strings"
	"time" // For http.Client timeout
)

const (
	defaultMinLen    = 4
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36" // A generic UA
	defaultTimeout   = 10                                                                                                           // seconds
)

var httpClient *http.Client

func main() {
	var minLen int
	flag.IntVar(&minLen, "min-len", defaultMinLen, "Minimum length of reflected value to report")

	var userAgent string
	flag.StringVar(&userAgent, "ua", defaultUserAgent, "User-Agent string for requests")

	var skipVerify bool
	flag.BoolVar(&skipVerify, "k", false, "Skip TLS certificate verification (insecure)")
	
	var timeoutSeconds int
	flag.IntVar(&timeoutSeconds, "t", defaultTimeout, "Request timeout in seconds")


	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Detects query string values reflected in HTTP response bodies.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [url...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "If no URLs are provided as arguments, URLs are read from stdin.\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
		// Add other transport settings if needed, e.g., DialContext timeouts
	}
	httpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not follow redirects
		},
	}

	var input io.Reader
	if flag.NArg() > 0 {
		input = strings.NewReader(strings.Join(flag.Args(), "\n"))
	} else {
		fmt.Fprintln(os.Stderr, "No URL arguments provided, reading URLs from stdin...")
		input = os.Stdin
	}

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		rawURL := sc.Text()
		if rawURL == "" {
			continue
		}
		u, err := url.Parse(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing URL %s: %v\n", rawURL, err)
			continue
		}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for %s: %v\n", u.String(), err)
			continue
		}
		req.Header.Set("User-Agent", userAgent)

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching %s: %v\n", u.String(), err)
			continue
		}
		
		if resp.Body == nil {
			if resp.Body != nil { // Should always be true if resp.Body is nil, but good practice
				resp.Body.Close()
			}
			fmt.Fprintf(os.Stderr, "Response body is nil for %s\n", u.String())
			continue
		}
		
		b, err := io.ReadAll(resp.Body) // Changed ioutil.ReadAll to io.ReadAll
		resp.Body.Close() // Close body immediately after reading
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading body from %s: %v\n", u.String(), err)
			continue
		}
		body := string(b)

		for k, vv := range u.Query() {
			for _, v := range vv {
				if len(v) < minLen {
					continue
				}

				// A fairly shonky way to get a few chars of context either side of the match
				// but it helps avoid trying to find the locations of all the matches in the
				// body, and then getting the context on either side, with all the bounds
				// checking etc that would need to be done for that.
				re, err := regexp.Compile("(.{0,6}" + regexp.QuoteMeta(v) + ".{0,6})")
				if err != nil {
					fmt.Fprintf(os.Stderr, "regexp compile error: %s", err)
				}

				matches := re.FindAllStringSubmatch(body, -1)

				for _, m := range matches {
					fmt.Printf("%s: '%s=%s' reflected in response body (...%s...)\n", u, k, v, m[0])
				}
			}

		}
	}
}
