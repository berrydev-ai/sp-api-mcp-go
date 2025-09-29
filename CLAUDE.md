# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Build and verify the project
go build ./...

# Run tests
go test ./...

# Run the MCP server (default: stdio transport)
go run ./cmd/server

# Build a distributable binary
go build -o bin/sp-api-mcp ./cmd/server

# Run with specific transports (requires PORT env var)
PORT=8080 MCP_TRANSPORT=sse go run ./cmd/server
PORT=9090 MCP_TRANSPORT=streamablehttp go run ./cmd/server
```

## Architecture Overview

This is a Go-based Model Context Protocol (MCP) server that wraps Amazon's Selling Partner API. The codebase follows standard Go project structure with clear separation of concerns:

### Core Structure
- `cmd/server/main.go` - Entry point that initializes configuration, SP-API client, and starts the MCP server
- `internal/app/server.go` - MCP server construction and dependency wiring
- `internal/config/config.go` - Environment-based configuration with validation
- `internal/spapi/` - Selling Partner API client abstraction
- `internal/tools/` - MCP tool implementations (orders, sales, placeholders)
- `internal/resources/` - MCP resource definitions for documentation

### Tool System
The MCP tools are organized by SP-API domain:
- **Orders** (`internal/tools/orders.go`) - Live SP-API integration for listing/retrieving orders
- **Sales** (`internal/tools/sales.go`) - Sales analytics and reporting tools  
- **Placeholders** (`internal/tools/specs.go`) - Template tools for catalog, inventory, feeds, etc.

New tools are registered in `internal/tools/registry.go` via the `BuildAll()` function.

### Configuration Requirements
The server requires SP-API credentials via environment variables:
- `SP_API_CLIENT_ID` (required)
- `SP_API_CLIENT_SECRET` (required) 
- `SP_API_REFRESH_TOKEN` (required)
- `SP_API_ENDPOINT` (defaults to NA region)

Optional `.env` file support is available for local development.

### Transport Options
Supports three MCP transport modes:
- `stdio` (default) - For desktop clients like Cursor/Claude
- `sse` - Server-sent events over HTTP
- `streamablehttp` - HTTP-based streaming transport

## Code Patterns

- Follow Go standard formatting with `gofmt` (tabs, trailing newlines)
- SP-API client interactions are abstracted through the `spapi.Client` interface
- MCP tools use the `mark3labs/mcp-go` framework for tool definitions and JSON schema validation
- Configuration uses environment variables with reasonable defaults for development
- The `__reference/` directory contains generated SP-API SDK code for reference but is not used directly

## Testing

Run the full test suite before making changes:
```bash
go test ./...
```

The current test coverage focuses on tool functionality, particularly the sales analytics components in `internal/tools/sales_test.go`.