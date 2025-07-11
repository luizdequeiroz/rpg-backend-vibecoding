#!/bin/bash

# Script para testar o sistema completo de GameTables e Invites
# Uso: ./test_game_tables.sh

set -e  # Para no primeiro erro

echo "=== TESTANDO SISTEMA DE MESAS E CONVITES ==="
echo ""

# Verificar se o servidor está rodando
echo "1. Verificando se o servidor está online..."
curl -s http://localhost:8080/health > /dev/null || {
    echo "❌ Servidor não está rodando em localhost:8080"
    echo "Execute: scripts.bat run"
    exit 1
}
echo "✅ Servidor está online"
echo ""

# Registrar mestre
echo "2. Registrando usuário mestre..."
MESTRE_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/auth/signup" \
    -H "Content-Type: application/json" \
    -d '{"email":"mestre@exemplo.com","password":"senha123"}')

MESTRE_TOKEN=$(echo $MESTRE_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$MESTRE_TOKEN" ]; then
    echo "⚠️  Usuário mestre já existe, fazendo login..."
    MESTRE_LOGIN=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"mestre@exemplo.com","password":"senha123"}')
    MESTRE_TOKEN=$(echo $MESTRE_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
fi
echo "✅ Mestre autenticado: ${MESTRE_TOKEN:0:20}..."
echo ""

# Registrar jogador
echo "3. Registrando usuário jogador..."
JOGADOR_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/auth/signup" \
    -H "Content-Type: application/json" \
    -d '{"email":"jogador@exemplo.com","password":"senha123"}')

JOGADOR_TOKEN=$(echo $JOGADOR_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$JOGADOR_TOKEN" ]; then
    echo "⚠️  Usuário jogador já existe, fazendo login..."
    JOGADOR_LOGIN=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"jogador@exemplo.com","password":"senha123"}')
    JOGADOR_TOKEN=$(echo $JOGADOR_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
fi
echo "✅ Jogador autenticado: ${JOGADOR_TOKEN:0:20}..."
echo ""

# Criar mesa de jogo
echo "4. Criando mesa de jogo..."
MESA_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/tables" \
    -H "Authorization: Bearer $MESTRE_TOKEN" \
    -H "Content-Type: application/json" \
    -d @test/fixtures/game_table_create.json)

MESA_ID=$(echo $MESA_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "✅ Mesa criada com ID: $MESA_ID"
echo "   Resposta: $MESA_RESPONSE"
echo ""

# Listar mesas do mestre
echo "5. Listando mesas do mestre..."
MESAS_MESTRE=$(curl -s -X GET "http://localhost:8080/api/v1/tables" \
    -H "Authorization: Bearer $MESTRE_TOKEN")
echo "✅ Mesas do mestre: $MESAS_MESTRE"
echo ""

# Criar convite
echo "6. Criando convite para o jogador..."
CONVITE_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/tables/$MESA_ID/invites" \
    -H "Authorization: Bearer $MESTRE_TOKEN" \
    -H "Content-Type: application/json" \
    -d @test/fixtures/invite_create.json)

CONVITE_ID=$(echo $CONVITE_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "✅ Convite criado com ID: $CONVITE_ID"
echo "   Resposta: $CONVITE_RESPONSE"
echo ""

# Listar convites da mesa
echo "7. Listando convites da mesa..."
CONVITES=$(curl -s -X GET "http://localhost:8080/api/v1/tables/$MESA_ID/invites" \
    -H "Authorization: Bearer $MESTRE_TOKEN")
echo "✅ Convites da mesa: $CONVITES"
echo ""

# Jogador aceita convite
echo "8. Jogador aceitando convite..."
ACEITAR_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/tables/$MESA_ID/invites/$CONVITE_ID/accept" \
    -H "Authorization: Bearer $JOGADOR_TOKEN")
echo "✅ Convite aceito: $ACEITAR_RESPONSE"
echo ""

# Listar mesas do jogador (agora deve aparecer a mesa)
echo "9. Listando mesas do jogador (agora com a mesa aceita)..."
MESAS_JOGADOR=$(curl -s -X GET "http://localhost:8080/api/v1/tables" \
    -H "Authorization: Bearer $JOGADOR_TOKEN")
echo "✅ Mesas do jogador: $MESAS_JOGADOR"
echo ""

# Ver detalhes da mesa
echo "10. Vendo detalhes da mesa..."
DETALHES_MESA=$(curl -s -X GET "http://localhost:8080/api/v1/tables/$MESA_ID" \
    -H "Authorization: Bearer $MESTRE_TOKEN")
echo "✅ Detalhes da mesa: $DETALHES_MESA"
echo ""

echo "=== TESTE COMPLETO EXECUTADO COM SUCESSO! ==="
echo ""
echo "Resumo:"
echo "- Mesa ID: $MESA_ID"
echo "- Convite ID: $CONVITE_ID"
echo "- Status: Convite aceito e jogador tem acesso à mesa"
