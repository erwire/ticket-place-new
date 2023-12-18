@echo off
set /p version="version: "

if not exist "build/release" mkdir build/release
if not exist "build/release/%version%" mkdir "build/release/%version%"



go build -ldflags -H=windowsgui -o ./build/release/%version%/ticket-place_windows_amd64.exe cmd/main.go
go build -ldflags -H=windowsgui -o ./build/release/%version%/updater_windows_amd64.exe cmd/update.go
7z a -tzip "./build/release/%version%/ticket-place_windows_amd64.zip" "./build/release/%version%/ticket-place_windows_amd64.exe" "./build/release/%version%/updater_windows_amd64.exe"