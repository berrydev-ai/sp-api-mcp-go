package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amzapi/selling-partner-api-sdk/productPricing"
	"github.com/mark3labs/mcp-go/mcp"
)

type productPricingDecoded struct {
	pricing        []interface{}
	apiErrors      *productPricing.ErrorList
	payloadPresent bool
}

type productPricingGetPricingResult struct {
	Operation     string        `json:"operation"`
	MarketplaceID string        `json:"marketplaceId"`
	ItemType      string        `json:"itemType"`
	ItemCount     int           `json:"itemCount"`
	PricePoints   []interface{} `json:"pricePoints"`
	RetrievedAt   time.Time     `json:"retrievedAt"`
}

type productPricingGetCompetitivePricingResult struct {
	Operation        string        `json:"operation"`
	MarketplaceID    string        `json:"marketplaceId"`
	ItemType         string        `json:"itemType"`
	ItemCount        int           `json:"itemCount"`
	CompetitiveItems []interface{} `json:"competitiveItems"`
	RetrievedAt      time.Time     `json:"retrievedAt"`
}

// DTO types for unmarshaling response data
type productPricingResponseDTO struct {
	Payload *json.RawMessage             `json:"payload,omitempty"`
	Errors  *productPricing.ErrorList    `json:"errors,omitempty"`
}

func decodeProductPricingResponse(resp *http.Response, operation string) (*mcp.CallToolResult, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to read response body", err), nil
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return handleProductPricingAPIError(resp.StatusCode, body, operation)
	}

	decoded, err := decodeProductPricingBody(body)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to decode response", err), nil
	}

	if decoded.apiErrors != nil && len(*decoded.apiErrors) > 0 {
		errorDetails := []string{}
		for _, apiErr := range *decoded.apiErrors {
			if strings.TrimSpace(apiErr.Message) != "" {
				errorDetails = append(errorDetails, apiErr.Message)
			}
		}
		return mcp.NewToolResultError(fmt.Sprintf("API returned errors: %s", strings.Join(errorDetails, "; "))), nil
	}

	if !decoded.payloadPresent {
		return mcp.NewToolResultError(fmt.Sprintf("%s response payload is empty", operation)), nil
	}

	switch operation {
	case "GetPricing":
		result := buildGetPricingResult(decoded)
		fallback := buildGetPricingFallback(result)
		return mcp.NewToolResultStructured(result, fallback), nil
	case "GetCompetitivePricing":
		result := buildGetCompetitivePricingResult(decoded)
		fallback := buildGetCompetitivePricingFallback(result)
		return mcp.NewToolResultStructured(result, fallback), nil
	default:
		return mcp.NewToolResultError(fmt.Sprintf("unknown operation: %s", operation)), nil
	}
}

func decodeProductPricingBody(body []byte) (productPricingDecoded, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return productPricingDecoded{}, fmt.Errorf("response body is empty")
	}

	var dto productPricingResponseDTO
	if err := json.Unmarshal(trimmed, &dto); err != nil {
		return productPricingDecoded{}, err
	}

	decoded := productPricingDecoded{
		apiErrors: dto.Errors,
	}

	if dto.Payload != nil {
		decoded.payloadPresent = true
		var payload map[string]interface{}
		if err := json.Unmarshal(*dto.Payload, &payload); err == nil {
			// Extract pricing items from various possible array fields
			if items, ok := payload["offers"].([]interface{}); ok {
				decoded.pricing = items
			} else if items, ok := payload["pricing"].([]interface{}); ok {
				decoded.pricing = items
			} else if items, ok := payload["competitivePricing"].([]interface{}); ok {
				decoded.pricing = items
			} else {
				// If no recognized array, store the raw payload
				decoded.pricing = []interface{}{payload}
			}
		}
	}

	return decoded, nil
}

func buildGetPricingResult(decoded productPricingDecoded) productPricingGetPricingResult {
	return productPricingGetPricingResult{
		Operation:   "GetPricing",
		ItemCount:   len(decoded.pricing),
		PricePoints: decoded.pricing,
		RetrievedAt: time.Now().UTC(),
	}
}

func buildGetCompetitivePricingResult(decoded productPricingDecoded) productPricingGetCompetitivePricingResult {
	return productPricingGetCompetitivePricingResult{
		Operation:        "GetCompetitivePricing",
		ItemCount:        len(decoded.pricing),
		CompetitiveItems: decoded.pricing,
		RetrievedAt:      time.Now().UTC(),
	}
}

func buildGetPricingFallback(result productPricingGetPricingResult) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Product pricing retrieved for %d items", result.ItemCount))
	
	if result.ItemCount > 0 {
		summary.WriteString(fmt.Sprintf(" at %s", result.RetrievedAt.Format("15:04:05 UTC")))
		
		// Try to extract some key pricing info for the summary
		for i, item := range result.PricePoints {
			if i >= 3 { // Limit to first 3 items in summary
				break
			}
			if itemMap, ok := item.(map[string]interface{}); ok {
				if i > 0 {
					summary.WriteString("; ")
				}
				summary.WriteString(" - ")
				
				// Try to get identifier
				if asin, ok := itemMap["ASIN"].(string); ok {
					summary.WriteString(fmt.Sprintf("ASIN %s", asin))
				} else if sku, ok := itemMap["SellerSKU"].(string); ok {
					summary.WriteString(fmt.Sprintf("SKU %s", sku))
				} else {
					summary.WriteString("Item")
				}
				
				// Try to get price info
				if buyingPrice, ok := itemMap["BuyingPrice"].(map[string]interface{}); ok {
					if listPrice, ok := buyingPrice["ListingPrice"].(map[string]interface{}); ok {
						if amount, ok := listPrice["Amount"].(float64); ok {
							if currency, ok := listPrice["CurrencyCode"].(string); ok {
								summary.WriteString(fmt.Sprintf(": %s %.2f", currency, amount))
							}
						}
					}
				}
			}
		}
	}
	
	return summary.String()
}

func buildGetCompetitivePricingFallback(result productPricingGetCompetitivePricingResult) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Competitive pricing retrieved for %d items", result.ItemCount))
	
	if result.ItemCount > 0 {
		summary.WriteString(fmt.Sprintf(" at %s", result.RetrievedAt.Format("15:04:05 UTC")))
		
		// Try to extract some key competitive pricing info for the summary
		for i, item := range result.CompetitiveItems {
			if i >= 3 { // Limit to first 3 items in summary
				break
			}
			if itemMap, ok := item.(map[string]interface{}); ok {
				if i > 0 {
					summary.WriteString("; ")
				}
				summary.WriteString(" - ")
				
				// Try to get identifier
				if identifier, ok := itemMap["identifier"].(map[string]interface{}); ok {
					if asin, ok := identifier["ASIN"].(string); ok {
						summary.WriteString(fmt.Sprintf("ASIN %s", asin))
					} else if sku, ok := identifier["SellerSKU"].(string); ok {
						summary.WriteString(fmt.Sprintf("SKU %s", sku))
					}
				}
				
				// Try to get competitive price count
				if competitivePricing, ok := itemMap["competitivePricing"].(map[string]interface{}); ok {
					if prices, ok := competitivePricing["CompetitivePrices"].([]interface{}); ok {
						summary.WriteString(fmt.Sprintf(": %d competitive prices", len(prices)))
					}
				}
			}
		}
	}
	
	return summary.String()
}

func handleProductPricingAPIError(statusCode int, body []byte, operation string) (*mcp.CallToolResult, error) {
	message := fmt.Sprintf("%s failed with status %d", operation, statusCode)
	if len(body) > 0 {
		message += fmt.Sprintf(": %s", string(body))
	}
	return mcp.NewToolResultError(message), nil
}