package tools

import "github.com/mark3labs/mcp-go/mcp"

var placeholderSpecs = []placeholderSpec{
	{
		Name:        "auth.beginAuthorization",
		Title:       "Authentication",
		Description: "Initiate Login with Amazon workflow for the Selling Partner API.",
		Guidance:    "Provide Login with Amazon authorization URLs and exchange refresh tokens using the Tokens API when implementing this tool.",
		Options: []mcp.ToolOption{
			mcp.WithString("marketplaceId", mcp.Description("Optional marketplace to scope the authorization.")),
		},
	},
	{
		Name:        "catalog.lookupItem",
		Title:       "Catalog",
		Description: "Retrieve catalog metadata for a specific identifier or keyword.",
		Guidance:    "Call the Catalog Items API (2022-04-01) to return attributes, dimensions, and relationships for the requested identifier.",
		Options: []mcp.ToolOption{
			mcp.WithString("identifier", mcp.Title("Identifier"), mcp.Description("ASIN, seller SKU, or keyword to search for."), mcp.Required()),
			mcp.WithString("identifierType", mcp.Title("Identifier Type"), mcp.Enum("ASIN", "SKU", "Keyword"), mcp.Description("Controls how the identifier value is interpreted.")),
		},
	},
	{
		Name:        "inventory.getSummary",
		Title:       "Inventory Management",
		Description: "Return inventory summaries for a given SKU across marketplaces.",
		Guidance:    "Use the FBA Inventory API to fetch current availability and inbound quantities.",
		Options: []mcp.ToolOption{
			mcp.WithString("sku", mcp.Required(), mcp.Description("Seller SKU to summarise.")),
			mcp.WithString("marketplaceId", mcp.Description("Optional marketplace override when looking up inventory.")),
		},
	},
	{
		Name:        "orders.getOrder",
		Title:       "Order Processing",
		Description: "Fetch order details and order items for a specific Amazon order ID.",
		Guidance:    "Combine the Orders API v0 and v2 endpoints to retrieve order headers and order items for downstream workflows.",
		Options: []mcp.ToolOption{
			mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
		},
	},
	{
		Name:        "reports.createReport",
		Title:       "Reports",
		Description: "Submit an asynchronous report creation request.",
		Guidance:    "Leverage the Reports API createReport operation and monitor the resulting report document until it is ready for download.",
		Options: []mcp.ToolOption{
			mcp.WithString("reportType", mcp.Required(), mcp.Description("Report type identifier, e.g. GET_FLAT_FILE_ALL_ORDERS_DATA_BY_LAST_UPDATE.")),
			mcp.WithArray("marketplaceIds", mcp.Description("Optional list of marketplaces to include."), mcp.WithStringItems()),
		},
	},
	{
		Name:        "feeds.submitFeed",
		Title:       "Feed Submission",
		Description: "Upload and submit a processing feed to Amazon.",
		Guidance:    "Use the Feeds API createFeed operation with the appropriate content type and optional encryption metadata.",
		Options: []mcp.ToolOption{
			mcp.WithString("feedType", mcp.Required(), mcp.Description("Feed type identifier, e.g. POST_PRODUCT_DATA.")),
		},
	},
	{
		Name:        "finance.listFinancialEvents",
		Title:       "Financial Data",
		Description: "List financial events related to orders, refunds, and shipments.",
		Guidance:    "Invoke the Finances API to page through financial events and reconcile transactions.",
		Options: []mcp.ToolOption{
			mcp.WithString("amazonOrderId", mcp.Description("Optional Amazon order identifier to filter results.")),
		},
	},
	{
		Name:        "notifications.subscribe",
		Title:       "Notification Management",
		Description: "Register a subscription for an SP-API notification type.",
		Guidance:    "Implement this tool using the Notifications API to create or update subscriptions tied to your destination resources.",
		Options: []mcp.ToolOption{
			mcp.WithString("notificationType", mcp.Required(), mcp.Description("Notification type to subscribe to, e.g. ANY_OFFER_CHANGED.")),
		},
	},
	{
		Name:        "pricing.getPricing",
		Title:       "Product Pricing",
		Description: "Retrieve competitive pricing and offers for a catalog item.",
		Guidance:    "Use the Product Pricing API getPricing endpoint to evaluate buy box competition and price recommendations.",
		Options: []mcp.ToolOption{
			mcp.WithString("asin", mcp.Description("Optional ASIN when pricing by catalog identifier.")),
			mcp.WithString("sku", mcp.Description("Optional seller SKU when pricing by offer.")),
			mcp.WithString("marketplaceId", mcp.Required(), mcp.Description("Marketplace identifier for the pricing request.")),
		},
	},
	{
		Name:        "listings.updateListing",
		Title:       "Listings",
		Description: "Create or update a listing for a given SKU.",
		Guidance:    "Tie into the Listings Items API to patch attributes, images, and compliance details for existing offers.",
		Options: []mcp.ToolOption{
			mcp.WithString("sku", mcp.Required(), mcp.Description("Seller SKU whose listing should be updated.")),
			mcp.WithString("marketplaceId", mcp.Required(), mcp.Description("Marketplace identifier for the listing.")),
		},
	},
	{
		Name:        "fba.createInboundShipmentPlan",
		Title:       "FBA Operations",
		Description: "Plan inbound shipments to Amazon fulfillment centers.",
		Guidance:    "Coordinate with the FBA Inbound Eligibility and Inbound Shipment APIs to generate labels and routing information.",
		Options: []mcp.ToolOption{
			mcp.WithString("shipFromAddressId", mcp.Required(), mcp.Description("Identifier for the ship-from address resource.")),
		},
	},
}
