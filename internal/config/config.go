package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultServerName    = "Selling Partner MCP Server"
	defaultServerVersion = "0.1.0"
	defaultInstructions  = "This MCP server organizes Amazon Selling Partner API domains into discrete tools and documentation resources. The current implementation exposes placeholder tools; they should be replaced with real SP-API integrations as capabilities mature."
	defaultEndpoint      = "https://sellingpartnerapi-na.amazon.com"
	defaultTransport     = "stdio"
	defaultHost          = "localhost"
	defaultPort          = "8080"
)

// Transport is the mechanism used to expose the MCP server.
type Transport string

const (
	// TransportSTDIO exposes the server over stdio. This is the default and works well for local development.
	TransportSTDIO Transport = "stdio"
	// TransportSSE exposes the server via server-sent events.
	TransportSSE Transport = "sse"
	// TransportStreamableHTTP exposes the server using the Streamable HTTP transport.
	TransportStreamableHTTP Transport = "streamablehttp"
)

// Credentials encapsulates SP-API credentials sourced from the environment.
type Credentials struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

// IsEmpty returns true when no credential fields are populated.
func (c Credentials) IsEmpty() bool {
	return c.ClientID == "" && c.ClientSecret == "" && c.RefreshToken == ""
}

// IsComplete returns true when every credential field has a value.
func (c Credentials) IsComplete() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RefreshToken != ""
}

// Config represents all runtime configuration required by the MCP server.
type Config struct {
	ServerName    string
	ServerVersion string
	Instructions  string
	SPAPIEndpoint string
	Credentials   Credentials
	Transport     Transport
	Host          string
	Port          string
}

// Load constructs a Config from environment variables, applying defaults and validation.
func Load() (Config, error) {
	transport, err := parseTransport(envOrDefault("MCP_TRANSPORT", defaultTransport))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ServerName:    envOrDefault("MCP_SERVER_NAME", defaultServerName),
		ServerVersion: envOrDefault("MCP_SERVER_VERSION", defaultServerVersion),
		Instructions:  envOrDefault("MCP_SERVER_INSTRUCTIONS", defaultInstructions),
		SPAPIEndpoint: envOrDefault("SP_API_ENDPOINT", defaultEndpoint),
		Credentials: Credentials{
			ClientID:     strings.TrimSpace(os.Getenv("SP_API_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("SP_API_CLIENT_SECRET")),
			RefreshToken: strings.TrimSpace(os.Getenv("SP_API_REFRESH_TOKEN")),
		},
		Transport: transport,
		Host:      envOrDefault("HOST", defaultHost),
		Port:      envOrDefault("PORT", defaultPort),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if !c.Credentials.IsEmpty() && !c.Credentials.IsComplete() {
		return fmt.Errorf("SP-API credentials are partially configured; provide all values or none")
	}

	if (c.Transport == TransportSSE || c.Transport == TransportStreamableHTTP) && strings.TrimSpace(c.Port) == "" {
		return fmt.Errorf("PORT must be set for transport %q", c.Transport)
	}

	if c.Port != "" {
		if _, err := strconv.Atoi(c.Port); err != nil {
			return fmt.Errorf("PORT must be numeric: %w", err)
		}
	}

	return nil
}

func parseTransport(raw string) (Transport, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", string(TransportSTDIO):
		return TransportSTDIO, nil
	case string(TransportSSE):
		return TransportSSE, nil
	case string(TransportStreamableHTTP):
		return TransportStreamableHTTP, nil
	default:
		return Transport(""), fmt.Errorf("unsupported MCP transport %q", raw)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
