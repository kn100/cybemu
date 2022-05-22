fmt :
	go fmt ./...
generate:
	go generate ./...
build: fmt generate
	go build -o cybemu
test: fmt generate
	go test ./... --cover
deps:
	go mod tidy
