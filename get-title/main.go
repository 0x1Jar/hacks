package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tomnomnom/gahttp"
	"golang.org/x/net/html"
)

func extractTitle(req *http.Request, resp *http.Response, err error) {
	if err != nil {
		return
	}

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.StartTagToken && t.Data == "title" {
			if z.Next() == html.TextToken {
				title := strings.TrimSpace(z.Token().Data)
				fmt.Printf("%s (%s)\n", title, req.URL)
				break
			}
		}

	}
}

func main() {
	// Use a different name for the flag variable to avoid redeclaration if 'concurrency' is needed elsewhere.
	// However, since it's only used to set p.SetConcurrency, we can use it directly.
	var concurrencyLevel int
	flag.IntVar(&concurrencyLevel, "c", 20, "Concurrency level")

	var skipVerifyFlag bool
	flag.BoolVar(&skipVerifyFlag, "k", false, "Skip TLS certificate verification (insecure)")
	flag.BoolVar(&skipVerifyFlag, "skip-verify", false, "Skip TLS certificate verification (insecure) (long form)")


	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Fetch and get the title of HTML pages from URLs provided on stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var p *gahttp.Pipeline
	if skipVerifyFlag {
		client := gahttp.NewClient(gahttp.SkipVerify)
		p = gahttp.NewPipelineWithClient(client)
	} else {
		// For default client (with verification), NewPipeline() is sufficient
		// as it creates a default client internally.
		p = gahttp.NewPipeline()
	}

	p.SetConcurrency(concurrencyLevel)
	extractFn := gahttp.Wrap(extractTitle, gahttp.CloseBody)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		p.Get(sc.Text(), extractFn)
	}
	p.Done()

	p.Wait()

}
