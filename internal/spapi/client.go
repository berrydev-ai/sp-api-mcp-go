package spapi

import (
	"fmt"
	"net/http"

	sp "github.com/amzapi/selling-partner-api-sdk/pkg/selling-partner"
)

// Client defines the behaviour expected by MCP tools that need to call SP-API endpoints.
type Client interface {
	AuthorizeRequest(req *http.Request) error
	Endpoint() string
	Status() Status
}

// Status captures the readiness of the underlying SP-API integration.
type Status struct {
	Ready   bool
	Message string
}

// Config contains all runtime settings required to initialise the Selling Partner API client.
type Config struct {
	Endpoint    string
	Credentials Credentials
}

// Credentials mirrors the SP-API secrets required to sign requests.
type Credentials struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

// IsComplete reports whether every credential field has been supplied.
func (c Credentials) IsComplete() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RefreshToken != ""
}

// NewClient builds either a fully-initialised SP-API client or a noop placeholder when credentials are absent.
func NewClient(cfg Config) (Client, error) {
	if !cfg.Credentials.IsComplete() {
		return &noopClient{endpoint: cfg.Endpoint, reason: "selling partner credentials are not configured"}, nil
	}

	spClient, err := sp.NewSellingPartner(&sp.Config{
		ClientID:     cfg.Credentials.ClientID,
		ClientSecret: cfg.Credentials.ClientSecret,
		RefreshToken: cfg.Credentials.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("initialising selling partner client: %w", err)
	}

	return &sellingPartnerClient{endpoint: cfg.Endpoint, client: spClient}, nil
}

type sellingPartnerClient struct {
	endpoint string
	client   *sp.SellingPartner
}

func (c *sellingPartnerClient) AuthorizeRequest(req *http.Request) error {
	return c.client.AuthorizeRequest(req)
}

func (c *sellingPartnerClient) Endpoint() string {
	return c.endpoint
}

func (c *sellingPartnerClient) Status() Status {
	return Status{Ready: true}
}

type noopClient struct {
	endpoint string
	reason   string
}

func (c *noopClient) AuthorizeRequest(_ *http.Request) error {
	return fmt.Errorf("selling partner client unavailable: %s", c.reason)
}

func (c *noopClient) Endpoint() string {
	return c.endpoint
}

func (c *noopClient) Status() Status {
	return Status{Ready: false, Message: c.reason}
}
