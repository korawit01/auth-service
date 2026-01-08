PRD_READY=true
ADR_READY=true
TESTPLAN_READY=true
IMPLEMENTED=true
AC_PASS=true
TEST_PASS=true
LINT_PASS=false
LINT_SKIPPED=true
COMMANDS:
- gofmt ./... (failed: open ...\\-p: The system cannot find the path specified.)
- gofmt -w (rg --files -g '*.go') (pass)
- go test ./... (failed: access denied to Go build cache under %LOCALAPPDATA%\\go-build)
- GOCACHE=D:\\Project\\go-micro-taskboard\\.gocache go test ./... (pass)
- golangci-lint run (skipped: not installed)
SUMMARY:
- replaced auth/task service config loading with Viper defaults + config.yaml + TASKBOARD_ env support and validation
- added config tests for defaults, file overrides, env overrides, and validation failures
- updated docker-compose env names and added config.example.yaml + run instructions
ISSUES:
- go test ./... initially failed due to access denied on Go build cache; reran with GOCACHE set to workspace to pass
