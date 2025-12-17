@echo off
setlocal

rem Always run from the project root (location of this script)
cd /d "%~dp0"

echo Starting stack with Docker Compose...
docker compose version >nul 2>&1
if errorlevel 1 (
    echo docker compose not found. Please install Docker Desktop (with Compose V2) and try again.
    exit /b 1
)

docker compose up --build
