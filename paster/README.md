# Paster - Simple Snippet Clipboard Tool

`Paster` is a small graphical Go application that provides a window with buttons for predefined text snippets. Clicking a button copies the corresponding snippet to the system clipboard.

## Features

*   Simple GUI with buttons for each snippet.
*   Cross-platform clipboard access using `github.com/atotto/clipboard`.
*   Cross-platform GUI using `github.com/andlabs/ui`.
*   Currently hardcoded snippets:
    *   "Upside Down": `🙃`
    *   "Eye Roll": `🙄`
    *   "Shrug": `¯\_(ツ)_/¯`

## Prerequisites

*   **Go**: Version 1.16 or newer is recommended.
*   **Platform-Specific GUI Dependencies**: The `github.com/andlabs/ui` library has dependencies that need to be installed on your system for the program to build and run.
    *   **Linux**: `libgtk-3-dev` (e.g., `sudo apt-get install libgtk-3-dev` on Debian/Ubuntu).
    *   **macOS**: Xcode Command Line Tools (usually already installed if you have Go set up; if not, `xcode-select --install`).
    *   **Windows**: `gcc` (e.g., via MinGW-w64) is required for Cgo. Ensure `gcc` is in your PATH.
    *   Refer to the [andlabs/ui prerequisites](https://github.com/andlabs/ui#requirements) for the most up-to-date information.

## Installation

1.  **Ensure GUI prerequisites are installed** for your operating system (see above).
2.  Install `paster` using `go install`:
    ```bash
    go install github.com/0x1Jar/new-hacks/paster@latest
    ```
    This command will fetch the module, compile it (which may take a moment due to Cgo and the UI library), and install the `paster` binary to your Go binary directory (e.g., `$GOPATH/bin` or `$HOME/go/bin`).

**For local development or building from source:**

1.  **Navigate to the `paster` project directory:**
    ```bash
    cd path/to/your/new-hacks/paster
    ```
2.  **Ensure GUI prerequisites are installed.**
3.  **Build the executable:**
    ```bash
    go build
    ```
    This creates an executable file named `paster` (or `paster.exe` on Windows) in the current directory.

## Usage

Run the compiled executable:
```bash
./paster
```
Or, if installed to your PATH:
```bash
paster
```
A small window will appear with buttons for each snippet. Clicking a button will copy its associated text to your clipboard.

## Customizing Snippets

Currently, the snippets are hardcoded in `main.go`. To change or add snippets, you will need to:
1.  Edit the `snippets` slice in the `main.go` file.
2.  Recompile the application using `go build` or `go install`.

Example `snippets` slice in `main.go`:
```go
snippets := [][]string{
    {"Upside Down", "🙃"},
    {"Eye Roll", "🙄"},
    {"Shrug", `¯\_(ツ)_/¯`},
    // Add new snippets here, e.g.:
    // {"My New Snippet", "Text to copy"},
}
```

## How it Works
The application uses `github.com/andlabs/ui` to create a native graphical window and buttons. Each button is associated with a text snippet. When a button is clicked, its `OnClicked` event handler uses `github.com/atotto/clipboard` to write the corresponding text snippet to the system clipboard.
