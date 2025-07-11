# Copilot Instructions

# Instruções para o Copilot
- Sempre responda e gere comentários de documentação em português, independentemente do idioma da pergunta ou do código fornecido.
- Sempre atualize este arquivo quando houver mudanças significativas no projeto ou na estrutura do código.

## Fase 1: Scaffold do Projeto

### Estrutura Criada
- ✅ Módulo Go inicializado: `github.com/luizdequeiroz/rpg-backend`
- ✅ Diretórios criados:
  - `cmd/api/` - Ponto de entrada da aplicação
  - `internal/bff/` - Backend For Frontend
  - `internal/app/` - Lógica de negócio
  - `pkg/config/` - Configurações
  - `pkg/db/` - Banco de dados
  - `docs/` - Documentação
- ✅ Arquivos base:
  - `go.mod` - Definição do módulo
  - `README.md` - Documentação do projeto

### Próximos Passos
- [x] Implementar estrutura básica da API
- [x] Configurar conexão com banco de dados
- [x] Definir modelos de dados
- [ ] Implementar autenticação
- [ ] Criar endpoints básicos

## Fase 2: Banco de Dados e Migrações

### Funcionalidades Implementadas
- ✅ **Conexão de Banco**: 
  - Abstração com sqlx para fácil troca de driver
  - Suporte a SQLite (padrão) e configuração via `DATABASE_URL`
  - Pool de conexões configurado
  - Health check da conexão
  
- ✅ **Sistema de Migrações**:
  - Integração com goose para gerenciamento de migrações
  - Execução automática na inicialização da aplicação
  - CLI para gerenciamento manual de migrações
  
- ✅ **Migrações Criadas**:
  - `00001_create_users_table.sql` - Tabela de usuários
  - `00002_create_campaigns_table.sql` - Tabela de campanhas
  - `00003_create_characters_table.sql` - Tabela de personagens
  - `20250710190310_update_users_for_auth.sql` - Ajustes na tabela de usuários para autenticação
  - `20250710190617_create_sheet_templates_table.sql` - Tabela de templates de fichas

### Como Executar Migrações

#### Via Makefile (Recomendado)
```bash
# Executar todas as migrações pendentes
make migrate-up

# Desfazer a última migração
make migrate-down

# Ver status das migrações
make migrate-status

# Resetar todas as migrações
make migrate-reset

# Criar nova migração (modo interativo)
make migrate-create
```

#### Via CLI Direta
```bash
# Executar migrações
go run cmd/migrate/main.go -action=up

# Desfazer migração
go run cmd/migrate/main.go -action=down

# Status das migrações
go run cmd/migrate/main.go -action=status

# Criar nova migração
go run cmd/migrate/main.go -action=create -name=nome_da_migracao

# Resetar migrações
go run cmd/migrate/main.go -action=reset
```

#### Via Script Windows (scripts.bat)
```batch
# Executar migrações
scripts.bat migrate-up

# Desfazer migração
scripts.bat migrate-down

# Status das migrações
scripts.bat migrate-status

# Criar nova migração
scripts.bat migrate-create nome_da_migracao

# Resetar migrações
scripts.bat migrate-reset
```

### Configuração do Banco
- **Padrão**: SQLite local em `./data/rpg.db`
- **Configuração**: Via variável de ambiente `DATABASE_URL`
- **Exemplo PostgreSQL**: `DATABASE_URL=postgres://user:password@localhost/rpg_db`

### Estrutura de Dados
- **Users**: Usuários do sistema (jogadores e mestres)
- **Campaigns**: Campanhas de RPG
- **Characters**: Personagens dos jogadores
- **Sheet_Templates**: Templates de fichas de personagens
- **Relações**: FK constraints entre tabelas para integridade

### Próximos Passos
- [x] Implementar modelos de domínio (structs Go)
- [x] Criar repositories para acesso aos dados
- [x] Implementar autenticação JWT
- [x] Criar endpoints REST para CRUD básico
- [x] Adicionar validação de dados

## Fase 3: API HTTP e BFF Layer

### Funcionalidades Implementadas
- ✅ **Sistema de Configuração**:
  - Configuração via variáveis de ambiente
  - Configurações para servidor, banco, auth e logs
  - Valores padrão para desenvolvimento
  
- ✅ **Servidor HTTP**:
  - Framework Gin para alta performance
  - Middlewares de logging e recovery
  - CORS configurado para desenvolvimento
  - Graceful shutdown implementado
  
- ✅ **BFF Layer (Backend For Frontend)**:
  - Endpoints RESTful em `/api/v1`
  - Estrutura organizada por recursos
  - Handlers temporários para todos os CRUDs
  
- ✅ **Endpoints Implementados**:
  - `GET /health` - Healthcheck com status de serviços
  - `GET /` - Informações da API e endpoints disponíveis
  - `GET|POST|PUT|DELETE /api/v1/users` - Gestão de usuários
  - `GET|POST|PUT|DELETE /api/v1/campaigns` - Gestão de campanhas
  - `GET|POST|PUT|DELETE /api/v1/characters` - Gestão de personagens
  - `GET|POST|PUT|DELETE /api/v1/sessions` - Gestão de sessões

### Como Executar a API

#### Iniciar Servidor
```bash
# Via script Windows
scripts.bat run

# Via Makefile
make run

# Via comando direto
go run cmd/api/main.go
```

#### Modo Desenvolvimento (Debug)
```bash
# Via script Windows
scripts.bat dev

# Via Makefile
make dev
```

#### Testar Endpoints
```bash
# Via script Windows (servidor deve estar rodando)
scripts.bat test-api

# Via Makefile
make test-api

# Manualmente
curl http://localhost:8080/health
curl http://localhost:8080/
curl http://localhost:8080/api/v1/users
```

### Configuração do Servidor
- **Padrão**: `localhost:8080`
- **Configuração**: Via `HOST` e `PORT` no ambiente
- **Logs**: Configurável via `LOG_LEVEL` (info, debug)
- **CORS**: Habilitado para desenvolvimento

### URLs da API
- **Servidor**: `http://localhost:8080`
- **Healthcheck**: `http://localhost:8080/health`
- **Documentação**: `http://localhost:8080/` (lista endpoints)
- **API v1**: `http://localhost:8080/api/v1`

### Próximos Passos
- [ ] Implementar modelos de domínio com validação
- [ ] Criar services na camada `internal/app`
- [ ] Implementar repositories com queries reais
- [ ] Adicionar middleware de autenticação JWT
- [ ] Implementar testes unitários e de integração
- [x] Adicionar documentação OpenAPI/Swagger

## Fase 4: Documentação API com Swagger

### Funcionalidades Implementadas
- ✅ **Integração com swaggo/swag**:
  - Biblioteca swag para geração automática de documentação
  - Integração com gin-swagger para servir a UI
  - Anotações Go para definir especificação OpenAPI
  
- ✅ **Documentação Swagger UI**:
  - Interface web em `/docs/index.html`
  - Documentação interativa da API
  - Especificação OpenAPI 2.0 gerada automaticamente
  
- ✅ **Endpoints Documentados**:
  - `GET /health` - Healthcheck com exemplos de resposta
  - Modelos de dados com exemplos (HealthResponse)
  - Tags organizadas por funcionalidade
  
- ✅ **Scripts de Automação**:
  - `make swagger-generate` - Gera documentação
  - `scripts.bat swagger-generate` - Versão Windows
  - Integração nos fluxos de build

### Como Gerar e Acessar Documentação

#### Gerar Documentação Swagger
```bash
# Via Makefile
make swagger-generate

# Via script Windows
scripts.bat swagger-generate

# Via comando direto
swag init -g cmd/api/main.go -o docs
```

#### Acessar Documentação
- **Swagger UI**: `http://localhost:8080/docs/index.html`
- **JSON OpenAPI**: `http://localhost:8080/docs/swagger.json`
- **YAML OpenAPI**: `http://localhost:8080/docs/swagger.yaml`

### Estrutura de Arquivos Swagger
- `docs/docs.go` - Código Go gerado automaticamente
- `docs/swagger.json` - Especificação OpenAPI em JSON
- `docs/swagger.yaml` - Especificação OpenAPI em YAML

### Anotações Implementadas
- **API Info**: Título, versão, descrição, contato
- **Endpoints**: Sumário, descrição, tags, parâmetros
- **Modelos**: Estruturas com exemplos e tipos
- **Respostas**: Códigos HTTP com descrições

### Próximos Passos
- [ ] Documentar todos os endpoints da API v1
- [ ] Adicionar autenticação Bearer Token no Swagger
- [ ] Incluir exemplos de requisições para POST/PUT
- [ ] Configurar validação de esquemas
- [ ] Adicionar testes da documentação

## Fase 5: Sistema de Autenticação e Autorização

### Funcionalidades Implementadas
- ✅ **Modelos de Autenticação**:
  - Modelo User com ID, email, password_hash, created_at, updated_at
  - Estruturas para requests de signup e login
  - Response models para autenticação com JWT
  
- ✅ **Serviço de Autenticação**:
  - Hash de senhas com bcrypt
  - Geração e validação de tokens JWT
  - Verificação de credenciais
  - Gerenciamento de usuários
  
- ✅ **Middleware JWT**:
  - Middleware obrigatório para rotas protegidas
  - Middleware opcional para rotas que podem ter auth
  - Extração de informações do usuário do contexto
  
- ✅ **Endpoints de Autenticação**:
  - `POST /api/v1/auth/signup` - Registro de novos usuários
  - `POST /api/v1/auth/login` - Login e geração de JWT
  - `GET /api/v1/auth/me` - Informações do usuário autenticado (protegido)
  
- ✅ **Documentação Swagger**:
  - Anotações completas para todos os endpoints de auth
  - Suporte a Bearer Token Authentication
  - Exemplos de requests e responses
  - Códigos de erro documentados

### Como Usar Autenticação

#### Registrar Novo Usuário
```bash
# Via PowerShell
$body = @{email='user@example.com'; password='senha123'} | ConvertTo-Json
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/signup' -Method POST -Body $body -ContentType 'application/json'
```

#### Fazer Login
```bash
# Via PowerShell
$body = @{email='user@example.com'; password='senha123'} | ConvertTo-Json
$response = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -Body $body -ContentType 'application/json'
$token = $response.token
```

#### Acessar Endpoints Protegidos
```bash
# Via PowerShell
$headers = @{Authorization="Bearer $token"}
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/me' -Headers $headers
```

### Configuração de Segurança
- **JWT Secret**: Configurável via `JWT_SECRET` no ambiente
- **Expiração**: Tokens expiram em 24 horas
- **Hash de Senhas**: bcrypt com custo padrão
- **Validação**: Email único e senha mínima de 6 caracteres

### Migração de Banco
- ✅ **Migração 20250710190310**: Ajustou tabela users para modelo de auth
- ✅ **Campos Opcionais**: username e display_name tornaram-se opcionais
- ✅ **Integridade**: Mantida compatibilidade com dados existentes

### Arquivos Criados/Modificados
- `internal/app/models/user.go` - Modelos de autenticação
- `internal/app/services/auth.go` - Serviço de autenticação
- `internal/app/middleware/auth.go` - Middleware JWT
- `internal/bff/auth.go` - Handlers de autenticação
- `migrations/20250710190310_update_users_for_auth.sql` - Migração de ajuste

### Próximos Passos
- [ ] Proteger endpoints de usuários, campanhas e personagens com JWT
- [ ] Implementar roles/permissões (admin, player, master)
- [ ] Adicionar refresh tokens
- [ ] Implementar reset de senha via email
- [ ] Adicionar rate limiting nos endpoints de auth
- [ ] Implementar logout com blacklist de tokens

## Fase 6: CRUD Completo para SheetTemplate

### Funcionalidades Implementadas
- ✅ **Modelo SheetTemplate**:
  - ID autoincremental único
  - Nome obrigatório (validação mínimo 3 caracteres)
  - Descrição opcional
  - Definition (JSON armazenado como string no banco)
  - Timestamps automáticos (created_at, updated_at)
  
- ✅ **Repositório de Dados**:
  - CRUD completo com sqlx
  - Queries SQL otimizadas
  - Tratamento de erros de banco
  - Conversão JSON automática para definition
  
- ✅ **Serviço de Negócio**:
  - Validação de payloads de entrada
  - Regras de negócio centralizadas
  - Tratamento de casos de erro específicos
  
- ✅ **Endpoints REST**:
  - `GET /api/v1/templates` - Listar todos os templates
  - `POST /api/v1/templates` - Criar novo template
  - `GET /api/v1/templates/{id}` - Buscar template por ID
  - `PUT /api/v1/templates/{id}` - Atualizar template
  - `DELETE /api/v1/templates/{id}` - Deletar template
  
- ✅ **Documentação Swagger**:
  - Modelos completos com exemplos
  - Códigos de resposta documentados (200, 201, 400, 404, 500)
  - Parâmetros de path e body documentados
  - Exemplos de payloads JSON

### Como Usar Templates

#### Listar Todos os Templates
```bash
# Via PowerShell
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/templates' -Method GET

# Via curl
curl -X GET http://localhost:8080/api/v1/templates
```

#### Criar Novo Template
```bash
# Via PowerShell
$template = @{
    name = "Ficha D&D 5e"
    description = "Template para personagens de Dungeons & Dragons 5ª edição"
    definition = @{
        sections = @(
            @{
                name = "Atributos"
                fields = @(
                    @{name = "Força"; type = "number"; min = 1; max = 20}
                    @{name = "Destreza"; type = "number"; min = 1; max = 20}
                    @{name = "Constituição"; type = "number"; min = 1; max = 20}
                )
            }
        )
    }
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/templates' -Method POST -Body $template -ContentType 'application/json'
```

#### Buscar Template por ID
```bash
# Via PowerShell
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/templates/1' -Method GET

# Via curl
curl -X GET http://localhost:8080/api/v1/templates/1
```

#### Atualizar Template
```bash
# Via PowerShell
$update = @{
    name = "Ficha D&D 5e Atualizada"
    description = "Template atualizado para D&D 5e com mais campos"
    definition = @{
        sections = @(
            @{
                name = "Atributos Básicos"
                fields = @(
                    @{name = "Força"; type = "number"; min = 1; max = 20}
                    @{name = "Destreza"; type = "number"; min = 1; max = 20}
                )
            }
        )
    }
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/templates/1' -Method PUT -Body $update -ContentType 'application/json'
```

#### Deletar Template
```bash
# Via PowerShell
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/templates/1' -Method DELETE

# Via curl
curl -X DELETE http://localhost:8080/api/v1/templates/1
```

### Estrutura do Campo Definition
O campo `definition` é um JSON flexível que pode conter:
```json
{
  "sections": [
    {
      "name": "Nome da Seção",
      "fields": [
        {
          "name": "Nome do Campo",
          "type": "text|number|boolean|select",
          "required": true,
          "options": ["Opção 1", "Opção 2"],
          "min": 1,
          "max": 20,
          "default": "Valor Padrão"
        }
      ]
    }
  ],
  "rules": {
    "calculations": [],
    "validations": []
  }
}
```

### Migração de Banco
- ✅ **Migração 20250710190617**: Criou tabela sheet_templates
- ✅ **Campos**: id, name, description, definition (TEXT para JSON), timestamps
- ✅ **Índices**: Criado índice único no campo name

### Validações Implementadas
- **Nome**: Obrigatório, mínimo 3 caracteres
- **Description**: Opcional
- **Definition**: Deve ser JSON válido
- **ID**: Validação de existência para operações de update/delete

### Tratamento de Erros
- **400 Bad Request**: Payload inválido ou dados incorretos
- **404 Not Found**: Template não encontrado
- **500 Internal Server Error**: Erros de banco ou servidor
- **Mensagens**: Descritivas em português para facilitar debugging

### Arquivos Criados/Modificados
- `internal/app/models/sheet_template.go` - Modelos e estruturas
- `internal/app/repositories/sheet_template.go` - Camada de dados
- `internal/app/services/sheet_template.go` - Lógica de negócio
- `internal/bff/sheet_template.go` - Handlers HTTP
- `migrations/20250710190617_create_sheet_templates_table.sql` - Migração
- `docs/` - Documentação Swagger atualizada

### Próximos Passos
- [ ] Implementar autenticação nos endpoints de templates (se necessário)
- [ ] Criar testes unitários para repositories e services
- [ ] Adicionar testes de integração para endpoints
- [ ] Implementar versionamento de templates
- [ ] Adicionar validação avançada de schemas de definition
- [ ] Implementar importação/exportação de templates