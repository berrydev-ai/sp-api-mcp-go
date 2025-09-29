# Repository Guidelines

## Project Structure & Module Organization
The Go module lives at `github.com/berrydev-ai/sp-api-mcp-go`. Active source code sits in `example.go`, which bootstraps the Selling Partner API client and samples the Sellers "GetMarketplaceParticipations" endpoint. Generated artifacts (such as the `sp-api-mcp-go` binary) should remain untracked; rebuild locally when needed. Configuration and credentials load from an optional `.env` file in the root, keeping SP-API secrets out of source control.

## Build, Test, and Development Commands
Run `go build ./...` to compile the module and confirm dependencies resolve. Use `go run .` for a quick local execution against your environment. Execute `go test ./...` before submitting changes; even if suites are sparse today, the command guards against regressions as tests grow.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (tabs for indentation, trailing newlines required). Prefer descriptive, PascalCase export names that mirror Selling Partner API terminology, and use lowerCamelCase for locals. Group imports by standard library, third-party, then internal packages. Keep logging via the `log` package consistent with existing request/response dumps.

## Testing Guidelines
New features should ship with `_test.go` files under the same package, using Go's built-in `testing` framework. Favor table-driven tests for API helpers and mock external calls where possible. Target full coverage for new logic touching Amazon endpoints, and validate with `go test -run TestName` while iterating.

## Commit & Pull Request Guidelines
Write commit subjects in the imperative mood (e.g., `Add seller participation client`) and limit to ~72 characters. Reference linked issues in the body when applicable. Pull requests should summarize behavior changes, note testing performed (`go test ./...`), and include any screenshots or trace snippets that help reviewers verify API interactions.

## Environment & Security Notes
Populate `.env` with `SP_API_CLIENT_ID`, `SP_API_CLIENT_SECRET`, and `SP_API_REFRESH_TOKEN`; never commit these values. Rotate credentials regularly and prefer scoped IAM policies. When sharing logs, scrub request IDs or dumps that may leak sensitive payloads.
