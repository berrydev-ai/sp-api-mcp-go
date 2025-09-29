package tools

import "github.com/mark3labs/mcp-go/server"

// BuildAll assembles every tool the server should expose. Future tools can be registered by extending the relevant specs.
func BuildAll(deps Dependencies) []server.ServerTool {
	orders := newOrdersTools(deps)
	sales := newSalesTools(deps)
	fbaInventory := newFBAInventoryTools(deps)
	all := make([]server.ServerTool, 0, len(orders)+len(sales)+len(fbaInventory)+len(placeholderSpecs))

	all = append(all, orders...)
	all = append(all, sales...)
	all = append(all, fbaInventory...)

	for _, spec := range placeholderSpecs {
		all = append(all, newPlaceholderTool(spec, deps))
	}

	return all
}
