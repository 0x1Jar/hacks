# Security.txt Parser (PHP)

This PHP script (`parse.php`) parses `security.txt` files. It reads a file path as a command-line argument, processes the content, and then prints out encountered errors, comments, and the parsed fields.

The script specifically looks for and validates the following fields (case-insensitive):
-   `Contact`
-   `Encryption`
-   `Disclosure`
-   `Acknowledgement`

## Features

-   Parses fields defined in the script.
-   Validates `Contact` entries (accepts email, URL, or phone numbers starting with `+`).
-   Validates `Encryption` and `Acknowledgement` fields as URIs.
-   Validates `Disclosure` field against a predefined list: `full`, `partial`, `none`.
-   Collects comments (lines starting with `#`).
-   Reports parsing errors and validation errors.
-   Requires at least one `Contact` field to be present.

## Files

-   `parse.php`: The main PHP parser script.
-   `example.txt`: An example `security.txt`-like file that can be used for testing the parser.

## Prerequisites

-   PHP installed on your system and accessible via your system's PATH.

## Usage

You can run the script from the command line, providing the path to a `security.txt` file as an argument.

1.  **Make the script executable (optional but recommended for direct execution):**
    ```bash
    chmod +x parse.php
    ```

2.  **Run the script:**
    Using direct execution (if made executable and shebang `#!/usr/bin/env php` works):
    ```bash
    ./parse.php /path/to/your/security.txt
    ```
    Or by explicitly invoking PHP:
    ```bash
    php parse.php /path/to/your/security.txt
    ```

    To test with the provided example:
    ```bash
    php parse.php example.txt
    ```

## Example Output (using `example.txt`)

Running `php parse.php example.txt` will produce output similar to this (line numbers for errors may vary based on the exact `example.txt` used):

```
errors:
	invalid value '+44-INVALID-NUMBER' for option 'contact' on line 18
	invalid URI 'not a URL' for option 'encryption' on line 14
	invalid URI 'not a URL' for option 'acknowledgement' on line 15
	invalid value 'foo' for option 'disclosure' on line 12; must be one of [full, partial, none]
comments:
	# This is an example file
	# This is another comment
	# This field has different casing
contact:
	mail@tomnomnom.com
	https://tomnomnom.com
	+44 7555 555 555
	+44-7555-555-555
encryption:
	https://tomnomnom.uk/pgpkey
disclosure:
	Full
	full
acknowledgement:
	https://tomnomnom.uk/hof
```

**Note on "Disclosure" field:** The `Disclosure` field with allowed values `full`, `partial`, `none` (case-insensitive) is a custom implementation in this script. Values not matching these (like "foo" in `example.txt`) will generate an error and will not be stored as valid `Disclosure` entries, which is the correct behavior according to the script's validation logic.

## Script Logic Overview

-   The `SecurityTxt` class encapsulates the parsing and validation logic.
-   `parse()` method:
    -   Splits input into lines.
    -   Ignores empty lines and processes comments.
    -   Splits lines by the first colon into field name and value.
    -   Validates field names and their corresponding values.
    -   Stores valid fields and any encountered errors.
-   Validation methods (`validateContact`, `validateDisclosure`, `validateUri`) check the format and content of field values.
-   The script requires at least one valid `Contact` field.
-   After parsing, it prints errors (if any), comments, and then all successfully parsed fields and their values.

## Potential Improvements (Not Implemented)

-   Alignment with RFC 9116: This would involve adding mandatory fields like `Expires`, supporting other standard fields (`Acknowledgments` (plural), `Preferred-Languages`, `Canonical`, `Policy`, `Hiring`), and potentially re-evaluating the custom `Disclosure` field.
-   Stricter validation for `Disclosure` values if only `full`, `partial`, `none` are meant to be stored.
-   More detailed phone number validation.
