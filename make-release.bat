@echo off
if not exist "release" mkdir release
set /p version="Version: "
set current_path=%cd%
set destination_path=\release\
set source=%current_path%\build
set filename=Ticket^ Place^ v
set ext=.zip
go build -o ./build/Ticket-Place.exe -ldflags -H=windowsgui cmd/main.go
7z a -tzip "%current_path%%destination_path%%filename%%version%%ext%" "%source%"