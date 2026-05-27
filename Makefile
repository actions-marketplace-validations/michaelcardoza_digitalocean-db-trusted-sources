BIN := dist/do-db-trusted-sources

.PHONY: build run clean vet test

build:
	@mkdir -p dist
	go build -o $(BIN) .

run: build
	./$(BIN)

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -rf dist
