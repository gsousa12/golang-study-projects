@echo off

REM Iniciar api-gateway
start cmd /k "cd /d %~dp0api-gateway && go run main.go"

REM Iniciar conversion-service
start cmd /k "cd /d %~dp0conversion-service && go run main.go"

REM Iniciar transaction-service
start cmd /k "cd /d %~dp0transaction-service && go run main.go"
