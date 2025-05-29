package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io" // For io.ReadAll
	// "io/ioutil" // Replaced by io
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultConcurrency = 20
const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36"

type paramCheck struct {
	url   string
	param string
}

var transport = &http.Transport{ // Keep as global or pass to newClient
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO: Make this configurable
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, // TODO: Make this configurable
		KeepAlive: time.Second,      // TODO: Make this configurable
		DualStack: true,
	}).DialContext,
}

var httpClient = &http.Client{ // Keep as global or pass to newClient
	Transport: transport,
}
var userAgent string // Global to be set by flag

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Number of concurrent workers per stage")
	flag.StringVar(&userAgent, "ua", defaultUserAgent, "User-Agent string for requests")
	// TODO: Add flags for timeouts, TLS skip verify, etc.

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "kxss - XSS parameter reflection checker.\n\n")
		fmt.Fprintf(os.Stderr, "Reads URLs from stdin and checks for reflected parameters.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()


	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	sc := bufio.NewScanner(os.Stdin)

	initialChecks := make(chan paramCheck, concurrency*2) // Buffer size related to concurrency

	appendChecks := makePool(initialChecks, concurrency, func(c paramCheck, output chan paramCheck) {
		reflected, err := checkReflected(c.url)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "Error from checkReflected for %s: %v\n", c.url, err)
			return
		}

		// if len(reflected) == 0 {
		// 	// This can be very verbose if many URLs don't reflect anything.
		// 	// fmt.Printf("No params initially reflected in %s\n", c.url)
		// 	return
		// }

		for _, param := range reflected {
			output <- paramCheck{c.url, param}
		}
	})

	charChecks := makePool(appendChecks, concurrency, func(c paramCheck, output chan paramCheck) {
		// Using a more unique random-like string for append check
		wasReflected, err := checkAppend(c.url, c.param, "kXssRand0mStr1ng")
		if err != nil {
			// fmt.Fprintf(os.Stderr, "Error from checkAppend for url %s with param %s: %v\n", c.url, c.param, err)
			return
		}

		if wasReflected {
			// fmt.Printf("Confirmed reflection of param %s on %s after append check\n", c.param, c.url) // Verbose
			output <- paramCheck{c.url, c.param}
		}
	})

	done := makePool(charChecks, concurrency, func(c paramCheck, output chan paramCheck) {
		// Test a common set of XSS probe characters
		// TODO: Make this list configurable
		probeChars := []string{"\"", "'", "<", ">", "(", ")", "`", ";", "{", "}"}
		for _, char := range probeChars {
			testPayload := "kXssT3st" + char + "P4yL0ad" // Prefix/suffix to make it more unique
			wasReflected, err := checkAppend(c.url, c.param, testPayload)
			if err != nil {
				// fmt.Fprintf(os.Stderr, "Error from checkAppend for url %s with param %s with char '%s': %v\n", c.url, c.param, char, err)
				continue
			}

			if wasReflected {
				fmt.Printf("param %s is reflected and allows %s on %s\n", c.param, char, c.url)
			}
		}
	})

	for sc.Scan() {
		initialChecks <- paramCheck{url: sc.Text()}
	}

	close(initialChecks)
	<-done
}

func checkReflected(targetURL string) ([]string, error) {

	out := make([]string, 0)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return out, err
	}

	// temporary. Needs to be an option
	req.Header.Set("User-Agent", userAgent) // Use the configurable User-Agent

	resp, err := httpClient.Do(req)
	if err != nil {
		return out, err
	}
	if resp.Body == nil {
		// Ensure body is closed even if it's nil to prevent resource leaks with some HTTP client behaviors
		if resp.Body != nil {
			resp.Body.Close()
		}
		return out, fmt.Errorf("response body is nil")
	}
	defer resp.Body.Close()

	// always read the full body so we can re-use the tcp connection
	b, err := io.ReadAll(resp.Body) // Changed ioutil.ReadAll to io.ReadAll
	if err != nil {
		return out, err
	}

	// nope (:
	if strings.HasPrefix(resp.Status, "3") {
		return out, nil
	}

	// also nope
	ct := resp.Header.Get("Content-Type")
	if ct != "" && !strings.Contains(ct, "html") {
		return out, nil
	}

	body := string(b)

	u, err := url.Parse(targetURL)
	if err != nil {
		return out, err
	}

	for key, vv := range u.Query() {
		for _, v := range vv {
			if !strings.Contains(body, v) {
				continue
			}

			out = append(out, key)
		}
	}

	return out, nil
}

func checkAppend(targetURL, param, suffix string) (bool, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	qs := u.Query()
	val := qs.Get(param)
	//if val == "" {
	//return false, nil
	//return false, fmt.Errorf("can't append to non-existant param %s", param)
	//}

	qs.Set(param, val+suffix)
	u.RawQuery = qs.Encode()

	reflected, err := checkReflected(u.String())
	if err != nil {
		return false, err
	}

	for _, r := range reflected {
		if r == param {
			return true, nil
		}
	}

	return false, nil
}

type workerFunc func(paramCheck, chan paramCheck)

func makePool(input chan paramCheck, poolSize int, fn workerFunc) chan paramCheck {
	var wg sync.WaitGroup
	output := make(chan paramCheck)

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range input {
				fn(item, output)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
