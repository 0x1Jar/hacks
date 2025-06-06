# chromeredir - Chrome Redirect Checker

`chromeredir` is a Node.js command-line tool that uses Puppeteer (headless Chrome) to check if a list of URLs redirect and to what destination they redirect.

## Features

*   Reads URLs from standard input.
*   Launches a headless Chrome browser to navigate to each URL.
*   Reports if a URL redirects and its final destination.
*   Reports if a URL does not redirect.
*   Handles errors during navigation.
*   Processes URLs with a degree of concurrency.

## Prerequisites

*   **Node.js**: Version 12.x or newer recommended.
*   **npm**: Comes with Node.js.

## Setup

1.  **Navigate to the `chromeredir` project directory:**
    ```bash
    cd path/to/your/new-hacks/chromeredir
    ```

2.  **Install dependencies (including Puppeteer, which will download a version of Chromium):**
    ```bash
    npm install
    ```
    This command reads the `package.json` file and installs the necessary packages (primarily Puppeteer) into a `node_modules` directory.

## Usage

You can run the script by providing URLs via standard input or from a file.

```bash
node checkredir.js [options]
```

### Options

*   `-f, --file <filepath>`: Specify an input file containing URLs, one URL per line. If not provided, the script reads from standard input.
*   `-c, --concurrency <number>`: Set the number of concurrent Puppeteer pages to use for checking URLs. (Default: 10)
*   `-h, --help`: Show help information.

### Examples

**1. Using standard input:**
```bash
cat list_of_urls.txt | node checkredir.js
```
Or with `echo`:
```bash
echo -e "http://example.com\nhttps://google.com" | node checkredir.js
```

**2. Using an input file:**
```bash
node checkredir.js -f list_of_urls.txt
```

**3. Using an input file with custom concurrency:**
```bash
node checkredir.js -f list_of_urls.txt -c 5
```

### Input Format

The tool expects one URL per line, whether from stdin or an input file. Each URL should include the scheme (e.g., `http://` or `https://`).

### Output Format

The script will output one line per processed URL:

*   If a URL redirects: `<original_url> redirects to <final_url>`
*   If a URL does not redirect (the domain remains the same): `<original_url> does not redirect`
*   If an error occurs: `error checking <original_url>`

**Example Output:**
```
http://example.com redirects to https://www.example.com/
https://google.com does not redirect
http://nonexistentdomain123abc.com error checking http://nonexistentdomain123abc.com
```

## How it Works

The script reads URLs either from standard input or a specified file. It launches a single headless Chrome browser instance. A queue of URLs is processed, with a configurable number of Puppeteer pages operating concurrently. For each URL, a new page navigates to it (ignoring HTTPS errors and using a common User-Agent), and then checks the final `document.domain` and `document.location.href` against the original URL's host and path to determine if a redirect occurred. Errors during navigation are caught and reported. The browser instance is closed once all URLs have been processed.
