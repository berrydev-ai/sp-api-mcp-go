package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

// Dependencies carries the external clients that tool handlers can leverage.
type Dependencies struct {
	SellingPartner spapi.Client
}

type toolSpec struct {
	Name        string
	Title       string
	Description string
	Guidance    string
	Options     []mcp.ToolOption
}

func serverToolFromSpec(spec toolSpec, handler server.ToolHandlerFunc) server.ServerTool {
	options := []mcp.ToolOption{
		mcp.WithDescription(spec.Description),
		mcp.WithTitleAnnotation(spec.Title),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(true),
	}

	options = append(options, spec.Options...)

	tool := mcp.NewTool(spec.Name, options...)

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

func newPlaceholderTool(spec toolSpec, deps Dependencies) server.ServerTool {
	return serverToolFromSpec(spec, placeholderHandler(spec, deps))
}

func placeholderHandler(spec toolSpec, _ Dependencies) server.ToolHandlerFunc {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message := spec.Guidance
		if message == "" {
			message = spec.Description
		}

		var argsBlock string
		if args := req.GetArguments(); len(args) > 0 {
			pretty, err := json.MarshalIndent(args, "", "  ")
			if err != nil {
				pretty = []byte(fmt.Sprintf("%v", args))
			}
			argsBlock = fmt.Sprintf("\n\nReceived arguments:\n%s", pretty)
		}

		text := fmt.Sprintf(
			"%s is a placeholder for future %s capabilities. %s%s",
			spec.Name,
			spec.Title,
			message,
			argsBlock,
		)

		return mcp.NewToolResultText(text), nil
	}
}
