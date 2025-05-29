# unimap - Unicode Character Map Viewer

`unimap` is a simple web-based tool to explore Unicode characters. It displays a grid of characters based on a draggable range slider, allowing users to quickly view blocks of Unicode code points.

## Features

-   Displays a 256-character grid at a time.
-   Uses a range slider to navigate through Unicode code points (0 to 65535, i.e., the Basic Multilingual Plane).
-   Dynamically updates the character grid as the slider moves, using `String.fromCodePoint()` to render characters.
-   Built with HTML, PHP (for a constant), and JavaScript.

## File

-   `index.php`: The single file containing the HTML structure, CSS styling, PHP constant definition, and JavaScript logic.

## How it Works

-   The page sets up an HTML range input (`<input type=range>`) with a min of 0 and a max of 65535.
-   PHP is used to define a constant `CHAR_COUNT` (set to 256) and to initially generate 256 empty `<span>` elements that will serve as cells in the character grid.
-   JavaScript listens to the range slider's value.
-   When the slider's value changes, the JavaScript calculates the starting code point for the 256-character block to display.
-   It then iterates from this starting code point for `CHAR_COUNT` items, converting each code point to its corresponding character using `String.fromCodePoint(i)` and updating the content of the respective `<span>` cell.
-   `window.requestAnimationFrame` is used to efficiently update the display.

## Usage

1.  **Ensure you have a PHP-enabled web server.** Examples include Apache with `mod_php`, Nginx with PHP-FPM, or using PHP's built-in development server.
2.  Place `index.php` in a directory served by your web server.
3.  Access `index.php` through your web browser (e.g., `http://localhost/path/to/unimap/index.php`).

Alternatively, for local testing/development using PHP's built-in server:
1.  Navigate to the `unimap/` directory in your terminal.
2.  Start the PHP development server:
    ```bash
    php -S localhost:8000
    ```
3.  Open your web browser and go to `http://localhost:8000/index.php` (or just `http://localhost:8000/` if `index.php` is the default).

Once loaded, use the slider at the top of the page to scroll through different blocks of Unicode characters. Each cell in the grid will display the character corresponding to its code point.
