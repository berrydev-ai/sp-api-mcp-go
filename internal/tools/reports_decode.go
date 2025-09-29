package tools

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/amzapi/selling-partner-api-sdk/reports"
)

type reportsGetReportsDecoded struct {
	reports        []reports.Report
	nextToken      string
	apiErrors      *reports.ErrorList
	payloadPresent bool
}

type reportsCreateReportDecoded struct {
	reportID       string
	apiErrors      *reports.ErrorList
	payloadPresent bool
}

type reportsGetReportDecoded struct {
	report         reports.Report
	apiErrors      *reports.ErrorList
	payloadPresent bool
}

type reportsGetReportDocumentDecoded struct {
	reportDocumentID     string
	url                  string
	compressionAlgorithm string
	apiErrors            *reports.ErrorList
	payloadPresent       bool
}

func decodeReportsGetReports(body []byte) (reportsGetReportsDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return reportsGetReportsDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto reportsGetReportsResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return reportsGetReportsDecoded{}, err
	}

	decoded := reportsGetReportsDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.reports = dto.Payload.Reports
		decoded.nextToken = valueOrEmpty(dto.Payload.NextToken)
	}

	return decoded, nil
}

func decodeReportsCreateReport(body []byte) (reportsCreateReportDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return reportsCreateReportDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto reportsCreateReportResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return reportsCreateReportDecoded{}, err
	}

	decoded := reportsCreateReportDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.reportID = dto.Payload.ReportId
	}

	return decoded, nil
}

func decodeReportsGetReport(body []byte) (reportsGetReportDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return reportsGetReportDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto reportsGetReportResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return reportsGetReportDecoded{}, err
	}

	decoded := reportsGetReportDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.report = *dto.Payload
	}

	return decoded, nil
}

func decodeReportsGetReportDocument(body []byte) (reportsGetReportDocumentDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return reportsGetReportDocumentDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto reportsGetReportDocumentResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return reportsGetReportDocumentDecoded{}, err
	}

	decoded := reportsGetReportDocumentDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.reportDocumentID = valueOrEmpty(dto.ReportDocumentId)
		decoded.url = valueOrEmpty(dto.Url)
		decoded.compressionAlgorithm = valueOrEmpty(dto.CompressionAlgorithm)
	}

	return decoded, nil
}

type reportsGetReportsResponseDTO struct {
	Errors  *reports.ErrorList           `json:"errors,omitempty"`
	Payload *reportsGetReportsPayloadDTO `json:"payload,omitempty"`
}

type reportsGetReportsPayloadDTO struct {
	Reports   []reports.Report `json:"reports"`
	NextToken *string          `json:"nextToken,omitempty"`
}

type reportsCreateReportResponseDTO struct {
	Errors  *reports.ErrorList                 `json:"errors,omitempty"`
	Payload *reportsCreateReportResultDTO      `json:"payload,omitempty"`
}

type reportsCreateReportResultDTO struct {
	ReportId string `json:"reportId"`
}

type reportsGetReportResponseDTO struct {
	Errors  *reports.ErrorList `json:"errors,omitempty"`
	Payload *reports.Report    `json:"payload,omitempty"`
}

type reportsGetReportDocumentResponseDTO struct {
	Errors               *reports.ErrorList `json:"errors,omitempty"`
	Payload              *reports.ReportDocument `json:"payload,omitempty"`
	ReportDocumentId     *string            `json:"reportDocumentId,omitempty"`
	Url                  *string            `json:"url,omitempty"`
	CompressionAlgorithm *string            `json:"compressionAlgorithm,omitempty"`
}