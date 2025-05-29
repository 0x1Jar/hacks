package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

const defaultWorkers = 10 // Default number of concurrent workers for DNS lookups

func main() {
	var depth int
	flag.IntVar(&depth, "depth", 4, "max recursion depth (alias -d)")
	flag.IntVar(&depth, "d", 4, "max recursion depth")

	var workers int
	flag.IntVar(&workers, "w", defaultWorkers, "number of concurrent workers for DNS lookups")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: target domain not specified.")
		fmt.Fprintln(os.Stderr, "usage: ettu [--depth=<int>] [--w=<int>] <domain> [<wordfile>|-]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	suffix := flag.Arg(0)
	wordListFile := ""
	if flag.NArg() > 1 {
		wordListFile = flag.Arg(1)
	}
	var wordlistInput io.Reader
	var err error

	if wordListFile == "" || wordListFile == "-" {
		wordlistInput = os.Stdin
		if wordListFile == "" {
			fmt.Fprintln(os.Stderr, "No wordlist file specified, reading from stdin...")
		}
	} else {
		wordlistInput, err = os.Open(wordListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open word list '%s': %s\n", wordListFile, err)
			os.Exit(1) // Exit if specified wordlist file cannot be opened
		}
		// defer wordlistInput.Close() // This would be an issue as it's an io.Reader
	}

	sc := bufio.NewScanner(wordlistInput)

	words := make([]string, 0)
	for sc.Scan() {
		words = append(words, sc.Text())
	}

	out := make(chan string, 1000)

	// Output goroutine
	go func() {
		for o := range out {
			fmt.Println(o)
		}
	}()

	// Semaphore for controlling concurrency of DNS lookups within each brute level
	sem := make(chan struct{}, workers)

	brute(suffix, words, out, 1, depth, sem)

	close(out) // Close output channel after brute returns and all its goroutines are done
}

func brute(suffix string, words []string, out chan string, currentDepth, maxDepth int, sem chan struct{}) {
	if currentDepth > maxDepth {
		return
	}

	var wg sync.WaitGroup

	for _, word := range words {
		candidate := fmt.Sprintf("%s.%s", word, suffix)

		// Acquire semaphore slot
		sem <- struct{}{}
		wg.Add(1)

		go func(cand string) {
			defer func() {
				<-sem // Release semaphore slot
				wg.Done()
			}()

			_, err := net.LookupHost(cand)

			canRecurse := false
			if err == nil {
				// If it resolves, output it and we can recurse
				out <- cand
				canRecurse = true
			} else {
				// If it's not a "no such host" or timeout, we might still recurse
				// as per original logic (dead-end avoidance)
				if dnsErr, ok := err.(*net.DNSError); ok {
					if !dnsErr.IsTimeout && dnsErr.Err != "no such host" {
						canRecurse = true
					}
				} else {
					// Non-DNS error, might be worth logging or handling differently
					// For now, assume we don't recurse on unknown errors
				}
			}

			if canRecurse && currentDepth < maxDepth {
				// Note: The semaphore 'sem' is passed down. Each level of recursion
				// will share this semaphore, effectively limiting total active DNS lookups.
				// If per-level concurrency is desired, a new semaphore would be created here.
				// For now, global concurrency for lookups.
				brute(cand, words, out, currentDepth+1, maxDepth, sem)
			}

		}(candidate)
	}
	wg.Wait()
}
