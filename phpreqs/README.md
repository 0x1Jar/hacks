# phpreqs - PHP Raw HTTP/S Request Library

`phpreqs` is a simple PHP library for crafting and sending raw HTTP/HTTPS requests and parsing the responses. It provides a low-level way to interact with HTTP servers, similar to tools like `netcat` but with some HTTP-specific conveniences.

## Features

*   Construct and send custom, raw HTTP/S request strings.
*   Supports `tcp` and `tls` (for HTTPS) transports.
*   Fluent interface for building requests (`Request` class).
*   Parses HTTP response status line, headers, and body (`Response` class).
*   Helper methods on the `Response` object for checking headers and body content.
*   Includes `main.php` as a runnable example demonstrating library usage.

## Files

*   `request.php`: Contains the `Request` class, which extends `RawRequest` to provide a higher-level interface for building HTTP requests.
*   `rawrequest.php`: Contains the `RawRequest` class, responsible for the actual socket communication and sending the raw request string.
*   `response.php`: Contains the `Response` class for parsing and interacting with the server's response.
*   `main.php`: An example script showing how to use the library.
*   `composer.json`: Project metadata file for Composer.

## Installation & Usage

This library is designed to be used by including its files directly or via a PSR-4 autoloader if you integrate it into a Composer-managed project.

**1. Direct Inclusion:**

You can clone this repository or download the files and include them in your PHP project using `require_once`:
```php
<?php
require_once __DIR__.'/path/to/phpreqs/request.php'; // Also loads rawrequest.php and response.php due to internal requires

// Your code using the Request class
$req = new Request("https://httpbin.org/get");
$req->addHeader("X-My-Header: Test");
$response = $req->send();

echo $response->toString();
?>
```

**2. Using Composer:**

If you have Composer in your project, you can add this library.
If `phpreqs` were published on Packagist, you would run:
`composer require 0x1jar/phpreqs`

Alternatively, for a local path or Git repository, you might configure your project's `composer.json` to include it. The provided `composer.json` in this directory sets up classmap autoloading, so if you copy this directory into your `vendor` folder (or manage it via Composer), the classes should be autoloadable.

After setting up with Composer and running `composer install` or `composer dump-autoload`, you can use the classes directly without manual `require_once` calls if your project's autoloader is included.

## Core Classes

*   **`Request` (extends `RawRequest`)**
    *   Constructor: `new Request(string $url = "", array $headers = [])`
    *   Methods: `setTransport()`, `setHost()`, `setPort()`, `setMethod()`, `setPath()`, `setQuery()`, `setFragment()`, `setProto()`, `addHeader()`, `setEol()`, `send()`, `toString()`.
*   **`RawRequest`**
    *   Used internally by `Request` or can be used directly for sending completely arbitrary data over a socket.
*   **`Response`**
    *   Methods: `readHeaders()`, `readBody()`, `getHeader()`, `hasHeader()`, `getBody()`, `bodyMatches()`, `toString()`, `close()`.

## Example (`main.php`)

The `main.php` script provides a comprehensive example:
```bash
php main.php
```
This script constructs a POST request to `https://httpbin.org/anything`, sets various properties and headers, sends the request, and then prints the request and response.

Example output snippet from `main.php`:
```
POST /anything?param1=value1&param2=value2#thefragment HTTP/1.1
Host: httpbin.org
Origin: http://evil.com

HTTP/1.1 200 OK
Connection: keep-alive
Server: meinheld/0.6.1
...
{
  "args": {
    "param1": "value1", 
    "param2": "value2"
  }, 
  ...
}
Found ACAO: http://evil.com
It's probably JSON
```

## Notes & Future Enhancements

*   **Error Handling**: The library uses `@` to suppress some socket errors. For production use, more robust error handling (e.g., custom exceptions, more detailed return values) would be beneficial.
*   **Namespacing**: The classes are currently in the global namespace. For better integration into larger projects, they should be moved into a PHP namespace (e.g., `PhpReqs\\Http`).
*   **Dependencies**: The library currently relies only on built-in PHP extensions.
