package tools

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/amzapi/selling-partner-api-sdk/fbaInventory"
)

type fbaInventoryGetInventorySummariesDecoded struct {
	granularityType    string
	granularityID      string
	inventorySummaries []fbaInventory.InventorySummary
	nextToken          string
	apiErrors          *fbaInventory.ErrorList
	payloadPresent     bool
}

func decodeFBAInventoryGetInventorySummaries(body []byte) (fbaInventoryGetInventorySummariesDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return fbaInventoryGetInventorySummariesDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto fbaInventoryGetInventorySummariesResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return fbaInventoryGetInventorySummariesDecoded{}, err
	}

	decoded := fbaInventoryGetInventorySummariesDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		decoded.inventorySummaries = dto.Payload.InventorySummaries
		
		if dto.Payload.Granularity.GranularityType != nil {
			decoded.granularityType = *dto.Payload.Granularity.GranularityType
		}
		if dto.Payload.Granularity.GranularityId != nil {
			decoded.granularityID = *dto.Payload.Granularity.GranularityId
		}
		
		if dto.Pagination != nil && dto.Pagination.NextToken != nil {
			decoded.nextToken = *dto.Pagination.NextToken
		}
	}

	return decoded, nil
}

type fbaInventoryGetInventorySummariesResponseDTO struct {
	Errors     *fbaInventory.ErrorList                    `json:"errors,omitempty"`
	Pagination *fbaInventory.Pagination                   `json:"pagination,omitempty"`
	Payload    *fbaInventory.GetInventorySummariesResult  `json:"payload,omitempty"`
}