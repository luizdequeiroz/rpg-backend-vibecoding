#!/bin/bash

# Script para executar todos os testes do projeto
# Uso: ./run-tests.sh [opcoes]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funções de output
print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Configurar variáveis de ambiente para testes
export DATABASE_URL="file::memory:?cache=shared"
export JWT_SECRET="test-jwt-secret-for-testing-only"
export LOG_LEVEL="error"

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    print_error "Go não está instalado ou não está no PATH"
    exit 1
fi

print_header "Iniciando execução dos testes"
echo "Go version: $(go version)"
echo "Database: $DATABASE_URL"
echo ""

# Limpar cache de testes
print_header "Limpando cache de testes"
go clean -testcache
print_success "Cache limpo"
echo ""

# Verificar formatação do código
print_header "Verificando formatação do código"
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    print_error "Arquivos não formatados encontrados:"
    echo "$UNFORMATTED"
    print_warning "Execute: go fmt ./..."
    exit 1
fi
print_success "Código está formatado corretamente"
echo ""

# Executar go vet
print_header "Executando análise estática (go vet)"
if go vet ./...; then
    print_success "Análise estática passou"
else
    print_error "Análise estática falhou"
    exit 1
fi
echo ""

# Executar testes unitários
print_header "Executando testes unitários"
if CGO_ENABLED=0 go test ./pkg/... ./internal/app/services/... -v -coverprofile=unit-coverage.out; then
    print_success "Testes unitários passaram"
else
    print_error "Testes unitários falharam"
    exit 1
fi
echo ""

# Executar testes de integração
print_header "Executando testes de integração"
if CGO_ENABLED=0 go test ./tests/integration/... -v -coverprofile=integration-coverage.out; then
    print_success "Testes de integração passaram"
else
    print_error "Testes de integração falharam"
    exit 1
fi
echo ""

# Executar todos os testes com cobertura
print_header "Executando todos os testes com cobertura"
if CGO_ENABLED=0 go test ./... -v -coverprofile=coverage.out -covermode=atomic; then
    print_success "Todos os testes passaram"
else
    print_error "Alguns testes falharam"
    exit 1
fi
echo ""

# Verificar cobertura de código
print_header "Verificando cobertura de código"
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "Cobertura total: ${COVERAGE}%"

# Verificar se atende o mínimo de 80%
COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')
if (( $(echo "$COVERAGE_NUM < 80.0" | bc -l) )); then
    print_error "Cobertura abaixo do mínimo (80%): ${COVERAGE}%"
    exit 1
fi
print_success "Cobertura aprovada: ${COVERAGE}%"
echo ""

# Gerar relatório HTML de cobertura
print_header "Gerando relatório HTML de cobertura"
go tool cover -html=coverage.out -o coverage.html
print_success "Relatório gerado: coverage.html"
echo ""

# Executar benchmarks
print_header "Executando benchmarks"
go test ./pkg/roll/... -bench=. -benchmem -run=^$ > benchmark-results.txt
print_success "Benchmarks executados - resultados em benchmark-results.txt"
echo ""

# Mostrar resumo dos arquivos gerados
print_header "Arquivos gerados"
echo "📊 coverage.out - Dados de cobertura"
echo "📊 unit-coverage.out - Cobertura dos testes unitários"
echo "📊 integration-coverage.out - Cobertura dos testes de integração"
echo "📊 coverage.html - Relatório HTML de cobertura"
echo "📊 benchmark-results.txt - Resultados dos benchmarks"
echo ""

print_success "Todos os testes foram executados com sucesso!"
print_success "Cobertura de código: ${COVERAGE}%"

# Abrir relatório HTML se estiver em ambiente desktop
if command -v xdg-open &> /dev/null; then
    print_warning "Abrindo relatório de cobertura no navegador..."
    xdg-open coverage.html
elif command -v open &> /dev/null; then
    print_warning "Abrindo relatório de cobertura no navegador..."
    open coverage.html
fi
