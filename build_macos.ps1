# Script de build pour macOS
# Nécessite un cross-compiler C pour macOS

Write-Host "Construction de Bashcord pour macOS..." -ForegroundColor Green

# Vérifier si nous avons les outils nécessaires
Write-Host "Vérification des prérequis..." -ForegroundColor Yellow

# Obtenir le hash git
$gitHash = "dev"
$gitTag = "dev-build"

try {
    $gitHash = (git rev-parse --short HEAD 2>$null) -replace "`n", "" -replace "`r", ""
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Attention: Impossible d'obtenir le hash git, utilisation de 'dev'" -ForegroundColor Yellow
        $gitHash = "dev"
    }
} catch {
    Write-Host "Attention: Git non disponible, utilisation de 'dev'" -ForegroundColor Yellow
    $gitHash = "dev"
}

Write-Host "Hash Git: $gitHash" -ForegroundColor Cyan
Write-Host "Tag: $gitTag" -ForegroundColor Cyan

# Essayer de compiler sans CGO d'abord (ne fonctionnera probablement pas pour la GUI)
Write-Host "`nTentative de compilation sans CGO (peut échouer)..." -ForegroundColor Yellow
$env:CGO_ENABLED = "0"
$env:GOOS = "darwin"
$env:GOARCH = "amd64"

$ldflags = "-s -w -X 'vencord/buildinfo.InstallerGitHash=$gitHash' -X 'vencord/buildinfo.InstallerTag=$gitTag'"

try {
    go build -v -tags "static gui" -ldflags $ldflags -o "Bashcord-macOS-nocgo"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Succès! Binaire créé: Bashcord-macOS-nocgo" -ForegroundColor Green
        Write-Host "Attention: Ce binaire peut ne pas fonctionner correctement car il manque les dépendances GUI" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Échec de la compilation sans CGO (attendu pour la GUI)" -ForegroundColor Red
}

# Essayer de compiler la version CLI qui ne nécessite pas CGO
Write-Host "`nTentative de compilation de la version CLI..." -ForegroundColor Yellow
try {
    go build -v -tags "static cli" -ldflags $ldflags -o "Bashcord-macOS-cli"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Succès! Version CLI créée: Bashcord-macOS-cli" -ForegroundColor Green
    }
} catch {
    Write-Host "Échec de la compilation CLI" -ForegroundColor Red
}

Write-Host "`nPour compiler la version GUI complète, vous avez besoin:" -ForegroundColor Yellow
Write-Host "1. Un cross-compiler C pour macOS (comme zig ou un toolchain macOS)" -ForegroundColor White
Write-Host "2. Ou compiler directement sur macOS" -ForegroundColor White
Write-Host "`nCommande pour macOS natif:" -ForegroundColor Cyan
Write-Host "CGO_ENABLED=1 go build -v -tags `"static gui`" -ldflags `"$ldflags`" -o Bashcord-macOS" -ForegroundColor Gray 