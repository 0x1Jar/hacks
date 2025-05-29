# pupjs - Puppeteer JS/JSON Resource Extractor

`pupjs` is a Node.js command-line tool that uses Puppeteer to launch a headless Chrome browser, navigate to a specified URL, and identify and print the URLs of any JavaScript or JSON resources loaded by the page.

## Features

*   Uses Puppeteer to control a headless Chrome browser.
*   Navigates to a user-provided URL.
*   Intercepts network responses to check their content types.
*   Prints the URLs of resources identified as `javascript` or `json`.
*   Sets a default User-Agent.
*   Ignores HTTPS errors by default.
*   Includes basic error handling for page navigation and Puppeteer launch.

## Prerequisites

*   **Node.js**: Version 10.x or newer is generally recommended for modern Puppeteer versions.
*   **npm**: Comes with Node.js.
*   **Chromium**: Puppeteer downloads a compatible version of Chromium by default when installed via npm.

## Setup

1.  **Navigate to the `pupjs` project directory:**
    ```bash
    cd path/to/your/new-hacks/pupjs
    ```

2.  **Install dependencies:**
    This command reads `package.json` and installs `puppeteer` (which includes downloading Chromium).
    ```bash
    npm install
    ```
    If you want to skip the automatic Chromium download (e.g., if you want to use an existing Chrome/Chromium installation), you can set an environment variable before running `npm install`:
    ```bash
    PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true npm install
    # Then, you might need to specify the executable path:
    # export PUPPETEER_EXECUTABLE_PATH=/path/to/your/chrome
    ```
    Refer to the [Puppeteer documentation](https://pptr.dev/guides/installation) for more details on environment variables.

## Usage

Run the script with the target URL as a command-line argument:
```bash
node main.js <url>
```

### Example

```bash
node main.js https://example.com
```

Example Output (will vary based on the resources loaded by example.com):
```
https://example.com/some/script.js
https://example.com/api/data.json
https://some-cdn.com/library.min.js
...
```

If no URL is provided, or if Puppeteer fails to launch or navigate, an error message will be printed to stderr.

## How it Works
The script launches a Puppeteer instance. It opens a new page, sets a User-Agent, and enables request interception (though it currently continues all requests without modification). It then listens for `response` events. For each response, it checks the `content-type` header. If the content type indicates JavaScript or JSON, the URL of that response is printed to the console. The script navigates to the URL provided as a command-line argument, waits for network activity to settle (`networkidle2`), and then closes the browser.
