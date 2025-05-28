# concurl

Concurrently request URLs provided on `stdin` using the `curl` command line utility, with per-domain rate limiting. The output of each `curl` command is saved to a file.

## Installation

To install the `concurl` command-line tool, ensure you have Go installed (version 1.16 or newer is recommended) and `curl` is available in your system's PATH.

You can install `concurl` using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/concurl@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `concurl` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Make sure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `concurl` project directory:**
    ```bash
    cd path/to/your/new-hacks/concurl
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `concurl` in the current directory. You would then run it as `./concurl`.

## Usage

`concurl` reads URLs from standard input, one URL per line.

```bash
cat urls.txt | concurl [options] [-- <curl_options>]
```

### Options

You can view the options by running `concurl -h`:
```
Usage of concurl:
  -c int
    	Concurrency level (default 20)
  -d int
    	Delay between requests to the same domain in milliseconds (default 5000)
  -o string
    	Output directory (default "out")
```

### Passing Options to `curl`

Any arguments provided after a double dash (`--`) will be passed directly to the underlying `curl` command for each request.

## Examples

**Basic usage:**

Suppose `urls.txt` contains:
```
https://example.com/path?one=1&two=2
https://example.com/pathtwo?two=2&one=1
https://example.net/a/path?two=2&one=1
```

Then run:
```bash
cat urls.txt | concurl -c 3
```
Output to stdout (example):
```
out/example.com/6ad33f150c6a17b4d51bb3a5425036160e18643c https://example.com/path?one=1&two=2
out/example.net/33ce069e645b0cb190ef0205af9200ae53b57e53 https://example.net/a/path?two=2&one=1
out/example.com/5657622dd56a6c64da72459132d576a8f89576e2 https://example.com/pathtwo?two=2&one=1
```
This also creates files in the `out/` directory (or the directory specified by `-o`). For example, `out/example.net/33ce069e645b0cb190ef0205af9200ae53b57e53` would contain:
```
cmd: curl --silent https://example.net/a/path?two=2&one=1
------

<!doctype html>
<html>
<head>
    <title>Example Domain</title>
    ...
```

**Supplying options to `curl`:**
```bash
echo "https://httpbin.org/anything" | concurl -c 5 -- -H "User-Agent: MyConcurlClient/1.0" -H "X-Custom-Header: TestValue"
```
Output to stdout (example):
```
out/httpbin.org/somehashvalue https://httpbin.org/anything
```
The file `out/httpbin.org/somehashvalue` would contain:
```
cmd: curl --silent https://httpbin.org/anything -H User-Agent: MyConcurlClient/1.0 -H X-Custom-Header: TestValue
------

{
  "args": {},
  "headers": {
    "User-Agent": "MyConcurlClient/1.0",
    "X-Custom-Header": "TestValue",
    ...
  },
  ...
}
```

## How it Works

`concurl` reads URLs from stdin and distributes them to a pool of worker goroutines. Each worker:
1.  Parses the URL to extract the domain for rate-limiting and directory structure.
2.  Respects a per-domain rate limit (configurable with `-d`) to avoid overwhelming servers.
3.  Constructs and executes a `curl` command. The `--silent` flag is always added. Any arguments after `--` in the `concurl` command are appended.
4.  Saves the `curl` command executed and its combined output (stdout + stderr) to a file. The filepath is `[outputDir]/[domain]/[sha1_hash_of_url_and_args]`.
5.  Prints the path to the output file and the original URL to stdout.
