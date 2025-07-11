#!/bin/bash

# Script para testar fluxo completo da API
echo "=== TESTE DE FLUXO COMPLETO DA API RPG BACKEND ==="
echo ""

BASE_URL="http://localhost:8080"

# Função para extrair token do JSON
extract_token() {
    echo "$1" | grep -o '"token":"[^"]*"' | cut -d'"' -f4
}

# Função para extrair ID do JSON  
extract_id() {
    echo "$1" | grep -o '"id":"[^"]*"' | cut -d'"' -f4
}

# Função para extrair ID numérico do JSON
extract_numeric_id() {
    echo "$1" | grep -o '"id":[0-9]*' | cut -d':' -f2
}

echo "1. Registrando Mestre João..."
JOAO_SIGNUP=$(curl -s -X POST $BASE_URL/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"joao.mestre@rpg.com","password":"senha123"}')
echo "Resposta signup João: $JOAO_SIGNUP"
echo ""

echo "2. Login do Mestre João..."
JOAO_LOGIN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"joao.mestre@rpg.com","password":"senha123"}')
echo "Resposta login João: $JOAO_LOGIN"
JOAO_TOKEN=$(extract_token "$JOAO_LOGIN")
echo "Token João: $JOAO_TOKEN"
echo ""

echo "3. Registrando Player Maria..."
MARIA_SIGNUP=$(curl -s -X POST $BASE_URL/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"maria.player@rpg.com","password":"senha123"}')
echo "Resposta signup Maria: $MARIA_SIGNUP"
echo ""

echo "4. Login da Player Maria..."
MARIA_LOGIN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"maria.player@rpg.com","password":"senha123"}')
echo "Resposta login Maria: $MARIA_LOGIN"
MARIA_TOKEN=$(extract_token "$MARIA_LOGIN")
echo "Token Maria: $MARIA_TOKEN"
echo ""

echo "5. João criando template de ficha D&D 5e..."
TEMPLATE_CREATE=$(curl -s -X POST $BASE_URL/api/v1/templates \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JOAO_TOKEN" \
  -d '{
    "name": "Ficha D&D 5e",
    "description": "Template para personagens de Dungeons & Dragons 5ª edição",
    "definition": {
      "sections": [
        {
          "name": "Atributos",
          "fields": [
            {"name": "Força", "type": "number", "min": 1, "max": 20},
            {"name": "Destreza", "type": "number", "min": 1, "max": 20},
            {"name": "Constituição", "type": "number", "min": 1, "max": 20},
            {"name": "Inteligência", "type": "number", "min": 1, "max": 20},
            {"name": "Sabedoria", "type": "number", "min": 1, "max": 20},
            {"name": "Carisma", "type": "number", "min": 1, "max": 20}
          ]
        }
      ]
    }
  }')
echo "Resposta criação template: $TEMPLATE_CREATE"
TEMPLATE_ID=$(extract_numeric_id "$TEMPLATE_CREATE")
echo "Template ID: $TEMPLATE_ID"
echo ""

echo "6. João criando mesa de jogo..."
TABLE_CREATE=$(curl -s -X POST $BASE_URL/api/v1/tables/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JOAO_TOKEN" \
  -d '{
    "name": "A Busca pelo Artefato Perdido",
    "system": "D&D 5e"
  }')
echo "Resposta criação mesa: $TABLE_CREATE"
TABLE_ID=$(extract_id "$TABLE_CREATE")
echo "Mesa ID: $TABLE_ID"
echo ""

echo "7. João convidando Maria para a mesa..."
INVITE_CREATE=$(curl -s -X POST $BASE_URL/api/v1/tables/$TABLE_ID/invites/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JOAO_TOKEN" \
  -d '{
    "invitee_email": "maria.player@rpg.com"
  }')
echo "Resposta convite: $INVITE_CREATE"
INVITE_ID=$(extract_id "$INVITE_CREATE")
echo "Convite ID: $INVITE_ID"
echo ""

echo "8. Maria aceitando convite..."
INVITE_ACCEPT=$(curl -s -X POST $BASE_URL/api/v1/tables/$TABLE_ID/invites/$INVITE_ID/accept \
  -H "Authorization: Bearer $MARIA_TOKEN")
echo "Resposta aceitar convite: $INVITE_ACCEPT"
echo ""

echo "9. Maria criando ficha de personagem..."
SHEET_CREATE=$(curl -s -X POST $BASE_URL/api/v1/sheets/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MARIA_TOKEN" \
  -d '{
    "table_id": "'$TABLE_ID'",
    "template_id": '$TEMPLATE_ID',
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
      }
    }
  }')
echo "Resposta criação ficha: $SHEET_CREATE"
SHEET_ID=$(extract_id "$SHEET_CREATE")
echo "Ficha ID: $SHEET_ID"
echo ""

echo "10. Maria rolando dados (teste de Arcana)..."
ROLL_1=$(curl -s -X POST $BASE_URL/api/v1/rolls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MARIA_TOKEN" \
  -d '{
    "sheet_id": "'$SHEET_ID'",
    "expression": "1d20+8"
  }')
echo "Resposta rolagem 1: $ROLL_1"
echo ""

echo "11. Maria rolando dados (baseado em campo da ficha)..."
ROLL_2=$(curl -s -X POST $BASE_URL/api/v1/rolls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MARIA_TOKEN" \
  -d '{
    "sheet_id": "'$SHEET_ID'",
    "field_name": "skills.arcana"
  }')
echo "Resposta rolagem 2: $ROLL_2"
echo ""

echo "12. João visualizando histórico de rolagens da mesa..."
ROLLS_HISTORY=$(curl -s -X GET $BASE_URL/api/v1/rolls/table/$TABLE_ID \
  -H "Authorization: Bearer $JOAO_TOKEN")
echo "Histórico de rolagens: $ROLLS_HISTORY"
echo ""

echo "13. Listando fichas da mesa..."
SHEETS_LIST=$(curl -s -X GET "$BASE_URL/api/v1/sheets/?table_id=$TABLE_ID" \
  -H "Authorization: Bearer $JOAO_TOKEN")
echo "Lista de fichas: $SHEETS_LIST"
echo ""

echo "=== TESTE CONCLUÍDO ==="
