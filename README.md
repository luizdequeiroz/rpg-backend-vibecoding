# RPG Session Backend

## Descrição

Este é um backend desenvolvido em Go para gerenciamento de sessões de RPG (Role-Playing Game). O sistema fornece APIs para criar, gerenciar e facilitar sessões de jogos de RPG, incluindo funcionalidades para:

- Gerenciamento de campanhas
- Controle de personagens
- Sistema de combate
- Gerenciamento de inventário
- Chat e comunicação entre jogadores
- Sistema de dados (dice rolling)

## Estrutura do Projeto

```
rpg-backend/
├── cmd/api/          # Ponto de entrada da aplicação
├── internal/bff/     # Backend For Frontend - camada de orquestração
├── internal/app/     # Lógica de negócio e casos de uso
├── pkg/config/       # Configurações da aplicação
├── pkg/db/           # Conexões e operações de banco de dados
└── docs/             # Documentação do projeto
```

## Tecnologias

- **Go** - Linguagem de programação principal
- **Gin** - Framework web (a ser adicionado)
- **PostgreSQL** - Banco de dados principal (a ser configurado)
- **Redis** - Cache e sessões (a ser configurado)

## Como Executar

```bash
# Clone o repositório
git clone https://github.com/luizdequeiroz/rpg-backend.git

# Entre no diretório
cd rpg-backend

# Instale as dependências
go mod tidy

# Execute a aplicação
go run cmd/api/main.go
```

## Status do Projeto

🚧 **Em Desenvolvimento** - Este projeto está em fase inicial de desenvolvimento.

## Contribuição

Contribuições são bem-vindas! Por favor, abra uma issue ou envie um pull request.

## Licença

Este projeto está sob a licença MIT.
