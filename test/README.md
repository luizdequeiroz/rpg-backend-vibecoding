# Documentação de Testes

Esta pasta contém todos os arquivos e scripts necessários para testar o sistema RPG Backend.

## Estrutura de Pastas

```
test/
├── fixtures/           # Arquivos JSON com dados de teste
└── scripts/           # Scripts de automação de testes
```

## Fixtures (Dados de Teste)

### Autenticação
- `auth_signup_master.json` - Dados para registro do mestre
- `auth_signup_player.json` - Dados para registro do jogador
- `auth_signup_player2.json` - Dados para registro de segundo jogador
- `auth_signup_admin.json` - Dados para registro de administrador

### Mesas de Jogo
- `game_table_create.json` - Dados para criar mesa D&D
- `game_table_vampiro.json` - Dados para criar mesa de Vampiro
- `game_table_cthulhu.json` - Dados para criar mesa de Call of Cthulhu
- `game_table_update.json` - Dados para atualizar mesa existente

### Convites
- `invite_create.json` - Dados para criar convite para jogador principal
- `invite_create_player2.json` - Dados para criar convite para segundo jogador

## Scripts de Teste

### PowerShell (Recomendado para Windows)
```bash
# Executar teste completo
./test/scripts/test_game_tables.ps1
```

### Bash (Linux/macOS)
```bash
# Executar teste completo
./test/scripts/test_game_tables.sh
```

### Batch (Windows CMD)
```bash
# Executar teste básico
./test/scripts/test_game_tables.bat
```

## Como Usar

1. **Iniciar o Servidor**
   ```bash
   go run cmd/api/main.go
   ```

2. **Executar Testes**
   ```bash
   # PowerShell (Recomendado)
   cd test/scripts
   ./test_game_tables.ps1
   
   # ou usar scripts.bat na raiz
   scripts.bat test-api
   ```

3. **Teste Manual com curl**
   ```bash
   # Registrar usuário
   curl -X POST "http://localhost:8080/api/v1/auth/signup" \
        -H "Content-Type: application/json" \
        -d @test/fixtures/auth_signup_master.json
   
   # Criar mesa
   curl -X POST "http://localhost:8080/api/v1/tables/" \
        -H "Authorization: Bearer {TOKEN}" \
        -H "Content-Type: application/json" \
        -d @test/fixtures/game_table_create.json
   ```

## Cenários de Teste Disponíveis

### Fluxo Completo (test_game_tables.ps1)
- ✅ Verificação de servidor online
- ✅ Registro de mestre e jogador
- ✅ Autenticação JWT
- ✅ Criação de mesa de jogo
- ✅ Sistema de convites
- ✅ Aceitação de convites
- ✅ Verificação de permissões
- ✅ Listagem de mesas por role

### Diferentes Sistemas de RPG
- **D&D 5e** - Mesa clássica de fantasia medieval
- **Vampiro: A Máscara** - Mesa de horror urbano
- **Call of Cthulhu** - Mesa de horror cósmico

### Múltiplos Usuários
- **Mestre** - Criador e administrador de mesas
- **Jogador Principal** - Primeiro convidado
- **Jogador Secundário** - Segundo convidado
- **Administrador** - Usuário com privilégios especiais

## Códigos de Resposta Esperados

- `200` - Operação bem-sucedida
- `201` - Recurso criado com sucesso
- `401` - Token inválido ou ausente
- `403` - Sem permissão para operação
- `404` - Recurso não encontrado
- `409` - Conflito (ex: convite duplicado)

## Troubleshooting

### Servidor não está rodando
```bash
# Verificar se porta 8080 está livre
netstat -ano | findstr :8080

# Iniciar servidor
go run cmd/api/main.go
```

### Problemas com PowerShell
```powershell
# Definir política de execução
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope CurrentUser
```

### Problemas de Token
- Tokens JWT expiram em 24 horas
- Regenerar token fazendo novo login
- Verificar se Bearer está incluído no header

## Logs e Debugging

O servidor exibe logs detalhados no console:
```
[GIN] 2025/07/11 - 12:05:55 | 201 | 55.9036ms | 127.0.0.1 | POST /api/v1/auth/signup
[GIN] 2025/07/11 - 12:06:16 | 201 |  3.141ms | 127.0.0.1 | POST /api/v1/tables/
```

Para debug detalhado, usar `curl -v` ou verificar logs do servidor.
