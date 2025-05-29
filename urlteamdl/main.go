package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os" // Added os import
	"time"
)

const searchURL = "https://archive.org/advancedsearch.php"
const metaURL = "http://archive.org/metadata/%s"

type file struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func main() {
	flag.Parse()

	sinceStr := flag.Arg(0)
	if sinceStr == "" {
		fmt.Fprintln(os.Stderr, "usage: urlteamdl <sinceISODate>")
		os.Exit(1)
	}

	sinceTime, err := time.Parse("2006-01-02", sinceStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid date format for '%s'. Please use YYYY-MM-DD (e.g., 2017-10-26).\n", sinceStr)
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request for search URL: %v\n", err)
		os.Exit(1)
	}

	since := sinceTime.Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("collection:(UrlteamWebCrawls) AND addeddate:[%s TO %s]", since, today))
	q.Add("fl[]", "identifier")
	q.Add("sort[]", "addeddate desc")
	q.Add("rows", "500")
	q.Add("output", "json")
	req.URL.RawQuery = q.Encode()

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching search results from %s: %v\n", req.URL.String(), err)
		os.Exit(1)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: Received non-200 status code (%d) from search URL %s\n", res.StatusCode, req.URL.String())
		os.Exit(1)
	}

	dec := json.NewDecoder(res.Body)
	wrapper := &struct {
		Response struct {
			Docs []struct {
				Identifier string `json:"identifier"`
			} `json:"docs"`
		} `json:"response"`
	}{}

	err = dec.Decode(wrapper)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding search results JSON: %v\n", err)
		os.Exit(1)
	}

	for _, d := range wrapper.Response.Docs {
		files, err := getDownloadURLs(httpClient, d.Identifier) // Pass httpClient
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting download URLs for identifier %s: %v\n", d.Identifier, err)
			continue
		}

		for _, f := range files {
			if f.Format != "ZIP" {
				continue
			}
			fmt.Printf("https://archive.org/download/%s/%s\n", d.Identifier, f.Name)
		}
	}
}

func getDownloadURLs(client *http.Client, ident string) ([]file, error) { // Accept httpClient
	metaRequestURL := fmt.Sprintf(metaURL, ident)
	res, err := client.Get(metaRequestURL) // Use passed client
	if err != nil {
		return []file{}, fmt.Errorf("fetching metadata from %s: %w", metaRequestURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []file{}, fmt.Errorf("received non-200 status code (%d) from metadata URL %s", res.StatusCode, metaRequestURL)
	}

	wrapper := &struct {
		Files []file `json:"files"`
	}{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(wrapper)
	if err != nil {
		return []file{}, fmt.Errorf("decoding metadata JSON for %s: %w", ident, err)
	}

	return wrapper.Files, nil
}
