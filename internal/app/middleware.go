package app

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ErrorLoggingMiddleware wraps tool handlers to log detailed error information.
func ErrorLoggingMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := next(ctx, request)

		// Log any errors returned from the handler
		if err != nil {
			log.Printf("[ERROR] Tool %s failed with error: %v", request.Params.Name, err)
		}

		// Log tool results that indicate errors
		if result != nil && result.IsError {
			log.Printf("[ERROR] Tool %s returned error result:", request.Params.Name)
			for i, content := range result.Content {
				if textContent, ok := content.(mcp.TextContent); ok {
					log.Printf("  [%d] %s", i, textContent.Text)
				}
			}
		}

		return result, err
	}
}
