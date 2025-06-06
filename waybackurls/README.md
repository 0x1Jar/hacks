# waybackurls

Accept line-delimited domains on stdin, fetch known URLs from the Wayback Machine for the exact domain and output them on stdout.

Usage example:

```
▶ cat domains.txt | waybackurls > urls
```

Flags:
- `-c <number>`: Number of concurrent requests (default: 10)
- `-v`: Enable verbose output

Install:

Ensure you have Go (version 1.24.3 or later, as specified in `go.mod`) installed.
```bash
go install github.com/0x1Jar/new-hacks/waybackurls
```
This will place the compiled binary in your Go bin directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`).

Alternatively, to build from source:
```bash
# Clone the repository (if not already done)
# git clone https://github.com/0x1Jar/new-hacks.git
# cd new-hacks/waybackurls
go build
```
This will create a `waybackurls` executable in the current directory.
