#!/bin/bash

# Teste WebSocket Simples - Verifica endpoints disponíveis

echo "=== TESTE WEBSOCKET RPG BACKEND ==="

# Configurações
BASE_URL="http://localhost:8080/api/v1"

echo "1. Testando login..."
curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao.mestre@rpg.com", 
    "password": "senha123"
  }'

echo -e "\n\n2. Testando endpoint WebSocket stats (sem auth - deve dar erro 401)..."
curl -s -X GET "$BASE_URL/ws/stats"

echo -e "\n\n3. WebSocket endpoints disponíveis:"
echo "   - GET  /api/v1/ws?table_id=UUID (WebSocket connection)"
echo "   - GET  /api/v1/ws/stats (Connection statistics)" 
echo "   - POST /api/v1/ws/test (Test event broadcast)"

echo -e "\n\n4. Para testar conexão WebSocket real, use:"
echo "   wscat -c \"ws://localhost:8080/api/v1/ws?table_id=MESA_ID\" -H \"Authorization: Bearer TOKEN\""

echo -e "\n=== TESTE WEBSOCKET CONCLUÍDO ==="
