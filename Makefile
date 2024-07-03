TEMP_DIR=./tmp
BINARY_NAME=ctrl_plus_revise

clean:
	@echo "\n> Cleaning project...\n"
	go clean
	rm -rf ${TEMP_DIR}

test:
	@echo "\n> Run tests...\n"
	go test -v -cover -race ./...

build: clean test
	@echo "\n> Building project backend...\n"
	go build -o ${TEMP_DIR}/${BINARY_NAME} .

run: build
	@echo "\n> Running project...\n"
	${TEMP_DIR}/${BINARY_NAME}

lint:
	@echo "\n> Run linter...\n"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

stringer:
	@echo "\n> Run stringer...\n"
	go run golang.org/x/tools/cmd/stringer@latest -linecomment -type=PromptMsg
	go run golang.org/x/tools/cmd/stringer@latest -linecomment -type=processingDevice
	go run golang.org/x/tools/cmd/stringer@latest -linecomment -type=gpuModel
	go run golang.org/x/tools/cmd/stringer@latest -linecomment -type=ModelName