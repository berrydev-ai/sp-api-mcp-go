.PHONY: help inspector

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make dev-http    - Start development server with HTTP transport (port 8002)"
	@echo "  make dev-sse     - Start development server with SSE transport (port 8002)"
	@echo "  make inspector   - Run MCP Inspector debugging tool"
	@echo "  make help        - Show this help message"

dev-http:
	@echo "Starting development server"
	VERBOSE=true MCP_TRANSPORT=streamablehttp PORT=8002 go run ./cmd/server/main.go --config ./config.json --server sp-api

dev-sse:
	@echo "Starting development server"
	VERBOSE=true MCP_TRANSPORT=sse PORT=8002 go run ./cmd/server/main.go --config ./config.json --server sp-api

inspector:
	@echo "Running MCP Inspector"
	DANGEROUSLY_OMIT_AUTH=true mcp-inspector go run ./cmd/server/main.go --config ./config.json --server sp-api