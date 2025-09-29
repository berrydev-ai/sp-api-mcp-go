package tools

import "github.com/mark3labs/mcp-go/server"

// BuildAll assembles every tool the server should expose. Future tools can be registered by extending the relevant specs.
func BuildAll(deps Dependencies) []server.ServerTool {
	orders := newOrdersTools(deps)
	sales := newSalesTools(deps)
	reports := newReportsTools(deps)
	fbaInventory := newFBAInventoryTools(deps)
	productPricing := newProductPricingTools(deps)
	all := make([]server.ServerTool, 0, len(orders)+len(sales)+len(reports)+len(fbaInventory)+len(productPricing)+len(placeholderSpecs))

	all = append(all, orders...)
	all = append(all, sales...)
	all = append(all, reports...)
	all = append(all, fbaInventory...)
	all = append(all, productPricing...)

	for _, spec := range placeholderSpecs {
		all = append(all, newPlaceholderTool(spec, deps))
	}

	return all
}
