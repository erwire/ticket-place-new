@echo off
if not exist "build/release" mkdir build/release
go build -ldflags -H=windowsgui -o ./build/ticket-place_windows_amd64.exe cmd/main.go
go build -ldflags -H=windowsgui -o ./build/updater_windows_amd64.exe cmd/update.go
7z a -tzip "./build/release/ticket-place_windows_amd64.zip" "./build/ticket-place_windows_amd64.exe" "./build/updater_windows_amd64.exe"