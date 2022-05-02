fmt :
	go fmt ./...
build: fmt
	go build -o cybemu
test: fmt
	go test ./... --cover
deps:
	go mod tidy
