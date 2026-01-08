# Test Plan: Viper Standard Config

## Test cases mapped to acceptance criteria

| AC | Test case | Steps | Expected result |
| --- | --- | --- | --- |
| AC1 | Defaults -> optional file -> env precedence | 1) Ensure no `config.yaml` and no `TASKBOARD_` env vars. 2) Start app. 3) Add `config.yaml` with non-default values. 4) Start app. 5) Set `TASKBOARD_` env vars overlapping file values. 6) Start app. | 1) App uses defaults. 2) File overrides defaults. 3) Env overrides file. |
| AC2 | Missing `config.yaml` still starts | 1) Delete or rename `config.yaml`. 2) Start app with defaults present. | App starts successfully using defaults unless required values are missing. |
| AC3 | Env prefix and key mapping | 1) Set `TASKBOARD_DB_URL`. 2) Set unrelated env `DB_URL` and `TASKBOARD__DB__URL`. 3) Start app. | Only `TASKBOARD_DB_URL` is read and maps to `db.url`. Unprefixed env is ignored. |
| AC4 | Env wins over file when both set | 1) Set `db.url` in `config.yaml`. 2) Set `TASKBOARD_DB_URL` to a different value. 3) Start app. | Effective config uses env value. |
| AC5 | Required config validated at startup | 1) Remove required `db.url` from defaults, file, and env (or set empty). 2) Start app. | App validates at startup and fails before serving requests. |
| AC6 | Clear error and non-zero exit | 1) Trigger validation failure (missing or invalid key). 2) Observe exit code and logs. | Non-zero exit code and clear error message indicating missing/invalid key. |
| AC7 | Run instructions documented | 1) Open docs for feature. | Docs include how to run with defaults, optional `config.yaml`, and `TASKBOARD_` env vars. |

## Negative cases

- Invalid `db.url` format (empty string, whitespace, or malformed URL) -> startup fails with actionable error.
- `config.yaml` contains unknown keys -> app ignores or logs warning without crashing.
- `config.yaml` present but unreadable (permission denied) -> startup fails with clear error.
- `config.yaml` has invalid YAML -> startup fails with clear error.
- Env values with wrong type (e.g., non-numeric port) -> startup fails with validation error.
- Conflicting env variables with different casing -> only correct `TASKBOARD_` keys are applied.

## Automated test suggestions (Go)

- `internal/config` unit tests using table-driven cases for precedence (defaults < file < env).
- Viper env mapping tests using `t.Setenv` and key replacer for `TASKBOARD_DB_URL` -> `db.url`.
- Validation tests for required fields and type constraints, assert error contents.
- Integration test that invokes `config.Load()` and verifies failure before server start when required fields missing.

## Standard checks

- `gofmt ./...`
- `go test ./...`
- `golangci-lint run` (optional, if available)
