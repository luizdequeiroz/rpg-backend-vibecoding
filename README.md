# RPG Session Backend

Backend desenvolvido em Go para gerenciamento de mesas de RPG (Role-Playing Game). Sistema completo com autenticação JWT, criação de mesas, sistema de convites e autorização baseada em roles.

## 🎮 Funcionalidades

- ✅ **Autenticação JWT** - Sistema completo de login/signup
- ✅ **Gerenciamento de Mesas** - CRUD completo para mesas de RPG
- ✅ **Sistema de Convites** - Convide jogadores e gerencie permissões
- ✅ **Autorização por Roles** - Mestres (owners) e Jogadores (players)
- ✅ **Múltiplos Sistemas** - D&D 5e, Vampiro, Call of Cthulhu, etc.
- ✅ **API RESTful** - Endpoints padronizados com Swagger
- ✅ **Banco SQLite** - Desenvolvimento local simplificado
- ✅ **Testes Automatizados** - Scripts PowerShell e Bash

## 🏗️ Arquitetura

```
rpg-backend/
├── cmd/api/                    # Ponto de entrada da aplicação
├── internal/
│   ├── bff/                   # Backend For Frontend (handlers HTTP)
│   └── app/
│       ├── models/            # Modelos de domínio
│       ├── repositories/      # Camada de dados
│       └── services/          # Lógica de negócio
├── migrations/                # Migrações do banco de dados
├── test/
│   ├── fixtures/             # Dados de teste (JSON)
│   └── scripts/              # Scripts de automação
├── pkg/                      # Pacotes utilitários
└── docs/                     # Documentação e diagramas
```

## 🚀 Como Executar

### Pré-requisitos
- Go 1.21+
- Git

### Instalação e Execução

```bash
# Clone o repositório
git clone https://github.com/luizdequeiroz/rpg-backend-vibecoding.git
cd rpg-backend

# Instale dependências
go mod tidy

# Execute migrações
go run cmd/migrate/main.go up

# Inicie o servidor
go run cmd/api/main.go
```

### Usando Scripts de Automação

```bash
# Windows
scripts.bat run                    # Inicia o servidor
scripts.bat test-integration       # Executa testes automatizados
scripts.bat test-full              # Executa testes com verbose

# Linux/macOS
./scripts.sh run
./scripts.sh test
```

## 🧪 Testes

### Teste Rápido
```bash
# Verificação básica da API
.\test\scripts\quick_test.ps1
```

### Teste Completo
```bash
# Teste completo com todos os cenários
.\test\scripts\test_game_tables_v2.ps1
```

### Cenários de Teste Disponíveis
- ✅ Registro e autenticação de usuários
- ✅ Criação de mesas de diferentes sistemas
- ✅ Sistema completo de convites
- ✅ Autorização baseada em roles
- ✅ Múltiplos jogadores e permissões

## 📡 API Endpoints

### Autenticação
- `POST /api/v1/auth/signup` - Registrar usuário
- `POST /api/v1/auth/login` - Login

### Mesas de Jogo
- `GET /api/v1/tables/` - Listar mesas do usuário
- `POST /api/v1/tables/` - Criar nova mesa
- `GET /api/v1/tables/{id}` - Detalhes da mesa
- `PUT /api/v1/tables/{id}` - Atualizar mesa
- `DELETE /api/v1/tables/{id}` - Excluir mesa

### Convites
- `POST /api/v1/tables/{id}/invites/` - Criar convite
- `GET /api/v1/tables/{id}/invites/` - Listar convites
- `POST /api/v1/tables/{id}/invites/{invite_id}/accept` - Aceitar convite
- `POST /api/v1/tables/{id}/invites/{invite_id}/decline` - Recusar convite

### Utilitários
- `GET /health` - Health check
- `GET /` - Informações da API

## 🎯 Tecnologias Utilizadas

- **Go 1.21+** - Linguagem principal
- **Gin** - Framework HTTP
- **SQLite** - Banco de dados
- **JWT** - Autenticação
- **UUID** - Identificadores únicos
- **Goose** - Migrações de banco
- **Swagger** - Documentação da API

## 📊 Estrutura de Dados

### GameTable
```json
{
  "id": "uuid",
  "name": "string",
  "system": "string",
  "owner_id": "integer",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Invite
```json
{
  "id": "uuid",
  "table_id": "uuid",
  "inviter_id": "integer", 
  "invitee_id": "integer",
  "status": "pending|accepted|declined",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## 🔧 Configuração

### Variáveis de Ambiente
```bash
# Servidor
HOST=localhost
PORT=8080

# Banco de dados
DB_PATH=./data/rpg.db

# JWT
JWT_SECRET=your-secret-key

# Log
LOG_LEVEL=info
```

## 📖 Documentação Adicional

- [Testes](./test/README.md) - Guia completo de testes
- [Arquitetura](./docs/architecture.md) - Diagramas e decisões arquiteturais
- [Copilot Instructions](./docs/copilot-instructions.md) - Guia para desenvolvimento

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Status do Projeto

🎯 **Funcional e Pronto para Produção**

- ✅ Sistema de autenticação completo
- ✅ CRUD de mesas implementado
- ✅ Sistema de convites funcional
- ✅ Testes automatizados
- ✅ Documentação completa
- ✅ Scripts de deploy prontos

🚧 **Em Desenvolvimento** - Este projeto está em fase inicial de desenvolvimento.

## Contribuição

Contribuições são bem-vindas! Por favor, abra uma issue ou envie um pull request.

## Licença

Este projeto está sob a licença MIT.
