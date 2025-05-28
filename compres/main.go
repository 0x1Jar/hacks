package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"flag"
	"os"

	"github.com/aofei/mimesniffer"
	"github.com/nsf/jsondiff"
)

// customHeader for parsing -H "Key: Value" flags
type customHeader map[string]string

func (h *customHeader) String() string {
	return "custom header"
}

func (h *customHeader) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("header must be in format Key:Value")
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	(*h)[key] = val
	return nil
}

func main() {
	baseUrlStr := flag.String("burl", "", "Base URL (required)")
	candUrlStr := flag.String("curl", "", "Candidate URL (required)")

	baseMethod := flag.String("bmethod", "GET", "HTTP method for base URL")
	candMethod := flag.String("cmethod", "GET", "HTTP method for candidate URL")

	var baseHeaders customHeader = make(map[string]string)
	flag.Var(&baseHeaders, "bH", "Header for base URL (e.g., \"User-Agent: test\"). Can be used multiple times.")
	var candHeaders customHeader = make(map[string]string)
	flag.Var(&candHeaders, "cH", "Header for candidate URL. Can be used multiple times.")

	proxyStr := flag.String("proxy", "", "Proxy URL (e.g., http://127.0.0.1:8080)")
	timeoutSec := flag.Int("timeout", 10, "HTTP client timeout in seconds")
	keepAlives := flag.Bool("keepalive", false, "Enable HTTP keep-alives")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -burl <base_url> -curl <candidate_url> [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Compares two HTTP responses and outputs their differences.\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *baseUrlStr == "" || *candUrlStr == "" {
		fmt.Fprintln(os.Stderr, "Error: Base URL (-burl) and Candidate URL (-curl) are required.")
		flag.Usage()
		os.Exit(1)
	}

	mimesniffer.Register("application/json", func(bs []byte) bool {
		return json.Valid(bs)
	})

	baseReq, err := http.NewRequest(strings.ToUpper(*baseMethod), *baseUrlStr, nil)
	fatalErr(err, "creating base request")
	for k, v := range baseHeaders {
		baseReq.Header.Add(k, v)
	}

	candReq, err := http.NewRequest(strings.ToUpper(*candMethod), *candUrlStr, nil)
	fatalErr(err, "creating candidate request")
	for k, v := range candHeaders {
		candReq.Header.Add(k, v)
	}

	client := newHTTPClient(*keepAlives, time.Duration(*timeoutSec), *proxyStr)

	fmt.Printf("Fetching base URL: %s %s\n", baseReq.Method, baseReq.URL)
	baseResp, err := client.Do(baseReq)
	fatalErr(err, "fetching base URL")
	defer baseResp.Body.Close()

	fmt.Printf("Fetching candidate URL: %s %s\n", candReq.Method, candReq.URL)
	candidateResp, err := client.Do(candReq)
	fatalErr(err, "fetching candidate URL")
	defer candidateResp.Body.Close()

	differences, err := compare(baseResp, candidateResp)
	fatalErr(err, "comparing responses")

	if len(differences) == 0 {
		fmt.Println("No differences found.")
		return
	}

	fmt.Println("\nDifferences found:")
	for _, d := range differences {
		fmt.Println(d.String())
	}
}

type diff struct {
	Kind         string `json:"kind"`
	Key          string `json:"key"`
	BaseVal      string `json:"base_value,omitempty"`
	CandidateVal string `json:"candidate_value,omitempty"`
	Description  string `json:"description,omitempty"`
}

// compare function remains largely the same but might need minor adjustments
// if response bodies are closed prematurely by defer in main.
// The defer statements in main will close bodies after compare returns.
func compare(b, c *http.Response) ([]diff, error) {
	out := make([]diff, 0)

	// Status Code
	if b.StatusCode != c.StatusCode { // Corrected from b.Status to b.StatusCode for int comparison
		out = append(out, diff{
			Kind:         "status-code",
			Key:          "code",
			BaseVal:      fmt.Sprintf("%d %s", b.StatusCode, http.StatusText(b.StatusCode)),
			CandidateVal: fmt.Sprintf("%d %s", c.StatusCode, http.StatusText(c.StatusCode)),
		})
	}

	// HTTP Protocol Version
	if b.Proto != c.Proto {
		out = append(out, diff{Kind: "status-line", Key: "protocol", BaseVal: b.Proto, CandidateVal: c.Proto})
	}

	// Headers
	// Check headers present in base but different or missing in candidate
	for name, baseValues := range b.Header {
		candValues, ok := c.Header[name]
		if !ok {
			out = append(out, diff{Kind: "header-missing", Key: name, BaseVal: strings.Join(baseValues, ", ")})
			continue
		}
		if slicesDiffer(baseValues, candValues) {
			out = append(out, diff{
				Kind:         "header-values",
				Key:          name,
				BaseVal:      strings.Join(baseValues, ", "),
				CandidateVal: strings.Join(candValues, ", "),
			})
		}
	}
	// Check headers present in candidate but missing in base
	for name, candValues := range c.Header {
		if _, ok := b.Header[name]; !ok {
			out = append(out, diff{Kind: "header-added", Key: name, CandidateVal: strings.Join(candValues, ", ")})
		}
	}

	// Body comparisons
	bBody, errB := ioutil.ReadAll(b.Body)
	if errB != nil {
		return out, fmt.Errorf("reading base body: %w", errB)
	}
	cBody, errC := ioutil.ReadAll(c.Body)
	if errC != nil {
		return out, fmt.Errorf("reading candidate body: %w", errC)
	}

	// Length difference
	if len(bBody) != len(cBody) {
		out = append(out, diff{
			Kind:         "body",
			Key:          "length",
			BaseVal:      fmt.Sprintf("%d", len(bBody)),
			CandidateVal: fmt.Sprintf("%d", len(cBody)),
		})
	}

	// Hash difference
	bHashStr := fmt.Sprintf("%x", sha256.Sum256(bBody))
	cHashStr := fmt.Sprintf("%x", sha256.Sum256(cBody))

	if bHashStr != cHashStr {
		out = append(out, diff{Kind: "body", Key: "hash", BaseVal: bHashStr, CandidateVal: cHashStr})
	}

	// MIME sniff
	bMIME := mimesniffer.Sniff(bBody)
	cMIME := mimesniffer.Sniff(cBody)
	if bMIME != cMIME {
		out = append(out, diff{Kind: "body", Key: "mime", BaseVal: bMIME, CandidateVal: cMIME})
	}

	// JSON diff if both are JSON
	if bMIME == "application/json" && cMIME == "application/json" {
		opts := jsondiff.DefaultJSONOptions()
		opts.SkipMatches = true // Report only differences
		difference, desc := jsondiff.Compare(bBody, cBody, &opts)

		// jsondiff.Compare returns a description of the first difference.
		// For a more comprehensive diff, one might need to parse and walk the JSON.
		// This provides a basic "are they different and how"
		if difference != jsondiff.FullMatch && difference != jsondiff.SupersetMatch { // SupersetMatch means C is a superset of B
			out = append(out, diff{
				Kind:        "body-json-diff",
				Key:         "content",
				Description: desc, // The description from jsondiff
			})
		}
	} else if bMIME == "application/json" && cMIME != "application/json" {
		out = append(out, diff{Kind: "body-type-mismatch", Key: "mime", BaseVal: bMIME, CandidateVal: cMIME, Description: "Base is JSON, Candidate is not."})
	} else if bMIME != "application/json" && cMIME == "application/json" {
		out = append(out, diff{Kind: "body-type-mismatch", Key: "mime", BaseVal: bMIME, CandidateVal: cMIME, Description: "Candidate is JSON, Base is not."})
	}
	// Consider adding other content-type specific diffs here (e.g., HTML, XML)

	// Keywords
	// This keyword counting can be noisy if bodies are large and different.
	// Consider making this optional or more sophisticated.
	keywords := []string{"error", "warn", "debug", "exception", "failed", "trace"}
	bBodyLower := strings.ToLower(string(bBody))
	cBodyLower := strings.ToLower(string(cBody))

	for _, k := range keywords {
		bCount := strings.Count(bBodyLower, k)
		cCount := strings.Count(cBodyLower, k)
		if bCount != cCount {
			out = append(out, diff{
				Kind:         "keyword-count",
				Key:          k,
				BaseVal:      fmt.Sprintf("%d", bCount),
				CandidateVal: fmt.Sprintf("%d", cCount),
			})
		}
	}
	return out, nil
}

func slicesDiffer(a, b []string) bool { // Simplified for string slices (headers)
	if len(a) != len(b) {
		return true
	}
	// Order might matter for headers with multiple identical values, but for simplicity,
	// we'll consider them different if counts of unique values differ or values themselves differ.
	// This basic check assumes order matters and values are directly comparable.
	for i := range a {
		if a[i] != b[i] {
			return true
		}
	}
	return false
}

// String representation of a diff
func (d diff) String() string {
	if d.Description != "" {
		return fmt.Sprintf("[%s] %s: %s", d.Kind, d.Key, d.Description)
	}
	return fmt.Sprintf("[%s] %s: Base=\"%s\", Candidate=\"%s\"", d.Kind, d.Key, d.BaseVal, d.CandidateVal)
}

func fatalErr(err error, context ...string) {
	if err != nil {
		if len(context) > 0 {
			log.Fatalf("Error %s: %v", strings.Join(context, " "), err)
		}
		log.Fatalf("Error: %v", err)
	}
}

func newHTTPClient(keepAlives bool, timeout time.Duration, proxyURLStr string) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:      100, // Increased for potentially more connections
		IdleConnTimeout:   90 * time.Second,
		DisableKeepAlives: !keepAlives,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // Common for security tools, make configurable if sensitive
		DialContext: (&net.Dialer{
			Timeout:   timeout, // Use the passed timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	if proxyURLStr != "" {
		proxyURL, err := url.Parse(proxyURLStr)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyURL)
		} else {
			log.Printf("Warning: Invalid proxy URL '%s': %v", proxyURLStr, err)
		}
	}

	// Prevent following redirects automatically, as we want to compare the direct response.
	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: redirectPolicy,
		Timeout:       timeout + (5 * time.Second), // Client timeout slightly larger than dialer timeout
	}
}
