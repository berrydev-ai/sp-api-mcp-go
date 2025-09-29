package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	sales "github.com/amzapi/selling-partner-api-sdk/sales"
)

type salesOrderMetricsDecoded struct {
	metrics        []sales.OrderMetricsInterval
	apiErrors      *sales.ErrorList
	payloadPresent bool
}

func decodeSalesOrderMetrics(body []byte) (salesOrderMetricsDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return salesOrderMetricsDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto salesOrderMetricsResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return salesOrderMetricsDecoded{}, err
	}

	decoded := salesOrderMetricsDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.metrics = convertSalesOrderMetrics(*dto.Payload)
	}

	return decoded, nil
}

type salesOrderMetricsResponseDTO struct {
	Errors  *sales.ErrorList                `json:"errors,omitempty"`
	Payload *[]salesOrderMetricsIntervalDTO `json:"payload"`
}

type salesOrderMetricsIntervalDTO struct {
	AverageUnitPrice salesMoneyDTO `json:"averageUnitPrice"`
	Interval         string        `json:"interval"`
	OrderCount       int           `json:"orderCount"`
	OrderItemCount   int           `json:"orderItemCount"`
	TotalSales       salesMoneyDTO `json:"totalSales"`
	UnitCount        int           `json:"unitCount"`
}

type salesMoneyDTO struct {
	Amount       decimalString `json:"amount"`
	CurrencyCode string        `json:"currencyCode"`
}

func (dto salesMoneyDTO) toSalesMoney() sales.Money {
	return sales.Money{
		Amount:       sales.Decimal(dto.Amount.String()),
		CurrencyCode: dto.CurrencyCode,
	}
}

func convertSalesOrderMetrics(items []salesOrderMetricsIntervalDTO) []sales.OrderMetricsInterval {
	if len(items) == 0 {
		return nil
	}

	metrics := make([]sales.OrderMetricsInterval, 0, len(items))
	for _, item := range items {
		metrics = append(metrics, sales.OrderMetricsInterval{
			AverageUnitPrice: item.AverageUnitPrice.toSalesMoney(),
			Interval:         item.Interval,
			OrderCount:       item.OrderCount,
			OrderItemCount:   item.OrderItemCount,
			TotalSales:       item.TotalSales.toSalesMoney(),
			UnitCount:        item.UnitCount,
		})
	}

	return metrics
}

type decimalString string

func (d *decimalString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*d = ""
		return nil
	}

	if trimmed[0] == '"' {
		var raw string
		if err := json.Unmarshal(trimmed, &raw); err != nil {
			return err
		}
		*d = decimalString(strings.TrimSpace(raw))
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(trimmed, &number); err != nil {
		return err
	}

	*d = decimalString(number.String())
	return nil
}

func (d decimalString) String() string {
	return string(d)
}
