Goal: Integrate Viper for config + standardize Go codebase
Constraints: minimal refactor, must run gofmt + go test ./...
Config requirements:
  - defaults < config.yaml (optional) < env overrides
  - env prefix TASKBOARD_
  - fail-fast validation for required values (ex: db.url)
Out of scope: unrelated refactors
