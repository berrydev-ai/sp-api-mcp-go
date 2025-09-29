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
	intervals := []sales.OrderMetricsInterval{
		{
			Interval: "2025-09-01T00:00Z--2025-09-02T00:00Z",
			AverageUnitPrice: sales.Money{Amount: sales.Decimal("73.57"), CurrencyCode: "USD"},
			OrderCount: 100,
			TotalSales: sales.Money{Amount: sales.Decimal("8975.55"), CurrencyCode: "USD"},
		},
		{
			Interval: "2025-09-02T00:00Z--2025-09-03T00:00Z", 
			AverageUnitPrice: sales.Money{Amount: sales.Decimal("73.72"), CurrencyCode: "USD"},
			OrderCount: 80,
			TotalSales: sales.Money{Amount: sales.Decimal("6856.09"), CurrencyCode: "USD"},
		},
	}
	result := salesGetOrderMetricsResult{
		MarketplaceIDs: []string{"ATVPDKIKX0DER"},
		Interval:       "2025-09-01T00:00:00Z--2025-09-03T00:00:00Z",
		Granularity:    "Day",
		BuyerType:      "B2B",
		Metrics:        intervals,
	}

	message := buildSalesMetricsFallback(result)
	if !strings.Contains(message, "Retrieved 2 day interval(s)") {
		t.Fatalf("unexpected summary: %s", message)
	}
	if !strings.Contains(message, "ATVPDKIKX0DER") {
		t.Fatalf("expected marketplace ID: %s", message)
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

func TestDecodeSalesOrderMetricsHandlesNumericAmounts(t *testing.T) {
	body := []byte(`{
		"payload": [
			{
				"averageUnitPrice": {"amount": 12.34, "currencyCode": "USD"},
				"interval": "2025-01-01T00:00:00Z--2025-01-02T00:00:00Z",
				"orderCount": 5,
				"orderItemCount": 6,
				"totalSales": {"amount": 61.7, "currencyCode": "USD"},
				"unitCount": 10
			}
		]
	}`)

	decoded, err := decodeSalesOrderMetrics(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !decoded.payloadPresent {
		t.Fatalf("expected payload to be present")
	}

	if len(decoded.metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(decoded.metrics))
	}

	metric := decoded.metrics[0]
	if metric.AverageUnitPrice.Amount != sales.Decimal("12.34") {
		t.Fatalf("unexpected average unit price amount: %s", metric.AverageUnitPrice.Amount)
	}
	if metric.TotalSales.Amount != sales.Decimal("61.7") {
		t.Fatalf("unexpected total sales amount: %s", metric.TotalSales.Amount)
	}
}

func TestDecodeSalesOrderMetricsHandlesStringAmounts(t *testing.T) {
	body := []byte(`{
		"payload": [
			{
				"averageUnitPrice": {"amount": "73.57", "currencyCode": "USD"},
				"interval": "2025-09-01T00:00Z--2025-09-02T00:00Z",
				"orderCount": 100,
				"orderItemCount": 113,
				"totalSales": {"amount": "8975.55", "currencyCode": "USD"},
				"unitCount": 122
			}
		]
	}`)

	decoded, err := decodeSalesOrderMetrics(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !decoded.payloadPresent {
		t.Fatalf("expected payload to be present")
	}

	if len(decoded.metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(decoded.metrics))
	}

	metric := decoded.metrics[0]
	if metric.AverageUnitPrice.Amount != sales.Decimal("73.57") {
		t.Fatalf("unexpected average unit price amount: %s", metric.AverageUnitPrice.Amount)
	}
	if metric.TotalSales.Amount != sales.Decimal("8975.55") {
		t.Fatalf("unexpected total sales amount: %s", metric.TotalSales.Amount)
	}
	if metric.AverageUnitPrice.CurrencyCode != "USD" {
		t.Fatalf("unexpected currency code: %s", metric.AverageUnitPrice.CurrencyCode)
	}
	if metric.Interval != "2025-09-01T00:00Z--2025-09-02T00:00Z" {
		t.Fatalf("unexpected interval: %s", metric.Interval)
	}
	if metric.OrderCount != 100 {
		t.Fatalf("unexpected order count: %d", metric.OrderCount)
	}
	if metric.OrderItemCount != 113 {
		t.Fatalf("unexpected order item count: %d", metric.OrderItemCount)
	}
	if metric.UnitCount != 122 {
		t.Fatalf("unexpected unit count: %d", metric.UnitCount)
	}
}

func TestDecodeSalesOrderMetricsRealApiResponse(t *testing.T) {
	// Test with actual API response structure (multiple intervals)
	body := []byte(`{
		"payload": [
			{
				"averageUnitPrice": {"amount": "73.57", "currencyCode": "USD"},
				"interval": "2025-09-01T00:00Z--2025-09-02T00:00Z",
				"orderCount": 100,
				"orderItemCount": 113,
				"totalSales": {"amount": "8975.55", "currencyCode": "USD"},
				"unitCount": 122
			},
			{
				"averageUnitPrice": {"amount": "73.72", "currencyCode": "USD"},
				"interval": "2025-09-02T00:00Z--2025-09-03T00:00Z",
				"orderCount": 80,
				"orderItemCount": 90,
				"totalSales": {"amount": "6856.09", "currencyCode": "USD"},
				"unitCount": 93
			}
		]
	}`)

	decoded, err := decodeSalesOrderMetrics(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !decoded.payloadPresent {
		t.Fatalf("expected payload to be present")
	}

	if len(decoded.metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(decoded.metrics))
	}

	// Test first metric
	metric1 := decoded.metrics[0]
	if metric1.AverageUnitPrice.Amount != sales.Decimal("73.57") {
		t.Fatalf("unexpected first metric average unit price: %s", metric1.AverageUnitPrice.Amount)
	}
	if metric1.OrderCount != 100 {
		t.Fatalf("unexpected first metric order count: %d", metric1.OrderCount)
	}

	// Test second metric
	metric2 := decoded.metrics[1]
	if metric2.AverageUnitPrice.Amount != sales.Decimal("73.72") {
		t.Fatalf("unexpected second metric average unit price: %s", metric2.AverageUnitPrice.Amount)
	}
	if metric2.OrderCount != 80 {
		t.Fatalf("unexpected second metric order count: %d", metric2.OrderCount)
	}
}

func TestDecodeSalesOrderMetricsMissingPayload(t *testing.T) {
	body := []byte(`{"errors": [{"code": "InvalidInput", "message": "bad"}]}`)

	decoded, err := decodeSalesOrderMetrics(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decoded.payloadPresent {
		t.Fatalf("expected payload to be absent")
	}

	if decoded.apiErrors == nil || len(*decoded.apiErrors) != 1 {
		t.Fatalf("expected one api error, got %v", decoded.apiErrors)
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
