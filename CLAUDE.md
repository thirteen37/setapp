# Setapp CLI

Go CLI for managing Setapp applications on macOS. Reads from Setapp's local SQLite database (read-only).

## Build & Run

- `go build -o setapp .` — build the binary
- `go run .` — run without building
- `go install github.com/thirteen37/setapp@latest` — install globally

## Architecture

- `cmd/` — Cobra command definitions (one file per command)
- `internal/db/` — SQLite database access (read-only, via modernc.org/sqlite)
- `internal/model/` — Data types and helpers
- Database path: `~/Library/Application Support/Setapp/Default/Databases/Apps.sqlite`

## Notes

- No tests yet — use TDD when adding new features
- Binary `setapp` is gitignored
- Sandbox must be disabled for git operations
