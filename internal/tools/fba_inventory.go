package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amzapi/selling-partner-api-sdk/fbaInventory"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

type fbaInventoryGetInventorySummariesArgs struct {
	GranularityType string   `json:"granularityType"`
	GranularityID   string   `json:"granularityId"`
	StartDateTime   string   `json:"startDateTime"`
	SellerSkus      []string `json:"sellerSkus"`
	NextToken       string   `json:"nextToken"`
	Details         *bool    `json:"details"`
}

type fbaInventoryGetInventorySummariesResult struct {
	GranularityType    string                           `json:"granularityType"`
	GranularityID      string                           `json:"granularityId,omitempty"`
	InventorySummaries []fbaInventory.InventorySummary  `json:"inventorySummaries"`
	NextToken          string                           `json:"nextToken,omitempty"`
	RetrievedAt        time.Time                        `json:"retrievedAt"`
}

func newFBAInventoryTools(deps Dependencies) []server.ServerTool {
	spClient := deps.SellingPartner

	getInventorySummariesHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args fbaInventoryGetInventorySummariesArgs) (*mcp.CallToolResult, error) {
		return executeFBAInventoryGetInventorySummaries(ctx, args, spClient)
	})

	return []server.ServerTool{
		serverToolFromSpec(fbaInventoryGetInventorySummariesSpec, getInventorySummariesHandler),
	}
}

func executeFBAInventoryGetInventorySummaries(ctx context.Context, args fbaInventoryGetInventorySummariesArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureFBAInventoryClient(spClient)
	if failure != nil {
		return failure, nil
	}

	granularityType := strings.TrimSpace(args.GranularityType)
	if granularityType == "" {
		return mcp.NewToolResultError("granularityType is required"), nil
	}

	params := &fbaInventory.GetInventorySummariesParams{
		GranularityType: granularityType,
	}

	if granularityID := strings.TrimSpace(args.GranularityID); granularityID != "" {
		params.GranularityId = granularityID
	}

	if startDateTime := strings.TrimSpace(args.StartDateTime); startDateTime != "" {
		if parsedTime, err := time.Parse(time.RFC3339, startDateTime); err == nil {
			params.StartDateTime = &parsedTime
		} else {
			return mcp.NewToolResultError("startDateTime must be in ISO 8601 format"), nil
		}
	}

	if len(args.SellerSkus) > 0 {
		sellerSkus := trimStringSlice(args.SellerSkus)
		if len(sellerSkus) > 0 {
			params.SellerSkus = &sellerSkus
		}
	}

	if nextToken := strings.TrimSpace(args.NextToken); nextToken != "" {
		params.NextToken = &nextToken
	}

	if args.Details != nil {
		params.Details = args.Details
	}

	httpResp, err := client.GetInventorySummaries(ctx, params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("fbaInventory.getInventorySummaries request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("fbaInventory.getInventorySummaries returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read fbaInventory.getInventorySummaries response", readErr), nil
	}

	decoded, decodeErr := decodeFBAInventoryGetInventorySummaries(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode fbaInventory.getInventorySummaries response", decodeErr), nil
	}

	if err := ensureFBAInventoryAPIResponse("getInventorySummaries", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("fbaInventory.getInventorySummaries response payload is empty"), nil
	}

	result := fbaInventoryGetInventorySummariesResult{
		GranularityType:    decoded.granularityType,
		GranularityID:      decoded.granularityID,
		InventorySummaries: decoded.inventorySummaries,
		NextToken:          decoded.nextToken,
		RetrievedAt:        time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved %d inventory summaries", len(result.InventorySummaries))
	if result.NextToken != "" {
		fallback = fmt.Sprintf("%s, more available via nextToken", fallback)
	}
	if result.GranularityType != "" {
		fallback = fmt.Sprintf("%s (granularity: %s)", fallback, result.GranularityType)
	}

	return mcp.NewToolResultStructured(result, fallback), nil
}

func ensureFBAInventoryClient(spClient spapi.Client) (*fbaInventory.Client, *mcp.CallToolResult) {
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

	client, err := buildFBAInventoryClient(spClient)
	if err != nil {
		return nil, mcp.NewToolResultErrorFromErr("failed to create fbaInventory client", err)
	}

	return client, nil
}

func buildFBAInventoryClient(spClient spapi.Client) (*fbaInventory.Client, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	return &fbaInventory.Client{
		Endpoint:      spClient.Endpoint(),
		Client:        httpClient,
		RequestBefore: buildFBAInventoryRequestBefore(spClient),
	}, nil
}

func buildFBAInventoryRequestBefore(spClient spapi.Client) fbaInventory.RequestBeforeFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Amzn-Requestid", uuid.NewString())
		req.Header.Set("Accept", "application/json")
		if err := spClient.AuthorizeRequest(req); err != nil {
			return fmt.Errorf("authorize request: %w", err)
		}
		return nil
	}
}

func ensureFBAInventoryAPIResponse(operation string, resp *http.Response, body []byte, errors *fbaInventory.ErrorList) error {
	if resp == nil {
		return fmt.Errorf("%s: no HTTP response returned", operation)
	}

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), sanitizeBodySnippet(body))
	}

	if errors != nil && len(*errors) > 0 {
		return fmt.Errorf("%s: %s", operation, formatFBAInventoryErrors(*errors))
	}

	return nil
}

func formatFBAInventoryErrors(list fbaInventory.ErrorList) string {
	segments := make([]string, 0, len(list))
	for _, apiErr := range list {
		var builder strings.Builder
		if apiErr.Message != nil {
			builder.WriteString(strings.TrimSpace(*apiErr.Message))
		}
		if apiErr.Code != "" {
			builder.WriteString(" (" + apiErr.Code + ")")
		}
		if apiErr.Details != nil {
			detail := strings.TrimSpace(*apiErr.Details)
			if detail != "" {
				builder.WriteString(": " + detail)
			}
		}
		if text := builder.String(); text != "" {
			segments = append(segments, text)
		}
	}
	return strings.Join(segments, "; ")
}