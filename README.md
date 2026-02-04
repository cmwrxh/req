# req

Minimalist HTTP client CLI written in Go — fast, colorful, and simple.

Inspired by httpie, but even lighter.

## Features
- GET and POST requests
- JSON body support (`-d` / `--data`)
- Custom headers (`-H` / `--header`, multiple allowed)
- Colored output (method, URL, status codes)
- Automatic pretty-printing of JSON responses
- Clear error messages

## Installation

```bash
# Install globally (recommended)
go install github.com/cmwrxh/req@latest

# Or build from source
git clone https://github.com/cmwrxh/req.git
cd req
go build -o req
# Move ./req to a directory in your $PATH (optional)
