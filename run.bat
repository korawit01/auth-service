@echo off
setlocal

rem Always run from the project root (location of this script)
cd /d "%~dp0"

start "task-service" cmd /k "go run services\task-service\cmd\task-service\main.go"
start "auth-service" cmd /k "go run services\auth-service\cmd\auth-service\main.go"
start "api-gateway" cmd /k "go run api-gateway\cmd\api-gateway\main.go"

echo Started task-service, auth-service, and api-gateway in separate windows.
echo Close each window or hit Ctrl+C in it to stop the service.
pause >nul
