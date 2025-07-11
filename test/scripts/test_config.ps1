# Configuração de Testes
# Este arquivo define configurações padrão para os testes automatizados

# URLs e Endpoints
$script:BaseUrl = "http://localhost:8080"
$script:HealthEndpoint = "/health"
$script:AuthEndpoint = "/api/v1/auth"
$script:TablesEndpoint = "/api/v1/tables"

# Timeouts
$script:RequestTimeout = 10
$script:HealthTimeout = 5

# Fixtures padrão
$script:DefaultMaster = "test/fixtures/auth_signup_master.json"
$script:DefaultPlayer = "test/fixtures/auth_signup_player.json"
$script:DefaultTable = "test/fixtures/game_table_create.json"
$script:DefaultInvite = "test/fixtures/invite_create.json"

# Funções utilitárias
function Get-TestConfig {
    return @{
        BaseUrl = $script:BaseUrl
        Timeout = $script:RequestTimeout
        Fixtures = @{
            Master = $script:DefaultMaster
            Player = $script:DefaultPlayer
            Table = $script:DefaultTable
            Invite = $script:DefaultInvite
        }
    }
}

function Test-ServerOnline {
    param([string]$Url = $script:BaseUrl)
    
    try {
        $response = Invoke-RestMethod -Uri "$Url$script:HealthEndpoint" -TimeoutSec $script:HealthTimeout
        return $response.status -eq "ok"
    }
    catch {
        return $false
    }
}

function Write-TestHeader {
    param([string]$Title)
    
    Write-Host ""
    Write-Host "=" * 60 -ForegroundColor Cyan
    Write-Host "  $Title" -ForegroundColor White
    Write-Host "=" * 60 -ForegroundColor Cyan
    Write-Host ""
}

function Write-TestStep {
    param([string]$Step, [string]$Description)
    
    Write-Host "$Step. $Description" -ForegroundColor Yellow
}

function Write-TestSuccess {
    param([string]$Message)
    
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-TestError {
    param([string]$Message)
    
    Write-Host "❌ $Message" -ForegroundColor Red
}

function Write-TestWarning {
    param([string]$Message)
    
    Write-Host "⚠️ $Message" -ForegroundColor Yellow
}

function Write-TestInfo {
    param([string]$Message)
    
    Write-Host "ℹ️ $Message" -ForegroundColor Blue
}
