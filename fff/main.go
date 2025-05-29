package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultConcurrency = 20

func init() {
	flag.Usage = func() {
		h := []string{
			"Request URLs provided on stdin fairly frickin' fast",
			"",
			"Options:",
			"  -c, --concurrency <num>   Number of concurrent requests (default: 20)",
			"  -b, --body <data>         Request body",
			"  -d, --delay <delay>       Delay between issuing requests (ms) (applied per worker, not globally before each request)",
			"  -H, --header <header>     Add a header to the request (can be specified multiple times)",
			"  -k, --keep-alive          Use HTTP Keep-Alive",
			"  -m, --method              HTTP method to use (default: GET, or POST if body is specified)",
			"  -o, --output <dir>        Directory to save responses in (will be created, default: out)",
			"  -s, --save-status <code>  Save responses with given status code (can be specified multiple times)",
			"  -S, --save                Save all responses",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "concurrency", defaultConcurrency, "")
	flag.IntVar(&concurrency, "c", defaultConcurrency, "")

	var body string
	flag.StringVar(&body, "body", "", "")
	flag.StringVar(&body, "b", "", "")

	var keepAlives bool
	flag.BoolVar(&keepAlives, "keep-alive", false, "")
	flag.BoolVar(&keepAlives, "k", false, "") // Removed duplicate -keep-alives

	var saveResponses bool
	flag.BoolVar(&saveResponses, "save", false, "")
	flag.BoolVar(&saveResponses, "S", false, "")

	var delayMs int
	flag.IntVar(&delayMs, "delay", 0, "") // Default delay 0, applied per worker if > 0
	flag.IntVar(&delayMs, "d", 0, "")

	var method string
	flag.StringVar(&method, "method", "GET", "") // Default GET, auto POST if body
	flag.StringVar(&method, "m", "GET", "")

	var outputDir string
	flag.StringVar(&outputDir, "output", "out", "")
	flag.StringVar(&outputDir, "o", "out", "")

	var headers headerArgs
	flag.Var(&headers, "header", "")
	flag.Var(&headers, "H", "")

	var saveStatus saveStatusArgs
	flag.Var(&saveStatus, "save-status", "")
	flag.Var(&saveStatus, "s", "")

	flag.Parse()

	if body != "" && method == "GET" { // Auto-set method to POST if body is provided and method is still GET
		method = "POST"
	}

	delay := time.Duration(delayMs) * time.Millisecond // Corrected delay to use Millisecond
	client := newClient(keepAlives, concurrency) // Pass concurrency to newClient
	// prefix := outputDir // outputDir is used directly now

	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rawURL := range jobs {
				if delay > 0 {
					time.Sleep(delay)
				}

				// create the request
				var reqBody io.Reader
				if body != "" {
					reqBody = strings.NewReader(body)
				}
				req, err := http.NewRequest(method, rawURL, reqBody)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create request: %s\n", err)
				return
			}

			// add headers to the request
			for _, h := range headers {
				parts := strings.SplitN(h, ":", 2)

				if len(parts) != 2 {
					continue
				}
				req.Header.Set(parts[0], parts[1])
			}

			// send the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "request failed: %s\n", err)
				return
			}
			defer resp.Body.Close()

			shouldSave := saveResponses || len(saveStatus) > 0 && saveStatus.Includes(resp.StatusCode)

			if !shouldSave {
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				fmt.Printf("%s %d\n", rawURL, resp.StatusCode)
				return
			}

			// output files are stored in outputDir/domain/normalisedpath/hash.(body|headers)
			normalisedPath := normalisePath(req.URL)
			// Hash should ideally include method, URL, body, and specific headers for uniqueness if they vary per request.
			// For simplicity, if headers are global, this is okay.
			hash := sha1.Sum([]byte(req.Method + req.URL.String() + body + headers.String()))
			p := path.Join(outputDir, req.URL.Hostname(), normalisedPath, fmt.Sprintf("%x.body", hash))
			err = os.MkdirAll(path.Dir(p), 0750)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create dir: %s\n", err)
				return
			}

			// create the body file
			f, err := os.Create(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create file: %s\n", err)
				return
			}
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write file contents: %s\n", err)
				return
			}

			// create the headers file
			headersPath := path.Join(outputDir, req.URL.Hostname(), normalisedPath, fmt.Sprintf("%x.headers", hash))
			headersFile, err := os.Create(headersPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create file: %s\n", err)
				return
			}
			defer headersFile.Close()

			var buf strings.Builder

			// put the request URL and method at the top
			buf.WriteString(fmt.Sprintf("%s %s\n\n", method, rawURL))

			// add the request headers
			for _, h := range headers {
				buf.WriteString(fmt.Sprintf("> %s\n", h))
			}
			buf.WriteRune('\n')

			// add the request body (if any)
			// For GET/HEAD etc. body is nil. For POST/PUT it's present.
			// The current 'body' variable is global. If body could change per request, this needs adjustment.
			// For now, assuming 'body' is a global static body for all POST/PUT.
			if req.Body != nil && body != "" { // Check if original body was set
				buf.WriteString(body) // Write the original string body
				buf.WriteString("\n\n")
			}


			// add the proto and status
			buf.WriteString(fmt.Sprintf("< %s %s\n", resp.Proto, resp.Status))

			// add the response headers
			for k, vs := range resp.Header {
				for _, v := range vs {
					buf.WriteString(fmt.Sprintf("< %s: %s\n", k, v))
				}
			}
			// No need to write response body to headers file. It's in the .body file.

			_, err = headersFile.WriteString(buf.String()) // Use WriteString for efficiency
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write headers file contents: %s\n", err)
				return
			}

			// output the body filename for each URL
			fmt.Printf("%s: %s %d\n", p, rawURL, resp.StatusCode)

			} // End of worker goroutine scope for 'rawURL'
		}() // Removed passing rawURL, it's captured by the closure
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)
	wg.Wait()
}

func newClient(keepAlives bool, numWorkers int) *http.Client { // Added numWorkers parameter
	// TODO: Make timeout and other transport settings configurable via flags
	tr := &http.Transport{
		MaxIdleConns:      numWorkers * 2, // Use passed numWorkers
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: !keepAlives,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // Consider making this a flag
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Default dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // Do not follow redirects
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: redirectPolicy,
		Timeout:       20 * time.Second, // Overall request timeout
	}
}

type headerArgs []string

func (h *headerArgs) Set(val string) error {
	*h = append(*h, val)
	return nil
}

func (h headerArgs) String() string {
	return strings.Join(h, ", ")
}

type saveStatusArgs []int

func (s *saveStatusArgs) Set(val string) error {
	i, _ := strconv.Atoi(val)
	*s = append(*s, i)
	return nil
}

func (s saveStatusArgs) String() string {
	return "string"
}

func (s saveStatusArgs) Includes(search int) bool {
	for _, status := range s {
		if status == search {
			return true
		}
	}
	return false
}

func normalisePath(u *url.URL) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9/._-]+`)
	return re.ReplaceAllString(u.Path, "-")
}
