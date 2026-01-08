# Run Instructions

## Defaults
Defaults load first, so you can run with no config file and no env overrides.

Auth service:
```
go run ./services/auth-service/cmd/auth-service
```

Task service:
```
go run ./services/task-service/cmd/task-service
```

## Optional config.yaml
Copy `config.example.yaml` into the working directory where you run the service,
rename it to `config.yaml`, and adjust values as needed (for example, update
`http.addr` to `:8082` when running the task service).

## Environment overrides
Set `TASKBOARD_`-prefixed variables to override file/defaults.

Examples:
```
set TASKBOARD_DB_URL=postgres://postgres:postgres@localhost:5432/go_micro_task_board?sslmode=disable
set TASKBOARD_HTTP_ADDR=:8082
set TASKBOARD_JWT_SECRET=dev-secret
```
