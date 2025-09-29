package tools

import "github.com/mark3labs/mcp-go/mcp"

var ordersListOrdersSpec = toolSpec{
	Name:        "orders.listOrders",
	Title:       "Order Processing",
	Description: "List orders created or updated within a time window, optionally filtered by status and fulfillment details.",
	Guidance:    "Leverage the Orders API GetOrders operation to page through orders by marketplace and timeframe. When supplying a next token, omit other filters.",
	Options: []mcp.ToolOption{
		mcp.WithArray("marketplaceIds", mcp.Required(), mcp.WithStringItems(), mcp.Description("One or more marketplace identifiers. Required unless using nextToken.")),
		mcp.WithString("createdAfter", mcp.Description("ISO 8601 timestamp filter for order creation time.")),
		mcp.WithString("createdBefore", mcp.Description("ISO 8601 timestamp upper bound for creation time.")),
		mcp.WithString("lastUpdatedAfter", mcp.Description("ISO 8601 timestamp filter for last update time.")),
		mcp.WithString("lastUpdatedBefore", mcp.Description("ISO 8601 timestamp upper bound for last update time.")),
		mcp.WithArray("orderStatuses", mcp.WithStringItems(), mcp.Description("Optional list of order status values to include.")),
		mcp.WithArray("fulfillmentChannels", mcp.WithStringItems(), mcp.Description("Optional fulfillment channel filters (FBA or SellerFulfilled).")),
		mcp.WithArray("paymentMethods", mcp.WithStringItems(), mcp.Description("Optional payment method filters (COD, CVS, Other).")),
		mcp.WithString("buyerEmail", mcp.Description("Filter orders by buyer email.")),
		mcp.WithString("sellerOrderId", mcp.Description("Filter by seller-defined order identifier.")),
		mcp.WithNumber("maxResultsPerPage", mcp.Description("Optional page size between 1 and 100.")),
		mcp.WithArray("easyShipShipmentStatuses", mcp.WithStringItems(), mcp.Description("Optional Amazon Easy Ship status filters.")),
		mcp.WithArray("amazonOrderIds", mcp.WithStringItems(), mcp.Description("Optional list of specific Amazon order IDs to retrieve.")),
		mcp.WithString("nextToken", mcp.Description("Pagination token returned from a previous listOrders call.")),
	},
}

var ordersGetOrderSpec = toolSpec{
	Name:        "orders.getOrder",
	Title:       "Order Processing",
	Description: "Fetch order details and order items for a specific Amazon order ID.",
	Guidance:    "Combine the Orders API v0 and v2 endpoints to retrieve order headers and order items for downstream workflows.",
	Options: []mcp.ToolOption{
		mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
	},
}

var ordersGetOrderAddressSpec = toolSpec{
	Name:        "orders.getOrderAddress",
	Title:       "Order Processing",
	Description: "Retrieve the shipping address for a specific Amazon order.",
	Guidance:    "Call the Orders API getOrderAddress operation to return the buyer-facing shipping address for fulfilment and customer service flows.",
	Options: []mcp.ToolOption{
		mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
	},
}

var ordersGetOrderBuyerInfoSpec = toolSpec{
	Name:        "orders.getOrderBuyerInfo",
	Title:       "Order Processing",
	Description: "Retrieve buyer contact details for a specific Amazon order.",
	Guidance:    "Use the Orders API getOrderBuyerInfo operation to obtain anonymised buyer contact data with the required SP-API scope.",
	Options: []mcp.ToolOption{
		mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
	},
}

var ordersGetOrderItemsSpec = toolSpec{
	Name:        "orders.getOrderItems",
	Title:       "Order Processing",
	Description: "List the line items for a specific Amazon order, supporting pagination via next tokens.",
	Guidance:    "Call the Orders API getOrderItems operation to retrieve order line items and handle pagination using next tokens for large orders.",
	Options: []mcp.ToolOption{
		mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
		mcp.WithString("nextToken", mcp.Description("Pagination token returned from a previous getOrderItems call.")),
	},
}

var ordersGetOrderItemsBuyerInfoSpec = toolSpec{
	Name:        "orders.getOrderItemsBuyerInfo",
	Title:       "Order Processing",
	Description: "Retrieve buyer information for each order item, including gift notes and customization data.",
	Guidance:    "Use the Orders API getOrderItemsBuyerInfo operation to fetch buyer-specific details (gift messages, customization URLs) for each order line item.",
	Options: []mcp.ToolOption{
		mcp.WithString("amazonOrderId", mcp.Required(), mcp.Description("Amazon order identifier (e.g. 123-1234567-1234567).")),
		mcp.WithString("nextToken", mcp.Description("Pagination token returned from a previous getOrderItemsBuyerInfo call.")),
	},
}

var salesGetOrderMetricsSpec = toolSpec{
	Name:        "sales.getOrderMetrics",
	Title:       "Sales Performance",
	Description: "Aggregate order metrics over a requested interval with configurable granularity and filters.",
	Guidance:    "Use the Sales API getOrderMetrics operation to analyse order, unit, and revenue trends. Provide an ISO-8601 interval separated by '--' and align the time zone when aggregating beyond hourly granularity.",
	Options: []mcp.ToolOption{
		mcp.WithArray("marketplaceIds", mcp.Required(), mcp.WithStringItems(), mcp.Description("One or more marketplace identifiers (for example ATVPDKIKX0DER).")),
		mcp.WithString("interval", mcp.Required(), mcp.Description("Inclusive/exclusive ISO 8601 interval formatted as start--end (e.g. 2024-01-01T00:00:00Z--2024-01-08T00:00:00Z).")),
		mcp.WithString("granularity", mcp.Required(), mcp.Enum("Hour", "Day", "Week", "Month", "Year", "Total"), mcp.Description("Time bucket granularity for the metrics.")),
		mcp.WithString("granularityTimeZone", mcp.Description("IANA time zone identifier required when granularity is Day or higher (e.g. UTC, US/Pacific).")),
		mcp.WithString("buyerType", mcp.Enum("B2B", "B2C"), mcp.Description("Optional buyer segment filter.")),
		mcp.WithString("fulfillmentNetwork", mcp.Enum("AFN", "MFN"), mcp.Description("Optional fulfillment network filter.")),
		mcp.WithString("firstDayOfWeek", mcp.Enum("Monday", "Sunday"), mcp.Description("Override the first weekday when granularity is Week.")),
		mcp.WithString("asin", mcp.Description("Optional ASIN filter. Cannot be combined with sku.")),
		mcp.WithString("sku", mcp.Description("Optional seller SKU filter. Cannot be combined with asin.")),
	},
}

var placeholderSpecs = []toolSpec{
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
