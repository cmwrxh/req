# req

A minimalist, fast HTTP client CLI written in Go — inspired by httpie but even simpler.

Commands:
- GET / POST / PUT / DELETE
- JSON body support (-d or --data)
- Headers (-H)
- Pretty-printed JSON output
- Save/load request templates (future)

Built in exactly 9 clean commits.

## Build & Run (once cloned locally)
```bash
go build -o req
./req get https://jsonplaceholder.typicode.com/todos/1
