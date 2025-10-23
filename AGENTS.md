# Agent Guidelines for sp-api-mcp-go

## Build & Test Commands
- `go build ./...` - compile all packages
- `go test ./...` - run all tests
- `go test -run TestName ./internal/tools` - run single test by name
- `go run ./cmd/server` - start MCP server (reads config from .env)
- `make dev-http` - dev server (HTTP transport, port 8002, verbose logging)
- `make inspector` - MCP Inspector debugging tool

## Code Style & Conventions
- **Formatting**: Use `gofmt` (tabs for indentation, trailing newlines required)
- **Imports**: Group as: standard library, third-party, internal (`github.com/berrydev-ai/sp-api-mcp-go/...`)
- **Naming**: PascalCase for exports matching SP-API terms; lowerCamelCase for locals
- **Error handling**: Wrap errors with `fmt.Errorf("context: %w", err)` for stack traces
- **Logging**: Use `log` package; enable verbose output with `VERBOSE=true` env var
- **Types**: Prefer explicit struct types for API request/response payloads; use pointer fields for optional values
- **No comments**: Do not add code comments unless explicitly requested

## Project Structure
- `cmd/server/main.go` - entrypoint loading config and starting MCP server
- `internal/config/` - environment-based configuration
- `internal/spapi/` - SP-API client wrapper with credential management
- `internal/tools/` - MCP tool implementations (orders.go, sales.go, etc.)
- `internal/app/` - MCP server setup and middleware

## Environment & Security
- Store credentials in `.env` (never commit): `SP_API_CLIENT_ID`, `SP_API_CLIENT_SECRET`, `SP_API_REFRESH_TOKEN`
- Generated binaries (`bin/sp-api-mcp`) should remain gitignored
