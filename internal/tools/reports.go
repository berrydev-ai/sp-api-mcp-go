package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amzapi/selling-partner-api-sdk/reports"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/berrydev-ai/sp-api-mcp-go/internal/spapi"
)

type reportsGetReportsArgs struct {
	ReportTypes              []string `json:"reportTypes"`
	ProcessingStatuses       []string `json:"processingStatuses"`
	MarketplaceIDs           []string `json:"marketplaceIds"`
	PageSize                 *int     `json:"pageSize"`
	CreatedSince             string   `json:"createdSince"`
	CreatedUntil             string   `json:"createdUntil"`
	NextToken                string   `json:"nextToken"`
}

type reportsGetReportsResult struct {
	Reports     []reports.Report `json:"reports"`
	NextToken   string           `json:"nextToken,omitempty"`
	RetrievedAt time.Time        `json:"retrievedAt"`
}

type reportsCreateReportArgs struct {
	ReportType      string            `json:"reportType"`
	MarketplaceIDs  []string          `json:"marketplaceIds"`
	DataStartTime   string            `json:"dataStartTime"`
	DataEndTime     string            `json:"dataEndTime"`
	ReportOptions   map[string]string `json:"reportOptions"`
}

type reportsCreateReportResult struct {
	ReportID    string    `json:"reportId"`
	ReportType  string    `json:"reportType"`
	RetrievedAt time.Time `json:"retrievedAt"`
}

type reportsGetReportArgs struct {
	ReportID string `json:"reportId"`
}

type reportsGetReportResult struct {
	ReportID        string    `json:"reportId"`
	ReportType      string    `json:"reportType"`
	ProcessingStatus string   `json:"processingStatus"`
	CreatedTime     string    `json:"createdTime"`
	ProcessingStartTime string `json:"processingStartTime,omitempty"`
	ProcessingEndTime   string `json:"processingEndTime,omitempty"`
	ReportDocumentID    string `json:"reportDocumentId,omitempty"`
	Report          reports.Report `json:"report"`
	RetrievedAt     time.Time      `json:"retrievedAt"`
}

type reportsGetReportDocumentArgs struct {
	ReportDocumentID string `json:"reportDocumentId"`
}

type reportsGetReportDocumentResult struct {
	ReportDocumentID     string    `json:"reportDocumentId"`
	URL                  string    `json:"url"`
	CompressionAlgorithm string    `json:"compressionAlgorithm,omitempty"`
	RetrievedAt          time.Time `json:"retrievedAt"`
}

func newReportsTools(deps Dependencies) []server.ServerTool {
	spClient := deps.SellingPartner

	getReportsHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args reportsGetReportsArgs) (*mcp.CallToolResult, error) {
		return executeReportsGetReports(ctx, args, spClient)
	})

	createReportHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args reportsCreateReportArgs) (*mcp.CallToolResult, error) {
		return executeReportsCreateReport(ctx, args, spClient)
	})

	getReportHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args reportsGetReportArgs) (*mcp.CallToolResult, error) {
		return executeReportsGetReport(ctx, strings.TrimSpace(args.ReportID), spClient)
	})

	getReportDocumentHandler := mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, args reportsGetReportDocumentArgs) (*mcp.CallToolResult, error) {
		return executeReportsGetReportDocument(ctx, strings.TrimSpace(args.ReportDocumentID), spClient)
	})

	return []server.ServerTool{
		serverToolFromSpec(reportsGetReportsSpec, getReportsHandler),
		serverToolFromSpec(reportsCreateReportSpec, createReportHandler),
		serverToolFromSpec(reportsGetReportSpec, getReportHandler),
		serverToolFromSpec(reportsGetReportDocumentSpec, getReportDocumentHandler),
	}
}

func executeReportsGetReports(ctx context.Context, args reportsGetReportsArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureReportsClient(spClient)
	if failure != nil {
		return failure, nil
	}

	params := &reports.GetReportsParams{}
	
	if nextToken := strings.TrimSpace(args.NextToken); nextToken != "" {
		params.NextToken = &nextToken
	}
	
	if len(args.ReportTypes) > 0 {
		params.ReportTypes = &args.ReportTypes
	}
	
	if len(args.ProcessingStatuses) > 0 {
		params.ProcessingStatuses = &args.ProcessingStatuses
	}
	
	if len(args.MarketplaceIDs) > 0 {
		params.MarketplaceIds = &args.MarketplaceIDs
	}
	
	if args.PageSize != nil && *args.PageSize > 0 && *args.PageSize <= 100 {
		params.PageSize = args.PageSize
	}
	
	if createdSince := strings.TrimSpace(args.CreatedSince); createdSince != "" {
		if parsedTime, err := time.Parse(time.RFC3339, createdSince); err == nil {
			params.CreatedSince = &parsedTime
		}
	}
	
	if createdUntil := strings.TrimSpace(args.CreatedUntil); createdUntil != "" {
		if parsedTime, err := time.Parse(time.RFC3339, createdUntil); err == nil {
			params.CreatedUntil = &parsedTime
		}
	}

	httpResp, err := client.GetReports(ctx, params)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("reports.getReports request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("reports.getReports returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read reports.getReports response", readErr), nil
	}

	decoded, decodeErr := decodeReportsGetReports(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode reports.getReports response", decodeErr), nil
	}

	if err := ensureReportsAPIResponse("getReports", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("reports.getReports response payload is empty"), nil
	}

	result := reportsGetReportsResult{
		Reports:     decoded.reports,
		NextToken:   decoded.nextToken,
		RetrievedAt: time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved %d reports", len(result.Reports))
	if result.NextToken != "" {
		fallback = fmt.Sprintf("%s, more available via nextToken", fallback)
	}

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeReportsCreateReport(ctx context.Context, args reportsCreateReportArgs, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureReportsClient(spClient)
	if failure != nil {
		return failure, nil
	}

	reportType := strings.TrimSpace(args.ReportType)
	if reportType == "" {
		return mcp.NewToolResultError("reportType is required"), nil
	}

	marketplaces := trimStringSlice(args.MarketplaceIDs)
	if len(marketplaces) == 0 {
		return mcp.NewToolResultError("marketplaceIds is required"), nil
	}

	spec := reports.CreateReportSpecification{
		ReportType:     reportType,
		MarketplaceIds: marketplaces,
	}

	if dataStartTime := strings.TrimSpace(args.DataStartTime); dataStartTime != "" {
		if parsedTime, err := time.Parse(time.RFC3339, dataStartTime); err == nil {
			spec.DataStartTime = &parsedTime
		} else {
			return mcp.NewToolResultError("dataStartTime must be in ISO 8601 format"), nil
		}
	}

	if dataEndTime := strings.TrimSpace(args.DataEndTime); dataEndTime != "" {
		if parsedTime, err := time.Parse(time.RFC3339, dataEndTime); err == nil {
			spec.DataEndTime = &parsedTime
		} else {
			return mcp.NewToolResultError("dataEndTime must be in ISO 8601 format"), nil
		}
	}

	if len(args.ReportOptions) > 0 {
		reportOptions := reports.ReportOptions{}
		for k, v := range args.ReportOptions {
			reportOptions.Set(k, v)
		}
		spec.ReportOptions = &reportOptions
	}

	httpResp, err := client.CreateReport(ctx, reports.CreateReportJSONRequestBody(spec))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("reports.createReport request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("reports.createReport returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read reports.createReport response", readErr), nil
	}

	decoded, decodeErr := decodeReportsCreateReport(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode reports.createReport response", decodeErr), nil
	}

	if err := ensureReportsAPIResponse("createReport", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("reports.createReport response payload is empty"), nil
	}

	result := reportsCreateReportResult{
		ReportID:    decoded.reportID,
		ReportType:  reportType,
		RetrievedAt: time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Created %s report with ID %s", reportType, result.ReportID)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeReportsGetReport(ctx context.Context, reportID string, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureReportsClient(spClient)
	if failure != nil {
		return failure, nil
	}

	if reportID == "" {
		return mcp.NewToolResultError("reportId is required"), nil
	}

	httpResp, err := client.GetReport(ctx, reportID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("reports.getReport request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("reports.getReport returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read reports.getReport response", readErr), nil
	}

	decoded, decodeErr := decodeReportsGetReport(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode reports.getReport response", decodeErr), nil
	}

	if err := ensureReportsAPIResponse("getReport", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("reports.getReport response payload is empty"), nil
	}

	report := decoded.report
	result := reportsGetReportResult{
		ReportID:            report.ReportId,
		ReportType:          report.ReportType,
		ProcessingStatus:    report.ProcessingStatus,
		CreatedTime:         report.CreatedTime.Format(time.RFC3339),
		ProcessingStartTime: formatTimePtr(report.ProcessingStartTime),
		ProcessingEndTime:   formatTimePtr(report.ProcessingEndTime),
		ReportDocumentID:    valueOrEmpty(report.ReportDocumentId),
		Report:              report,
		RetrievedAt:         time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Report %s (%s) - Status: %s", result.ReportID, result.ReportType, result.ProcessingStatus)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func executeReportsGetReportDocument(ctx context.Context, reportDocumentID string, spClient spapi.Client) (*mcp.CallToolResult, error) {
	client, failure := ensureReportsClient(spClient)
	if failure != nil {
		return failure, nil
	}

	if reportDocumentID == "" {
		return mcp.NewToolResultError("reportDocumentId is required"), nil
	}

	httpResp, err := client.GetReportDocument(ctx, reportDocumentID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("reports.getReportDocument request failed", err), nil
	}
	if httpResp == nil {
		return mcp.NewToolResultError("reports.getReportDocument returned no response"), nil
	}

	body, readErr := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	if readErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to read reports.getReportDocument response", readErr), nil
	}

	decoded, decodeErr := decodeReportsGetReportDocument(body)
	if decodeErr != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode reports.getReportDocument response", decodeErr), nil
	}

	if err := ensureReportsAPIResponse("getReportDocument", httpResp, body, decoded.apiErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError("reports.getReportDocument response payload is empty"), nil
	}

	result := reportsGetReportDocumentResult{
		ReportDocumentID:     decoded.reportDocumentID,
		URL:                  decoded.url,
		CompressionAlgorithm: decoded.compressionAlgorithm,
		RetrievedAt:          time.Now().UTC(),
	}

	fallback := fmt.Sprintf("Retrieved download URL for report document %s", result.ReportDocumentID)

	return mcp.NewToolResultStructured(result, fallback), nil
}

func ensureReportsClient(spClient spapi.Client) (*reports.Client, *mcp.CallToolResult) {
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

	client, err := buildReportsClient(spClient)
	if err != nil {
		return nil, mcp.NewToolResultErrorFromErr("failed to create reports client", err)
	}

	return client, nil
}

func buildReportsClient(spClient spapi.Client) (*reports.Client, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	return &reports.Client{
		Endpoint:      spClient.Endpoint(),
		Client:        httpClient,
		RequestBefore: buildReportsRequestBefore(spClient),
	}, nil
}

func buildReportsRequestBefore(spClient spapi.Client) reports.RequestBeforeFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Amzn-Requestid", uuid.NewString())
		req.Header.Set("Accept", "application/json")
		if err := spClient.AuthorizeRequest(req); err != nil {
			return fmt.Errorf("authorize request: %w", err)
		}
		return nil
	}
}

func ensureReportsAPIResponse(operation string, resp *http.Response, body []byte, errors *reports.ErrorList) error {
	if resp == nil {
		return fmt.Errorf("%s: no HTTP response returned", operation)
	}

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s: request failed with status %d %s: %s", operation, statusCode, http.StatusText(statusCode), sanitizeBodySnippet(body))
	}

	if errors != nil && len(*errors) > 0 {
		return fmt.Errorf("%s: %s", operation, formatReportsErrors(*errors))
	}

	return nil
}

func formatReportsErrors(list reports.ErrorList) string {
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

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}