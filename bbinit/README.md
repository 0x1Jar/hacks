# bbinit - Bug Bounty Program Scope Initializer

`bbinit` is a command-line tool that fetches scope information (in-scope and out-of-scope assets) for a given bug bounty program from HackerOne using its GraphQL API. It then generates three files:

*   `.scope`: A file compatible with Burp Suite's scope definition, including regex for domains.
*   `domains`: A plain list of all in-scope domain identifiers.
*   `wildcards`: A plain list of in-scope wildcard domain identifiers (e.g., `*.example.com` would list `example.com`).

## Prerequisites

*   **Go**: Version 1.16 or newer.
*   **HackerOne GraphQL Token**: You need a HackerOne API token.
    *   Set it as an environment variable: `export H1_GRAPHQL_TOKEN="your_token_here"`
    *   You can obtain a token from [https://hackerone.com/current_user/graphql_token.json](https://hackerone.com/current_user/graphql_token.json).

## Installation

You can install `bbinit` using `go install`:

```bash
go install github.com/0x1Jar/new-hacks/bbinit@latest
```
This will install the `bbinit` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`). Ensure this directory is in your system's `PATH`.

The output files (`.scope`, `domains`, `wildcards`) will be created in the directory where you run the `bbinit` command.

**For local development or building from source:**

1.  **Navigate to the `bbinit` project directory:**
    ```bash
    cd path/to/your/new-hacks/bbinit
    ```
2.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `bbinit` in the current directory. You would then run it as `./bbinit`.

## Usage

```bash
bbinit [options] <program_handle>
```

*   `<program_handle>`: The handle of the HackerOne program (e.g., `hackerone`, `github`). This is a required argument.

### Options

*   `-risky`: (boolean, default: `false`)
    If set, treats all in-scope domains as wildcards when generating the `.scope` file and adding to the `wildcards` file. For example, if `example.com` is in scope, with `-risky` it will be treated as `*.example.com`.
*   `-append-scope`: (boolean, default: `false`)
    If set, appends to the existing `.scope` file instead of overwriting it. The `domains` and `wildcards` files are always appended to.

### Examples

**1. Fetch scope for the "hackerone" program:**

```bash
bbinit hackerone
```
This will create/overwrite `.scope` and create/append to `domains` and `wildcards` in the current directory.

**2. Fetch scope for "hackerone" and treat all domains as wildcards:**

```bash
bbinit -risky hackerone
```

**3. Fetch scope for "anotherprogram" and append to an existing `.scope` file:**

```bash
bbinit -append-scope anotherprogram
```

## Output Files

*   **.scope**:
    *   Contains regex patterns for Burp Suite.
    *   In-scope domains: `^exactdomain\.com$` or `.*\.wildcarddomain\.com$`
    *   Out-of-scope domains: `!^exactoutofscope\.com$` or `!.*\.wildcardoutofscope\.com$`
*   **domains**:
    *   A list of all unique in-scope domain identifiers (e.g., `example.com`, `sub.example.org`).
*   **wildcards**:
    *   A list of the base domains for in-scope wildcards (e.g., if `*.example.com` is in scope, `example.com` is listed here).
    *   If `-risky` is used, all in-scope domains will have their base domain listed here.

These files will be created in the directory from which `bbinit` is executed.
