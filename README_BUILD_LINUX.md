# Compilation Bashcord pour Linux

Ce guide explique comment compiler Bashcord pour Linux à partir d'un environnement Windows ou Linux.

## Scripts disponibles

### PowerShell (Windows)
```powershell
.\build_linux.ps1
```

### Batch (Windows/Cross-compilation)
```cmd
.\build_linux.bat
```

## Versions créées

- **Bashcord-Linux-cli** : Version ligne de commande (sans interface graphique)
- **Bashcord-Linux-gui** : Version avec interface graphique (nécessite dépendances Linux)

## Compilation sur Linux natif

### Prérequis Ubuntu/Debian
```bash
sudo apt update
sudo apt install -y git golang-go pkg-config \
  libsdl2-dev libx11-dev libxcursor-dev libxrandr-dev \
  libxinerama-dev libxi-dev libglx-dev libgl1-mesa-dev \
  libxxf86vm-dev libwayland-dev libxkbcommon-dev \
  wayland-protocols extra-cmake-modules
```

### Prérequis Arch Linux
```bash
sudo pacman -S git go pkg-config sdl2 libx11 libxcursor \
  libxrandr libxinerama libxi mesa wayland wayland-protocols
```

### Compilation CLI (sans interface)
```bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

GIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS="-s -w -X 'vencord/buildinfo.InstallerGitHash=$GIT_HASH' -X 'vencord/buildinfo.InstallerTag=dev-build'"

go build -v -tags "static cli" -ldflags "$LDFLAGS" -o "Bashcord-Linux-cli"
```

### Compilation GUI (avec interface)
```bash
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

GIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS="-s -w -X 'vencord/buildinfo.InstallerGitHash=$GIT_HASH' -X 'vencord/buildinfo.InstallerTag=dev-build'"

go build -v -tags "static gui" -ldflags "$LDFLAGS" -o "Bashcord-Linux-gui"
```

## Cross-compilation depuis Windows

La version CLI peut être compilée depuis Windows sans problème.
La version GUI nécessite CGO et les bibliothèques Linux, donc elle échouera probablement sur Windows.

## Dépannage

### Erreur "CGO required"
- Installez les dépendances de développement listées ci-dessus
- Vérifiez que `pkg-config` est installé
- Essayez de compiler uniquement la version CLI en premier

### Erreur "go command not found"
```bash
# Ubuntu/Debian
sudo apt install golang-go

# Arch Linux  
sudo pacman -S go
```

### Version Go trop ancienne
Bashcord nécessite Go 1.19 ou plus récent :
```bash
go version
```

## Utilisation

### CLI
```bash
./Bashcord-Linux-cli
```

### GUI
```bash
./Bashcord-Linux-gui
```

## Tailles typiques

- CLI : ~8-10 MB
- GUI : ~15-20 MB (avec les bibliothèques graphiques) 