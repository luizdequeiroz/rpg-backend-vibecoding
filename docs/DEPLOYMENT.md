# Deployment & Production Guide - RPG Backend

## 🚀 Visão Geral

Este guia cobre o processo de deployment e configuração para produção do RPG Backend, incluindo configurações de banco de dados, variáveis de ambiente, monitoramento e boas práticas de segurança.

---

## 📋 Pré-requisitos de Produção

### Sistema Operacional
- **Linux**: Ubuntu 20.04+ / CentOS 8+ / Debian 11+ (recomendado)
- **Windows**: Windows Server 2019+ (com suporte a Go)
- **Docker**: Suporte completo com containers

### Runtime
- **Go**: versão 1.21 ou superior
- **Banco de Dados**: PostgreSQL 13+ (recomendado) ou SQLite
- **Proxy Reverso**: Nginx ou Apache (recomendado)
- **SSL/TLS**: Certificado válido para HTTPS

### Hardware Mínimo
- **CPU**: 2 cores
- **RAM**: 2GB (4GB recomendado)
- **Armazenamento**: 10GB SSD
- **Rede**: 100Mbps

---

## 🔧 Configuração de Ambiente

### Variáveis de Ambiente Produção

```bash
# Server Configuration
HOST=0.0.0.0
PORT=8080
GIN_MODE=release

# Database
DATABASE_URL=postgres://rpg_user:secure_password@localhost:5432/rpg_production?sslmode=require

# Security
JWT_SECRET=sua_chave_super_secreta_minimo_32_caracteres
JWT_EXPIRATION=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Performance
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s
MAX_HEADER_BYTES=1048576

# Database Connection Pool
MAX_IDLE_CONNS=10
MAX_OPEN_CONNS=100
CONN_MAX_LIFETIME=1h

# CORS (produção - domínios específicos)
CORS_ORIGINS=https://seudominio.com,https://www.seudominio.com
CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_HEADERS=Content-Type,Authorization

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
```

### Arquivo de Configuração (.env)

```bash
# Criar arquivo .env para produção
cat > .env << 'EOF'
# RPG Backend Production Configuration

# Server
HOST=0.0.0.0
PORT=8080
GIN_MODE=release

# Database
DATABASE_URL=postgres://rpg_user:${DB_PASSWORD}@${DB_HOST}:5432/rpg_production?sslmode=require

# Security
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRATION=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Performance
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s
MAX_HEADER_BYTES=1048576

# Connection Pool
MAX_IDLE_CONNS=10
MAX_OPEN_CONNS=100
CONN_MAX_LIFETIME=1h

# CORS
CORS_ORIGINS=https://seudominio.com
CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_HEADERS=Content-Type,Authorization

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
EOF
```

---

## 🗄️ Configuração de Banco de Dados

### PostgreSQL (Recomendado para Produção)

#### Instalação Ubuntu/Debian
```bash
# Instalar PostgreSQL
sudo apt update
sudo apt install -y postgresql postgresql-contrib

# Iniciar serviço
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Configurar usuário e banco
sudo -u postgres psql << 'EOF'
CREATE USER rpg_user WITH PASSWORD 'sua_senha_segura';
CREATE DATABASE rpg_production OWNER rpg_user;
GRANT ALL PRIVILEGES ON DATABASE rpg_production TO rpg_user;
\q
EOF
```

#### Configuração de Segurança PostgreSQL
```bash
# Editar configuração PostgreSQL
sudo nano /etc/postgresql/13/main/postgresql.conf

# Configurações recomendadas
listen_addresses = 'localhost'
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100

# Configurar autenticação
sudo nano /etc/postgresql/13/main/pg_hba.conf

# Adicionar linha para aplicação
local   rpg_production   rpg_user                     md5
host    rpg_production   rpg_user     127.0.0.1/32    md5

# Reiniciar PostgreSQL
sudo systemctl restart postgresql
```

#### Backup Automatizado
```bash
#!/bin/bash
# Script: /opt/rpg-backend/scripts/backup.sh

set -e

BACKUP_DIR="/opt/rpg-backend/backups"
DB_NAME="rpg_production"
DB_USER="rpg_user"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/rpg_backup_$TIMESTAMP.sql"

# Criar diretório se não existir
mkdir -p "$BACKUP_DIR"

# Fazer backup
PGPASSWORD="$DB_PASSWORD" pg_dump -h localhost -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"

# Comprimir backup
gzip "$BACKUP_FILE"

# Manter apenas últimos 7 backups
find "$BACKUP_DIR" -name "rpg_backup_*.sql.gz" -mtime +7 -delete

echo "Backup concluído: ${BACKUP_FILE}.gz"
```

#### Crontab para Backups
```bash
# Adicionar ao crontab
crontab -e

# Backup diário às 2:00 AM
0 2 * * * /opt/rpg-backend/scripts/backup.sh >> /var/log/rpg-backup.log 2>&1
```

---

## 🐳 Deployment com Docker

### Dockerfile Otimizado
```dockerfile
# Multi-stage build para produção
FROM golang:1.21-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache git ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -g '' appuser

# Definir diretório de trabalho
WORKDIR /build

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Download de dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Stage final - imagem mínima
FROM scratch

# Copiar certificados CA
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar usuário
COPY --from=builder /etc/passwd /etc/passwd

# Copiar binário
COPY --from=builder /build/main /app/main

# Copiar migrations
COPY --from=builder /build/migrations /app/migrations

# Criar diretório de dados
COPY --from=builder /build/data /app/data

# Definir usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "--health-check"]

# Comando de inicialização
ENTRYPOINT ["/app/main"]
```

### Docker Compose para Produção
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  rpg-backend:
    build:
      context: .
      dockerfile: Dockerfile
    image: rpg-backend:latest
    container_name: rpg-backend-prod
    restart: unless-stopped
    ports:
      - "127.0.0.1:8080:8080"
    environment:
      - HOST=0.0.0.0
      - PORT=8080
      - GIN_MODE=release
      - DATABASE_URL=postgres://rpg_user:${DB_PASSWORD}@postgres:5432/rpg_production?sslmode=require
      - JWT_SECRET=${JWT_SECRET}
      - LOG_LEVEL=info
      - LOG_FORMAT=json
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - rpg-network
    healthcheck:
      test: ["CMD", "/app/main", "--health-check"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: rpg-postgres-prod
    restart: unless-stopped
    environment:
      - POSTGRES_DB=rpg_production
      - POSTGRES_USER=rpg_user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_INITDB_ARGS=--auth-host=md5
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - rpg-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U rpg_user -d rpg_production"]
      interval: 10s
      timeout: 5s
      retries: 5

  nginx:
    image: nginx:alpine
    container_name: rpg-nginx-prod
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/sites-available:/etc/nginx/sites-available:ro
      - ./ssl:/etc/nginx/ssl:ro
      - ./logs/nginx:/var/log/nginx
    depends_on:
      - rpg-backend
    networks:
      - rpg-network

volumes:
  postgres_data:
    driver: local

networks:
  rpg-network:
    driver: bridge
```

### Script de Deploy
```bash
#!/bin/bash
# Script: deploy.sh

set -e

echo "🚀 Iniciando deployment do RPG Backend..."

# Verificar se .env existe
if [ ! -f ".env" ]; then
    echo "❌ Arquivo .env não encontrado!"
    exit 1
fi

# Carregar variáveis de ambiente
source .env

# Fazer backup antes do deploy
echo "📦 Fazendo backup..."
./scripts/backup.sh

# Parar containers antigos
echo "🛑 Parando containers..."
docker-compose -f docker-compose.prod.yml down

# Fazer pull das imagens base mais recentes
echo "📥 Atualizando imagens base..."
docker-compose -f docker-compose.prod.yml pull postgres nginx

# Build da nova versão
echo "🔨 Fazendo build da aplicação..."
docker-compose -f docker-compose.prod.yml build --no-cache rpg-backend

# Executar migrações
echo "🗄️ Executando migrações..."
docker-compose -f docker-compose.prod.yml run --rm rpg-backend \
    /app/main --migrate-up

# Iniciar serviços
echo "▶️ Iniciando serviços..."
docker-compose -f docker-compose.prod.yml up -d

# Aguardar serviços ficarem saudáveis
echo "⏳ Aguardando serviços..."
timeout 120 bash -c '
  until docker-compose -f docker-compose.prod.yml ps | grep -q "healthy"; do
    sleep 5
    echo "Aguardando serviços ficarem saudáveis..."
  done
'

# Verificar health check
echo "🏥 Verificando saúde da aplicação..."
sleep 10
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Deploy concluído com sucesso!"
    echo "🌐 API disponível em: http://localhost:8080"
    echo "📚 Documentação: http://localhost:8080/docs/"
else
    echo "❌ Deploy falhou - aplicação não responde ao health check!"
    echo "📋 Logs da aplicação:"
    docker-compose -f docker-compose.prod.yml logs rpg-backend
    exit 1
fi
```

---

## 🌐 Configuração Nginx

### Configuração Principal
```nginx
# /etc/nginx/nginx.conf
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                   '$status $body_bytes_sent "$http_referer" '
                   '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    # Performance
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 10M;

    # Gzip
    gzip on;
    gzip_vary on;
    gzip_min_length 10240;
    gzip_proxied expired no-cache no-store private must-revalidate auth;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;

    # Include site configurations
    include /etc/nginx/sites-available/*;
}
```

### Site Configuration
```nginx
# /etc/nginx/sites-available/rpg-backend
server {
    listen 80;
    server_name seudominio.com www.seudominio.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name seudominio.com www.seudominio.com;

    # SSL Configuration
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # API Rate Limiting
    location /api/v1/auth {
        limit_req zone=auth burst=10 nodelay;
        proxy_pass http://127.0.0.1:8080;
        include /etc/nginx/proxy_params;
    }

    location /api {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://127.0.0.1:8080;
        include /etc/nginx/proxy_params;
    }

    # Health check (sem rate limiting)
    location /health {
        proxy_pass http://127.0.0.1:8080;
        include /etc/nginx/proxy_params;
    }

    # Documentation
    location /docs {
        proxy_pass http://127.0.0.1:8080;
        include /etc/nginx/proxy_params;
    }

    # Root
    location / {
        proxy_pass http://127.0.0.1:8080;
        include /etc/nginx/proxy_params;
    }

    # Logging
    access_log /var/log/nginx/rpg-backend.access.log main;
    error_log /var/log/nginx/rpg-backend.error.log warn;
}
```

### Proxy Parameters
```nginx
# /etc/nginx/proxy_params
proxy_set_header Host $http_host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host $server_name;
proxy_redirect off;

# Timeouts
proxy_connect_timeout 30s;
proxy_send_timeout 30s;
proxy_read_timeout 30s;

# Buffering
proxy_buffering on;
proxy_buffer_size 4k;
proxy_buffers 8 4k;
proxy_busy_buffers_size 8k;
```

---

## 📊 Monitoramento e Logs

### Configuração de Logs Estruturados
```go
// internal/pkg/logger/logger.go
package logger

import (
    "os"
    "github.com/sirupsen/logrus"
)

func SetupLogger() *logrus.Logger {
    logger := logrus.New()
    
    // Formato baseado na variável de ambiente
    if os.Getenv("LOG_FORMAT") == "json" {
        logger.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: "2006-01-02T15:04:05.000Z",
            FieldMap: logrus.FieldMap{
                logrus.FieldKeyTime:  "timestamp",
                logrus.FieldKeyLevel: "level",
                logrus.FieldKeyMsg:   "message",
            },
        })
    } else {
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:   true,
            TimestampFormat: "2006-01-02 15:04:05",
        })
    }
    
    // Nível de log
    level := os.Getenv("LOG_LEVEL")
    switch level {
    case "debug":
        logger.SetLevel(logrus.DebugLevel)
    case "warn":
        logger.SetLevel(logrus.WarnLevel)
    case "error":
        logger.SetLevel(logrus.ErrorLevel)
    default:
        logger.SetLevel(logrus.InfoLevel)
    }
    
    return logger
}
```

### Health Check Avançado
```go
// internal/app/health/health.go
package health

import (
    "context"
    "database/sql"
    "time"
)

type HealthChecker struct {
    db *sql.DB
}

type HealthStatus struct {
    Status    string                 `json:"status"`
    Timestamp time.Time             `json:"timestamp"`
    Version   string                `json:"version"`
    Services  map[string]string     `json:"services"`
    Metrics   map[string]interface{} `json:"metrics,omitempty"`
}

func (h *HealthChecker) Check(ctx context.Context) HealthStatus {
    status := HealthStatus{
        Status:    "ok",
        Timestamp: time.Now(),
        Version:   "1.0.0",
        Services:  make(map[string]string),
        Metrics:   make(map[string]interface{}),
    }
    
    // Check database
    if err := h.checkDatabase(ctx); err != nil {
        status.Services["database"] = "error"
        status.Status = "unhealthy"
    } else {
        status.Services["database"] = "ok"
    }
    
    // Add metrics if enabled
    if os.Getenv("METRICS_ENABLED") == "true" {
        status.Metrics = h.getMetrics()
    }
    
    return status
}

func (h *HealthChecker) checkDatabase(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return h.db.PingContext(ctx)
}

func (h *HealthChecker) getMetrics() map[string]interface{} {
    return map[string]interface{}{
        "uptime":           time.Since(startTime).Seconds(),
        "goroutines":       runtime.NumGoroutine(),
        "memory_alloc":     runtime.MemStats{}.Alloc,
        "requests_total":   requestCounter.Value(),
        "response_time_avg": responseTimeAvg.Value(),
    }
}
```

### Script de Monitoramento
```bash
#!/bin/bash
# Script: monitor.sh

LOG_FILE="/var/log/rpg-backend-monitor.log"
API_URL="http://localhost:8080/health"
SLACK_WEBHOOK="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

check_health() {
    local response=$(curl -s -w "%{http_code}" -o /tmp/health_response.json "$API_URL")
    local http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        local status=$(jq -r '.status' /tmp/health_response.json 2>/dev/null)
        if [ "$status" = "ok" ]; then
            echo "$(date): API healthy" >> "$LOG_FILE"
            return 0
        else
            echo "$(date): API unhealthy - status: $status" >> "$LOG_FILE"
            return 1
        fi
    else
        echo "$(date): API not responding - HTTP $http_code" >> "$LOG_FILE"
        return 1
    fi
}

send_alert() {
    local message="$1"
    curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"🚨 RPG Backend Alert: $message\"}" \
        "$SLACK_WEBHOOK"
}

# Verificar saúde
if ! check_health; then
    send_alert "API health check failed at $(date)"
    
    # Tentar restart se configurado
    if [ "$AUTO_RESTART" = "true" ]; then
        echo "$(date): Attempting restart..." >> "$LOG_FILE"
        docker-compose -f /opt/rpg-backend/docker-compose.prod.yml restart rpg-backend
        sleep 30
        
        if check_health; then
            send_alert "API restarted successfully at $(date)"
        else
            send_alert "API restart failed at $(date) - manual intervention required"
        fi
    fi
fi
```

### Crontab para Monitoramento
```bash
# Adicionar ao crontab
crontab -e

# Verificar saúde a cada 5 minutos
*/5 * * * * /opt/rpg-backend/scripts/monitor.sh

# Rotação de logs diária
0 1 * * * /opt/rpg-backend/scripts/rotate_logs.sh
```

---

## 🔒 Configurações de Segurança

### Firewall (UFW)
```bash
# Configurar firewall
sudo ufw --force reset
sudo ufw default deny incoming
sudo ufw default allow outgoing

# SSH
sudo ufw allow 22/tcp

# HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# PostgreSQL (apenas local)
sudo ufw allow from 127.0.0.1 to any port 5432

# Ativar firewall
sudo ufw --force enable

# Verificar status
sudo ufw status verbose
```

### SSL/TLS com Let's Encrypt
```bash
# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obter certificado
sudo certbot --nginx -d seudominio.com -d www.seudominio.com

# Renovação automática
sudo crontab -e

# Adicionar linha para renovação automática
0 3 * * * certbot renew --quiet --deploy-hook "systemctl reload nginx"
```

### Configurações do Sistema
```bash
# Limites de arquivo
echo "fs.file-max = 65536" >> /etc/sysctl.conf

# Network tuning
echo "net.core.somaxconn = 1024" >> /etc/sysctl.conf
echo "net.ipv4.tcp_max_syn_backlog = 1024" >> /etc/sysctl.conf

# Aplicar configurações
sudo sysctl -p
```

---

## 📈 Performance e Escalabilidade

### Load Balancing com Nginx
```nginx
# Configuração upstream para múltiplas instâncias
upstream rpg_backend {
    least_conn;
    server 127.0.0.1:8080 weight=1 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8081 weight=1 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8082 weight=1 max_fails=3 fail_timeout=30s;
    
    # Health check
    keepalive 32;
}

server {
    # ... outras configurações ...
    
    location / {
        proxy_pass http://rpg_backend;
        proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
        proxy_next_upstream_tries 3;
        proxy_next_upstream_timeout 10s;
        
        include /etc/nginx/proxy_params;
    }
}
```

### Docker Compose Escalável
```yaml
# docker-compose.scale.yml
version: '3.8'

services:
  rpg-backend:
    build: .
    restart: unless-stopped
    environment:
      - HOST=0.0.0.0
      - PORT=8080
      - DATABASE_URL=postgres://rpg_user:${DB_PASSWORD}@postgres:5432/rpg_production?sslmode=require
    depends_on:
      - postgres
    networks:
      - rpg-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/load-balancer.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - rpg-backend
    networks:
      - rpg-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=rpg_production
      - POSTGRES_USER=rpg_user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - rpg-network

volumes:
  postgres_data:

networks:
  rpg-network:
```

### Comando para Escalar
```bash
# Escalar para 3 instâncias da aplicação
docker-compose -f docker-compose.scale.yml up -d --scale rpg-backend=3

# Verificar instâncias rodando
docker-compose -f docker-compose.scale.yml ps
```

---

## 🔄 CI/CD Pipeline

### GitHub Actions
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [ main ]
  release:
    types: [ published ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      run: |
        go mod download
        go test -v ./...
    
    - name: Run linting
      uses: golangci/golangci-lint-action@v3

  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Deploy to server
      uses: appleboy/ssh-action@v0.1.5
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SERVER_SSH_KEY }}
        script: |
          cd /opt/rpg-backend
          git pull origin main
          ./deploy.sh
```

---

## 🚨 Troubleshooting

### Problemas Comuns

#### 1. Aplicação não inicia
```bash
# Verificar logs
docker-compose -f docker-compose.prod.yml logs rpg-backend

# Verificar configuração de rede
docker network ls
docker network inspect rpg-backend_rpg-network

# Verificar variáveis de ambiente
docker-compose -f docker-compose.prod.yml config
```

#### 2. Banco de dados não conecta
```bash
# Testar conexão direta
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U rpg_user -d rpg_production -c "SELECT version();"

# Verificar logs do PostgreSQL
docker-compose -f docker-compose.prod.yml logs postgres

# Verificar configuração de rede
docker-compose -f docker-compose.prod.yml exec rpg-backend \
    nslookup postgres
```

#### 3. Alto uso de CPU/Memória
```bash
# Monitorar recursos
docker stats

# Verificar logs por erros
docker-compose -f docker-compose.prod.yml logs rpg-backend | grep ERROR

# Verificar conexões de banco
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U rpg_user -d rpg_production -c "SELECT * FROM pg_stat_activity;"
```

### Scripts de Diagnóstico
```bash
#!/bin/bash
# Script: diagnose.sh

echo "🔍 RPG Backend Diagnostic Report"
echo "================================"

echo "📅 Data: $(date)"
echo "🖥️ Host: $(hostname)"
echo "👤 User: $(whoami)"

echo -e "\n🐳 Docker Status:"
docker --version
docker-compose --version
docker system df

echo -e "\n📦 Containers:"
docker-compose -f docker-compose.prod.yml ps

echo -e "\n🌐 Network:"
ss -tlnp | grep ':80\|:443\|:8080\|:5432'

echo -e "\n💾 Disk Usage:"
df -h

echo -e "\n🧠 Memory:"
free -h

echo -e "\n⚡ CPU:"
top -bn1 | head -5

echo -e "\n🏥 API Health:"
curl -s http://localhost:8080/health | jq . || echo "API não disponível"

echo -e "\n📋 Recent Logs (last 50 lines):"
docker-compose -f docker-compose.prod.yml logs --tail=50 rpg-backend
```

---

## 📚 Recursos Adicionais

### Documentação
- [README.md](README.md) - Visão geral do projeto
- [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitetura e design
- [DEVELOPMENT.md](DEVELOPMENT.md) - Guia de desenvolvimento
- [API.md](API.md) - Referência completa da API

### Ferramentas Recomendadas
- **Monitoramento**: Prometheus + Grafana
- **Logs**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **APM**: Jaeger para tracing distribuído
- **Backup**: pgBackRest para PostgreSQL
- **Secrets**: HashiCorp Vault
- **Load Testing**: k6 ou Apache Bench

### Próximos Passos
1. Implementar métricas com Prometheus
2. Configurar alertas avançados
3. Implementar tracing distribuído
4. Adicionar cache com Redis
5. Configurar CDN para assets estáticos
6. Implementar backup incremental
7. Configurar disaster recovery

---

**Este guia garante um deployment robusto e escalável do RPG Backend em produção. Para suporte adicional, consulte a documentação técnica ou entre em contato com a equipe de DevOps.**
