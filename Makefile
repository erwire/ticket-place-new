run:
	go run cmd/main.go
test-run:
	go run cmd/test.go
builder:
	SET GOOS=windows
	SET GOARCH=386
	go build -o ./build/Ticket-Place_windows_386.exe -ldflags -H=windowsgui cmd/main.go
	SET GOOS=windows
	SET GOARCH=amd64
	go build -o ./build/Ticket-Place_windows_amd64.exe -ldflags -H=windowsgui cmd/main.go
build-test:
	go build -ldflags -H=windowsgui -o ./build/ticket-place_windows_amd64.exe cmd/main.go
	go build -ldflags -H=windowsgui -o ./build/updater_windows_amd64.exe cmd/update.go