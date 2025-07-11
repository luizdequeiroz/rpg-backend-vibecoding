#!/bin/bash

# Teste WebSocket - Verifica conectividade e notificações em tempo real

echo "=== TESTE WEBSOCKET RPG BACKEND ==="

# Configurações
BASE_URL="http://localhost:8080/api/v1"
WS_URL="ws://localhost:8080/api/v1/ws"

# Fazer login para obter token
echo "1. Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao.mestre@rpg.com",
    "password": "senha123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Erro: Não foi possível obter token"
  exit 1
fi

# Buscar uma mesa existente
echo "2. Buscando mesas..."
TABLES_RESPONSE=$(curl -s -X GET "$BASE_URL/tables/my" \
  -H "Authorization: Bearer $TOKEN")

TABLE_ID=$(echo $TABLES_RESPONSE | jq -r '.tables[0].id // empty')
echo "Mesa ID: $TABLE_ID"

if [ -z "$TABLE_ID" ]; then
  echo "❌ Erro: Nenhuma mesa encontrada"
  exit 1
fi

# Testar estatísticas WebSocket
echo "3. Verificando estatísticas WebSocket..."
WS_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/ws/stats" \
  -H "Authorization: Bearer $TOKEN")

echo "Estatísticas WebSocket: $WS_STATS_RESPONSE"

# Testar evento de teste
echo "4. Enviando evento de teste..."
TEST_EVENT_RESPONSE=$(curl -s -X POST "$BASE_URL/ws/test" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_id\": \"$TABLE_ID\",
    \"event_type\": \"roll_performed\",
    \"data\": {
      \"message\": \"Teste de evento WebSocket\",
      \"timestamp\": \"$(date -Iseconds)\"
    }
  }")

echo "Resposta evento teste: $TEST_EVENT_RESPONSE"

# Verificar estatísticas novamente
echo "5. Verificando estatísticas após teste..."
WS_STATS_RESPONSE2=$(curl -s -X GET "$BASE_URL/ws/stats" \
  -H "Authorization: Bearer $TOKEN")

echo "Estatísticas finais: $WS_STATS_RESPONSE2"

echo "=== TESTE WEBSOCKET CONCLUÍDO ==="
echo "💡 Para testar conexão WebSocket real, use um cliente WebSocket:"
echo "   URL: $WS_URL?table_id=$TABLE_ID"
echo "   Header: Authorization: Bearer $TOKEN"
