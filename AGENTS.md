# Project rules (Go)

## Commands (must run)

- gofmt ./...
- go test ./...

## Optional lint

- golangci-lint run (if available)

## Definition of Done

- All AC in PRD are satisfied
- go test ./... passes
- gofmt applied
- STATUS.md updated with PASS flags

## Safety

- Prefer sandbox workspace-write
- No unrelated refactors
