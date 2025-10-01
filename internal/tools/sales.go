package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sales "github.com/amzapi/selling-partner-api-sdk/sales"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

type salesGetOrderMetricsArgs struct {
	MarketplaceIDs      []string `json:"marketplaceIds"`
	Interval            string   `json:"interval"`
	Granularity         string   `json:"granularity"`
	GranularityTimeZone string   `json:"granularityTimeZone"`
	BuyerType           string   `json:"buyerType"`
	FulfillmentNetwork  string   `json:"fulfillmentNetwork"`
	FirstDayOfWeek      string   `json:"firstDayOfWeek"`
	ASIN                string   `json:"asin"`
	SKU                 string   `json:"sku"`
}

type salesGetOrderMetricsResult struct {
	MarketplaceIDs      []string                     `json:"marketplaceIds"`
	Interval            string                       `json:"interval"`
	Granularity         string                       `json:"granularity"`
	GranularityTimeZone string                       `json:"granularityTimeZone,omitempty"`
	BuyerType           string                       `json:"buyerType,omitempty"`
	FulfillmentNetwork  string                       `json:"fulfillmentNetwork,omitempty"`
	FirstDayOfWeek      string                       `json:"firstDayOfWeek,omitempty"`
	ASIN                string                       `json:"asin,omitempty"`
	SKU                 string                       `json:"sku,omitempty"`
	Metrics             []sales.OrderMetricsInterval `json:"metrics"`
	RetrievedAt         time.Time                    `json:"retrievedAt"`
}

func newSalesTools(deps Dependencies) []server.ServerTool {
	spClient := deps.SellingPartner

	orderMetricsHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args salesGetOrderMetricsArgs) (*mcp.CallToolResult, error) {
		return executeSalesGetOrderMetrics(ctx, args, spClient)
	})

	return []server.ServerTool{
		serverToolFromSpec(salesGetOrderMetricsSpec, orderMetricsHandler),
	}
}

func executeSalesGetOrderMetrics(ctx context.Context, args salesGetOrderMetricsArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureSalesClient(spClient)
	if failure != nil {
		return failure, nil
	}

	params, failure := prepareSalesGetOrderMetricsParams(args)
	if failure != nil {
		return failure, nil
	}

	httpResp, err := client.GetOrderMetrics(ctx, params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("sales.getOrderMetrics request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("sales.getOrderMetrics returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read sales.getOrderMetrics response", readErr), nil
	}

	decoded, decodeErr := decodeSalesOrderMetrics(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode sales.getOrderMetrics response", decodeErr), nil
	}

	if err := ensureSalesAPIResponse("getOrderMetrics", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("sales.getOrderMetrics response payload is empty"), nil
	}

	result := buildSalesMetricsResult(params, decoded.metrics)
	fallback := buildSalesMetricsFallback(result)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func prepareSalesGetOrderMetricsParams(args salesGetOrderMetricsArgs) (*sales.GetOrderMetricsParams, *mcp.CallToolResult) {
	marketplaces := trimStringSlice(args.MarketplaceIDs)
	if len(marketplaces) == 0 {
		return nil, mcp.NewToolResultError("marketplaceIds is required")
	}

	interval := strings.TrimSpace(args.Interval)
	if interval == "" {
		return nil, mcp.NewToolResultError("interval is required")
	}

	granularity := strings.TrimSpace(args.Granularity)
	if granularity == "" {
		return nil, mcp.NewToolResultError("granularity is required")
	}

	asin := strings.TrimSpace(args.ASIN)
	sku := strings.TrimSpace(args.SKU)
	if asin != "" && sku != "" {
		return nil, mcp.NewToolResultError("provide either asin or sku, not both")
	}

	granularityTimeZone := strings.TrimSpace(args.GranularityTimeZone)
	switch strings.ToUpper(granularity) {
	case "HOUR", "TOTAL":
		// optional timezone
	default:
		if granularityTimeZone == "" {
			return nil, mcp.NewToolResultError("granularityTimeZone is required when granularity is Day or greater")
		}
	}

	params := &sales.GetOrderMetricsParams{
		MarketplaceIds: marketplaces,
		Interval:       interval,
		Granularity:    granularity,
	}

	params.GranularityTimeZone = stringPtr(granularityTimeZone)
	params.BuyerType = stringPtr(args.BuyerType)
	params.FulfillmentNetwork = stringPtr(args.FulfillmentNetwork)
	params.FirstDayOfWeek = stringPtr(args.FirstDayOfWeek)
	params.Asin = stringPtr(asin)
	params.Sku = stringPtr(sku)

	return params, nil
}

func buildSalesMetricsResult(params *sales.GetOrderMetricsParams, metrics []sales.OrderMetricsInterval) salesGetOrderMetricsResult {
	result := salesGetOrderMetricsResult{
		MarketplaceIDs:      append([]string(nil), params.MarketplaceIds...),
		Interval:            params.Interval,
		Granularity:         params.Granularity,
		GranularityTimeZone: valueOrEmpty(params.GranularityTimeZone),
		BuyerType:           valueOrEmpty(params.BuyerType),
		FulfillmentNetwork:  valueOrEmpty(params.FulfillmentNetwork),
		FirstDayOfWeek:      valueOrEmpty(params.FirstDayOfWeek),
		ASIN:                valueOrEmpty(params.Asin),
		SKU:                 valueOrEmpty(params.Sku),
		Metrics:             make([]sales.OrderMetricsInterval, 0),
		RetrievedAt:         time.Now().UTC(),
	}

	if len(metrics) > 0 {
		result.Metrics = append(result.Metrics, metrics...)
	}

	return result
}

func buildSalesMetricsFallback(result salesGetOrderMetricsResult) string {
	var summary strings.Builder
	count := len(result.Metrics)
	granularity := strings.ToLower(strings.TrimSpace(result.Granularity))
	if granularity == "" {
		granularity = "requested"
	}

	marketplaces := "requested marketplaces"
	if len(result.MarketplaceIDs) > 0 {
		marketplaces = strings.Join(result.MarketplaceIDs, ", ")
	}

	summary.WriteString(fmt.Sprintf("Retrieved %d %s interval(s) for %s", count, granularity, marketplaces))

	if interval := strings.TrimSpace(result.Interval); interval != "" {
		summary.WriteString(" within ")
		summary.WriteString(interval)
	}

	if buyer := strings.TrimSpace(result.BuyerType); buyer != "" {
		summary.WriteString(" (buyer type: ")
		summary.WriteString(buyer)
		summary.WriteString(")")
	}

	return summary.String()
}

func ensureSalesClient(spClient spapi.Client) (*sales.Client, *mcp.CallToolResult) {
	if spClient == nil {
		return nil, mcp.NewToolResultError("Selling Partner API client is not initialised")
	}

	if status := spClient.Status(); !status.Ready {
		message := strings.TrimSpace(status.Message)
		if message == "" {
			message = "Selling Partner API client is not ready"
		}
		return nil, mcp.NewToolResultError(message)
	}

	client, err := buildSalesClient(spClient)
	if err != nil {
		return nil, mcp.NewToolResultErrorFromErr("failed to create sales client", err)
	}

	return client, nil
}

func buildSalesClient(spClient spapi.Client) (*sales.Client, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	return &sales.Client{
		Endpoint:      spClient.Endpoint(),
		Client:        httpClient,
		RequestBefore: buildSalesRequestBefore(spClient),
	}, nil
}

func buildSalesRequestBefore(spClient spapi.Client) sales.RequestBeforeFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Amzn-Requestid", uuid.NewString())
		req.Header.Set("Accept", "application/json")
		if err := spClient.AuthorizeRequest(req); err != nil {
			return fmt.Errorf("authorize request: %w", err)
		}
		return nil
	}
}

func ensureSalesAPIResponse(operation string, resp *http.Response, body []byte, errors *sales.ErrorList) error {
	if resp == nil {
		return fmt.Errorf("%s: no HTTP response returned", operation)
	}

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		// Try to extract detailed error messages from the ErrorList first
		if errors != nil && len(*errors) > 0 {
			return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), formatSalesErrors(*errors))
		}
		// Fall back to body snippet if no structured errors available
		return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), sanitizeBodySnippet(body))
	}

	if errors != nil && len(*errors) > 0 {
		return fmt.Errorf("%s: %s", operation, formatSalesErrors(*errors))
	}

	return nil
}

func formatSalesErrors(list sales.ErrorList) string {
	segments := make([]string, 0, len(list))
	for _, apiErr := range list {
		var builder strings.Builder
		builder.WriteString(strings.TrimSpace(apiErr.Message))
		if apiErr.Code != "" {
			builder.WriteString(" (" + apiErr.Code + ")")
		}
		if apiErr.Details != nil {
			detail := strings.TrimSpace(*apiErr.Details)
			if detail != "" {
				builder.WriteString(": " + detail)
			}
		}
		segments = append(segments, builder.String())
	}
	return strings.Join(segments, "; ")
}
