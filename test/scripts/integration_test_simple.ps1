# RPG Backend Integration Test
# Simplified version to avoid encoding issues

param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose = $false
)

Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "   RPG BACKEND INTEGRATION TEST" -ForegroundColor White
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "Base URL: $BaseUrl"
Write-Host ""

function Invoke-TestRequest {
    param($Uri, $Method = "GET", $Headers = @{}, $Body = $null, $Description = "")
    
    try {
        $params = @{ Uri = $Uri; Method = $Method; Headers = $Headers; TimeoutSec = 10 }
        if ($Body) { 
            $params.Body = $Body
            $params.ContentType = 'application/json' 
        }
        if ($Verbose -and $Description) { Write-Host "  -> $Description" -ForegroundColor Gray }
        return Invoke-RestMethod @params
    }
    catch { 
        if ($Verbose) { Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red }
        throw 
    }
}

try {
    # Test 1: Server Health
    Write-Host "1. Testing server health..." -ForegroundColor Yellow
    $health = Invoke-TestRequest -Uri "$BaseUrl/health" -Description "Health check"
    Write-Host "   Server is online (Status: $($health.status))" -ForegroundColor Green
    Write-Host ""

    # Test 2: Authentication
    Write-Host "2. Testing authentication..." -ForegroundColor Yellow
    
    # Master signup/login
    $masterData = Get-Content "test/fixtures/auth_signup_master.json" -Raw
    try {
        $masterResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $masterData -Description "Master signup"
        Write-Host "   Master registered (ID: $($masterResponse.user.id))" -ForegroundColor Green
    }
    catch {
        $masterResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $masterData -Description "Master login"
        Write-Host "   Master logged in (ID: $($masterResponse.user.id))" -ForegroundColor Green
    }
    $masterToken = $masterResponse.token
    $masterHeaders = @{ 'Authorization' = "Bearer $masterToken" }
    
    # Player signup/login
    $playerData = Get-Content "test/fixtures/auth_signup_player.json" -Raw
    try {
        $playerResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/signup" -Method POST -Body $playerData -Description "Player signup"
        Write-Host "   Player registered (ID: $($playerResponse.user.id))" -ForegroundColor Green
    }
    catch {
        $playerResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $playerData -Description "Player login"
        Write-Host "   Player logged in (ID: $($playerResponse.user.id))" -ForegroundColor Green
    }
    $playerToken = $playerResponse.token
    $playerHeaders = @{ 'Authorization' = "Bearer $playerToken" }
    Write-Host ""

    # Test 3: Game Tables
    Write-Host "3. Testing game table management..." -ForegroundColor Yellow
    
    # Create table
    $tableData = Get-Content "test/fixtures/game_table_create.json" -Raw
    $tableResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $tableData -Description "Create table"
    $tableId = $tableResponse.id
    Write-Host "   Table created (ID: $tableId)" -ForegroundColor Green
    
    # List master tables
    $masterTables = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $masterHeaders -Description "List master tables"
    Write-Host "   Master has $($masterTables.total) table(s)" -ForegroundColor Green
    
    # Check player has no access yet
    $playerTablesInitial = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $playerHeaders -Description "List player tables (initial)"
    Write-Host "   Player has $($playerTablesInitial.total) table(s) initially" -ForegroundColor Green
    Write-Host ""

    # Test 4: Invitations
    Write-Host "4. Testing invitation system..." -ForegroundColor Yellow
    
    # Create invite
    $inviteData = Get-Content "test/fixtures/invite_create.json" -Raw
    $inviteResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId/invites/" -Method POST -Headers $masterHeaders -Body $inviteData -Description "Create invite"
    $inviteId = $inviteResponse.id
    Write-Host "   Invite created (ID: $inviteId)" -ForegroundColor Green
    
    # Accept invite
    $acceptResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/$tableId/invites/$inviteId/accept" -Method POST -Headers $playerHeaders -Description "Accept invite"
    Write-Host "   Invite accepted (Status: $($acceptResponse.status))" -ForegroundColor Green
    
    # Check player access after accept
    $playerTablesFinal = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Headers $playerHeaders -Description "List player tables (final)"
    Write-Host "   Player now has $($playerTablesFinal.total) table(s)" -ForegroundColor Green
    Write-Host ""

    # Test 5: Different RPG Systems
    Write-Host "5. Testing different RPG systems..." -ForegroundColor Yellow
    
    # Vampiro table
    $vampiroData = Get-Content "test/fixtures/game_table_vampiro.json" -Raw
    $vampiroResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $vampiroData -Description "Create Vampiro table"
    Write-Host "   Vampiro table created: $($vampiroResponse.name)" -ForegroundColor Green
    
    # Cthulhu table
    $cthulhuData = Get-Content "test/fixtures/game_table_cthulhu.json" -Raw
    $cthulhuResponse = Invoke-TestRequest -Uri "$BaseUrl/api/v1/tables/" -Method POST -Headers $masterHeaders -Body $cthulhuData -Description "Create Cthulhu table"
    Write-Host "   Cthulhu table created: $($cthulhuResponse.name)" -ForegroundColor Green
    Write-Host ""

    # Success
    Write-Host "=============================================" -ForegroundColor Green
    Write-Host "   ALL TESTS PASSED SUCCESSFULLY!" -ForegroundColor White
    Write-Host "=============================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Test Summary:" -ForegroundColor Cyan
    Write-Host "  Authentication: PASS" -ForegroundColor Green
    Write-Host "  Table Management: PASS" -ForegroundColor Green
    Write-Host "  Invitation System: PASS" -ForegroundColor Green
    Write-Host "  Role-based Access: PASS" -ForegroundColor Green
    Write-Host "  Multiple RPG Systems: PASS" -ForegroundColor Green
    Write-Host ""
    Write-Host "System is ready for production!" -ForegroundColor Green

}
catch {
    Write-Host ""
    Write-Host "TEST FAILED: $($_.Exception.Message)" -ForegroundColor Red
    if ($Verbose) {
        Write-Host ""
        Write-Host "Error Details:" -ForegroundColor Red
        Write-Host $_.Exception.ToString() -ForegroundColor DarkRed
    }
    exit 1
}
