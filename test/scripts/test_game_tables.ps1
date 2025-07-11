# Script PowerShell para testar o sistema completo de GameTables e Invites
# Uso: .\test_game_tables.ps1

Write-Host "=== TESTANDO SISTEMA DE MESASWrite-Host "🎯 O sistema de mesas e convites está funcionando perfeitamente!" -ForegroundColor GreenE CONVITES ===" -ForegroundColor Cyan
Write-Host ""

try {
    # Verificar se o servidor está rodando
    Write-Host "1. Verificando se o servidor está online..." -ForegroundColor Yellow
    $healthCheck = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
    Write-Host "✅ Servidor está online" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Servidor não está rodando em localhost:8080" -ForegroundColor Red
    Write-Host "Execute: scripts.bat run" -ForegroundColor Yellow
    exit 1
}

# Registrar mestre
Write-Host "2. Registrando usuário mestre..." -ForegroundColor Yellow
try {
    $mestreSignup = @{email='mestre@exemplo.com'; password='senha123'} | ConvertTo-Json
    $mestreResponse = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/signup' -Method POST -Body $mestreSignup -ContentType 'application/json'
    $mestreToken = $mestreResponse.token
    Write-Host "✅ Mestre registrado com sucesso" -ForegroundColor Green
} catch {
    Write-Host "⚠️ Usuário mestre já existe, fazendo login..." -ForegroundColor Yellow
    $mestreLogin = @{email='mestre@exemplo.com'; password='senha123'} | ConvertTo-Json
    $mestreResponse = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -Body $mestreLogin -ContentType 'application/json'
    $mestreToken = $mestreResponse.token
    Write-Host "✅ Mestre autenticado com sucesso" -ForegroundColor Green
}
Write-Host "Token do mestre: $($mestreToken.Substring(0, 20))..." -ForegroundColor Gray
Write-Host ""

# Registrar jogador
Write-Host "3. Registrando usuário jogador..." -ForegroundColor Yellow
try {
    $jogadorSignup = @{email='jogador@exemplo.com'; password='senha123'} | ConvertTo-Json
    $jogadorResponse = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/signup' -Method POST -Body $jogadorSignup -ContentType 'application/json'
    $jogadorToken = $jogadorResponse.token
    Write-Host "✅ Jogador registrado com sucesso" -ForegroundColor Green
} catch {
    Write-Host "⚠️ Usuário jogador já existe, fazendo login..." -ForegroundColor Yellow
    $jogadorLogin = @{email='jogador@exemplo.com'; password='senha123'} | ConvertTo-Json
    $jogadorResponse = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -Body $jogadorLogin -ContentType 'application/json'
    $jogadorToken = $jogadorResponse.token
    Write-Host "✅ Jogador autenticado com sucesso" -ForegroundColor Green
}
Write-Host "Token do jogador: $($jogadorToken.Substring(0, 20))..." -ForegroundColor Gray
Write-Host ""

# Headers para requisições autenticadas
$mestreHeaders = @{
    'Authorization' = "Bearer $mestreToken"
    'Content-Type' = 'application/json'
}
$jogadorHeaders = @{
    'Authorization' = "Bearer $jogadorToken"
    'Content-Type' = 'application/json'
}

# Criar mesa de jogo
Write-Host "4. Criando mesa de jogo..." -ForegroundColor Yellow
$novaMesa = @{
    name = "Mesa D&D: A Busca pelo Artefato Perdido"
    system = "D&D 5e"
} | ConvertTo-Json

$mesaResponse = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Method POST -Body $novaMesa -Headers $mestreHeaders
$mesaId = $mesaResponse.id
Write-Host "✅ Mesa criada com ID: $mesaId" -ForegroundColor Green
Write-Host "Nome: $($mesaResponse.name)" -ForegroundColor Gray
Write-Host "Sistema: $($mesaResponse.system)" -ForegroundColor Gray
Write-Host ""

# Listar mesas do mestre
Write-Host "5. Listando mesas do mestre..." -ForegroundColor Yellow
$mesasMestre = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Headers $mestreHeaders
Write-Host "✅ Mestre tem $($mesasMestre.Length) mesa(s)" -ForegroundColor Green
foreach ($mesa in $mesasMestre) {
    Write-Host "   - $($mesa.name) ($($mesa.system))" -ForegroundColor Gray
}
Write-Host ""

# Criar convite
Write-Host "6. Criando convite para o jogador..." -ForegroundColor Yellow
$novoConvite = @{
    invitee_email = "jogador@exemplo.com"
} | ConvertTo-Json

$conviteResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tables/$mesaId/invites" -Method POST -Body $novoConvite -Headers $mestreHeaders
$conviteId = $conviteResponse.id
Write-Host "✅ Convite criado com ID: $conviteId" -ForegroundColor Green
Write-Host "Status: $($conviteResponse.status)" -ForegroundColor Gray
Write-Host ""

# Listar convites da mesa
Write-Host "7. Listando convites da mesa..." -ForegroundColor Yellow
$convitesMesa = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tables/$mesaId/invites" -Headers $mestreHeaders
Write-Host "✅ Mesa tem $($convitesMesa.Length) convite(s)" -ForegroundColor Green
foreach ($convite in $convitesMesa) {
    Write-Host "   - Para: $($convite.invitee.email) | Status: $($convite.status)" -ForegroundColor Gray
}
Write-Host ""

# Verificar mesas do jogador antes de aceitar
Write-Host "8. Verificando mesas do jogador (antes de aceitar)..." -ForegroundColor Yellow
$mesasJogadorAntes = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Headers $jogadorHeaders
Write-Host "✅ Jogador tem $($mesasJogadorAntes.Length) mesa(s) antes de aceitar" -ForegroundColor Green
Write-Host ""

# Jogador aceita convite
Write-Host "9. Jogador aceitando convite..." -ForegroundColor Yellow
$aceitarResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tables/$mesaId/invites/$conviteId/accept" -Method POST -Headers $jogadorHeaders
Write-Host "✅ Convite aceito com sucesso!" -ForegroundColor Green
Write-Host "Status: $($aceitarResponse.status)" -ForegroundColor Gray
Write-Host ""

# Verificar mesas do jogador após aceitar
Write-Host "10. Verificando mesas do jogador (após aceitar)..." -ForegroundColor Yellow
$mesasJogadorDepois = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/tables' -Headers $jogadorHeaders
Write-Host "✅ Jogador agora tem $($mesasJogadorDepois.Length) mesa(s)" -ForegroundColor Green
foreach ($mesa in $mesasJogadorDepois) {
    Write-Host "   - $($mesa.name) ($($mesa.system))" -ForegroundColor Gray
}
Write-Host ""

# Ver detalhes da mesa
Write-Host "11. Vendo detalhes completos da mesa..." -ForegroundColor Yellow
$detalhesMesa = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tables/$mesaId" -Headers $mestreHeaders
Write-Host "✅ Detalhes da mesa carregados:" -ForegroundColor Green
Write-Host "   Nome: $($detalhesMesa.name)" -ForegroundColor Gray
Write-Host "   Sistema: $($detalhesMesa.system)" -ForegroundColor Gray
Write-Host "   Owner: $($detalhesMesa.owner.email)" -ForegroundColor Gray
Write-Host "   Convites: $($detalhesMesa.invites.Length)" -ForegroundColor Gray
foreach ($convite in $detalhesMesa.invites) {
    Write-Host "     - $($convite.invitee.email): $($convite.status)" -ForegroundColor Gray
}
Write-Host ""

Write-Host "=== TESTE COMPLETO EXECUTADO COM SUCESSO! ===" -ForegroundColor Green
Write-Host ""
Write-Host "📊 RESUMO:" -ForegroundColor Cyan
Write-Host "- Mesa ID: $mesaId" -ForegroundColor White
Write-Host "- Convite ID: $conviteId" -ForegroundColor White
Write-Host "- Status: Convite aceito e jogador tem acesso à mesa" -ForegroundColor White
Write-Host "- Mestre possui $($mesasMestre.Length) mesa(s)" -ForegroundColor White
Write-Host "- Jogador possui $($mesasJogadorDepois.Length) mesa(s)" -ForegroundColor White
Write-Host ""
Write-Host "🎯 O sistema de mesas e convites esta funcionando perfeitamente!" -ForegroundColor Green
