APP_NAME   := app
SRC_DIR    := $(CURDIR)/src
BIN_DIR    := $(CURDIR)/bin
APP_FILE   := $(BIN_DIR)/$(APP_NAME)
CMD_DIR    := ./cmd/server
IMAGE_NAME := {{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}

.PHONY: build test test-race vet fmt fmt-check lint clean run container

build: fmt-check vet lint test clean
	mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && go build -o $(APP_FILE) $(CMD_DIR)

run:
	cd $(SRC_DIR) && go run $(CMD_DIR)

test:
	cd $(SRC_DIR) && go test ./...

test-race:
	cd $(SRC_DIR) && go test -race ./...

vet:
	cd $(SRC_DIR) && go vet ./...

fmt:
	cd $(SRC_DIR) && gofmt -w .

fmt-check:
	@cd $(SRC_DIR) && test -z "$$(gofmt -l .)" || (echo "files not formatted:"; gofmt -l .; exit 1)

lint: vet
	cd $(SRC_DIR) && golangci-lint run ./...

clean:
	@[ -d "$(BIN_DIR)" ] && rm -f $(BIN_DIR)/* || true

container: fmt-check vet lint test
	docker build --tag $(IMAGE_NAME) .
