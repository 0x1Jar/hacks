package main

import (
	"bufio"
	"fmt"
	"log" // Added log import
	"net"
	"os"
	"sync"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup

	for sc.Scan() {
		domain := sc.Text()

		wg.Add(1)
		go func() {

			if _, err := net.LookupHost(domain); err == nil {
				fmt.Println(domain)
			}

			wg.Done()
		}()
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	wg.Wait()
}
