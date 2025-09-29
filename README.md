# Selling Partner API MCP Server (Go)

This project packages Amazon's Selling Partner API into a Model Context Protocol (MCP) server written in Go. The server uses [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) for the transport/runtime layer and [`amzapi/selling-partner-api-sdk`](https://github.com/amzapi/selling-partner-api-sdk) for SP-API clients. Most tools still act as placeholders, but the Orders suite now calls the live SP-API so you can retrieve listings of orders, addresses, buyer info, and line items.

---

## Prerequisites

- Go 1.22 or newer
- Amazon Selling Partner API credentials (Login with Amazon) with access to the desired marketplaces
- Optional `.env` file in the repository root so secrets stay out of your shell history

Keep credentials secure; never commit them to source control. Rotate refresh tokens regularly and scope IAM roles to the minimum permissions your workflow requires.

---

## Configuration

The server reads configuration from environment variables (or a local `.env` file). Defaults are chosen for local development, but production deployments should set every value explicitly.

| Variable | Default | Description |
| --- | --- | --- |
| `SP_API_CLIENT_ID` | _required_ | Login with Amazon client identifier |
| `SP_API_CLIENT_SECRET` | _required_ | Login with Amazon client secret |
| `SP_API_REFRESH_TOKEN` | _required_ | Refresh token scoped to your SP-API role |
| `SP_API_ENDPOINT` | `https://sellingpartnerapi-na.amazon.com` | SP-API regional endpoint |
| `MCP_SERVER_NAME` | `Selling Partner MCP Server` | Name shown to MCP clients |
| `MCP_SERVER_VERSION` | `0.1.0` | Semantic-ish version string reported to clients |
| `MCP_SERVER_INSTRUCTIONS` | placeholder text | High-level instructions shared with the assistant |
| `MCP_TRANSPORT` | `stdio` | One of `stdio`, `sse`, `streamablehttp` |
| `PORT` | `8080` | Required when using `sse` or `streamablehttp` transports |

Example `.env` template:

```dotenv
SP_API_CLIENT_ID=amzn1.application-oa2-client....
SP_API_CLIENT_SECRET=super-secret
SP_API_REFRESH_TOKEN=Atzr|...
SP_API_ENDPOINT=https://sellingpartnerapi-na.amazon.com
MCP_TRANSPORT=stdio
```

---

## Quick Start

```bash
# install dependencies & verify the project builds
go build ./...

# run the MCP server over stdio (default)
go run ./cmd/server
```

When `MCP_TRANSPORT` is left as `stdio`, the process communicates via stdin/stdout and is ideal for desktop MCP clients like Cursor or Claude. To expose it as a network service, pick an alternate transport:

```bash
# SSE endpoint on localhost:8080
PORT=8080 MCP_TRANSPORT=sse go run ./cmd/server

# Streamable HTTP endpoint on localhost:9090
PORT=9090 MCP_TRANSPORT=streamablehttp go run ./cmd/server
```

Build a reusable binary if you plan to host it somewhere persistent:

```bash
go build -o bin/sp-api-mcp ./cmd/server
```

---

## Registering the Server with an MCP Client

The server can be launched by any MCP-compatible client by delegating to the `go run` command (or a compiled binary) and exporting the required environment variables.

### Cursor

Add an entry to `~/Library/Application Support/Cursor/mcp.json` (macOS) or the equivalent path on your platform:

```json
{
  "mcpServers": {
    "sp-api": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "cwd": "/absolute/path/to/sp-api-mcp-go",
      "env": {
        "SP_API_CLIENT_ID": "amzn1.application-oa2-client...",
        "SP_API_CLIENT_SECRET": "...",
        "SP_API_REFRESH_TOKEN": "Atzr|...",
        "SP_API_ENDPOINT": "https://sellingpartnerapi-na.amazon.com"
      }
    }
  }
}
```

Restart Cursor and enable the `sp-api` server in the MCP panel.

### Claude Desktop

Create or update `~/Library/Application Support/Claude/mcp/config.json` (macOS) with the server configuration:

```json
{
  "servers": {
    "sp-api": {
      "command": "/absolute/path/to/bin/sp-api-mcp",
      "transport": "stdio",
      "workingDirectory": "/absolute/path/to/sp-api-mcp-go",
      "env": {
        "SP_API_CLIENT_ID": "amzn1.application-oa2-client...",
        "SP_API_CLIENT_SECRET": "...",
        "SP_API_REFRESH_TOKEN": "Atzr|..."
      }
    }
  }
}
```

Restart Claude Desktop to pick up the new configuration. If you prefer to run over SSE or Streamable HTTP, point `command` to a process supervisor (e.g. `npm`, `systemd`, `forever`) that hosts the binary and then reference the network endpoint inside Claude's UI when adding the server.

### Other MCP Clients

- **Continue**: add a stanza under `~/.continue/mcpServers.json` using the same structure as above.
- **Runes/CLI tools**: invoke the binary directly and pipe stdin/stdout, or configure the client to connect to the SSE/HTTP endpoint depending on your transport choice.

Regardless of the client, make sure the process can read your `.env` or has the credentials in its environment before starting.

---

## Tools & Resources Exposed

The current build ships placeholder tools to help you scaffold real SP-API workflows:

- `auth.beginAuthorization` – Guides implementing Login with Amazon authorization.
- `catalog.lookupItem` – Placeholder for catalog metadata lookups.
- `inventory.getSummary` – Placeholder for inventory summaries across marketplaces.
- `orders.listOrders` – Lists orders for a marketplace window and returns Amazon next tokens for pagination.
- `orders.getOrder` – Fetches order metadata and line items via the Orders API when SP-API credentials are configured.
- `orders.getOrderAddress` – Returns the shipping address for an order.
- `orders.getOrderBuyerInfo` – Returns buyer contact details where scopes allow it.
- `orders.getOrderItems` – Lists order line items page by page.
- `orders.getOrderItemsBuyerInfo` – Lists buyer-specific data such as gift messages per line item.
- `reports.createReport` – Placeholder for asynchronous report generation.
- `feeds.submitFeed` – Placeholder for feed submission workflows.
- `finance.listFinancialEvents` – Placeholder for reconciling financial events.
- `notifications.subscribe` – Placeholder for managing notification subscriptions.
- `pricing.getPricing` – Placeholder for competitive pricing retrieval.
- `listings.updateListing` – Placeholder for listing patch operations.
- `fba.createInboundShipmentPlan` – Placeholder for FBA inbound shipment planning.

Documentation resources are available under URIs like `amazon-sp-api://overview`, providing structured notes you can expand with live references as integrations are implemented.

---

## Development Workflow

- Format code with `gofmt` (tabs, trailing newline).
- Run `go test ./...` before pushing changes.
- `go run ./cmd/server` exercises the server end-to-end against your environment.
- Generated binaries (`bin/sp-api-mcp` or similar) should remain untracked; rebuild locally when needed.

Feel free to replace placeholder tool implementations with real SP-API calls by extending the types under `internal/tools` and wiring additional dependencies through `internal/app`.

---

## Troubleshooting

- **Configuration errors**: the server validates environment variables on startup and will explain missing or partial credentials in the log output.
- **Authentication failures**: confirm your refresh token is still valid and that the IAM role has the scopes required by the requested SP-API endpoints.
- **Rate limits**: Amazon enforces per-endpoint throttling; cache responses and implement retries with backoff when you replace the placeholder logic with live calls.

For deeper debugging, enable MCP client logging (Cursor/Claude provide toggles) and inspect the JSON RPC traffic to trace tool invocations.

---

## TODO - SP-API Implementation Checklist

### Currently Implemented (12 functions)
- [x] **Authorization**: GetAuthorizationCode  
- [x] **Orders**: GetOrders, GetOrder, GetOrderAddress, GetOrderBuyerInfo, GetOrderItems, GetOrderItemsBuyerInfo
- [x] **Sales**: GetOrderMetrics
- [x] **Reports**: GetReports, CreateReport, GetReport, GetReportDocument

### Phase 1: Core READ Operations (High Priority)

#### Finances API
- [ ] **ListFinancialEventGroups** - List financial event groups [#4](https://github.com/berrydev-ai/sp-api-mcp-go/issues/4)
- [ ] **ListFinancialEventsByGroupId** - List financial events by group [#5](https://github.com/berrydev-ai/sp-api-mcp-go/issues/5)
- [ ] **ListFinancialEventsByOrderId** - List financial events by order [#6](https://github.com/berrydev-ai/sp-api-mcp-go/issues/6)
- [ ] **ListFinancialEvents** - List all financial events [#7](https://github.com/berrydev-ai/sp-api-mcp-go/issues/7)

#### FBA Inventory API
- [ ] **GetInventorySummaries** - Get inventory summaries [#8](https://github.com/berrydev-ai/sp-api-mcp-go/issues/8)

#### Product Pricing API
- [ ] **GetPricing** - Get pricing for products [#9](https://github.com/berrydev-ai/sp-api-mcp-go/issues/9)
- [ ] **GetCompetitivePricing** - Get competitive pricing [#10](https://github.com/berrydev-ai/sp-api-mcp-go/issues/10)
- [ ] **GetListingOffers** - Get listing offers
- [ ] **GetItemOffers** - Get item offers
- [ ] **GetItemOffersBatch** - Get item offers in batch

#### Product Fees API
- [ ] **GetMyFeesEstimateForSKU** - Get fees estimate for SKU [#11](https://github.com/berrydev-ai/sp-api-mcp-go/issues/11)
- [ ] **GetMyFeesEstimateForASIN** - Get fees estimate for ASIN [#12](https://github.com/berrydev-ai/sp-api-mcp-go/issues/12)

### Phase 2: Enhanced READ Functionality (Medium Priority)

#### Catalog Items API
- [ ] **SearchCatalogItems** - Search for catalog items [#13](https://github.com/berrydev-ai/sp-api-mcp-go/issues/13)
- [ ] **GetCatalogItem** - Get details for a specific catalog item [#14](https://github.com/berrydev-ai/sp-api-mcp-go/issues/14)

#### Listings Items API
- [ ] **GetListingsItem** - Get listing details [#15](https://github.com/berrydev-ai/sp-api-mcp-go/issues/15)

#### Sellers API
- [ ] **GetMarketplaceParticipations** - Get marketplace participations [#16](https://github.com/berrydev-ai/sp-api-mcp-go/issues/16)

#### FBA Inbound API (READ-only)
- [ ] **GetInboundGuidance** - Get inbound guidance for items [#17](https://github.com/berrydev-ai/sp-api-mcp-go/issues/17)
- [ ] **GetPreorderInfo** - Get preorder information [#18](https://github.com/berrydev-ai/sp-api-mcp-go/issues/18)
- [ ] **GetPrepInstructions** - Get prep instructions [#19](https://github.com/berrydev-ai/sp-api-mcp-go/issues/19)
- [ ] **GetTransportDetails** - Get transport details [#20](https://github.com/berrydev-ai/sp-api-mcp-go/issues/20)
- [ ] **GetLabels** - Get inbound shipment labels [#21](https://github.com/berrydev-ai/sp-api-mcp-go/issues/21)
- [ ] **GetBillOfLading** - Get bill of lading [#22](https://github.com/berrydev-ai/sp-api-mcp-go/issues/22)
- [ ] **GetShipments** - List inbound shipments [#23](https://github.com/berrydev-ai/sp-api-mcp-go/issues/23)
- [ ] **GetShipmentItemsByShipmentId** - Get shipment items [#24](https://github.com/berrydev-ai/sp-api-mcp-go/issues/24)
- [ ] **GetShipmentItems** - Get shipment items across shipments [#25](https://github.com/berrydev-ai/sp-api-mcp-go/issues/25)

#### FBA Outbound API (READ-only)
- [ ] **GetFulfillmentPreview** - Get fulfillment preview [#26](https://github.com/berrydev-ai/sp-api-mcp-go/issues/26)
- [ ] **GetFulfillmentOrder** - Get fulfillment order details [#27](https://github.com/berrydev-ai/sp-api-mcp-go/issues/27)
- [ ] **ListAllFulfillmentOrders** - List all fulfillment orders [#28](https://github.com/berrydev-ai/sp-api-mcp-go/issues/28)
- [ ] **GetPackageTrackingDetails** - Get package tracking details [#29](https://github.com/berrydev-ai/sp-api-mcp-go/issues/29)
- [ ] **ListReturnReasonCodes** - List return reason codes [#30](https://github.com/berrydev-ai/sp-api-mcp-go/issues/30)
- [ ] **GetFulfillmentReturn** - Get fulfillment return [#31](https://github.com/berrydev-ai/sp-api-mcp-go/issues/31)
- [ ] **GetFeatures** - Get available features [#32](https://github.com/berrydev-ai/sp-api-mcp-go/issues/32)
- [ ] **GetFeatureInventory** - Get feature inventory [#33](https://github.com/berrydev-ai/sp-api-mcp-go/issues/33)
- [ ] **GetFeatureSKU** - Get feature SKU [#34](https://github.com/berrydev-ai/sp-api-mcp-go/issues/34)

### Phase 3: Advanced READ Features (Lower Priority)

#### Feeds API (READ-only)
- [ ] **GetFeeds** - Get feed processing reports [#35](https://github.com/berrydev-ai/sp-api-mcp-go/issues/35)
- [ ] **GetFeed** - Get feed details [#36](https://github.com/berrydev-ai/sp-api-mcp-go/issues/36)
- [ ] **GetFeedDocument** - Get feed document [#37](https://github.com/berrydev-ai/sp-api-mcp-go/issues/37)

#### Reports API (Additional READ methods)
- [ ] **GetReportSchedules** - Get report schedules [#38](https://github.com/berrydev-ai/sp-api-mcp-go/issues/38)
- [ ] **GetReportSchedule** - Get report schedule details [#39](https://github.com/berrydev-ai/sp-api-mcp-go/issues/39)

#### Messaging API (READ-only)
- [ ] **GetMessagingActionsForOrder** - Get messaging actions for order [#40](https://github.com/berrydev-ai/sp-api-mcp-go/issues/40)
- [ ] **GetAttributes** - Get messaging attributes [#41](https://github.com/berrydev-ai/sp-api-mcp-go/issues/41)

#### Notifications API (READ-only)
- [ ] **GetSubscription** - Get subscription details [#42](https://github.com/berrydev-ai/sp-api-mcp-go/issues/42)
- [ ] **GetSubscriptionById** - Get subscription by ID [#43](https://github.com/berrydev-ai/sp-api-mcp-go/issues/43)
- [ ] **GetDestinations** - Get notification destinations [#44](https://github.com/berrydev-ai/sp-api-mcp-go/issues/44)
- [ ] **GetDestination** - Get destination details [#45](https://github.com/berrydev-ai/sp-api-mcp-go/issues/45)

#### Merchant Fulfillment API (READ-only)
- [ ] **GetEligibleShipmentServices** - Get eligible shipment services [#46](https://github.com/berrydev-ai/sp-api-mcp-go/issues/46)
- [ ] **GetShipment** - Get shipment details [#47](https://github.com/berrydev-ai/sp-api-mcp-go/issues/47)
- [ ] **GetAdditionalSellerInputs** - Get additional seller inputs [#48](https://github.com/berrydev-ai/sp-api-mcp-go/issues/48)

#### Service API (READ-only)
- [ ] **GetServiceJobs** - Get service jobs [#49](https://github.com/berrydev-ai/sp-api-mcp-go/issues/49)
- [ ] **GetServiceJobByServiceJobId** - Get service job details [#50](https://github.com/berrydev-ai/sp-api-mcp-go/issues/50)

#### Shipping API (READ-only)
- [ ] **GetShipment** - Get shipment details [#51](https://github.com/berrydev-ai/sp-api-mcp-go/issues/51)
- [ ] **GetRates** - Get shipping rates [#52](https://github.com/berrydev-ai/sp-api-mcp-go/issues/52)
- [ ] **GetAccount** - Get account information [#53](https://github.com/berrydev-ai/sp-api-mcp-go/issues/53)
- [ ] **GetTrackingInformation** - Get tracking information [#54](https://github.com/berrydev-ai/sp-api-mcp-go/issues/54)

#### Small and Light API (READ-only)
- [ ] **GetSmallAndLightEnrollmentBySellerSKU** - Get S&L enrollment by SKU [#55](https://github.com/berrydev-ai/sp-api-mcp-go/issues/55)
- [ ] **GetSmallAndLightEligibilityBySellerSKU** - Check S&L eligibility [#56](https://github.com/berrydev-ai/sp-api-mcp-go/issues/56)
- [ ] **GetSmallAndLightFeePreview** - Get S&L fee preview [#57](https://github.com/berrydev-ai/sp-api-mcp-go/issues/57)

#### Solicitations API (READ-only)
- [ ] **GetSolicitationActionsForOrder** - Get solicitation actions [#58](https://github.com/berrydev-ai/sp-api-mcp-go/issues/58)