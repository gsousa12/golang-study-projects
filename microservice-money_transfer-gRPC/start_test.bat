@echo off

REM Iniciar script de teste - /transfer
start cmd /k "cd /d %~dp0scripts && go run main.go"
