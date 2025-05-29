package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"flag"
	"os"
	"strings" // Added strings import
	"sync"
	"time"
)

const defaultConcurrency = 20
const defaultTimeout = 10 // seconds

func main() {
	concurrency := flag.Int("c", defaultConcurrency, "Concurrency level")
	timeout := flag.Int("t", defaultTimeout, "HTTP client timeout in seconds")
	// TODO: Add flag for custom origins list file
	// TODO: Add flag for custom patterns list file

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "CORS Misconfiguration Scanner\n\n")
		fmt.Fprintf(os.Stderr, "Reads URLs from stdin and tests for common CORS misconfigurations.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	urls := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		client := getClient(time.Duration(*timeout) * time.Second)
		go func(c *http.Client) {
			defer wg.Done()
			for u := range urls {
				testOrigins(c, u)
			}
		}(client)
	}

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		urls <- sc.Text()
	}
	close(urls)

	wg.Wait()
}

func getClient(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:    30, // TODO: Make configurable?
		IdleConnTimeout: time.Second, // TODO: Make configurable?
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO: Make configurable?
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second, // TODO: Make configurable?
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // Do not follow redirects
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       timeout,
	}
}

func testOrigins(c *http.Client, targetURL string) {
	// Validate targetURL before proceeding
	parsedTargetURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid target URL %s: %v\n", targetURL, err)
		return
	}
	if parsedTargetURL.Scheme == "" || parsedTargetURL.Host == "" {
		fmt.Fprintf(os.Stderr, "Target URL %s must be absolute\n", targetURL)
		return
	}


	permutations, err := getPermutations(targetURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating permutations for %s: %v\n", targetURL, err)
		return
	}

	for _, origin := range permutations {
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for %s with origin %s: %v\n", targetURL, origin, err)
			continue // Try next origin
		}
		req.Header.Set("Origin", origin)
		// Add other common headers that might influence CORS behavior?
		// req.Header.Set("User-Agent", "CORSBlimeyScanner/1.0")

		resp, err := c.Do(req)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "Error requesting %s with origin %s: %v\n", targetURL, origin, err)
			if resp != nil && resp.Body != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
			continue // Try next origin
		}

		acao := resp.Header.Get("Access-Control-Allow-Origin")
		acac := resp.Header.Get("Access-Control-Allow-Credentials") // true or omitted

		// Close body to allow connection reuse
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()

		// Check for interesting CORS headers
		if acao != "" { // If ACAO is set at all
			if acao == origin || acao == "*" {
				// If ACAO reflects the origin or is a wildcard
				// and if ACAC is true, it's particularly interesting
				output := fmt.Sprintf("Target: %s | Origin: %s | ACAO: %s", targetURL, origin, acao)
				if acac == "true" {
					output += " | ACAC: true (VULNERABLE!)"
				} else {
					output += " | ACAC: " + acac
				}
				fmt.Println(output)
			} else if originTld, err := getTldPlusOne(origin); err == nil {
				// Check if ACAO is a less specific TLD+1 of the tested origin
				// e.g. Origin: foo.bar.example.com, ACAO: example.com
				if acaoTld, errAcao := getTldPlusOne(acao); errAcao == nil && acaoTld == originTld && acao != origin {
					output := fmt.Sprintf("Target: %s | Origin: %s | ACAO: %s (Potentially interesting parent domain reflection)", targetURL, origin, acao)
					if acac == "true" {
						output += " | ACAC: true"
					}
					fmt.Println(output)
				}
			}
		}
	}
}


func getPermutations(rawTargetURL string) ([]string, error) {
	target, err := url.Parse(rawTargetURL)
	if err != nil {
		return nil, fmt.Errorf("parsing target URL for permutations %s: %w", rawTargetURL, err)
	}

	hostname := target.Hostname()
	if hostname == "" {
		return nil, fmt.Errorf("target URL %s has no hostname", rawTargetURL)
	}

	// Base set of origins to test
	permutations := []string{
		"null", // Common misconfiguration
		"https://evil.com",
		"http://evil.com",
		target.Scheme + "://" + hostname, // Self-reflection
		target.Scheme + "://sub." + hostname, // Subdomain of target
	}

	// Generate permutations based on the target's hostname
	// Example: if target is sub.example.com, try example.com, evilsub.example.com etc.
	// This part can be expanded significantly.
	// The original patterns were:
	// "https://%s.evil.com" -> https://targethostname.evil.com
	// "https://%sevil.com"  -> https://targethostnameevil.com (less common)

	patterns := []string{
		"https://%s.attacker.com", // Attacker controlled subdomain of attacker.com
		"http://%s.attacker.com",
		"https://sub.%s",          // Attacker controlled subdomain of target's domain
		"http://sub.%s",
		// "https://%s" + ".attacker.com", // Similar to first one
		// "https://" + hostname + ".evil.com", // Specific evil domain
	}

	for _, p := range patterns {
		permutations = append(permutations, fmt.Sprintf(p, hostname))
	}
	
	// Add permutations based on TLD+1 of the target, if possible
	tldPlusOne, err := getTldPlusOne(hostname)
	if err == nil && tldPlusOne != hostname { // only if tldPlusOne is different from full hostname
		permutations = append(permutations, target.Scheme + "://" + tldPlusOne) // e.g. https://example.com
		permutations = append(permutations, target.Scheme + "://sub." + tldPlusOne) // e.g. https://sub.example.com
		permutations = append(permutations, "https://"+tldPlusOne+".attacker.com")
	}


	// TODO: Implement the idea from README:
	// "i have also seen stuff like they don't allow target.com.evil.com or anything after the first .
	// but if you put target.anothertld it works. Then you just purcahse that domain and you are set to go"
	// This would involve trying to replace the TLD of the target with common other TLDs.
	// e.g., if target is example.com, try example.net, example.org, example.co.uk etc.
	// This requires a list of TLDs and careful construction.

	return uniqueStrings(permutations), nil
}

// getTldPlusOne attempts to get the top-level domain plus one more label.
// e.g., "sub.example.co.uk" -> "example.co.uk", "example.com" -> "example.com"
// This is a simplified approach and might not cover all edge cases perfectly.
func getTldPlusOne(hostname string) (string, error) {
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("hostname %s has too few parts to determine TLD+1", hostname)
	}
	// A very basic heuristic: assume the last two parts are TLD+1 for common cases like .com, .org
	// and last three for common ccSLDs like .co.uk, .com.au. This is not robust.
	// For simplicity here, we'll just take the last two parts if there are more than two.
	// If only two, it's already TLD+1 (e.g. example.com)
	if len(parts) == 2 {
		return hostname, nil
	}
	// This is a naive assumption for TLDs like .co.uk, .com.au etc.
	// A proper solution would use a public suffix list.
	if len(parts) > 2 && (parts[len(parts)-2] == "co" || parts[len(parts)-2] == "com" || parts[len(parts)-2] == "org" || parts[len(parts)-2] == "net" || parts[len(parts)-2] == "gov" || parts[len(parts)-2] == "ac") {
		if len(parts) >=3 {
			return strings.Join(parts[len(parts)-3:], "."), nil
		}
	}
	return strings.Join(parts[len(parts)-2:], "."), nil
}

func uniqueStrings(input []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, val := range input {
        if _, ok := seen[val]; !ok {
            seen[val] = true
            result = append(result, val)
        }
    }
    return result
}
