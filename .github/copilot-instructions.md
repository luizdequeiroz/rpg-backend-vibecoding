````instructions
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

## Fase 7: Sistema de Mesas de Jogo e Convites

### Funcionalidades Implementadas
- ✅ **Modelos de GameTable e Invite**:
  - GameTable com ID (UUID), Name, System, OwnerID, timestamps
  - Invite com ID (UUID), TableID, InviterID, InviteeID, Status, timestamps
  - Relacionamentos FK com usuários e constrains de integridade
  
- ✅ **Repositórios de Dados**:
  - GameTableRepository para operações CRUD de mesas
  - InviteRepository para gerenciamento de convites
  - Queries otimizadas com JOINs para dados relacionados
  - Verificações de permissão e validações de negócio
  
- ✅ **Serviços de Negócio**:
  - GameTableService com lógica de autorização
  - Validação de permissões (owner vs convidado)
  - Prevenção de auto-convites e convites duplicados
  - Controle de status de convites (pending/accepted/declined)
  
- ✅ **Endpoints REST Completos**:
  - CRUD completo para mesas com autorização
  - Sistema de convites com controle de acesso
  - Validação de permissões em todos os endpoints
  - Tratamento de erros específicos (403, 404, 409)
  
- ✅ **Migrações de Banco**:
  - `20250711120000_create_game_tables.sql` - Tabela de mesas
  - `20250711120100_create_invites.sql` - Tabela de convites
  - Índices para performance e constraints de integridade

### Endpoints de Mesas de Jogo

#### Gestão de Mesas
- `POST /api/v1/tables` - Criar mesa (autenticado, owner = JWT.UserID) 🔒
- `GET /api/v1/tables` - Lista mesas do usuário (owner ou convidado aceito) 🔒
- `GET /api/v1/tables/{id}` - Detalhes da mesa (inclui lista de invites) 🔒
- `PUT /api/v1/tables/{id}` - Atualiza nome/sistema (só owner) 🔒
- `DELETE /api/v1/tables/{id}` - Remove mesa (só owner) 🔒

#### Gestão de Convites
- `POST /api/v1/tables/{id}/invites` - Criar convite (body: invitee_email) 🔒
- `GET /api/v1/tables/{id}/invites` - Lista convites (owner e convidados) 🔒
- `POST /api/v1/tables/{id}/invites/{inviteId}/accept` - Aceitar convite (só invitee) 🔒
- `POST /api/v1/tables/{id}/invites/{inviteId}/decline` - Recusar convite (só invitee) 🔒

### Exemplos de Uso

#### Criar Mesa
```bash
# Autenticar primeiro
$token = "seu_jwt_token_aqui"
$headers = @{Authorization="Bearer $token"; "Content-Type"="application/json"}

# Criar nova mesa
$mesa = @{
    name = "Mesa D&D: A Busca pelo Artefato Perdido"
    system = "D&D 5e"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Method POST -Body $mesa -Headers $headers
```

#### Listar Mesas do Usuário
```bash
# Lista todas as mesas onde é owner ou convidado aceito
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Headers $headers
```

#### Convidar Jogador
```bash
# Convidar usuário por email
$convite = @{
    invitee_email = "jogador@exemplo.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables/{mesa_id}/invites' -Method POST -Body $convite -Headers $headers
```

#### Aceitar Convite
```bash
# Usuário convidado aceita o convite
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables/{mesa_id}/invites/{invite_id}/accept' -Method POST -Headers $headers
```

### Estrutura de Dados

#### GameTable
```json
{
  "id": "uuid-da-mesa",
  "name": "Mesa D&D: Aventura Épica",
  "system": "D&D 5e",
  "owner_id": 1,
  "owner": {
    "id": 1,
    "email": "mestre@exemplo.com"
  },
  "invites": [
    {
      "id": "uuid-do-convite",
      "invitee": {
        "id": 2,
        "email": "jogador@exemplo.com"
      },
      "status": "accepted",
      "created_at": "2025-07-11T12:00:00Z"
    }
  ],
  "created_at": "2025-07-11T10:00:00Z",
  "updated_at": "2025-07-11T10:00:00Z"
}
```

#### Invite
```json
{
  "id": "uuid-do-convite",
  "table_id": "uuid-da-mesa",
  "inviter_id": 1,
  "invitee_id": 2,
  "status": "pending", // pending, accepted, declined
  "inviter": {
    "id": 1,
    "email": "mestre@exemplo.com"
  },
  "invitee": {
    "id": 2,
    "email": "jogador@exemplo.com"
  },
  "created_at": "2025-07-11T11:00:00Z",
  "updated_at": "2025-07-11T11:00:00Z"
}
```

### Validações e Autorização

#### Permissões de Mesa
- **Owner**: Pode criar/editar/deletar mesa, criar convites, ver todos os convites
- **Convidado Aceito**: Pode ver detalhes da mesa, ver convites
- **Convidado Pendente**: Pode aceitar/recusar próprio convite
- **Outros**: Sem acesso (403 Forbidden)

#### Validações de Negócio
- **Nome da Mesa**: 3-100 caracteres, obrigatório
- **Sistema**: 2-50 caracteres, obrigatório
- **Email do Convidado**: Deve ser usuário existente
- **Auto-convite**: Não permitido
- **Convite Duplicado**: Não permitido
- **Status de Convite**: Só pode mudar de "pending"

#### Códigos de Erro
- **400**: Dados inválidos no request
- **401**: Token JWT inválido ou ausente
- **403**: Sem permissão para a operação
- **404**: Mesa, convite ou usuário não encontrado
- **409**: Conflito (convite duplicado, já respondido)
- **500**: Erro interno do servidor

### Arquivos Criados/Modificados
- `migrations/20250711120000_create_game_tables.sql` - Tabela mesas
- `migrations/20250711120100_create_invites.sql` - Tabela convites
- `internal/app/models/game_table.go` - Modelos e estruturas
- `internal/app/repositories/game_table.go` - Repositórios
- `internal/app/services/game_table.go` - Lógica de negócio
- `internal/bff/game_table.go` - Handlers HTTP
- `internal/bff/handler.go` - Integração das rotas

### Observações Importantes

#### UUID vs Integer ID
- **Mesas e Convites**: Usam UUID para evitar enumeration attacks
- **Usuários**: Mantêm ID integer por compatibilidade

#### Segurança
- **Todas as rotas protegidas**: Requerem JWT válido
- **Autorização granular**: Validação por operação
- **Prevenção de leaks**: Não expõe dados sensíveis

#### Performance
- **Índices otimizados**: Para queries frequentes
- **JOINs eficientes**: Para evitar N+1 queries
- **Paginação**: Implementada onde necessário

### Próximos Passos
- [ ] Adicionar notificações de convites por email
- [ ] Implementar websockets para updates em tempo real
- [ ] Adicionar roles dentro das mesas (player, co-master)
- [ ] Implementar sistema de sessões de jogo
- [ ] Adicionar logs de auditoria para ações importantes
- [ ] Implementar soft delete para mesas arquivadas

## Fase 8: Sistema PlayerSheet e Motor de Dados

### Funcionalidades Implementadas
- ✅ **Migrações de Banco**:
  - `20250711150000_create_player_sheets.sql` - Tabela de fichas de personagens
  - `20250711160100_create_rolls.sql` - Tabela de histórico de rolagens
  - Relacionamentos FK com game_tables, users, sheet_templates
  - Índices otimizados para performance

- ✅ **Modelos de Dados**:
  - PlayerSheet com dados JSON flexíveis para diferentes sistemas
  - PlayerSheetData com validação e conversão de tipos
  - CreateRollRequest com suporte a expressões e campos da ficha
  - RollResponse com detalhes completos da rolagem
  - Sistema de validação com tags Go

- ✅ **Motor de Dados Custom**:
  - Engine de parsing de expressões regex (pkg/roll/)
  - Suporte a expressões como "1d20+3", "2d6+STR", "3d8-1"
  - Sistema de critical/fumble configurável
  - Rolagem baseada em campos da ficha de personagem
  - Geração de números aleatórios criptograficamente seguros

- ✅ **Repositórios de Dados**:
  - PlayerSheetRepository com CRUD completo
  - Queries otimizadas com JOINs para relacionamentos
  - Sistema de paginação para listas grandes
  - Validações de permissão por mesa e usuário

- ✅ **Serviços de Negócio**:
  - PlayerSheetService com lógica de autorização
  - Validação de acesso à mesa antes de operações
  - Sistema de propriedade (owner vs membros da mesa)
  - Integração com motor de dados para rolagens

- ✅ **Endpoints REST**:
  - CRUD completo para fichas de personagens
  - Sistema de rolagem de dados independente
  - Histórico de rolagens por ficha e por mesa
  - Documentação Swagger completa

### Endpoints de PlayerSheet

#### Gestão de Fichas
- `POST /api/v1/sheets` - Criar ficha (body: table_id, template_id, name, data) 🔒
- `GET /api/v1/sheets?table_id={id}` - Listar fichas da mesa 🔒
- `GET /api/v1/sheets/{id}` - Detalhes da ficha 🔒
- `PUT /api/v1/sheets/{id}` - Atualizar ficha (name, data) 🔒
- `DELETE /api/v1/sheets/{id}` - Remover ficha 🔒

#### Sistema de Rolagem
- `POST /api/v1/rolls` - Rolar dados (body: sheet_id, expression ou field_name) 🔒
- `GET /api/v1/rolls/sheet/{sheetID}` - Histórico de rolagens da ficha 🔒
- `GET /api/v1/rolls/table/{tableID}` - Histórico de rolagens da mesa 🔒

### Exemplos de Uso

#### Criar Ficha de Personagem
```bash
# Autenticar primeiro
$token = "seu_jwt_token_aqui"
$headers = @{Authorization="Bearer $token"; "Content-Type"="application/json"}

# Criar ficha D&D 5e
$ficha = @{
    table_id = "uuid-da-mesa"
    template_id = 1
    name = "Elara, a Élfica Arcana"
    data = @{
        attributes = @{
            strength = 12
            dexterity = 16
            constitution = 14
            intelligence = 18
            wisdom = 13
            charisma = 15
        }
        skills = @{
            arcana = 8
            investigation = 6
            perception = 3
        }
        combat = @{
            armor_class = 15
            hit_points = 32
            speed = 30
        }
    }
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/sheets' -Method POST -Body $ficha -Headers $headers
```

#### Listar Fichas da Mesa
```bash
# Listar todas as fichas de uma mesa específica
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/sheets?table_id=uuid-da-mesa' -Headers $headers
```

#### Rolar Dados
```bash
# Rolagem de expressão livre
$rolagem = @{
    sheet_id = "uuid-da-ficha"
    expression = "1d20+5"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/rolls' -Method POST -Body $rolagem -Headers $headers

# Rolagem baseada em campo da ficha
$teste_arcana = @{
    sheet_id = "uuid-da-ficha"
    field_name = "skills.arcana"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/rolls' -Method POST -Body $teste_arcana -Headers $headers
```

#### Histórico de Rolagens
```bash
# Ver rolagens de uma ficha específica
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/rolls/sheet/uuid-da-ficha' -Headers $headers

# Ver todas as rolagens da mesa
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/rolls/table/uuid-da-mesa' -Headers $headers
```

### Estrutura de Dados

#### PlayerSheet
```json
{
  "id": "uuid-da-ficha",
  "table_id": "uuid-da-mesa",
  "template_id": 1,
  "owner_id": 2,
  "name": "Elara, a Élfica Arcana",
  "data": {
    "attributes": {
      "strength": 12,
      "dexterity": 16,
      "constitution": 14,
      "intelligence": 18,
      "wisdom": 13,
      "charisma": 15
    },
    "skills": {
      "arcana": 8,
      "investigation": 6,
      "perception": 3
    },
    "combat": {
      "armor_class": 15,
      "hit_points": 32,
      "speed": 30
    }
  },
  "template": {
    "id": 1,
    "name": "Ficha D&D 5e"
  },
  "owner": {
    "id": 2,
    "email": "player@exemplo.com"
  },
  "created_at": "2025-07-11T15:00:00Z",
  "updated_at": "2025-07-11T15:00:00Z"
}
```

#### Roll Response
```json
{
  "id": "uuid-da-rolagem",
  "sheet_id": "uuid-da-ficha",
  "expression": "1d20+8",
  "field_name": "skills.arcana",
  "result": 18,
  "details": {
    "dice": [10],
    "modifier": 8,
    "total": 18,
    "critical": false,
    "fumble": false,
    "success": true
  },
  "roller": {
    "id": 2,
    "email": "player@exemplo.com"
  },
  "created_at": "2025-07-11T15:30:00Z"
}
```

### Motor de Dados - Expressões Suportadas

#### Expressões Básicas
- `1d20` - Um dado de 20 faces
- `2d6` - Dois dados de 6 faces
- `3d8+5` - Três dados de 8 faces mais 5
- `1d100-10` - Um dado de 100 faces menos 10

#### Expressões com Campos
- `skills.arcana` - Usa valor do campo skills.arcana da ficha
- `attributes.strength` - Usa valor do atributo força
- `combat.armor_class` - Usa valor da classe de armadura

#### Recursos Avançados
- **Critical/Fumble**: Detecta 20 natural (critical) e 1 natural (fumble)
- **Validação**: Expressões malformadas retornam erro descritivo
- **Limites**: Máximo 20 dados, máximo d1000, modificador ±999
- **Segurança**: Geração criptograficamente segura de números

### Validações e Autorização

#### Permissões de Ficha
- **Owner da Ficha**: Pode editar/deletar própria ficha
- **Membros da Mesa**: Podem ver fichas da mesa, rolar dados
- **Owner da Mesa**: Pode deletar qualquer ficha da mesa
- **Outros**: Sem acesso (403 Forbidden)

#### Validações de Dados
- **Nome da Ficha**: 3-100 caracteres, obrigatório
- **TableID**: UUID válido, mesa deve existir
- **TemplateID**: Template deve existir no banco
- **Data**: JSON válido, campos opcionais
- **Expressão de Dados**: Sintaxe validada pelo motor

#### Códigos de Erro
- **400**: Dados inválidos, expressão malformada
- **401**: Token JWT inválido ou ausente
- **403**: Sem permissão para operação
- **404**: Ficha, mesa ou template não encontrado
- **500**: Erro interno do servidor

### Arquivos Implementados

#### Migrações
- `migrations/20250711150000_create_player_sheets.sql`
- `migrations/20250711160100_create_rolls.sql`

#### Motor de Dados
- `pkg/roll/engine.go` - Motor principal de rolagem
- `pkg/roll/parser.go` - Parser de expressões regex
- `pkg/roll/types.go` - Tipos e estruturas

#### Modelos e Domínio
- `internal/app/models/player_sheet.go` - Modelos completos
- `internal/app/repositories/player_sheet.go` - Camada de dados
- `internal/app/services/player_sheet.go` - Lógica de negócio

#### API e Handlers
- `internal/bff/player_sheet.go` - Handlers HTTP
- `internal/bff/handler.go` - Integração de rotas

#### Documentação
- `docs/` - Swagger atualizado com endpoints PlayerSheet

### Performance e Otimizações

#### Consultas Otimizadas
- **Índices**: Criados para table_id, owner_id, template_id
- **JOINs Eficientes**: Busca relacionamentos em uma query
- **Paginação**: Implementada para listas grandes
- **Cache**: Prepared statements para queries frequentes

#### Motor de Dados
- **Regex Compilado**: Patterns compilados uma vez na inicialização
- **Pool de Random**: Gerador único para toda aplicação
- **Validação Rápida**: Checks básicos antes de parsing completo

### Observações Importantes

#### Padrão de Rotas
- **Separação**: Fichas em `/sheets`, rolagens em `/rolls`
- **Evita Conflitos**: Não sobrepõe com rotas de `/tables/:id`
- **RESTful**: Seguindo convenções REST para recursos

#### Flexibilidade de Dados
- **JSON Livre**: Campo `data` aceita qualquer estrutura
- **Validação Opcional**: Não força schema rígido
- **Extensibilidade**: Fácil adicionar novos campos sem migração

#### Segurança
- **UUIDs**: Evita enumeration attacks em fichas
- **Autorização Granular**: Validação em cada operação
- **Sanitização**: Input validado antes de processamento

### Próximos Passos
- [ ] Implementar templates de rolagem personalizados
- [ ] Adicionar macros de dados complexas
- [ ] Implementar sistema de vantagem/desvantagem (D&D 5e)
- [ ] Adicionar modificadores temporários nas fichas
- [ ] Implementar iniciativa e ordem de turnos
- [ ] Criar sistema de notas e anotações nas fichas
````