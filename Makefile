run:
	go run cmd/main.go
test-run:
	go run cmd/test.go
builder:
	go build -o ./build/Ticket-Place.exe -ldflags -H=windowsgui cmd/main.go


