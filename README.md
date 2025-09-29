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

Happy building!
