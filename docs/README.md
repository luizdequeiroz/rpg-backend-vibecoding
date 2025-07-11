# 📚 Documentação do RPG Backend

Bem-vindo à documentação completa do RPG Backend! Este diretório contém todos os guias e referências necessários para entender, desenvolver, implantar e manter o sistema.

---

## 📋 Índice da Documentação

### 🚀 [README.md](../README.md)
**Ponto de partida do projeto**
- Visão geral do sistema
- Instalação e configuração inicial
- Guia de início rápido
- Funcionalidades implementadas
- Status do projeto

### 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md)
**Arquitetura e Design do Sistema**
- Princípios de Clean Architecture
- Diagramas de arquitetura (Mermaid)
- Responsabilidades das camadas
- Fluxo de dados e dependências
- Padrões de design utilizados

### 💻 [DEVELOPMENT.md](DEVELOPMENT.md)
**Guia de Desenvolvimento**
- Configuração do ambiente de desenvolvimento
- Padrões de codificação
- Estrutura de testes
- Debugging e troubleshooting
- Ferramentas de desenvolvimento

### 🌐 [API.md](API.md)
**Referência Completa da API**
- Documentação de todos os endpoints
- Exemplos de requests e responses
- Códigos de status HTTP
- Schemas de dados
- Guias de autenticação

### 🚀 [DEPLOYMENT.md](DEPLOYMENT.md)
**Guia de Deployment e Produção**
- Configuração para produção
- Docker e containerização
- Configuração de banco de dados
- Nginx e proxy reverso
- Monitoramento e logs
- Segurança e performance

---

## 🗂️ Organização dos Arquivos

```
docs/
├── README.md           # Este arquivo - índice da documentação
├── ARCHITECTURE.md     # Arquitetura e design do sistema
├── DEVELOPMENT.md      # Guia de desenvolvimento
├── API.md             # Referência completa da API
├── DEPLOYMENT.md      # Guia de deployment e produção
└── assets/            # Imagens e diagramas (futuro)
    ├── diagrams/      # Diagramas de arquitetura
    └── screenshots/   # Screenshots da aplicação
```

---

## 🎯 Como Usar Esta Documentação

### Para Novos Desenvolvedores
1. **Comece pelo [README.md](../README.md)** - Entenda o projeto e faça a configuração inicial
2. **Leia [ARCHITECTURE.md](ARCHITECTURE.md)** - Compreenda a estrutura e design do sistema
3. **Siga [DEVELOPMENT.md](DEVELOPMENT.md)** - Configure seu ambiente de desenvolvimento
4. **Consulte [API.md](API.md)** - Para entender os endpoints disponíveis

### Para DevOps/SRE
1. **Revise [ARCHITECTURE.md](ARCHITECTURE.md)** - Entenda a arquitetura do sistema
2. **Implemente com [DEPLOYMENT.md](DEPLOYMENT.md)** - Guia completo para produção
3. **Use [API.md](API.md)** - Para configurar health checks e monitoramento

### Para Frontend Developers
1. **Comece com [README.md](../README.md)** - Entenda o projeto
2. **Foque em [API.md](API.md)** - Referência completa dos endpoints
3. **Configure ambiente com [DEVELOPMENT.md](DEVELOPMENT.md)** - Para testes locais

### Para Product Managers
1. **Leia [README.md](../README.md)** - Visão geral das funcionalidades
2. **Consulte [API.md](API.md)** - Para entender capacidades da API
3. **Revise [ARCHITECTURE.md](ARCHITECTURE.md)** - Para entender limitações e possibilidades

---

## 📊 Status da Documentação

| Documento | Status | Última Atualização | Completude |
|-----------|--------|-------------------|------------|
| README.md | ✅ Completo | 2025-01-11 | 100% |
| ARCHITECTURE.md | ✅ Completo | 2025-01-11 | 100% |
| DEVELOPMENT.md | ✅ Completo | 2025-01-11 | 100% |
| API.md | ✅ Completo | 2025-01-11 | 100% |
| DEPLOYMENT.md | ✅ Completo | 2025-01-11 | 100% |

---

## 🔄 Atualizações da Documentação

### Princípios de Manutenção
- **Mantenha atualizado**: Toda mudança no código deve refletir na documentação
- **Seja claro**: Use linguagem simples e exemplos práticos
- **Inclua exemplos**: Sempre forneça exemplos funcionais
- **Teste instruções**: Verifique se os comandos funcionam antes de documentar

### Como Contribuir
1. **Para correções menores**: Edite diretamente os arquivos
2. **Para mudanças significativas**: Abra uma issue primeiro
3. **Para novos recursos**: Atualize múltiplos documentos conforme necessário
4. **Sempre teste**: Verifique se exemplos e comandos funcionam

### Responsabilidades
- **Desenvolvedores**: Atualizar DEVELOPMENT.md e ARCHITECTURE.md
- **DevOps**: Manter DEPLOYMENT.md atualizado
- **API Team**: Manter API.md sincronizado com código
- **Product**: Revisar README.md para clareza

---

## 🛠️ Ferramentas e Recursos

### Diagramas
- **Mermaid**: Para diagramas de arquitetura e fluxo
- **PlantUML**: Para diagramas UML detalhados (futuro)
- **Draw.io**: Para diagramas complexos (futuro)

### Documentação API
- **Swagger UI**: Documentação interativa em `/docs/index.html`
- **OpenAPI 3.0**: Especificação padrão para APIs REST
- **Postman Collection**: Para testes e exemplos (futuro)

### Validação
- **Markdown Linting**: Para consistência de formato
- **Link Checking**: Para verificar links quebrados
- **Code Examples**: Testes automatizados dos exemplos

---

## 📞 Suporte e Contato

### Para Dúvidas Técnicas
- **Issues**: Use GitHub Issues para bugs e questões
- **Discussions**: Para discussões gerais e ideias
- **Code Review**: Para mudanças na documentação

### Para Contribuições
1. Fork o repositório
2. Crie uma branch para sua mudança
3. Faça as alterações na documentação
4. Teste todos os exemplos
5. Submeta um Pull Request

### Contatos da Equipe
- **Tech Lead**: Para arquitetura e design
- **DevOps**: Para deployment e infraestrutura  
- **Product**: Para funcionalidades e roadmap

---

## 🏷️ Tags e Versioning

### Tags de Documentação
- `#architecture` - Tópicos relacionados à arquitetura
- `#api` - Documentação de endpoints
- `#deployment` - Guias de deploy
- `#development` - Configuração de ambiente
- `#security` - Tópicos de segurança
- `#performance` - Otimizações e performance

### Versionamento
- **Major**: Mudanças na arquitetura principal
- **Minor**: Novos endpoints ou funcionalidades
- **Patch**: Correções e melhorias na documentação

---

## 🔮 Roadmap da Documentação

### Próximas Adições
- [ ] **Cookbook**: Receitas para tarefas comuns
- [ ] **Troubleshooting**: Guia de resolução de problemas
- [ ] **Security Guide**: Guia detalhado de segurança
- [ ] **Performance Tuning**: Otimizações avançadas
- [ ] **Testing Guide**: Estratégias de teste
- [ ] **Migration Guide**: Guias de migração entre versões

### Melhorias Planejadas
- [ ] **Video Tutorials**: Tutoriais em vídeo
- [ ] **Interactive Examples**: Exemplos executáveis
- [ ] **Multilingual**: Suporte a outros idiomas
- [ ] **Search**: Sistema de busca na documentação
- [ ] **Feedback System**: Sistema de feedback dos usuários

---

**Esta documentação é um recurso vivo e deve evoluir junto com o projeto. Contribuições são sempre bem-vindas!**

---

## 🌟 Quick Links

- 🏠 [Voltar ao README principal](../README.md)
- 🏗️ [Ver Arquitetura](ARCHITECTURE.md)
- 💻 [Guia de Desenvolvimento](DEVELOPMENT.md)
- 🌐 [Referência da API](API.md)
- 🚀 [Guia de Deployment](DEPLOYMENT.md)
- 📊 [Swagger UI](http://localhost:8080/docs/index.html) *(quando servidor estiver rodando)*



```powershell
# Registro
$signup = @{email='usuario@exemplo.com'; password='senha123'} | ConvertTo-Json
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/signup' -Method POST -Body $signup -ContentType 'application/json'

# Login
$login = @{email='usuario@exemplo.com'; password='senha123'} | ConvertTo-Json
$response = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -Body $login -ContentType 'application/json'
$token = $response.token

# Endpoint protegido
$headers = @{Authorization="Bearer $token"}
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/me' -Headers $headers
```

## 📚 Documentação da API

A documentação completa da API está disponível via Swagger:

- **URL**: `http://localhost:8080/docs/index.html`
- **JSON**: `http://localhost:8080/docs/swagger.json`
- **YAML**: `http://localhost:8080/docs/swagger.yaml`

### Regenerar documentação Swagger
```bash
# Via Make
make swagger-generate

# Via script Windows
scripts.bat swagger-generate

# Via swag direto
swag init -g cmd/api/main.go -o docs
```

## 🗃️ Banco de Dados

### SQLite (padrão)
O projeto usa SQLite por padrão, armazenado em `./data/rpg.db`.

### PostgreSQL (opcional)
Para usar PostgreSQL, configure a variável `DATABASE_URL`:
```bash
export DATABASE_URL="postgres://user:password@localhost/rpg_db?sslmode=disable"
```

### Migrações
```bash
# Executar migrações
make migrate-up

# Desfazer última migração
make migrate-down

# Status das migrações
make migrate-status

# Criar nova migração
make migrate-create

# Resetar todas as migrações
make migrate-reset
```

## 🔧 Scripts Úteis

### Makefile (Linux/Mac)
```bash
make run              # Executar aplicação
make dev              # Executar em modo debug
make build            # Compilar aplicação
make test             # Executar testes
make swagger-generate # Gerar documentação Swagger
make migrate-up       # Executar migrações
make migrate-down     # Desfazer migração
make clean            # Limpar arquivos temporários
```

### Scripts Windows (scripts.bat)
```batch
scripts.bat run              # Executar aplicação
scripts.bat dev              # Executar em modo debug
scripts.bat build            # Compilar aplicação
scripts.bat test             # Executar testes
scripts.bat swagger-generate # Gerar documentação Swagger
scripts.bat migrate-up       # Executar migrações
scripts.bat migrate-down     # Desfazer migração
```

## 🔐 Autenticação

O sistema usa JWT (JSON Web Tokens) para autenticação:

1. **Registro/Login**: Receba um token JWT
2. **Requisições**: Inclua o header `Authorization: Bearer <token>`
3. **Expiração**: Tokens expiram em 24 horas

### Middleware de Autenticação

- **AuthMiddleware**: Obrigatório - bloqueia acesso sem token
- **OptionalAuthMiddleware**: Opcional - permite acesso com ou sem token

## 📋 Modelos de Dados

### User
```json
{
  "id": 1,
  "email": "usuario@exemplo.com",
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:26:29Z"
}
```

### SheetTemplate
```json
{
  "id": 1,
  "name": "Ficha D&D 5e",
  "system": "D&D 5e",
  "description": "Template para personagens de D&D 5ª edição",
  "definition": {
    "sections": [
      {
        "name": "Atributos",
        "fields": [
          {
            "name": "Força",
            "type": "number",
            "min": 1,
            "max": 20
          }
        ]
      }
    ]
  },
  "is_active": true,
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:26:29Z"
}
```

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📞 Contato

**Luiz de Queiroz**
- Email: luiz@example.com
- GitHub: [@luizdequeiroz](https://github.com/luizdequeiroz)

---

**Feito com ❤️ em Go**
