@echo off
echo Compilation de Bashcord CLI pour macOS...
set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -tags "static,cli" -o Bashcord-macOS-cli
echo Terminé! 