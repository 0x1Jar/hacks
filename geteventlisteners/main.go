package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	// "io/ioutil" // No longer needed
	"os"
	"regexp"
	"strconv"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
)

const defaultConcurrency = 5
const defaultTimeout = 20 // seconds
const defaultOutputDir = "eventlisteners_out"

type filterArgs []string

func (f *filterArgs) Set(val string) error {
	*f = append(*f, val)
	return nil
}

func (f filterArgs) String() string {
	return "string"
}

func (f filterArgs) Includes(search string) bool {
	search = strings.ToLower(search)
	for _, filter := range f {
		filter = strings.ToLower(filter)
		if filter == search {
			return true
		}
	}
	return false
}

func main() {
	var filters filterArgs
	flag.Var(&filters, "filter", "")
	flag.Var(&filters, "f", "Event type to filter for (e.g., click, mouseover). Can be used multiple times.")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose mode (prints current URL being processed)")

	var concurrency int
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Number of concurrent browser contexts")

	var timeout int
	flag.IntVar(&timeout, "t", defaultTimeout, "Timeout in seconds for each URL processing")
	
	var outputDir string
	flag.StringVar(&outputDir, "o", defaultOutputDir, "Output directory to save listener files")


	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Extracts JavaScript event listeners from URLs.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [url]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "If no URL is provided as an argument, URLs are read from stdin.\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var input io.Reader
	if flag.NArg() > 0 {
		// If a URL is provided as a command-line argument
		input = strings.NewReader(flag.Arg(0))
	} else {
		// Otherwise, read from stdin
		fmt.Fprintln(os.Stderr, "No URL argument provided, reading URLs from stdin...")
		input = os.Stdin
	}
	sc := bufio.NewScanner(input)

	// Create a cancellable context for the main browser allocation
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
	defer allocCancel()

	// Create a parent context for all browser tabs
	parentBrowserCtx, parentBrowserCancel := chromedp.NewContext(allocCtx)
	defer parentBrowserCancel()

	// It's often better to launch one browser and use multiple tabs (contexts) from it.
	// However, the original code creates a new context from parent for each URL, implying new tabs.
	// For true concurrency with separate browser instances, the allocator would be used per worker.
	// For now, let's stick to one main browser and manage concurrency of tab operations.

	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for requestURL := range jobs {
				if verbose {
					fmt.Printf("Requesting %s\n", requestURL)
				}

				// Create a new tab context from the parent browser context
				ctx, cancel := chromedp.NewContext(parentBrowserCtx)
				ctx, cancelTimeout := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
				
				var res map[string][]string


				err := chromedp.Run(ctx,
					chromedp.Navigate(requestURL),
					chromedp.EvaluateAsDevTools(`
			var listeners = getEventListeners(window)

			for (let i in listeners){
				listeners[i] = listeners[i].map(l => {
					return l.listener.toString()
				})
			}

			listeners`,
				&res),
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", requestURL, err)
			cancelTimeout()
			cancel() // Important to cancel the tab context
			continue
		}

		buf := &strings.Builder{}
		first := true
		for event, listeners := range res {

			if len(filters) > 0 && !filters.Includes(event) {
				continue
			}

			if first {
				fmt.Fprintf(buf, "// %s\n", requestURL)
				buf.WriteString("(function(){")
				first = false
			}

			seen := make(map[string]bool)

			for i, l := range listeners {
				if seen[l] {
					continue
				}
				seen[l] = true

				suffix := strconv.Itoa(i + 1)
				if suffix == "1" {
					suffix = ""
				}

				fmt.Fprintf(buf, "    let on%s%s = %s\n\n", event, suffix, l)
			}
		}

		if first {
			// we didn't find any matching event listeners
		if verbose {
			fmt.Printf("No matching listeners on %s\n", requestURL)
		}
		cancelTimeout()
		cancel()
		continue
	}

	buf.WriteString("})();") // Corrected to valid JS

	raw := buf.String()
	options := jsbeautifier.DefaultOptions()
	beautifiedJs, err := jsbeautifier.Beautify(&raw, options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error beautifying JS for %s: %v. Saving raw.\n", requestURL, err)
		beautifiedJs = raw // Save raw if beautify fails
	}
	
	// Organize files into outputDir/domain/filename.js
	parsedURL, urlParseErr := url.Parse(requestURL)
	if urlParseErr != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL for filename %s: %v. Saving to root of output dir.\n", requestURL, urlParseErr)
		// Fallback filename if URL parsing fails
		fn := genFilename(requestURL)
		filePath := filepath.Join(outputDir, fn)
		
		// Ensure base output directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating base output directory %s: %v\n", outputDir, err)
			cancelTimeout()
			cancel()
			continue
		}
		
		err = os.WriteFile(filePath, []byte(beautifiedJs), 0644) // Changed ioutil.WriteFile to os.WriteFile
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Saved listeners for %s to %s\n", requestURL, filePath)
		}
		cancelTimeout()
		cancel()
		continue
	}

	domain := parsedURL.Hostname()
	domainDir := filepath.Join(outputDir, domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating domain directory %s: %v\n", domainDir, err)
		cancelTimeout()
		cancel()
		continue
	}
	
	filename := genFilename(parsedURL.Path + "?" + parsedURL.RawQuery) // Use path and query for more uniqueness within domain
	if filename == ".js" || filename == "-.js" { // Handle cases where path is empty or just "/"
		filename = "index.js"
	}

	filePath := filepath.Join(domainDir, filename)
	err = os.WriteFile(filePath, []byte(beautifiedJs), 0644) // Changed ioutil.WriteFile to os.WriteFile
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", filePath, err)
	} else {
		fmt.Printf("Saved listeners for %s to %s\n", requestURL, filePath)
	}

	cancelTimeout()
	cancel() // Cancel the tab context
			} // End of for range jobs
		}() // End of worker goroutine
	} // End of worker creation loop


	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)
	wg.Wait()
}

func genFilename(uPath string) string {
	// Sanitize path and query part of the URL to create a filename
	safePath := strings.TrimLeft(uPath, "/") // Remove leading slash if any
	re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	fn := re.ReplaceAllString(safePath, "-")

	re = regexp.MustCompile("-+")
	fn = re.ReplaceAllString(fn, "-")
	fn = strings.Trim(fn, "-") // Trim leading/trailing dashes

	if fn == "" {
		fn = "index" // Default for empty paths
	}
	
	// Truncate if too long
	const maxFilenameLen = 100 
	if len(fn) > maxFilenameLen {
		fn = fn[:maxFilenameLen]
	}

	return fmt.Sprintf("%s.js", fn)
}
