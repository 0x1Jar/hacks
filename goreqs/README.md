# goreqs - Go Raw HTTP Request Library

`goreqs` is a simple Go library for crafting and sending raw HTTP/HTTPS requests and parsing the responses. It allows for low-level control over the request string sent to the server.

## Features

*   Send custom, raw HTTP/HTTPS request strings.
*   Handles both TCP and TLS connections.
*   Parses HTTP response status line, headers, and body.
*   Includes a basic `Response` type to access parsed response data.
*   The `main.go` in this package provides a simple example of its usage.
*   Includes a test suite (`full_test.go`) demonstrating functionality against an HTTP test server.

## Installation (as a library)

To use `goreqs` in your own Go project, you can add it as a dependency:
```bash
go get github.com/0x1Jar/new-hacks/goreqs
```
Then, import it in your Go code:
```go
import "github.com/0x1Jar/new-hacks/goreqs"
```

## Core Components

*   **`RawRequest` struct**:
    *   `transport`: "tcp" or "tls"
    *   `host`: Target hostname or IP address
    *   `port`: Target port
    *   `request`: The raw HTTP request string (e.g., "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n")

*   **`Response` struct**:
    *   `rawStatus`: The full status line (e.g., "HTTP/1.1 200 OK").
    *   `headers`: A slice of raw header strings.
    *   `body`: The response body as `[]byte`.
    *   `Header(name string) string`: Method to get a specific header value (case-insensitive name).

*   **`Do(req Request) (*Response, error)` function**:
    *   Takes a `Request` interface (satisfied by `RawRequest`).
    *   Establishes a connection (TCP or TLS).
    *   Sends the raw request.
    *   Parses the response into a `*Response` object.

## Basic Usage Example (in your Go code)

```go
package main

import (
	"fmt"
	"log"

	"github.com/0x1Jar/new-hacks/goreqs"
)

func main() {
	// Construct a raw HTTP GET request
	rawReqStr := "GET /anything HTTP/1.1\r\n" +
		"Host: httpbin.org\r\n" +
		"User-Agent: goreqs-example/1.0\r\n" +
		"Connection: close\r\n" // Important to close connection for simple example

	req := goreqs.RawRequest{
		Transport: "tls",          // or "tcp" for HTTP
		Host:      "httpbin.org",
		Port:      "443",            // or "80" for HTTP
		Request:   rawReqStr,
	}

	resp, err := goreqs.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.RawStatus()) // Assuming RawStatus() method exists or access rawStatus field
	fmt.Printf("Content-Type: %s\n", resp.Header("Content-Type"))
	fmt.Printf("Body length: %d\n", len(resp.Body())) // Assuming Body() method exists or access body field
	// fmt.Printf("Body: %s\n", string(resp.Body()))
}
```
*(Note: The example above assumes methods like `RawStatus()` and `Body()` on the `Response` type. Based on the provided `response.go`, you'd access `resp.rawStatus` and `resp.body` directly or add such methods.)*

## `main.go` Example

The `main.go` file included in this package serves as a runnable example:
```bash
go run main.go
```
It makes a predefined request to `httpbin.org` and prints the response structure.

## Testing

The package includes tests in `full_test.go`. You can run them from the `goreqs` directory:
```bash
go test
