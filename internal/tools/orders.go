package tools

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ordersv0 "github.com/amzapi/selling-partner-api-sdk/ordersV0"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

type ordersListOrdersArgs struct {
	MarketplaceIDs           []string `json:"marketplaceIds"`
	CreatedAfter             string   `json:"createdAfter"`
	CreatedBefore            string   `json:"createdBefore"`
	LastUpdatedAfter         string   `json:"lastUpdatedAfter"`
	LastUpdatedBefore        string   `json:"lastUpdatedBefore"`
	OrderStatuses            []string `json:"orderStatuses"`
	FulfillmentChannels      []string `json:"fulfillmentChannels"`
	PaymentMethods           []string `json:"paymentMethods"`
	BuyerEmail               string   `json:"buyerEmail"`
	SellerOrderID            string   `json:"sellerOrderId"`
	MaxResultsPerPage        *int     `json:"maxResultsPerPage"`
	EasyShipShipmentStatuses []string `json:"easyShipShipmentStatuses"`
	AmazonOrderIDs           []string `json:"amazonOrderIds"`
	NextToken                string   `json:"nextToken"`
}

type ordersListOrdersResult struct {
	Orders            []ordersv0.Order `json:"orders"`
	NextToken         string           `json:"nextToken,omitempty"`
	CreatedBefore     string           `json:"createdBefore,omitempty"`
	LastUpdatedBefore string           `json:"lastUpdatedBefore,omitempty"`
	RetrievedAt       time.Time        `json:"retrievedAt"`
}

type ordersGetOrderArgs struct {
	AmazonOrderID string `json:"amazonOrderId"`
}

type ordersGetOrderResult struct {
	AmazonOrderID  string               `json:"amazonOrderId"`
	OrderStatus    string               `json:"orderStatus"`
	PurchaseDate   string               `json:"purchaseDate"`
	LastUpdateDate string               `json:"lastUpdateDate"`
	ItemCount      int                  `json:"itemCount"`
	RetrievedAt    time.Time            `json:"retrievedAt"`
	Order          ordersv0.Order       `json:"order"`
	OrderItems     []ordersv0.OrderItem `json:"orderItems"`
}

type ordersGetOrderAddressResult struct {
	AmazonOrderID   string            `json:"amazonOrderId"`
	ShippingAddress *ordersv0.Address `json:"shippingAddress,omitempty"`
	RetrievedAt     time.Time         `json:"retrievedAt"`
}

type ordersGetOrderBuyerInfoResult struct {
	AmazonOrderID string                   `json:"amazonOrderId"`
	BuyerInfo     *ordersv0.OrderBuyerInfo `json:"buyerInfo,omitempty"`
	RetrievedAt   time.Time                `json:"retrievedAt"`
}

type ordersGetOrderItemsArgs struct {
	AmazonOrderID string `json:"amazonOrderId"`
	NextToken     string `json:"nextToken"`
}

type ordersGetOrderItemsResult struct {
	AmazonOrderID string               `json:"amazonOrderId"`
	Items         []ordersv0.OrderItem `json:"items"`
	NextToken     string               `json:"nextToken,omitempty"`
	RetrievedAt   time.Time            `json:"retrievedAt"`
}

type ordersGetOrderItemsBuyerInfoResult struct {
	AmazonOrderID string                        `json:"amazonOrderId"`
	Items         []ordersv0.OrderItemBuyerInfo `json:"items"`
	NextToken     string                        `json:"nextToken,omitempty"`
	RetrievedAt   time.Time                     `json:"retrievedAt"`
}

func newOrdersTools(deps Dependencies) []server.ServerTool {
	spClient := deps.SellingPartner

	listOrdersHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersListOrdersArgs) (*mcp.CallToolResult, error) {
		return executeOrdersListOrders(ctx, args, spClient)
	})

	getOrderHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersGetOrderArgs) (*mcp.CallToolResult, error) {
		return executeOrdersGetOrder(ctx, strings.TrimSpace(args.AmazonOrderID), spClient)
	})

	getOrderAddressHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersGetOrderArgs) (*mcp.CallToolResult, error) {
		return executeOrdersGetOrderAddress(ctx, strings.TrimSpace(args.AmazonOrderID), spClient)
	})

	getOrderBuyerInfoHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersGetOrderArgs) (*mcp.CallToolResult, error) {
		return executeOrdersGetOrderBuyerInfo(ctx, strings.TrimSpace(args.AmazonOrderID), spClient)
	})

	getOrderItemsHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersGetOrderItemsArgs) (*mcp.CallToolResult, error) {
		return executeOrdersGetOrderItems(ctx, args, spClient)
	})

	getOrderItemsBuyerInfoHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args ordersGetOrderItemsArgs) (*mcp.CallToolResult, error) {
		return executeOrdersGetOrderItemsBuyerInfo(ctx, args, spClient)
	})

	return []server.ServerTool{
		serverToolFromSpec(ordersListOrdersSpec, listOrdersHandler),
		serverToolFromSpec(ordersGetOrderSpec, getOrderHandler),
		serverToolFromSpec(ordersGetOrderAddressSpec, getOrderAddressHandler),
		serverToolFromSpec(ordersGetOrderBuyerInfoSpec, getOrderBuyerInfoHandler),
		serverToolFromSpec(ordersGetOrderItemsSpec, getOrderItemsHandler),
		serverToolFromSpec(ordersGetOrderItemsBuyerInfoSpec, getOrderItemsBuyerInfoHandler),
	}
}

func executeOrdersListOrders(ctx context.Context, args ordersListOrdersArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	nextToken := strings.TrimSpace(args.NextToken)
	params := ordersv0.GetOrdersParams{}

	if nextToken != "" {
		if hasListOrdersFilters(args) {
			return mcp.NewToolResultError("when nextToken is provided, omit additional filters"), nil
		}
		params.NextToken = stringPtr(nextToken)
	} else {
		marketplaces := trimStringSlice(args.MarketplaceIDs)
		if len(marketplaces) == 0 {
			return mcp.NewToolResultError("marketplaceIds is required unless nextToken is provided"), nil
		}
		params.MarketplaceIds = marketplaces
		params.CreatedAfter = stringPtr(args.CreatedAfter)
		params.CreatedBefore = stringPtr(args.CreatedBefore)
		params.LastUpdatedAfter = stringPtr(args.LastUpdatedAfter)
		params.LastUpdatedBefore = stringPtr(args.LastUpdatedBefore)
		params.OrderStatuses = stringSlicePtr(args.OrderStatuses)
		params.FulfillmentChannels = stringSlicePtr(args.FulfillmentChannels)
		params.PaymentMethods = stringSlicePtr(args.PaymentMethods)
		params.BuyerEmail = stringPtr(args.BuyerEmail)
		params.SellerOrderId = stringPtr(args.SellerOrderID)
		params.EasyShipShipmentStatuses = stringSlicePtr(args.EasyShipShipmentStatuses)
		params.AmazonOrderIds = stringSlicePtr(args.AmazonOrderIDs)
		params.MaxResultsPerPage = sanitizeMaxResults(args.MaxResultsPerPage)
	}

	resp, err := client.GetOrdersWithResponse(ctx, &params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.listOrders request failed", err), nil
	}
	if resp == nil {
		return mcp.NewToolResultError("orders.listOrders returned no response"), nil
	}

	var apiErrors *ordersv0.ErrorList
	if resp.Model != nil {
		apiErrors = resp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("listOrders", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
		log.Printf("[ERROR] orders.listOrders: %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	if resp.Model == nil || resp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.listOrders response payload is empty"), nil
	}

	payload := resp.Model.Payload
	result := ordersListOrdersResult{
		Orders:            payload.Orders,
		NextToken:         valueOrEmpty(payload.NextToken),
		CreatedBefore:     valueOrEmpty(payload.CreatedBefore),
		LastUpdatedBefore: valueOrEmpty(payload.LastUpdatedBefore),
		RetrievedAt:       time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved %d orders", len(result.Orders))
	if result.NextToken != "" {
		fallback = fmt.Sprintf("%s, more available via nextToken", fallback)
	}

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeOrdersGetOrder(ctx context.Context, orderID string, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	if orderID == "" {
		return mcp.NewToolResultError("amazonOrderId is required"), nil
	}

	orderResp, err := client.GetOrderWithResponse(ctx, orderID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.getOrder request failed", err), nil
	}
	if orderResp == nil {
		return mcp.NewToolResultError("orders.getOrder returned no response"), nil
	}

	var orderErrors *ordersv0.ErrorList
	if orderResp.Model != nil {
		orderErrors = orderResp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("getOrder", orderResp.HTTPResponse, orderResp.Body, orderErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if orderResp.Model == nil || orderResp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.getOrder response payload is empty"), nil
	}

	order := orderResp.Model.Payload

	items, err := fetchAllOrderItems(ctx, client, orderID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to retrieve order items", err), nil
	}

	result := ordersGetOrderResult{
		AmazonOrderID:  order.AmazonOrderId,
		OrderStatus:    strings.TrimSpace(order.OrderStatus),
		PurchaseDate:   strings.TrimSpace(order.PurchaseDate),
		LastUpdateDate: strings.TrimSpace(order.LastUpdateDate),
		ItemCount:      len(items),
		RetrievedAt:    time.Now().UTC(),
		Order:          *order,
		OrderItems:     items,
	}

	fallback := buildOrderFallback(result)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeOrdersGetOrderAddress(ctx context.Context, orderID string, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	if orderID == "" {
		return mcp.NewToolResultError("amazonOrderId is required"), nil
	}

	resp, err := client.GetOrderAddressWithResponse(ctx, orderID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.getOrderAddress request failed", err), nil
	}
	if resp == nil {
		return mcp.NewToolResultError("orders.getOrderAddress returned no response"), nil
	}

	var apiErrors *ordersv0.ErrorList
	if resp.Model != nil {
		apiErrors = resp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("getOrderAddress", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if resp.Model == nil || resp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.getOrderAddress response payload is empty"), nil
	}

	payload := resp.Model.Payload
	result := ordersGetOrderAddressResult{
		AmazonOrderID:   payload.AmazonOrderId,
		ShippingAddress: payload.ShippingAddress,
		RetrievedAt:     time.Now().UTC(),
	}

	summaryAddress := "unknown"
	if payload.ShippingAddress != nil && payload.ShippingAddress.Name != "" {
		summaryAddress = payload.ShippingAddress.Name
	}

	fallback := fmt.Sprintf("Retrieved shipping address for order %s (%s)", result.AmazonOrderID, summaryAddress)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeOrdersGetOrderBuyerInfo(ctx context.Context, orderID string, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	if orderID == "" {
		return mcp.NewToolResultError("amazonOrderId is required"), nil
	}

	resp, err := client.GetOrderBuyerInfoWithResponse(ctx, orderID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.getOrderBuyerInfo request failed", err), nil
	}
	if resp == nil {
		return mcp.NewToolResultError("orders.getOrderBuyerInfo returned no response"), nil
	}

	var apiErrors *ordersv0.ErrorList
	if resp.Model != nil {
		apiErrors = resp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("getOrderBuyerInfo", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if resp.Model == nil || resp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.getOrderBuyerInfo response payload is empty"), nil
	}

	payload := resp.Model.Payload
	result := ordersGetOrderBuyerInfoResult{
		AmazonOrderID: payload.AmazonOrderId,
		BuyerInfo:     payload,
		RetrievedAt:   time.Now().UTC(),
	}

	buyerName := "unknown"
	if payload != nil && payload.BuyerName != nil && strings.TrimSpace(*payload.BuyerName) != "" {
		buyerName = strings.TrimSpace(*payload.BuyerName)
	}

	fallback := fmt.Sprintf("Retrieved buyer info for order %s (%s)", result.AmazonOrderID, buyerName)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeOrdersGetOrderItems(ctx context.Context, args ordersGetOrderItemsArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	orderID := strings.TrimSpace(args.AmazonOrderID)
	if orderID == "" {
		return mcp.NewToolResultError("amazonOrderId is required"), nil
	}

	params := ordersv0.GetOrderItemsParams{}
	nextToken := strings.TrimSpace(args.NextToken)
	if nextToken != "" {
		params.NextToken = &nextToken
	}

	resp, err := client.GetOrderItemsWithResponse(ctx, orderID, &params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.getOrderItems request failed", err), nil
	}
	if resp == nil {
		return mcp.NewToolResultError("orders.getOrderItems returned no response"), nil
	}

	var apiErrors *ordersv0.ErrorList
	if resp.Model != nil {
		apiErrors = resp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("getOrderItems", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if resp.Model == nil || resp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.getOrderItems response payload is empty"), nil
	}

	payload := resp.Model.Payload
	result := ordersGetOrderItemsResult{
		AmazonOrderID: payload.AmazonOrderId,
		Items:         payload.OrderItems,
		NextToken:     valueOrEmpty(payload.NextToken),
		RetrievedAt:   time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved %d items for order %s", len(result.Items), result.AmazonOrderID)
	if result.NextToken != "" {
		fallback = fmt.Sprintf("%s, more available via nextToken", fallback)
	}

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeOrdersGetOrderItemsBuyerInfo(ctx context.Context, args ordersGetOrderItemsArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureOrdersClient(spClient)
	if failure != nil {
		return failure, nil
	}

	orderID := strings.TrimSpace(args.AmazonOrderID)
	if orderID == "" {
		return mcp.NewToolResultError("amazonOrderId is required"), nil
	}

	params := ordersv0.GetOrderItemsBuyerInfoParams{}
	nextToken := strings.TrimSpace(args.NextToken)
	if nextToken != "" {
		params.NextToken = &nextToken
	}

	resp, err := client.GetOrderItemsBuyerInfoWithResponse(ctx, orderID, &params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("orders.getOrderItemsBuyerInfo request failed", err), nil
	}
	if resp == nil {
		return mcp.NewToolResultError("orders.getOrderItemsBuyerInfo returned no response"), nil
	}

	var apiErrors *ordersv0.ErrorList
	if resp.Model != nil {
		apiErrors = resp.Model.Errors
	}

	if err := ensureOrdersAPIResponse("getOrderItemsBuyerInfo", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if resp.Model == nil || resp.Model.Payload == nil {
		return mcp.NewToolResultError("orders.getOrderItemsBuyerInfo response payload is empty"), nil
	}

	payload := resp.Model.Payload
	result := ordersGetOrderItemsBuyerInfoResult{
		AmazonOrderID: payload.AmazonOrderId,
		Items:         payload.OrderItems,
		NextToken:     valueOrEmpty(payload.NextToken),
		RetrievedAt:   time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved buyer info for %d items on order %s", len(result.Items), result.AmazonOrderID)
	if result.NextToken != "" {
		fallback = fmt.Sprintf("%s, more available via nextToken", fallback)
	}

	return mcp.NewToolResultStructured(result, fallback), nil
}

// buildOrdersClient constructs an Orders API client that reuses the shared Selling Partner authentication state.
func buildOrdersClient(spClient spapi.Client) (*ordersv0.ClientWithResponses, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	return ordersv0.NewClientWithResponses(
		spClient.Endpoint(),
		withOrdersRequestBefore(spClient),
		ordersv0.WithHTTPClient(httpClient),
	)
}

func withOrdersRequestBefore(spClient spapi.Client) ordersv0.ClientOption {
	return ordersv0.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Amzn-Requestid", uuid.NewString())
		req.Header.Set("Accept", "application/json")
		if err := spClient.AuthorizeRequest(req); err != nil {
			return fmt.Errorf("authorize request: %w", err)
		}
		return nil
	})
}

// fetchAllOrderItems walks the Orders API pagination to collect every line item.
func fetchAllOrderItems(ctx context.Context, client *ordersv0.ClientWithResponses, orderID string) ([]ordersv0.OrderItem, error) {
	var (
		items        []ordersv0.OrderItem
		nextToken    string
		hasNextToken bool
	)

	for {
		var params *ordersv0.GetOrderItemsParams
		if hasNextToken {
			params = &ordersv0.GetOrderItemsParams{NextToken: &nextToken}
		}

		resp, err := client.GetOrderItemsWithResponse(ctx, orderID, params)
		if err != nil {
			return nil, fmt.Errorf("calling getOrderItems: %w", err)
		}
		if resp == nil {
			return nil, fmt.Errorf("getOrderItems returned no response")
		}

		var apiErrors *ordersv0.ErrorList
		if resp.Model != nil {
			apiErrors = resp.Model.Errors
		}

		if err := ensureOrdersAPIResponse("getOrderItems", resp.HTTPResponse, resp.Body, apiErrors); err != nil {
			return nil, err
		}

		if resp.Model == nil || resp.Model.Payload == nil {
			return nil, fmt.Errorf("getOrderItems response payload is empty")
		}

		payload := resp.Model.Payload
		items = append(items, payload.OrderItems...)

		if payload.NextToken != nil {
			value := strings.TrimSpace(*payload.NextToken)
			hasNextToken = value != ""
			nextToken = value
		} else {
			hasNextToken = false
			nextToken = ""
		}

		if !hasNextToken {
			break
		}
	}

	return items, nil
}

func ensureOrdersAPIResponse(operation string, resp *http.Response, body []byte, errors *ordersv0.ErrorList) error {
	if resp == nil {
		return fmt.Errorf("%s: no HTTP response returned", operation)
	}

	verbose := os.Getenv("VERBOSE") == "true"

	if verbose {
		log.Printf("[DEBUG] %s - HTTP Status: %d %s", operation, resp.StatusCode, http.StatusText(resp.StatusCode))
		log.Printf("[DEBUG] %s - Response Headers: %v", operation, resp.Header)
		log.Printf("[DEBUG] %s - Response Body: %s", operation, string(body))
	}

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		// Try to extract detailed error messages from the ErrorList first
		if errors != nil && len(*errors) > 0 {
			return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), formatOrdersErrors(*errors))
		}
		// Fall back to body snippet if no structured errors available
		return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), sanitizeBodySnippet(body))
	}

	if errors != nil && len(*errors) > 0 {
		return fmt.Errorf("%s: %s", operation, formatOrdersErrors(*errors))
	}

	return nil
}

func ensureOrdersClient(spClient spapi.Client) (*ordersv0.ClientWithResponses, *mcp.CallToolResult) {
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

	client, err := buildOrdersClient(spClient)
	if err != nil {
		return nil, mcp.NewToolResultErrorFromErr("failed to create orders client", err)
	}

	return client, nil
}

func sanitizeBodySnippet(body []byte) string {
	snippet := strings.TrimSpace(string(body))
	if snippet == "" {
		return "no response body"
	}
	if len(snippet) > 512 {
		snippet = snippet[:512] + "..."
	}
	return snippet
}

func formatOrdersErrors(list ordersv0.ErrorList) string {
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

func buildOrderFallback(result ordersGetOrderResult) string {
	status := result.OrderStatus
	if status == "" {
		status = "unknown"
	}

	fallback := fmt.Sprintf("Fetched Amazon order %s with %d items (status: %s)", result.AmazonOrderID, result.ItemCount, status)

	if total := formatMoney(result.Order.OrderTotal); total != "" {
		fallback = fmt.Sprintf("%s, total %s", fallback, total)
	}

	if date := strings.TrimSpace(result.PurchaseDate); date != "" {
		fallback = fmt.Sprintf("%s, purchased %s", fallback, date)
	}

	return fallback
}

func formatMoney(money *ordersv0.Money) string {
	if money == nil || money.Amount == nil || money.CurrencyCode == nil {
		return ""
	}
	amount := strings.TrimSpace(*money.Amount)
	currency := strings.TrimSpace(*money.CurrencyCode)
	if amount == "" && currency == "" {
		return ""
	}
	return strings.TrimSpace(amount + " " + currency)
}

func stringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func stringSlicePtr(values []string) *[]string {
	trimmed := trimStringSlice(values)
	if len(trimmed) == 0 {
		return nil
	}
	return &trimmed
}

func trimStringSlice(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func valueOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return strings.TrimSpace(*ptr)
}

func hasListOrdersFilters(args ordersListOrdersArgs) bool {
	if len(trimStringSlice(args.MarketplaceIDs)) > 0 {
		return true
	}
	if strings.TrimSpace(args.CreatedAfter) != "" || strings.TrimSpace(args.CreatedBefore) != "" {
		return true
	}
	if strings.TrimSpace(args.LastUpdatedAfter) != "" || strings.TrimSpace(args.LastUpdatedBefore) != "" {
		return true
	}
	if len(trimStringSlice(args.OrderStatuses)) > 0 || len(trimStringSlice(args.FulfillmentChannels)) > 0 || len(trimStringSlice(args.PaymentMethods)) > 0 {
		return true
	}
	if strings.TrimSpace(args.BuyerEmail) != "" || strings.TrimSpace(args.SellerOrderID) != "" {
		return true
	}
	if args.MaxResultsPerPage != nil {
		return true
	}
	if len(trimStringSlice(args.EasyShipShipmentStatuses)) > 0 || len(trimStringSlice(args.AmazonOrderIDs)) > 0 {
		return true
	}
	return false
}

func sanitizeMaxResults(value *int) *int {
	if value == nil {
		return nil
	}
	if *value < 1 || *value > 100 {
		return nil
	}
	return value
}
