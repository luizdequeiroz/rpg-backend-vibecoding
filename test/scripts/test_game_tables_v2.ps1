# Script PowerShell para testar o sistema completo de GameTables e Convites
# Uso: .\test\scripts\test_game_tables_v2.ps1

param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose = $false
)

Write-Host "=== TESTANDO SISTEMA DE MESAS E CONVITES ===" -ForegroundColor Cyan
Write-Host "Base URL: $BaseUrl" -ForegroundColor Gray
Write-Host ""

# Função para fazer requisições com tratamento de erro
function Invoke-ApiRequest {
    param(
        [string]$Uri,
        [string]$Method = "GET",
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = 'application/json'
        }
        
        if ($Verbose) {
            Write-Host "→ $Method $Uri" -ForegroundColor DarkGray
        }
        
        return Invoke-RestMethod @params
    }
    catch {
        Write-Host "❌ Erro na requisição: $($_.Exception.Message)" -ForegroundColor Red
        throw
    }
}

try {
    # Verificar se o servidor está rodando
    Write-Host "1. Verificando se o servidor está online..." -ForegroundColor Yellow
    $healthCheck = Invoke-ApiRequest -Uri "$BaseUrl/health" -Method GET
    Write-Host "✅ Servidor está online (Versão: $($healthCheck.version))" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Servidor não está rodando em $BaseUrl" -ForegroundColor Red
    Write-Host "Execute: go run cmd/api/main.go" -ForegroundColor Yellow
    exit 1
}

# Registrar mestre
Write-Host "2. Registrando usuário mestre..." -ForegroundColor Yellow
try {
    $mestreData = Get-Content "test/fixtures/auth_signup_master.json" -Raw
    $mestreResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $mestreData
    $mestreToken = $mestreResponse.token
    Write-Host "✅ Mestre registrado com sucesso (ID: $($mestreResponse.user.id))" -ForegroundColor Green
} catch {
    Write-Host "⚠️ Usuário mestre já existe, fazendo login..." -ForegroundColor Yellow
    $mestreResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $mestreData
    $mestreToken = $mestreResponse.token
    Write-Host "✅ Mestre autenticado com sucesso (ID: $($mestreResponse.user.id))" -ForegroundColor Green
}
Write-Host "Token do mestre: $($mestreToken.Substring(0, 20))..." -ForegroundColor Gray
Write-Host ""

# Registrar jogador
Write-Host "3. Registrando usuário jogador..." -ForegroundColor Yellow
try {
    $jogadorData = Get-Content "test/fixtures/auth_signup_player.json" -Raw
    $jogadorResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $jogadorData
    $jogadorToken = $jogadorResponse.token
    Write-Host "✅ Jogador registrado com sucesso (ID: $($jogadorResponse.user.id))" -ForegroundColor Green
} catch {
    Write-Host "⚠️ Usuário jogador já existe, fazendo login..." -ForegroundColor Yellow
    $jogadorResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $jogadorData
    $jogadorToken = $jogadorResponse.token
    Write-Host "✅ Jogador autenticado com sucesso (ID: $($jogadorResponse.user.id))" -ForegroundColor Green
}
Write-Host "Token do jogador: $($jogadorToken.Substring(0, 20))..." -ForegroundColor Gray
Write-Host ""

# Headers para requisições autenticadas
$mestreHeaders = @{
    'Authorization' = "Bearer $mestreToken"
}
$jogadorHeaders = @{
    'Authorization' = "Bearer $jogadorToken"
}

# Criar mesa de jogo
Write-Host "4. Criando mesa de jogo..." -ForegroundColor Yellow
$mesaData = Get-Content "test/fixtures/game_table_create.json" -Raw
$mesaResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $mestreHeaders -Body $mesaData
$mesaId = $mesaResponse.id
Write-Host "✅ Mesa criada com ID: $mesaId" -ForegroundColor Green
Write-Host "   Nome: $($mesaResponse.name)" -ForegroundColor Gray
Write-Host "   Sistema: $($mesaResponse.system)" -ForegroundColor Gray
Write-Host ""

# Listar mesas do mestre
Write-Host "5. Listando mesas do mestre..." -ForegroundColor Yellow
$mesasMestre = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $mestreHeaders
Write-Host "✅ Mestre tem $($mesasMestre.total) mesa(s)" -ForegroundColor Green
foreach ($mesa in $mesasMestre.tables) {
    Write-Host "   - $($mesa.name) ($($mesa.system)) [Role: $($mesa.role)]" -ForegroundColor Gray
}
Write-Host ""

# Criar convite
Write-Host "6. Criando convite para o jogador..." -ForegroundColor Yellow
$conviteData = Get-Content "test/fixtures/invite_create.json" -Raw
$conviteResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/$mesaId/invites/" -Method POST -Headers $mestreHeaders -Body $conviteData
$conviteId = $conviteResponse.id
Write-Host "✅ Convite criado com ID: $conviteId" -ForegroundColor Green
Write-Host "   Para: $($conviteResponse.invitee.email)" -ForegroundColor Gray
Write-Host "   Status: $($conviteResponse.status)" -ForegroundColor Gray
Write-Host ""

# Listar convites da mesa
Write-Host "7. Listando convites da mesa..." -ForegroundColor Yellow
$convitesMesa = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/$mesaId/invites/" -Headers $mestreHeaders
Write-Host "✅ Mesa tem $($convitesMesa.Length) convite(s)" -ForegroundColor Green
foreach ($convite in $convitesMesa) {
    Write-Host "   - Para: $($convite.invitee.email) | Status: $($convite.status)" -ForegroundColor Gray
}
Write-Host ""

# Verificar mesas do jogador antes de aceitar
Write-Host "8. Verificando mesas do jogador (antes de aceitar)..." -ForegroundColor Yellow
$mesasJogadorAntes = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $jogadorHeaders
Write-Host "✅ Jogador tem $($mesasJogadorAntes.total) mesa(s) antes de aceitar" -ForegroundColor Green
Write-Host ""

# Jogador aceita convite
Write-Host "9. Jogador aceitando convite..." -ForegroundColor Yellow
$aceitarResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/$mesaId/invites/$conviteId/accept" -Method POST -Headers $jogadorHeaders
Write-Host "✅ Convite aceito com sucesso!" -ForegroundColor Green
Write-Host "   Status: $($aceitarResponse.status)" -ForegroundColor Gray
Write-Host ""

# Verificar mesas do jogador após aceitar
Write-Host "10. Verificando mesas do jogador (após aceitar)..." -ForegroundColor Yellow
$mesasJogadorDepois = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $jogadorHeaders
Write-Host "✅ Jogador agora tem $($mesasJogadorDepois.total) mesa(s)" -ForegroundColor Green
foreach ($mesa in $mesasJogadorDepois.tables) {
    Write-Host "   - $($mesa.name) ($($mesa.system)) [Role: $($mesa.role)]" -ForegroundColor Gray
}
Write-Host ""

# Ver detalhes da mesa
Write-Host "11. Vendo detalhes completos da mesa..." -ForegroundColor Yellow
$detalhesMesa = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/$mesaId" -Headers $mestreHeaders
Write-Host "✅ Detalhes da mesa carregados:" -ForegroundColor Green
Write-Host "   Nome: $($detalhesMesa.name)" -ForegroundColor Gray
Write-Host "   Sistema: $($detalhesMesa.system)" -ForegroundColor Gray
Write-Host "   Owner ID: $($detalhesMesa.owner_id)" -ForegroundColor Gray
Write-Host "   Convites: $($detalhesMesa.invites.Length)" -ForegroundColor Gray
foreach ($convite in $detalhesMesa.invites) {
    Write-Host "     - $($convite.invitee.email): $($convite.status)" -ForegroundColor Gray
}
Write-Host ""

# Teste adicional: Criar mesa de Vampiro
Write-Host "12. Teste adicional: Criando mesa de Vampiro..." -ForegroundColor Yellow
$vampiroData = Get-Content "test/fixtures/game_table_vampiro.json" -Raw
$vampiroResponse = Invoke-ApiRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $mestreHeaders -Body $vampiroData
Write-Host "✅ Mesa de Vampiro criada: $($vampiroResponse.name)" -ForegroundColor Green
Write-Host ""

Write-Host "=== TESTE COMPLETO EXECUTADO COM SUCESSO! ===" -ForegroundColor Green
Write-Host ""
Write-Host "📊 RESUMO:" -ForegroundColor Cyan
Write-Host "- Mesa D e D ID: $mesaId" -ForegroundColor White
Write-Host "- Mesa Vampiro ID: $($vampiroResponse.id)" -ForegroundColor White
Write-Host "- Convite ID: $conviteId" -ForegroundColor White
Write-Host "- Status: Convite aceito e jogador tem acesso à mesa" -ForegroundColor White
Write-Host "- Mestre possui $($mesasMestre.total) mesa(s)" -ForegroundColor White
Write-Host "- Jogador possui $($mesasJogadorDepois.total) mesa(s)" -ForegroundColor White
Write-Host ""
Write-Host "🎯 O sistema de mesas e convites está funcionando perfeitamente!" -ForegroundColor Green
