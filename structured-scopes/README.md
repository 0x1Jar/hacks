# HackerOne Structured Scopes Fetcher (PHP)

This PHP script (`fetch.php`) connects to the HackerOne GraphQL API to fetch all "URL" type assets from the structured scopes of programs, provided they are eligible for submission.

## Features

-   Paginates through all teams (programs) on HackerOne.
-   Extracts structured scope information for each team.
-   Filters for assets where `asset_type` is "URL" and `eligible_for_submission` is true.
-   Prints the `asset_identifier` (the URL) for each matching asset, one per line, to standard output.
-   Requires a HackerOne `X-Auth-Token` for API authentication.
-   Includes basic error handling for API requests, JSON parsing, and GraphQL errors.
-   Prints progress information (page fetching) to standard error.

## Prerequisites

-   PHP (version 7.0 or newer recommended, due to use of `??` operator).
-   PHP JSON extension enabled (usually by default).
-   PHP OpenSSL extension enabled (for making HTTPS requests, usually by default).
-   A valid HackerOne `X-Auth-Token`.

## How to get an X-Auth-Token

1.  Log in to your HackerOne account.
2.  Open your browser's developer tools (e.g., by pressing F12).
3.  Go to the "Network" tab.
4.  Perform an action that makes an API request, for example, navigate to `https://hackerone.com/programs` or any page that loads data dynamically.
5.  Look for requests to `https://hackerone.com/graphql`.
6.  Inspect the request headers for one of these GraphQL requests. You should find an `X-Auth-Token` header. Copy its value.

**Note:** This token is sensitive. Keep it secure and do not share it. It grants access to the HackerOne API as your user.

## Usage

Run the script from the command line, providing your `X-Auth-Token` as the first argument:

```bash
php fetch.php YOUR_X_AUTH_TOKEN_HERE
```

Or, if you've made the script executable (e.g., `chmod +x fetch.php`):
```bash
./fetch.php YOUR_X_AUTH_TOKEN_HERE
```

**Example:**
```bash
php fetch.php "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" > h1_urls.txt
```
This will save all fetched URLs into `h1_urls.txt`. Progress and errors will be printed to the console (standard error).

## Output

-   **Standard Output (STDOUT):** A list of asset identifiers (URLs), one per line.
    Example:
    ```
    www.example.com
    *.example.org
    dev.example.net
    ```
-   **Standard Error (STDERR):** Progress messages indicating which page is being fetched, and any error messages encountered during the process.
    Example:
    ```
    Fetching page 1 (cursor: none)
    Fetching page 2 (cursor: ABCDEF12345==)
    ...
    ```

## Error Handling

The script includes checks for:
-   Failure to connect to the HackerOne API.
-   Failure to read the API response.
-   Invalid JSON in the API response.
-   GraphQL API errors (e.g., authentication issues, malformed queries).
-   Unexpected API response structure.

If a critical error occurs, the script will `die()` with an error message printed to standard error. Some non-critical issues, like a team missing scope data, will print a warning to standard error and continue.
