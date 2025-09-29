package tools

import (
	"context"
	"strings"

	"github.com/amzapi/selling-partner-api-sdk/productPricing"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

func newProductPricingTools(deps Dependencies) []server.ServerTool {
	spClient := deps.SellingPartner
	
	getPricingHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args productPricingGetPricingArgs) (*mcp.CallToolResult, error) {
		return executeProductPricingGetPricing(ctx, args, spClient)
	})
	
	getCompetitivePricingHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args productPricingGetCompetitivePricingArgs) (*mcp.CallToolResult, error) {
		return executeProductPricingGetCompetitivePricing(ctx, args, spClient)
	})
	
	return []server.ServerTool{
		serverToolFromSpec(productPricingGetPricingSpec, getPricingHandler),
		serverToolFromSpec(productPricingGetCompetitivePricingSpec, getCompetitivePricingHandler),
	}
}

type productPricingGetPricingArgs struct {
	MarketplaceID string   `json:"marketplaceId"`
	ItemType      string   `json:"itemType"`
	Asins         []string `json:"asins"`
	Skus          []string `json:"skus"`
	ItemCondition string   `json:"itemCondition"`
}

type productPricingGetCompetitivePricingArgs struct {
	MarketplaceID string   `json:"marketplaceId"`
	ItemType      string   `json:"itemType"`
	Asins         []string `json:"asins"`
	Skus          []string `json:"skus"`
}

func executeProductPricingGetPricing(ctx context.Context, args productPricingGetPricingArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureProductPricingClient(spClient)
	if failure != nil {
		return failure, nil
	}

	// Validation
	if strings.TrimSpace(args.MarketplaceID) == "" {
		return mcp.NewToolResultError("marketplaceId is required"), nil
	}
	if strings.TrimSpace(args.ItemType) == "" {
		return mcp.NewToolResultError("itemType is required"), nil
	}

	// Validate ItemType
	if args.ItemType != "Asin" && args.ItemType != "Sku" {
		return mcp.NewToolResultError("itemType must be 'Asin' or 'Sku'"), nil
	}

	// Validate that the appropriate identifier list is provided
	if args.ItemType == "Asin" && len(args.Asins) == 0 {
		return mcp.NewToolResultError("asins list is required when itemType is 'Asin'"), nil
	}
	if args.ItemType == "Sku" && len(args.Skus) == 0 {
		return mcp.NewToolResultError("skus list is required when itemType is 'Sku'"), nil
	}

	// Validate list sizes (max 20 items)
	if len(args.Asins) > 20 {
		return mcp.NewToolResultError("maximum 20 ASINs allowed"), nil
	}
	if len(args.Skus) > 20 {
		return mcp.NewToolResultError("maximum 20 SKUs allowed"), nil
	}

	// Set up parameters
	params := &productPricing.GetPricingParams{
		MarketplaceId: args.MarketplaceID,
		ItemType:      args.ItemType,
	}

	if len(args.Asins) > 0 {
		params.Asins = &args.Asins
	}
	if len(args.Skus) > 0 {
		params.Skus = &args.Skus
	}
	if strings.TrimSpace(args.ItemCondition) != "" {
		params.ItemCondition = &args.ItemCondition
	}

	// Make API call
	httpResp, err := client.GetPricing(ctx, params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("GetPricing API call failed", err), nil
	}
	defer httpResp.Body.Close()

	return decodeProductPricingResponse(httpResp, "GetPricing")
}

func executeProductPricingGetCompetitivePricing(ctx context.Context, args productPricingGetCompetitivePricingArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureProductPricingClient(spClient)
	if failure != nil {
		return failure, nil
	}

	// Validation
	if strings.TrimSpace(args.MarketplaceID) == "" {
		return mcp.NewToolResultError("marketplaceId is required"), nil
	}
	if strings.TrimSpace(args.ItemType) == "" {
		return mcp.NewToolResultError("itemType is required"), nil
	}

	// Validate ItemType
	if args.ItemType != "Asin" && args.ItemType != "Sku" {
		return mcp.NewToolResultError("itemType must be 'Asin' or 'Sku'"), nil
	}

	// Validate that the appropriate identifier list is provided
	if args.ItemType == "Asin" && len(args.Asins) == 0 {
		return mcp.NewToolResultError("asins list is required when itemType is 'Asin'"), nil
	}
	if args.ItemType == "Sku" && len(args.Skus) == 0 {
		return mcp.NewToolResultError("skus list is required when itemType is 'Sku'"), nil
	}

	// Validate list sizes (max 20 items)
	if len(args.Asins) > 20 {
		return mcp.NewToolResultError("maximum 20 ASINs allowed"), nil
	}
	if len(args.Skus) > 20 {
		return mcp.NewToolResultError("maximum 20 SKUs allowed"), nil
	}

	// Set up parameters
	params := &productPricing.GetCompetitivePricingParams{
		MarketplaceId: args.MarketplaceID,
		ItemType:      args.ItemType,
	}

	if len(args.Asins) > 0 {
		params.Asins = &args.Asins
	}
	if len(args.Skus) > 0 {
		params.Skus = &args.Skus
	}

	// Make API call
	httpResp, err := client.GetCompetitivePricing(ctx, params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("GetCompetitivePricing API call failed", err), nil
	}
	defer httpResp.Body.Close()

	return decodeProductPricingResponse(httpResp, "GetCompetitivePricing")
}

func ensureProductPricingClient(spClient spapi.Client) (*productPricing.Client, *mcp.CallToolResult) {
	endpoint := spClient.Endpoint()
	client, err := productPricing.NewClient(endpoint)
	if err != nil {
		return nil, mcp.NewToolResultErrorFromErr("Failed to create Product Pricing client", err)
	}
	return client, nil
}