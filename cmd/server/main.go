package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/app"
	"github.com/berrydev-ai/sp-api-mcp-go/internal/config"
	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

func main() {
	// Loading environment variables from .env is optional but convenient during development.
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	spClient, err := spapi.NewClient(spapi.Config{
		Endpoint: cfg.SPAPIEndpoint,
		Credentials: spapi.Credentials{
			ClientID:     cfg.Credentials.ClientID,
			ClientSecret: cfg.Credentials.ClientSecret,
			RefreshToken: cfg.Credentials.RefreshToken,
		},
	})
	if err != nil {
		log.Fatalf("failed to initialize Selling Partner client: %v", err)
	}

	if status := spClient.Status(); status.Message != "" {
		log.Printf("Selling Partner client status: ready=%t detail=%s", status.Ready, status.Message)
	}

	srv := app.NewServer(cfg, app.Dependencies{SellingPartner: spClient})

	switch cfg.Transport {
	case config.TransportSSE:
		log.Printf("starting SSE MCP server on port %s", cfg.Port)
		sse := server.NewSSEServer(srv)
		if err := sse.Start(":" + cfg.Port); err != nil {
			log.Fatalf("sse server exited: %v", err)
		}
	case config.TransportStreamableHTTP:
		log.Printf("starting StreamableHTTP MCP server on port %s", cfg.Port)
		httpSrv := server.NewStreamableHTTPServer(srv)
		if err := httpSrv.Start(":" + cfg.Port); err != nil {
			log.Fatalf("streamable HTTP server exited: %v", err)
		}
	default:
		log.Printf("starting STDIO MCP server")
		if err := server.ServeStdio(srv); err != nil {
			log.Fatalf("stdio server exited: %v", err)
		}
	}
}
