package tools

import "github.com/mark3labs/mcp-go/server"

// BuildAll assembles every tool the server should expose. Future tools can be registered by extending placeholderSpecs.
func BuildAll(deps Dependencies) []server.ServerTool {
	all := make([]server.ServerTool, 0, len(placeholderSpecs))
	for _, spec := range placeholderSpecs {
		all = append(all, newPlaceholderTool(spec, deps))
	}
	return all
}
