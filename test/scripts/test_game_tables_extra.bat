@echo off
REM Script para testar o sistema completo de GameTables e Invites no Windows
REM Uso: test_game_tables.bat

echo === TESTANDO SISTEMA DE MESAS E CONVITES ===
echo.

REM Verificar se o servidor está rodando
echo 1. Verificando se o servidor está online...
curl -s http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    echo ❌ Servidor não está rodando em localhost:8080
    echo Execute: scripts.bat run
    exit /b 1
)
echo ✅ Servidor está online
echo.

REM Registrar mestre
echo 2. Registrando usuário mestre...
curl -s -X POST "http://localhost:8080/api/v1/auth/signup" -H "Content-Type: application/json" -d @test_signup.json > mestre_response.tmp
type mestre_response.tmp

REM Se usuário já existe, fazer login
curl -s -X POST "http://localhost:8080/api/v1/auth/login" -H "Content-Type: application/json" -d @test_signup.json > mestre_login.tmp
echo ✅ Mestre autenticado
echo.

REM Registrar jogador
echo 3. Registrando usuário jogador...
curl -s -X POST "http://localhost:8080/api/v1/auth/signup" -H "Content-Type: application/json" -d @test_signup_player.json > jogador_response.tmp
type jogador_response.tmp

REM Se usuário já existe, fazer login
curl -s -X POST "http://localhost:8080/api/v1/auth/login" -H "Content-Type: application/json" -d @test_signup_player.json > jogador_login.tmp
echo ✅ Jogador autenticado
echo.

echo === PRÓXIMOS PASSOS ===
echo 1. Extrair os tokens JWT dos arquivos .tmp
echo 2. Usar os tokens para criar mesa e convites
echo 3. Testar fluxo completo de convites
echo.
echo Arquivos gerados:
echo - mestre_response.tmp (ou mestre_login.tmp)
echo - jogador_response.tmp (ou jogador_login.tmp)
echo.
echo Use PowerShell para testes mais avançados:
echo $mestre = Get-Content mestre_login.tmp ^| ConvertFrom-Json
echo $token = $mestre.token

REM Limpeza opcional
REM del *.tmp
