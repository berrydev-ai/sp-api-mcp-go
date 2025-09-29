package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// documentationEntry captures metadata for an MCP documentation resource.
type documentationEntry struct {
	Category            string
	Title               string
	Summary             string
	ImplementationNotes []string
}

var documentationEntries = []documentationEntry{
	{
		Category: "overview",
		Title:    "Selling Partner API Overview",
		Summary:  "Orientation material for the Selling Partner API and how this MCP server organizes functionality.",
		ImplementationNotes: []string{
			"Link to the official SP-API developer documentation landing page.",
			"Describe how SP-API authorization, throttling, and marketplace routing impact every tool.",
			"Clarify environment variables and deployment recommendations for the MCP server.",
		},
	},
	{
		Category: "authentication",
		Title:    "Authentication",
		Summary:  "Guidance for Login with Amazon (LWA) and role-based authorization required by SP-API.",
		ImplementationNotes: []string{
			"Outline the LWA authorization code flow and exchange for refresh tokens.",
			"Document role-based permissions, IAM policy requirements, and AWS STS integration.",
			"Highlight token rotation schedules and secure storage strategies.",
		},
	},
	{
		Category: "catalog",
		Title:    "Catalog",
		Summary:  "Documentation for querying Amazon's product catalog and normalizing item attributes.",
		ImplementationNotes: []string{
			"List supported catalog endpoints and how to select versions (e.g. 2022-04-01).",
			"Explain model coverage for ASINs, seller SKUs, and keyword search.",
			"Call out pagination, locale, and attribute expansion patterns.",
		},
	},
	{
		Category: "orders",
		Title:    "Orders",
		Summary:  "Process order retrieval, acknowledgements, and fulfillment workflows.",
		ImplementationNotes: []string{
			"Cover the Orders API v0/v2 capabilities and when to use each version.",
			"Detail order item pagination, buyer info access, and fulfillment channel nuances.",
			"Describe how to integrate shipment confirmations and refunds with order updates.",
		},
	},
	{
		Category: "inventory",
		Title:    "Inventory",
		Summary:  "Centralize inventory availability, inbound shipments, and restock metrics.",
		ImplementationNotes: []string{
			"Document FBA Inventory APIs and MFN inventory sources supported by the server.",
			"Explain how marketplaces and warehouses affect quantity calculations.",
			"Call out long-term storage fees and restock limit insights.",
		},
	},
	{
		Category: "reports",
		Title:    "Reports",
		Summary:  "Generate, monitor, and download asynchronous reports.",
		ImplementationNotes: []string{
			"List the most common report types and prerequisites for requesting them.",
			"Explain document encryption, compression, and signed URL handling.",
			"Track polling cadence, retry logic, and rate limit strategies for report generation.",
		},
	},
	{
		Category: "feeds",
		Title:    "Feeds",
		Summary:  "Submit, monitor, and validate feeds for catalog and fulfillment updates.",
		ImplementationNotes: []string{
			"Enumerate supported feed types and corresponding content schemas.",
			"Describe staging files to S3 or other storage before invoking the Feeds API.",
			"Detail feed document encryption and result inspection after processing.",
		},
	},
	{
		Category: "finance",
		Title:    "Finance",
		Summary:  "Work with financial event groups, settlements, and chargebacks.",
		ImplementationNotes: []string{
			"Map Finances API resources to accounting events and ledger entries.",
			"Explain pagination through financial event groups and reconciliation best practices.",
			"Highlight tax, fee, and refund event coverage per marketplace.",
		},
	},
	{
		Category: "notifications",
		Title:    "Notifications",
		Summary:  "Manage notification subscriptions and destinations.",
		ImplementationNotes: []string{
			"Document destination creation, encryption keys, and SQS/SNS webhook patterns.",
			"Clarify notification type availability and throttling behaviour.",
			"Provide testing approaches for validating notifications end-to-end.",
		},
	},
	{
		Category: "productPricing",
		Title:    "Product Pricing",
		Summary:  "Access competitive pricing, fee previews, and offer details.",
		ImplementationNotes: []string{
			"List pricing endpoints and required scopes for each marketplace.",
			"Explain batch request limits and caching strategies for price intelligence.",
			"Discuss how to merge pricing with catalog attributes for decisioning.",
		},
	},
	{
		Category: "listings",
		Title:    "Listings",
		Summary:  "Create and manage offer details, compliance data, and images.",
		ImplementationNotes: []string{
			"Summarise Listings Items API patch semantics and conflict handling.",
			"Highlight hazard, compliance, and image validation requirements.",
			"Outline error handling and retries when updating offers at scale.",
		},
	},
	{
		Category: "fba",
		Title:    "Fulfillment by Amazon",
		Summary:  "Operate inbound shipments, inventory placement, and customer fulfillment via FBA.",
		ImplementationNotes: []string{
			"Describe creating inbound shipment plans, labels, and routing workflows.",
			"Explain Small and Light / Amazon Warehousing & Distribution considerations.",
			"Cover reconciliation of received inventory and discrepancy reports.",
		},
	},
}

// Documentation returns server resources that act as placeholders for SP-API reference material.
func Documentation() []server.ServerResource {
	resources := make([]server.ServerResource, 0, len(documentationEntries))
	for _, entry := range documentationEntries {
		uri := fmt.Sprintf("amazon-sp-api://%s", entry.Category)
		resource := mcp.NewResource(
			uri,
			entry.Title,
			mcp.WithResourceDescription(entry.Summary),
			mcp.WithMIMEType("text/markdown"),
		)

		resources = append(resources, server.ServerResource{
			Resource: resource,
			Handler:  documentationHandler(uri, entry),
		})
	}

	return resources
}

func documentationHandler(uri string, entry documentationEntry) server.ResourceHandlerFunc {
	return func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		body := buildDocumentationBody(entry)
		content := mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "text/markdown",
			Text:     body,
		}
		return []mcp.ResourceContents{content}, nil
	}
}

func buildDocumentationBody(entry documentationEntry) string {
	var sb strings.Builder

	sb.WriteString("# ")
	sb.WriteString(entry.Title)
	sb.WriteString("\n\n")
	sb.WriteString(entry.Summary)
	sb.WriteString("\n\n## Implementation Notes\n")

	for _, note := range entry.ImplementationNotes {
		sb.WriteString("- ")
		sb.WriteString(note)
		sb.WriteString("\n")
	}

	sb.WriteString("\n## Status\n")
	sb.WriteString("This resource is a placeholder. Replace it with live links to Amazon documentation or generated knowledge base content as the MCP server matures.\n")

	return sb.String()
}
