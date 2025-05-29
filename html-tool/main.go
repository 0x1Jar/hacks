package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ericchiang/css"
	"github.com/tomnomnom/gahttp"
	"golang.org/x/net/html"
)

const defaultConcurrency = 20

var (
	showSource   bool
	skipVerify   bool
	concurrency  int
)

func extractSelector(r io.Reader, selector string) ([]string, error) {
	out := []string{}
	sel, err := css.Parse(selector)
	if err != nil {
		return out, err
	}

	node, err := html.Parse(r)
	if err != nil {
		return out, err
	}

	// it's kind of tricky to actually know what to output
	// if the resulting tags contain more than just a text node
	for _, ele := range sel.Select(node) {
		if ele.FirstChild == nil {
			continue
		}
		out = append(out, ele.FirstChild.Data)
	}

	return out, nil
}

func extractComments(r io.Reader) []string {

	z := html.NewTokenizer(r)

	out := []string{}
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.CommentToken {
			d := strings.Replace(t.Data, "\n", " ", -1)
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			out = append(out, d)
		}

	}
	return out
}

func extractAttribs(r io.Reader, attribs []string) []string {
	z := html.NewTokenizer(r)

	out := []string{}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		for _, a := range t.Attr {

			if a.Val == "" {
				continue
			}

			for _, attrib := range attribs {
				if attrib == a.Key {
					out = append(out, a.Val)
				}
			}
		}
	}
	return out
}

func extractTags(r io.Reader, tags []string) []string {
	z := html.NewTokenizer(r)

	out := []string{}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.StartTagToken {

			for _, tag := range tags {
				if t.Data == tag {
					if z.Next() == html.TextToken {
						text := strings.TrimSpace(z.Token().Data)
						if text == "" {
							continue
						}
						out = append(out, text)
					}
				}
			}
		}
	}
	return out
}

type target struct {
	location string
	r        io.ReadCloser
}

func main() {
	flag.BoolVar(&showSource, "s", false, "Prepend each output line with the source URL or filename")
	flag.BoolVar(&showSource, "source", false, "Prepend each output line with the source URL or filename (long form)")
	flag.BoolVar(&skipVerify, "skip-verify", false, "Skip TLS certificate verification for HTTPS URLs")
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Set the concurrency level for fetching URLs")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Accept URLs or filenames for HTML documents on stdin and extract parts of them.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: html-tool [global_options] <mode> [<mode_args>]\n\n")
		fmt.Fprintf(os.Stderr, "Global Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nModes:\n")
		fmt.Fprintf(os.Stderr, "  tags <tag-name1> [<tag-name2> ...]    Extract text contained in specified tags\n")
		fmt.Fprintf(os.Stderr, "  attribs <attr-name1> [<attr-name2> ...] Extract values of specified attributes\n")
		fmt.Fprintf(os.Stderr, "  comments                             Extract all HTML comments\n")
		fmt.Fprintf(os.Stderr, "  query <css-selector>                 Extract text from elements matching CSS selector\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  cat urls.txt | html-tool tags title h1\n")
		fmt.Fprintf(os.Stderr, "  find . -name '*.html' | html-tool -s attribs href src\n")
		fmt.Fprintf(os.Stderr, "  echo \"http://example.com\" | html-tool -skip-verify -c 10 comments\n")
		fmt.Fprintf(os.Stderr, "  echo \"http://example.com\" | html-tool query \"div.main > p\"\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: Mode not specified.")
		flag.Usage()
		os.Exit(1)
	}

	mode := flag.Arg(0)
	var modeArgs []string
	if flag.NArg() > 1 {
		modeArgs = flag.Args()[1:]
	}

	// Validate mode and modeArgs
	switch mode {
	case "tags", "attribs":
		if len(modeArgs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Mode '%s' requires at least one argument (tag/attribute name).\n", mode)
			flag.Usage()
			os.Exit(1)
		}
	case "comments":
		if len(modeArgs) > 0 {
			fmt.Fprintf(os.Stderr, "Error: Mode 'comments' does not take arguments.\n")
			flag.Usage()
			os.Exit(1)
		}
	case "query":
		if len(modeArgs) != 1 {
			fmt.Fprintf(os.Stderr, "Error: Mode 'query' requires exactly one argument (CSS selector).\n")
			flag.Usage()
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported mode '%s'.\n", mode)
		flag.Usage()
		os.Exit(1)
	}

	targets := make(chan *target)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ { // Use a pool of workers for processing targets
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range targets {
				var vals []string
				var err error

				switch mode {
				case "tags":
					vals = extractTags(t.r, modeArgs)
				case "attribs":
					vals = extractAttribs(t.r, modeArgs)
				case "comments":
					vals = extractComments(t.r)
				case "query":
					vals, err = extractSelector(t.r, modeArgs[0]) // modeArgs[0] is the selector
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error processing selector for %s: %v\n", t.location, err)
						t.r.Close()
						continue
					}
				}
				t.r.Close() // Close the reader as soon as processing is done

				for _, v := range vals {
					if showSource {
						fmt.Printf("%s %s\n", t.location, v)
					} else {
						fmt.Println(v)
					}
				}
			}
		}()
	}

	p := gahttp.NewPipeline()
	if skipVerify {
		client := gahttp.NewClient(gahttp.SkipVerify)
		p.SetClient(client)
	}
	// If not skipVerify, the pipeline uses its default client.
	p.SetConcurrency(concurrency)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		// location can be a filename or a URL
		location := strings.TrimSpace(sc.Text())

		// if it's a URL request it with gahttp
		nl := strings.ToLower(location)
		if strings.HasPrefix(nl, "http:") || strings.HasPrefix(nl, "https:") {
			p.Get(location, func(req *http.Request, resp *http.Response, err error) {
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to fetch URL: %s\n", err)
				}
				if resp != nil && resp.Body != nil {
					targets <- &target{req.URL.String(), resp.Body}
				}
			})
			continue
		}

		// if it's a file just open it
		f, err := os.Open(location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
			continue
		}

		targets <- &target{location, f}
	}
	p.Done()
	p.Wait()

	close(targets)
	wg.Wait()
}
