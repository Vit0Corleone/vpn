build-server:
	@go build -o ./s ./server/main.go

build-client:
	@go build -o ./c ./client/main.go