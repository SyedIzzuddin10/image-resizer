### Required Software

1. **Go (Golang)**

   - Download from: https://golang.org/
   - Version: 1.21 or later

2. **Node.js & npm**

   - Download from: https://nodejs.org/

3. **Wails Framework**

   - Install command:
     ```bash
     go install github.com/wailsapp/wails/v2/cmd/wails@latest
     ```

### Required Packages

1. **SQLite**

   ```bash
   # For macOS (using Homebrew)
   brew install sqlite3

   # The Go SQLite package will be installed automatically with go mod tidy:
   github.com/mattn/go-sqlite3 v1.14.24
   ```

2. **Image Processing**
   ```bash
   # The image resizing package will be installed automatically with go mod tidy:
   github.com/nfnt/resize
   ```

Note: All Go dependencies will be automatically installed when running `go mod tidy` in the installation steps.

## Access Credentials

Default Account:

- Username: `admin`
- Password: `admin123`

## Live Development

To run in live development mode, run `wails dev` in the project directory.
