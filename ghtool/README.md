# ghtool - GitHub API Utility

`ghtool` is a command-line tool for interacting with the GitHub API. It allows you to search code, list repositories for a user or organization, and list members of an organization.

## Features

*   Search for code across GitHub.
*   List public repositories for a specified user or organization (excluding forks).
*   List members of a specified organization.
*   Uses the official `go-github` library.
*   Handles API pagination automatically to retrieve all results.

## Prerequisites

*   **Go**: Version 1.18 or newer (as specified in `go.mod`).
*   **GitHub Personal Access Token**: You must provide a GitHub Personal Access Token with appropriate permissions (e.g., `public_repo` for accessing public repositories, `read:org` for organization members) via the `GITHUB_TOKEN` environment variable.
    ```bash
    export GITHUB_TOKEN="your_github_pat_here"
    ```
    You can create a token at [https://github.com/settings/tokens](https://github.com/settings/tokens).

## Installation

You can install `ghtool` using `go install`:
```bash
go install github.com/0x1Jar/new-hacks/ghtool@latest
```
This command will fetch the latest version of the module from its repository, compile the source code, and install the `ghtool` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Ensure this directory is in your system's `PATH`.

**For local development or building from source:**

1.  **Navigate to the `ghtool` project directory:**
    ```bash
    cd path/to/your/new-hacks/ghtool
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `ghtool` in the current directory.

## Usage

The basic command structure is:
```bash
ghtool <mode> <query_or_target>
```

### Modes

*   `search <search_query>`
    *   Performs a code search across GitHub using the provided query.
    *   Outputs a list of repository clone URLs containing matching code.
    *   Example: `ghtool search "mycompany.com password"`

*   `repos <username_or_orgname>`
    *   Lists all public, non-forked repositories for the specified GitHub username or organization name.
    *   Outputs a list of repository clone URLs.
    *   Example (user): `ghtool repos octocat`
    *   Example (org): `ghtool repos github`

*   `members <orgname>`
    *   Lists all public members of the specified GitHub organization.
    *   Outputs a list of member login names.
    *   Example: `ghtool members github`

### Examples

**1. Search for code containing "api.secret.com":**
```bash
ghtool search "api.secret.com"
```
Example Output:
```
https://github.com/someuser/somerepo.git
https://github.com/anotheruser/anotherrepo.git
...
```

**2. List repositories for the user "torvalds":**
```bash
ghtool repos torvalds
```
Example Output:
```
https://github.com/torvalds/linux.git
https://github.com/torvalds/uemacs.git
...
```

**3. List members of the "golang" organization:**
```bash
ghtool members golang
```
Example Output:
```
gopher
bradfitz
...
```

Error messages (e.g., missing token, API errors) are printed to standard error.
