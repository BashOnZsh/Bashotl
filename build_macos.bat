@echo off
echo Construction de Bashcord pour macOS...

REM Obtenir le hash git
for /f %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_HASH=%%i
if "%GIT_HASH%"=="" set GIT_HASH=dev

echo Hash Git: %GIT_HASH%

REM Compiler la version CLI pour macOS
echo.
echo Compilation de la version CLI pour macOS...
set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -v -tags "static cli" -ldflags "-s -w -X 'vencord/buildinfo.InstallerGitHash=%GIT_HASH%' -X 'vencord/buildinfo.InstallerTag=dev-build'" -o Bashcord-macOS-cli

if %errorlevel% equ 0 (
    echo.
    echo Succès! Version CLI créée: Bashcord-macOS-cli
) else (
    echo.
    echo Échec de la compilation CLI
)

echo.
echo Pour compiler la version GUI complète, vous avez besoin:
echo 1. Un cross-compiler C pour macOS (comme zig ou un toolchain macOS)
echo 2. Ou compiler directement sur macOS
echo.
echo Commande pour macOS natif:
echo CGO_ENABLED=1 go build -v -tags "static gui" -ldflags "-s -w -X 'vencord/buildinfo.InstallerGitHash=%GIT_HASH%' -X 'vencord/buildinfo.InstallerTag=dev-build'" -o Bashcord-macOS

pause 