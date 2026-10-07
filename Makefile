# Ikou backend task runner. `make` or `make help` lists everything.

BIN       := dist/api
PKG       := ./...
PORT      := 9001
COVERFILE := coverage.out

.DEFAULT_GOAL := help
.PHONY: help build start run stop test test-verbose test-cover cover-html fmt fmt-check vet check docker-build docker-up docker-down docker-logs clean

help:
	@echo "Ikou backend"
	@echo ""
	@echo "  Run"
	@echo "    make start         build, then run locally in the foreground"
	@echo "    make stop          kill a locally running backend"
	@echo "    make docker-up     build and run in Docker (reads app.env)"
	@echo "    make docker-down   stop and remove the Docker containers"
	@echo "    make docker-logs   follow the container logs"
	@echo ""
	@echo "  Verify"
	@echo "    make check         fmt-check + vet + test  <- run before committing"
	@echo "    make test          go test ./..."
	@echo "    make test-cover    tests with a coverage total"
	@echo "    make cover-html    open the coverage report in a browser"
	@echo ""
	@echo "  Fix / build"
	@echo "    make fmt           gofmt -w ."
	@echo "    make build         compile to $(BIN)"
	@echo "    make clean         remove build output and coverage (deletes files)"

build:
	@echo "Building..."
	@go build -o $(BIN) ./cmd
	@echo "Built $(BIN)"

start: build
	@echo "Starting backend on port $(PORT)..."
	@./$(BIN)

run: start

stop:
	@pkill -f '$(BIN)' && echo "Backend stopped." || echo "No running backend found."

test:
	@go test $(PKG)

test-verbose:
	@go test -v $(PKG)

test-cover:
	@go test -coverprofile=$(COVERFILE) $(PKG)
	@go tool cover -func=$(COVERFILE) | tail -n 1

cover-html: test-cover
	@go tool cover -html=$(COVERFILE)

fmt:
	@gofmt -w .
	@echo "Formatted."

fmt-check:
	@test -z "$$(gofmt -l .)" || { echo "These files need 'make fmt':"; gofmt -l .; exit 1; }

vet:
	@go vet $(PKG)

check: fmt-check vet test
	@echo "fmt, vet and tests all pass."

docker-build:
	@docker compose build

docker-up:
	@docker compose up -d --build
	@echo "Backend running at http://localhost:$(PORT)"

docker-down:
	@docker compose down

docker-logs:
	@docker compose logs -f app

clean:
	@rm -rf dist $(COVERFILE)
	@echo "Cleaned."
