
build:
	go build -o bin/ ./cmd/...

webserver: build
	./bin/webserver

