@echo off
REM Script para executar todos os testes do projeto no Windows
REM Uso: run-tests.bat

setlocal enabledelayedexpansion

REM Configurar variáveis de ambiente para testes
set DATABASE_URL=file::memory:?cache=shared
set JWT_SECRET=test-jwt-secret-for-testing-only
set LOG_LEVEL=error
set CGO_ENABLED=0

echo === Iniciando execucao dos testes ===
go version
echo Database: %DATABASE_URL%
echo.

REM Verificar se Go está instalado
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go nao esta instalado ou nao esta no PATH
    exit /b 1
)

REM Limpar cache de testes
echo === Limpando cache de testes ===
go clean -testcache
echo ✅ Cache limpo
echo.

REM Verificar formatação do código
echo === Verificando formatacao do codigo ===
go fmt ./...
echo ✅ Codigo formatado
echo.

REM Executar go vet
echo === Executando analise estatica (go vet) ===
go vet ./...
if errorlevel 1 (
    echo ❌ Analise estatica falhou
    exit /b 1
)
echo ✅ Analise estatica passou
echo.

REM Executar testes unitários
echo === Executando testes unitarios ===
go test ./pkg/... ./internal/app/services/... -v -coverprofile=unit-coverage.out
if errorlevel 1 (
    echo ❌ Testes unitarios falharam
    exit /b 1
)
echo ✅ Testes unitarios passaram
echo.

REM Executar testes de integração
echo === Executando testes de integracao ===
go test ./tests/integration/... -v -coverprofile=integration-coverage.out
if errorlevel 1 (
    echo ❌ Testes de integracao falharam
    exit /b 1
)
echo ✅ Testes de integracao passaram
echo.

REM Executar todos os testes com cobertura
echo === Executando todos os testes com cobertura ===
go test ./... -v -coverprofile=coverage.out -covermode=atomic
if errorlevel 1 (
    echo ❌ Alguns testes falharam
    exit /b 1
)
echo ✅ Todos os testes passaram
echo.

REM Verificar cobertura de código
echo === Verificando cobertura de codigo ===
go tool cover -func=coverage.out | findstr "total"
echo.

REM Gerar relatório HTML de cobertura
echo === Gerando relatorio HTML de cobertura ===
go tool cover -html=coverage.out -o coverage.html
echo ✅ Relatorio gerado: coverage.html
echo.

REM Executar benchmarks
echo === Executando benchmarks ===
go test ./pkg/roll/... -bench=. -benchmem -run=^$ > benchmark-results.txt
echo ✅ Benchmarks executados - resultados em benchmark-results.txt
echo.

REM Mostrar resumo dos arquivos gerados
echo === Arquivos gerados ===
echo 📊 coverage.out - Dados de cobertura
echo 📊 unit-coverage.out - Cobertura dos testes unitarios
echo 📊 integration-coverage.out - Cobertura dos testes de integracao
echo 📊 coverage.html - Relatorio HTML de cobertura
echo 📊 benchmark-results.txt - Resultados dos benchmarks
echo.

echo ✅ Todos os testes foram executados com sucesso!

REM Abrir relatório HTML no navegador
if exist coverage.html (
    echo ⚠️  Abrindo relatorio de cobertura no navegador...
    start coverage.html
)

endlocal
