package tools

import (
	"strings"
	"testing"

	sales "github.com/amzapi/selling-partner-api-sdk/sales"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestPrepareSalesGetOrderMetricsParamsValidation(t *testing.T) {
	tests := []struct {
		name string
		args salesGetOrderMetricsArgs
		msg  string
	}{
		{
			name: "missing marketplaces",
			args: salesGetOrderMetricsArgs{},
			msg:  "marketplaceIds is required",
		},
		{
			name: "missing interval",
			args: salesGetOrderMetricsArgs{MarketplaceIDs: []string{"ATVPDKIKX0DER"}},
			msg:  "interval is required",
		},
		{
			name: "missing granularity",
			args: salesGetOrderMetricsArgs{MarketplaceIDs: []string{"ATVPDKIKX0DER"}, Interval: "2024-01-01T00:00:00Z--2024-01-02T00:00:00Z"},
			msg:  "granularity is required",
		},
		{
			name: "asin and sku provided",
			args: salesGetOrderMetricsArgs{
				MarketplaceIDs:      []string{"ATVPDKIKX0DER"},
				Interval:            "2024-01-01T00:00:00Z--2024-01-02T00:00:00Z",
				Granularity:         "Day",
				ASIN:                "B01234",
				SKU:                 "SKU123",
				GranularityTimeZone: "UTC",
			},
			msg: "provide either asin or sku, not both",
		},
		{
			name: "missing timezone",
			args: salesGetOrderMetricsArgs{
				MarketplaceIDs: []string{"ATVPDKIKX0DER"},
				Interval:       "2024-01-01T00:00:00Z--2024-01-02T00:00:00Z",
				Granularity:    "Day",
			},
			msg: "granularityTimeZone is required when granularity is Day or greater",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params, failure := prepareSalesGetOrderMetricsParams(tc.args)
			if failure == nil {
				t.Fatalf("expected failure, got params: %+v", params)
			}

			message := toolResultText(failure)
			if !strings.Contains(message, tc.msg) {
				t.Fatalf("expected error %q, got %q", tc.msg, message)
			}
		})
	}
}

func TestPrepareSalesGetOrderMetricsParamsSuccess(t *testing.T) {
	args := salesGetOrderMetricsArgs{
		MarketplaceIDs:      []string{" ATVPDKIKX0DER ", "A1PA6795UKMFR9"},
		Interval:            " 2024-01-01T00:00:00Z--2024-01-02T00:00:00Z ",
		Granularity:         "Day",
		GranularityTimeZone: " UTC ",
		BuyerType:           " B2B ",
		FulfillmentNetwork:  " AFN ",
		FirstDayOfWeek:      " Sunday ",
		ASIN:                " B01234 ",
	}

	params, failure := prepareSalesGetOrderMetricsParams(args)
	if failure != nil {
		t.Fatalf("unexpected failure: %s", toolResultText(failure))
	}

	if got := len(params.MarketplaceIds); got != 2 {
		t.Fatalf("expected 2 marketplaces, got %d", got)
	}

	if params.MarketplaceIds[0] != "ATVPDKIKX0DER" {
		t.Fatalf("unexpected first marketplace: %q", params.MarketplaceIds[0])
	}

	if params.Interval != "2024-01-01T00:00:00Z--2024-01-02T00:00:00Z" {
		t.Fatalf("unexpected interval: %q", params.Interval)
	}

	if params.Granularity != "Day" {
		t.Fatalf("unexpected granularity: %q", params.Granularity)
	}

	if params.GranularityTimeZone == nil || *params.GranularityTimeZone != "UTC" {
		t.Fatalf("unexpected timezone: %v", params.GranularityTimeZone)
	}

	if params.BuyerType == nil || *params.BuyerType != "B2B" {
		t.Fatalf("unexpected buyer type: %v", params.BuyerType)
	}

	if params.FulfillmentNetwork == nil || *params.FulfillmentNetwork != "AFN" {
		t.Fatalf("unexpected fulfillment network: %v", params.FulfillmentNetwork)
	}

	if params.FirstDayOfWeek == nil || *params.FirstDayOfWeek != "Sunday" {
		t.Fatalf("unexpected first day of week: %v", params.FirstDayOfWeek)
	}

	if params.Asin == nil || *params.Asin != "B01234" {
		t.Fatalf("unexpected asin: %v", params.Asin)
	}

	if params.Sku != nil {
		t.Fatalf("expected sku nil, got %v", params.Sku)
	}
}

func TestBuildSalesMetricsFallback(t *testing.T) {
	intervals := []sales.OrderMetricsInterval{{Interval: "2024-01-01T00:00:00Z--2024-01-02T00:00:00Z"}, {Interval: "2024-01-02T00:00:00Z--2024-01-03T00:00:00Z"}}
	result := salesGetOrderMetricsResult{
		MarketplaceIDs: []string{"ATVPDKIKX0DER"},
		Interval:       "2024-01-01T00:00:00Z--2024-01-03T00:00:00Z",
		Granularity:    "Day",
		BuyerType:      "B2B",
		Metrics:        intervals,
	}

	message := buildSalesMetricsFallback(result)
	if !strings.Contains(message, "Retrieved 2 day interval(s)") {
		t.Fatalf("unexpected summary: %s", message)
	}
	if !strings.Contains(message, "buyer type") {
		t.Fatalf("expected buyer type mention: %s", message)
	}
}

func TestFormatSalesErrors(t *testing.T) {
	message := formatSalesErrors(sales.ErrorList{
		{Code: "InvalidInput", Message: "Invalid request", Details: stringPtr("Start date is after end date")},
		{Code: "Unauthorized", Message: "Access denied"},
	})

	expected := "Invalid request (InvalidInput): Start date is after end date; Access denied (Unauthorized)"
	if message != expected {
		t.Fatalf("unexpected formatted error: %q", message)
	}
}

func toolResultText(result *mcp.CallToolResult) string {
	if result == nil {
		return ""
	}

	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			return textContent.Text
		}
	}

	return ""
}
