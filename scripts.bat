@echo off
REM Script de automação para RPG Backend - Windows

if "%1"=="" goto help

if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="dev" goto dev
if "%1"=="test" goto test
if "%1"=="test-api" goto test-api
if "%1"=="test-integration" goto test-integration
if "%1"=="test-full" goto test-full
if "%1"=="migrate-up" goto migrate-up
if "%1"=="migrate-down" goto migrate-down
if "%1"=="migrate-status" goto migrate-status
if "%1"=="migrate-reset" goto migrate-reset
if "%1"=="migrate-create" goto migrate-create
if "%1"=="deps" goto deps
if "%1"=="swagger-generate" goto swagger-generate
if "%1"=="fmt" goto fmt
if "%1"=="clean" goto clean
goto help

:build
echo Compilando o projeto...
go build -o bin\rpg-backend.exe cmd\api\main.go
goto end

:run
echo Executando o servidor...
go run cmd\api\main.go
goto end

:dev
echo Executando o servidor em modo desenvolvimento...
set LOG_LEVEL=debug
go run cmd\api\main.go
goto end

:test
echo Executando testes...
go test -v .\...
goto end

:test-api
echo Testando endpoints da API...
echo Healthcheck:
curl -s http://localhost:8080/health
echo.
echo API Info:
curl -s http://localhost:8080/
echo.
echo Users:
curl -s http://localhost:8080/api/v1/users
goto end

:test-integration
echo Executando testes de integração...
echo Verificando se servidor está rodando...
powershell -Command "try { Invoke-RestMethod -Uri 'http://localhost:8080/health' -TimeoutSec 5 | Out-Null; Write-Host 'Servidor online - executando testes...' } catch { Write-Host 'Erro: Servidor não está rodando. Execute: scripts.bat run'; exit 1 }"
if errorlevel 1 goto end
powershell -ExecutionPolicy Bypass -File "test\scripts\integration_test_simple.ps1"
goto end

:test-full
echo Executando teste completo com todos os cenários...
powershell -ExecutionPolicy Bypass -File "test\scripts\integration_test_simple.ps1" -Verbose
goto end

:migrate-up
echo Executando migrações...
go run cmd\migrate\main.go -action=up
goto end

:migrate-down
echo Desfazendo última migração...
go run cmd\migrate\main.go -action=down
goto end

:migrate-status
echo Verificando status das migrações...
go run cmd\migrate\main.go -action=status
goto end

:migrate-reset
echo Resetando todas as migrações...
go run cmd\migrate\main.go -action=reset
goto end

:migrate-create
if "%2"=="" (
    set /p name="Nome da migração: "
) else (
    set name=%2
)
echo Criando migração: %name%
go run cmd\migrate\main.go -action=create -name=%name%
goto end

:deps
echo Baixando e organizando dependências...
go mod download
go mod tidy
goto end

:swagger-generate
echo Gerando documentação Swagger...
swag init -g cmd\api\main.go -o docs
goto end

:fmt
echo Formatando código...
go fmt .\...
goto end

:clean
echo Limpando arquivos de build...
if exist bin rmdir /s /q bin
if exist coverage.out del coverage.out
if exist coverage.html del coverage.html
goto end

:help
echo Comandos disponíveis:
echo   build          - Compila o projeto
echo   run            - Executa o servidor
echo   dev            - Executa em modo desenvolvimento (debug)
echo   test           - Executa os testes
echo   test-api       - Testa endpoints da API (servidor deve estar rodando)
echo   test-api       - Testa endpoints básicos da API
echo   test-integration - Executa testes de integração automatizados
echo   test-full      - Executa testes completos com verbose
echo   migrate-up     - Executa todas as migrações
echo   migrate-down   - Desfaz a última migração
echo   migrate-status - Mostra status das migrações
echo   migrate-reset  - Remove todas as migrações
echo   migrate-create [nome] - Cria uma nova migração
echo   deps           - Baixa e organiza dependências
echo   swagger-generate - Gera documentação Swagger
echo   fmt            - Formata o código
echo   clean          - Remove arquivos de build
echo   help           - Mostra esta ajuda
echo.
echo Uso: scripts.bat [comando] [argumentos]
goto end

:end
