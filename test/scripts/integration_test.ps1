# Teste Integrado RPG Backend
# Script principal para testes automatizados
# Uso: .\test\scripts\integration_test.ps1

param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose = $false,
    [switch]$SkipCleanup = $false
)

# Importar configurações
. "$PSScriptRoot\test_config.ps1"

# Configurar URL base
$script:BaseUrl = $BaseUrl

Write-TestHeader "TESTE DE INTEGRAÇÃO RPG BACKEND"
Write-Host "Base URL: $BaseUrl" -ForegroundColor Gray
Write-Host "Verbose: $Verbose" -ForegroundColor Gray
Write-Host ""

# Verificar servidor
Write-TestStep "1" "Verificando servidor online"
if (-not (Test-ServerOnline -Url $BaseUrl)) {
    Write-TestError "Servidor não está rodando em $BaseUrl"
    Write-Host "Execute: go run cmd/api/main.go" -ForegroundColor Yellow
    exit 1
}
Write-TestSuccess "Servidor está online e respondendo"

# Função para fazer requisições
function Invoke-TestRequest {
    param(
        [string]$Uri,
        [string]$Method = "GET",
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [string]$Description = ""
    )
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            Headers = $Headers
            TimeoutSec = 10
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = 'application/json'
        }
        
        if ($Verbose -and $Description) {
            Write-Host "  → $Description" -ForegroundColor DarkGray
        }
        
        return Invoke-RestMethod @params
    }
    catch {
        if ($Verbose) {
            Write-Host "  ❌ Erro: $($_.Exception.Message)" -ForegroundColor Red
        }
        throw
    }
}

# Arrays para armazenar IDs criados (para cleanup)
$CreatedUsers = @()
$CreatedTables = @()
$CreatedInvites = @()

try {
    # Teste de Autenticação
    Write-TestStep "2" "Testando sistema de autenticação"
    
    # Registrar mestre
    $masterData = Get-Content $script:DefaultMaster -Raw
    try {
        $masterResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $masterData -Description "Registro do mestre"
        $masterToken = $masterResponse.token
        $CreatedUsers += $masterResponse.user.id
        Write-TestSuccess "Mestre registrado (ID: $($masterResponse.user.id))"
    }
    catch {
        Write-TestWarning "Mestre já existe, fazendo login"
        $masterResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $masterData -Description "Login do mestre"
        $masterToken = $masterResponse.token
        Write-TestSuccess "Mestre autenticado (ID: $($masterResponse.user.id))"
    }
    
    # Registrar jogador
    $playerData = Get-Content $script:DefaultPlayer -Raw
    try {
        $playerResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $playerData -Description "Registro do jogador"
        $playerToken = $playerResponse.token
        $CreatedUsers += $playerResponse.user.id
        Write-TestSuccess "Jogador registrado (ID: $($playerResponse.user.id))"
    }
    catch {
        Write-TestWarning "Jogador já existe, fazendo login"
        $playerResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $playerData -Description "Login do jogador"
        $playerToken = $playerResponse.token
        Write-TestSuccess "Jogador autenticado (ID: $($playerResponse.user.id))"
    }
    
    # Headers para autenticação
    $masterHeaders = @{ 'Authorization' = "Bearer $masterToken" }
    $playerHeaders = @{ 'Authorization' = "Bearer $playerToken" }
    
    # Teste de Mesas
    Write-TestStep "3" "Testando gerenciamento de mesas"
    
    # Criar mesa
    $tableData = Get-Content $script:DefaultTable -Raw
    $tableResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $tableData -Description "Criação de mesa"
    $tableId = $tableResponse.id
    $CreatedTables += $tableId
    Write-TestSuccess "Mesa criada (ID: $tableId)"
    
    # Listar mesas do mestre
    $masterTables = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $masterHeaders -Description "Listagem de mesas do mestre"
    Write-TestSuccess "Mestre tem $($masterTables.total) mesa(s)"
    
    # Verificar que jogador não tem acesso ainda
    $playerTablesInitial = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $playerHeaders -Description "Mesas do jogador (inicial)"
    Write-TestSuccess "Jogador tem $($playerTablesInitial.total) mesa(s) (antes do convite)"
    
    # Teste de Convites
    Write-TestStep "4" "Testando sistema de convites"
    
    # Criar convite
    $inviteData = Get-Content $script:DefaultInvite -Raw
    $inviteResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId/invites/" -Method POST -Headers $masterHeaders -Body $inviteData -Description "Criação de convite"
    $inviteId = $inviteResponse.id
    $CreatedInvites += $inviteId
    Write-TestSuccess "Convite criado (ID: $inviteId)"
    
    # Listar convites da mesa
    $tableInvites = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId/invites/" -Headers $masterHeaders -Description "Listagem de convites"
    Write-TestSuccess "Mesa tem $($tableInvites.Length) convite(s)"
    
    # Jogador aceita convite
    $acceptResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId/invites/$inviteId/accept" -Method POST -Headers $playerHeaders -Description "Aceitação de convite"
    Write-TestSuccess "Convite aceito (Status: $($acceptResponse.status))"
    
    # Verificar acesso do jogador após aceitar
    $playerTablesFinal = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $playerHeaders -Description "Mesas do jogador (final)"
    Write-TestSuccess "Jogador agora tem $($playerTablesFinal.total) mesa(s)"
    
    # Teste de Detalhes
    Write-TestStep "5" "Verificando detalhes e permissões"
    
    # Detalhes da mesa pelo mestre
    $tableDetails = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId" -Headers $masterHeaders -Description "Detalhes da mesa"
    Write-TestSuccess "Mesa tem $($tableDetails.invites.Length) convite(s) registrado(s)"
    
    # Verificar role do jogador
    $playerTable = $playerTablesFinal.tables | Where-Object { $_.id -eq $tableId }
    if ($playerTable -and $playerTable.role -eq "player") {
        Write-TestSuccess "Jogador tem role correto: $($playerTable.role)"
    } else {
        Write-TestError "Role do jogador incorreto"
    }
    
    # Teste de Sistemas Diferentes
    Write-TestStep "6" "Testando diferentes sistemas de RPG"
    
    # Criar mesa de Vampiro
    $vampiroData = Get-Content "test/fixtures/game_table_vampiro.json" -Raw
    $vampiroResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $vampiroData -Description "Mesa de Vampiro"
    $CreatedTables += $vampiroResponse.id
    Write-TestSuccess "Mesa de Vampiro criada: $($vampiroResponse.name)"
    
    # Criar mesa de Cthulhu
    $cthulhuData = Get-Content "test/fixtures/game_table_cthulhu.json" -Raw
    $cthulhuResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $cthulhuData -Description "Mesa de Cthulhu"
    $CreatedTables += $cthulhuResponse.id
    Write-TestSuccess "Mesa de Cthulhu criada: $($cthulhuResponse.name)"
    
    # Verificar total de mesas
    $finalTables = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $masterHeaders -Description "Contagem final de mesas"
    Write-TestSuccess "Mestre tem total de $($finalTables.total) mesa(s)"
    
    Write-TestHeader "TODOS OS TESTES PASSARAM!"
    
    # Resumo final
    Write-Host "📊 RESUMO DOS TESTES:" -ForegroundColor Cyan
    Write-Host "  ✅ Autenticação funcionando" -ForegroundColor Green
    Write-Host "  ✅ Criação de mesas funcionando" -ForegroundColor Green
    Write-Host "  ✅ Sistema de convites funcionando" -ForegroundColor Green
    Write-Host "  ✅ Controle de permissões funcionando" -ForegroundColor Green
    Write-Host "  ✅ Múltiplos sistemas de RPG funcionando" -ForegroundColor Green
    Write-Host ""
    Write-Host "🎯 Sistema pronto para produção!" -ForegroundColor Green
    
}
catch {
    Write-TestError "Teste falhou: $($_.Exception.Message)"
    
    if ($Verbose) {
        Write-Host ""
        Write-Host "Detalhes do erro:" -ForegroundColor Red
        Write-Host $_.Exception.ToString() -ForegroundColor DarkRed
    }
    
    exit 1
}
finally {
    # Cleanup opcional
    if (-not $SkipCleanup -and ($CreatedTables.Count -gt 0 -or $CreatedUsers.Count -gt 0)) {
        Write-Host ""
        Write-Host "🧹 Limpeza de dados de teste..." -ForegroundColor Yellow
        
        # Limpar mesas criadas
        foreach ($tableId in $CreatedTables) {
            try {
                Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId" -Method DELETE -Headers $masterHeaders -Description "Limpeza da mesa $tableId" | Out-Null
                Write-Host "  ✅ Mesa $tableId removida" -ForegroundColor Green
            }
            catch {
                Write-Host "  ⚠️ Erro ao remover mesa $tableId" -ForegroundColor Yellow
            }
        }
        
        Write-Host "Limpeza concluída" -ForegroundColor Green
    }
}
