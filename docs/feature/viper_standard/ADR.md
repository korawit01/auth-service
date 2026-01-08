# ADR: Viper Standard Config

## Status
Proposed

## Context
We need a consistent configuration strategy that loads defaults, then an optional
config file, then environment overrides. The system must fail fast with clear
errors when required values are missing or invalid.

## Decision
### Proposed structure: internal/config
Create `internal/config` to own configuration concerns (types, loading, and
validation). This keeps config logic centralized and prevents scattering Viper
setup across services.

### Typed Config struct + Load()
Define a typed `Config` struct and a `Load()` function responsible for:
- establishing defaults
- optionally reading `config.yaml`
- binding env overrides
- validating the resulting config

This keeps consumption strongly typed and makes loading behavior explicit.

### Viper setup
Use Viper with:
- defaults set before reading any files
- optional `config.yaml` in the working directory (do not error if missing)
- env prefix `TASKBOARD_`
- key replacer mapping `.` to `_` so `db.url` maps to `TASKBOARD_DB_URL`
- env overrides win over file values

### Validation approach + required fields
Validation occurs in `Load()` after Viper unmarshals into `Config`. If required
values are missing or invalid, return a descriptive error and exit before
serving requests. At minimum, require `db.url` (or equivalent field in
`Config`) and validate any other critical fields (e.g., port range, timeouts)
as they are added.

### Where main/service should read config
The entrypoint (`main` or the service startup package) should call
`config.Load()` at process start, handle errors by logging a clear message, and
exit with non-zero status before wiring the server.

## Consequences
- Config logic is centralized and testable.
- The app has deterministic precedence: defaults < optional file < env.
- Misconfiguration fails fast with actionable errors.
