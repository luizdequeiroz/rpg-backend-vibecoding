# Makefile para RPG Backend

# Variáveis
BINARY_NAME=rpg-backend
MIGRATION_DIR=./migrations

# Comandos de build
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: dev
dev:
	LOG_LEVEL=debug go run cmd/api/main.go

.PHONY: clean
clean:
	rm -rf bin/

# Comandos de teste
.PHONY: test
test:
	go test -v ./...

.PHONY: test-api
test-api:
	@echo "Testando endpoints da API..."
	@echo "Healthcheck:"
	@curl -s http://localhost:8080/health | jq .
	@echo "\nAPI Info:"
	@curl -s http://localhost:8080/ | jq .
	@echo "\nUsers:"
	@curl -s http://localhost:8080/api/v1/users | jq .

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Comandos de migração
.PHONY: migrate-up
migrate-up:
	go run cmd/migrate/main.go -action=up

.PHONY: migrate-down
migrate-down:
	go run cmd/migrate/main.go -action=down

.PHONY: migrate-status
migrate-status:
	go run cmd/migrate/main.go -action=status

.PHONY: migrate-reset
migrate-reset:
	go run cmd/migrate/main.go -action=reset

.PHONY: migrate-create
migrate-create:
	@read -p "Nome da migração: " name; \
	go run cmd/migrate/main.go -action=create -name=$$name

# Comandos de desenvolvimento
.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: swagger-generate
swagger-generate:
	swag init -g cmd/api/main.go -o docs

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

# Comando completo de verificação
.PHONY: check
check: fmt vet test

# Comando de ajuda
.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  build          - Compila o projeto"
	@echo "  run            - Executa o servidor"
	@echo "  dev            - Executa em modo desenvolvimento (debug)"
	@echo "  clean          - Remove arquivos de build"
	@echo "  test           - Executa os testes"
	@echo "  test-api       - Testa endpoints da API (servidor deve estar rodando)"
	@echo "  test-coverage  - Executa testes com coverage"
	@echo "  migrate-up     - Executa todas as migrações"
	@echo "  migrate-down   - Desfaz a última migração"
	@echo "  migrate-status - Mostra status das migrações"
	@echo "  migrate-reset  - Remove todas as migrações"
	@echo "  migrate-create - Cria uma nova migração"
	@echo "  deps           - Baixa e organiza dependências"
	@echo "  swagger-generate - Gera documentação Swagger"
	@echo "  fmt            - Formata o código"
	@echo "  vet            - Verifica o código"
	@echo "  check          - Executa fmt, vet e test"
	@echo "  help           - Mostra esta ajuda"
