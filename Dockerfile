# Multi-stage build para aplicação em produção
FROM golang:1.24-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação com otimizações
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o rpg-backend \
    cmd/api/main.go

# Build do utilitário de migração
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o rpg-migrate \
    cmd/migrate/main.go

# Stage final - imagem mínima
FROM alpine:latest

# Instalar certificados SSL e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -s /bin/sh rpguser

# Criar diretórios necessários
RUN mkdir -p /app/data /app/migrations
RUN chown -R rpguser:rpguser /app

# Definir diretório de trabalho
WORKDIR /app

# Copiar binários do stage de build
COPY --from=builder /app/rpg-backend .
COPY --from=builder /app/rpg-migrate .

# Copiar migrations se existirem
COPY --chown=rpguser:rpguser migrations/ ./migrations/

# Definir permissões dos binários
RUN chmod +x rpg-backend rpg-migrate

# Mudar para usuário não-root
USER rpguser

# Expor porta
EXPOSE 8080

# Variáveis de ambiente padrão
ENV DATABASE_URL="file:./data/rpg.db?cache=shared&mode=rwc"
ENV JWT_SECRET=""
ENV LOG_LEVEL="info"
ENV HOST="0.0.0.0"
ENV PORT="8080"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando padrão
CMD ["./rpg-backend"]
