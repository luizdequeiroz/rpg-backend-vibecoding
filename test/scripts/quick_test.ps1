# Quick Test Script - Teste Rápido de API
# Uso: .\test\scripts\quick_test.ps1

param(
    [string]$BaseUrl = "http://localhost:8080"
)

Write-Host "🚀 TESTE RÁPIDO DE API" -ForegroundColor Cyan
Write-Host "Base URL: $BaseUrl" -ForegroundColor Gray
Write-Host ""

# Função para teste simples
function Test-Endpoint {
    param(
        [string]$Uri,
        [string]$Name,
        [string]$Method = "GET",
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            Headers = $Headers
            TimeoutSec = 5
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = 'application/json'
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "✅ $Name" -ForegroundColor Green
        
        # Mostrar informação relevante
        if ($response.PSObject.Properties.Name -contains "status") {
            Write-Host "   Status: $($response.status)" -ForegroundColor Gray
        }
        if ($response.PSObject.Properties.Name -contains "version") {
            Write-Host "   Versão: $($response.version)" -ForegroundColor Gray
        }
        if ($response.PSObject.Properties.Name -contains "token") {
            Write-Host "   Token gerado: $($response.token.Substring(0, 20))..." -ForegroundColor Gray
        }
        if ($response.PSObject.Properties.Name -contains "id") {
            Write-Host "   ID: $($response.id)" -ForegroundColor Gray
        }
        
        return $response
    }
    catch {
        Write-Host "❌ $Name - $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# Testes básicos
Write-Host "1. Testando endpoints básicos..." -ForegroundColor Yellow

$health = Test-Endpoint -Uri "$BaseUrl/health" -Name "Health Check"
$root = Test-Endpoint -Uri "$BaseUrl/" -Name "Root Endpoint"

Write-Host ""

# Teste de autenticação
Write-Host "2. Testando autenticação..." -ForegroundColor Yellow

$testUser = @{
    email = "quicktest@exemplo.com"
    password = "senha123"
} | ConvertTo-Json

$signup = Test-Endpoint -Uri "$BaseUrl/api/v1/auth/signup" -Name "Signup" -Method "POST" -Body $testUser

if ($signup) {
    $token = $signup.token
    $authHeaders = @{ 'Authorization' = "Bearer $token" }
    
    Write-Host ""
    Write-Host "3. Testando endpoints autenticados..." -ForegroundColor Yellow
    
    # Testar criação de mesa
    $mesaData = @{
        name = "Mesa de Teste Rápido"
        system = "Sistema de Teste"
    } | ConvertTo-Json
    
    $mesa = Test-Endpoint -Uri "$BaseUrl/api/v1/tables/" -Name "Criar Mesa" -Method "POST" -Headers $authHeaders -Body $mesaData
    
    # Testar listagem de mesas
    Test-Endpoint -Uri "$BaseUrl/api/v1/tables/" -Name "Listar Mesas" -Headers $authHeaders
    
    if ($mesa) {
        # Testar detalhes da mesa
        Test-Endpoint -Uri "$BaseUrl/api/v1/tables/$($mesa.id)" -Name "Detalhes da Mesa" -Headers $authHeaders
    }
}

Write-Host ""
Write-Host "🎯 TESTE RÁPIDO CONCLUÍDO!" -ForegroundColor Green
Write-Host ""
Write-Host "Para teste completo, execute:" -ForegroundColor Yellow
Write-Host "  scripts.bat test-integration" -ForegroundColor White
Write-Host "  ou" -ForegroundColor Gray
Write-Host "  .\test\scripts\test_game_tables_v2.ps1" -ForegroundColor White
