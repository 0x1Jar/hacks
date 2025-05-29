# Remove Subdomains (PHP Scripts)

This directory contains PHP scripts designed to manipulate domain names by reading from standard input and using a local `suffixes.txt` file to determine public suffixes.

## Scripts

1.  **`stripsubs.php`**:
    This script takes a list of domain names (one per line) from STDIN. For each domain, it attempts to identify the registrable domain (e.g., `example.com`, `example.co.uk`) and then outputs that registrable domain prepended by one preceding subdomain label, if such a label exists.
    -   If the input is `a.b.example.com`, it outputs `b.example.com`.
    -   If the input is `example.com`, it outputs `example.com`.
    -   The logic relies on the `suffixes.txt` file in the same directory.

2.  **`getsubs.php`**:
    This script also takes a list of domain names (one per line) from STDIN. For each domain, it outputs only the subdomain part that precedes the registrable domain (as determined by `suffixes.txt`).
    -   If the input is `a.b.example.com`, it outputs `a.b`.
    -   If the input is `example.com` (no subdomains relative to the registrable part), it outputs an empty line.

## Dependencies

-   PHP installed on your system.
-   A `suffixes.txt` file in the same directory as the scripts. This file should contain a list of public domain suffixes, one per line. The effectiveness of the scripts depends on the completeness and accuracy of this file.

## Usage

**Prerequisites:**
-   Ensure PHP is installed on your system and accessible via your system's PATH.
-   Make the scripts executable. You can do this by running the command: `chmod +x stripsubs.php getsubs.php`.

Once these prerequisites are met, the scripts can be run directly. The shebang line at the top of each script (e.g., `#!/usr/bin/env php`) tells the system to use the PHP interpreter.

Alternatively, you can always run them by explicitly invoking the PHP interpreter:
`cat list_of_domains.txt | php ./stripsubs.php`
`cat list_of_domains.txt | php ./getsubs.php`

**Direct execution (after `chmod +x`):**

### `stripsubs.php`

```bash
cat list_of_domains.txt | ./stripsubs.php
```

**Example:**
Input (`list_of_domains.txt`):
```
www.google.com
sub1.sub2.example.co.uk
example.org
```

Output of `cat list_of_domains.txt | ./stripsubs.php`:
```
www.google.com
sub2.example.co.uk
example.org
```
*(Note: The behavior of `stripsubs.php` is to return the effective Top-Level Domain Plus One (eTLD+1) along with the single label that immediately precedes it, if one exists. For `www.google.com`, `google.com` is the eTLD+1 and `www` is the preceding label. For `sub1.sub2.example.co.uk`, `example.co.uk` is the eTLD+1, `sub2` is the label immediately preceding it that gets included.)*

### `getsubs.php`

```bash
cat list_of_domains.txt | ./getsubs.php
```

**Example:**
Input (`list_of_domains.txt`):
```
www.google.com
sub1.sub2.example.co.uk
example.org
```

Output of `cat list_of_domains.txt | ./getsubs.php`:
```
www
sub1.sub2
# (empty line for example.org)
```

## `suffixes.txt`

This file is crucial for the correct operation of the scripts. It should list known public suffixes (like `.com`, `.co.uk`, `.org`). The scripts use this list to determine the boundary between the registrable domain and its subdomains.
An outdated or incomplete `suffixes.txt` may lead to incorrect domain parsing.
