.PHONY: inspector

dev-http:
	@echo "Starting development server"
	VERBOSE=true MCP_TRANSPORT=streamablehttp PORT=8002 go run ./cmd/server/main.go --config ./config.json --server sp-api

dev-sse:
	@echo "Starting development server"
	VERBOSE=true MCP_TRANSPORT=sse PORT=8002 go run ./cmd/server/main.go --config ./config.json --server sp-api

inspector:
	@echo "Running MCP Inspector"
	DANGEROUSLY_OMIT_AUTH=true mcp-inspector go run ./cmd/server/main.go --config ./config.json --server sp-api