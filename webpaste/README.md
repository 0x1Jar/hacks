# Webpaste

Webpaste is a Go application that runs a local web server to receive lines of text (e.g., URLs, notes) sent from a browser extension. It then prints these lines to standard output on the terminal where the server is running. It's useful for quickly collecting data from web pages.

## Features

-   Runs a local HTTP server to listen for incoming data.
-   Authenticates requests using an environment variable `WEBPASTE_TOKEN`.
-   Receives data as JSON payloads containing lines of text.
-   Prints received lines to standard output.
-   Optional `-u` flag to only print unique lines (deduplication).
-   Configurable listening address and port.
-   Includes a companion browser extension (in the `extension/` subdirectory) for sending data.

## Server Setup (Go Program)

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/0x1Jar/new-hacks.git
    cd new-hacks/webpaste
    ```
2.  **Build the server:**
    ```bash
    go build
    ```
    This creates a `webpaste` executable in the current directory.
    Alternatively, install it to your Go bin path:
    ```bash
    go install github.com/0x1Jar/new-hacks/webpaste
    ```

3.  **Set the authentication token:**
    Before starting `webpaste`, set the `WEBPASTE_TOKEN` environment variable. This token must match the one configured in the browser extension.
    ```bash
    export WEBPASTE_TOKEN="your_secret_token_here" 
    # Example: export WEBPASTE_TOKEN=iloveweb
    ```

4.  **Run the server:**
    By default, `webpaste` runs on `0.0.0.0:8080`.
    ```bash
    ./webpaste 
    # or if installed to PATH:
    # webpaste
    ```
    You will see a message: `Listening on 0.0.0.0:8080`.

    **Command-line flags for the server:**

```
$ ./webpaste -h
Usage of webpaste:
  -a string
        address to listen on (default "0.0.0.0")
  -p string
        port to listen on (default "8080")
  -u    only print unique lines (deduplicates received lines)
```

## Browser Extension Setup

The browser extension is located in the `extension/` subdirectory.

1.  **Open Extension Management Page:**
    In Chrome or a Chromium-based browser, go to the extensions page (e.g., `chrome://extensions` or Menu -> More Tools -> Extensions).

2.  **Enable Developer Mode:**
    If not already enabled, turn on "Developer mode" (usually a toggle in the top-right corner).

3.  **Load Unpacked Extension:**
    Click the "Load unpacked" button.
    Navigate to and select the `extension` folder inside your cloned `webpaste` directory (e.g., `path/to/new-hacks/webpaste/extension`).

4.  **Configure Extension Options:**
    *   The Webpaste extension icon should now appear in your browser's toolbar.
    *   Right-click the Webpaste extension icon and select "Options".
    *   **Server URL**: Enter the address and port where your `webpaste` server is running (e.g., `http://localhost:8080` or `http://your-ip:port` if running on a different machine/IP).
    *   **Token**: Enter the same secret token you set for the `WEBPASTE_TOKEN` environment variable (e.g., `iloveweb`).
    *   **Snippets (Optional)**: You can define JavaScript snippets to extract data from pages. The format is an array of objects, each with `name`, `code` (JavaScript to execute), and optionally `onsuccess` (JavaScript to execute after successful data sending, e.g., to navigate to the next page).
        Example snippets are provided in the original `snippets.js` (not part of the extension load, you'd copy-paste into options):
        ```javascript
        [
            {
                "name": "Google URLs",
                "code": "[...document.querySelectorAll('div.r>a:first-child')].map(n=>n.href)",
                "onsuccess": "document.location=document.querySelectorAll('a#pnnext')[0].href;"
            },
            {
                "name": "GitHub Code Results",
                "code": "[...document.querySelectorAll('#code_search_results a.text-bold')].map(n=>n.href)",
                "onsuccess": "document.location=document.querySelectorAll('a.next_page')[0].href;"
            }
        ]
        ```
    *   Save the options.

5.  **Usage Example with Extension:**
    *   Ensure your `webpaste` server is running in a terminal.
    *   Navigate to a web page (e.g., a Google search results page).
    *   Click the Webpaste extension icon in your browser.
    *   If you have snippets configured, click the name of the snippet (e.g., "Google URLs").
    *   The extracted data (URLs, in this example) should appear in the terminal where the `webpaste` server is running.
    *   If an `onsuccess` script is defined for the snippet (like navigating to the next page of results), it will execute.
```

## How it Works (Server)

-   The Go program starts an HTTP server.
-   The `/` endpoint listens for `POST` requests.
-   It expects a JSON payload in the request body with a `token` and an array of `lines`.
-   It validates the `token` against the `WEBPASTE_TOKEN` environment variable.
-   If authenticated, the `lines` are sent to an internal channel (`bus`).
-   A separate goroutine reads arrays of lines from the `bus`.
-   It iterates through these lines, optionally deduplicates them if the `-u` flag is set, and prints each line to standard output.
-   CORS headers are set to allow requests from any origin (`Access-Control-Allow-Origin: *`).

1. Open the Extension Manager by following:

Kebab menu(three vertical dots) -> More Tools -> Extensions

2. If the developer mode is not turned on, turn it on by clicking the toggle in the top right corner.

3. Now click on Load unpacked button on the top left

4. Go the directory where you have webpaste/extension folder and select it.

5. Extension is loaded now.

6. Right click on chrome extension, go to "Options"

Put Server name:

http://localhost:8080 or http://ip:port

Same token as set above example: iloveweb

For Snippets, cloned directory has Google and JS extraction snippets.
```
