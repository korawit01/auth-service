# PRD: Viper Standard Config

## User stories
- As a developer, I want configuration to load from defaults, then an optional config file, then environment overrides so I can run locally and in deployment with predictable precedence.
- As a developer, I want environment variables prefixed with TASKBOARD_ to map to config keys so I can configure the app in CI and containers without editing files.
- As an operator, I want the app to fail fast with a clear error when required configuration is missing or invalid so misconfiguration is caught early.
- As a developer, I want concise run instructions so I can start the app with defaults, a config file, or environment variables.

## Acceptance criteria
1. When the application starts, it uses built-in defaults, then overrides with values from config.yaml if present, then overrides with environment variables.
2. When config.yaml is missing, the application still starts using defaults unless a required value is missing.
3. Environment variables are read only when prefixed with TASKBOARD_ and map to config keys using Viper-compatible key mapping (for example, TASKBOARD_DB_URL maps to db.url).
4. When both config.yaml and TASKBOARD_ environment variables specify the same setting, the environment value takes precedence.
5. The application validates required configuration (for example, db.url) at startup and fails fast before serving requests.
6. On validation failure, the application exits with a non-zero status and prints a clear, actionable error message identifying the missing or invalid key.
7. Documentation includes how to run the app with defaults, with an optional config.yaml file, and with TASKBOARD_ environment variables.
