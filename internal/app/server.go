package app

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/config"
	"github.com/berrydev-ai/sp-api-mcp-go/internal/resources"
	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
	"github.com/berrydev-ai/sp-api-mcp-go/internal/tools"
)

// Dependencies bundles runtime clients the MCP server relies on.
type Dependencies struct {
	SellingPartner spapi.Client
}

// NewServer constructs the MCP server, wiring tools and resources so additional capabilities can be added in one place.
func NewServer(cfg config.Config, deps Dependencies) *server.MCPServer {
	srv := server.NewMCPServer(
		cfg.ServerName,
		cfg.ServerVersion,
		server.WithInstructions(cfg.Instructions),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, true),
		server.WithRecovery(),
		server.WithResourceRecovery(),
		server.WithLogging(),
		server.WithToolHandlerMiddleware(ErrorLoggingMiddleware),
	)

	srv.AddTools(tools.BuildAll(tools.Dependencies{SellingPartner: deps.SellingPartner})...)
	srv.AddResources(resources.Documentation()...)

	return srv
}
