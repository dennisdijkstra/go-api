# Go API Project

Modern backend API written in Go with a PostgreSQL database.

## Tech Stack

- Go
- PostgreSQL
- SQLC
- REST-style HTTP API

## Project Structure

- `main.go` – application entrypoint
- `server/` – HTTP server setup
- `internal/database/` – generated database access code
- `sql/schema/` – database schema migrations
- `sql/queries/` – SQL queries used by SQLC

## Getting Started

1. Install Go and PostgreSQL.
2. Create a PostgreSQL database.
3. Apply migrations from `sql/schema/`.
4. Run the app:

```bash
go run .
```

## Development

- Generate/update database code with SQLC after changing SQL queries or schema.
- Keep API handlers in the root files grouped by feature (users, chirps, auth, tokens, webhooks).
