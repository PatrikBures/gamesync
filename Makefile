ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BIN_NAME := gamesync
BIN_NAME_DEV := $(BIN_NAME)-dev
BIN_STATE_NAME := gamesync-state
BIN_SERVER_NAME := $(BIN_NAME)-server
CONTAINER_NAME := $(BIN_NAME)
VERSION ?= dev
PREFIX ?= /usr/local
BIN_DIR := $(DESTDIR)$(PREFIX)/bin
MAN1_DIR := $(DESTDIR)$(PREFIX)/share/man/man1
LICENSE_DIR := $(DESTDIR)$(PREFIX)/share/licenses/$(BIN_NAME)

all: build man

### build
build: mkbin
	@echo "Building $(BIN_NAME)..."
	CGO_ENABLED=0 go build -ldflags "-X gamesync/internal/vars.Version=$(VERSION)" -o bin/$(BIN_NAME) ./cmd/gamesync

build-server: mkbin
	@echo "Building $(BIN_SERVER_NAME)..."
	go build -ldflags="-s -w" -trimpath -o bin/$(BIN_SERVER_NAME) ./cmd/gamesync-server


### docs
man: build
	@echo "Generating man pages..."
	mkdir -p manpages
	bin/$(BIN_NAME) gen-man

### install client
install: all
	@echo "Installing binary to $(BIN_DIR)..."
	install -d $(BIN_DIR)
	install -m 755 bin/$(BIN_NAME) $(BIN_DIR)/$(BIN_NAME)

	@echo "Installing man pages..."
	install -d $(MAN1_DIR)
	install -m 644 manpages/*.1 $(MAN1_DIR)/

	@echo "Updating man database..."
	-mandb >/dev/null 2>&1 || true

	@echo "Installing license..."
	install -Dm644 LICENSE $(LICENSE_DIR)/LICENSE

uninstall:
	@echo "Removing $(BIN_NAME)..."
	rm -f $(BIN_DIR)/$(BIN_NAME)
	rm -f $(MAN1_DIR)/gamesync*.1
	rm -rf $(LICENSE_DIR)

clean:
	@echo "Cleaning up..."
	rm -f bin/$(BIN_NAME)
	rm -rf manpages

go-install:
	@echo "installing..."
	go build -ldflags "-X gamesync/internal/vars.Version=$(VERSION)" -o $${GOPATH}/bin/$(BIN_NAME_DEV) ./cmd/gamesync

### test client
go-test:
	@echo "testing..."
	go test -v ./...



### container
build-container: build-state
	@echo "building container..."
	docker build ./ -t $(CONTAINER_NAME):$(VERSION)

psql:
	docker compose exec -it db psql $(DB_NAME) $(DB_USER) 

up:
	docker compose up --build --remove-orphans -d
down:
	docker compose down -v

dev:
	docker compose up --build --watch


### generate code
gen: gen-sql gen-api gen-strings

gen-api:
	go generate .

gen-strings:
	go generate ./internal/server/permissions

gen-sql:
	rm -rf ./internal/server/dbm
	sqlc generate


### linting
lint: lint-go lint-api

lint-go:
	golangci-lint run
lint-api:
	vacuum --ext-refs lint --min-score 100 api/*
lint-apid:
	vacuum --ext-refs dashboard api/*


### format
fmt:
	go fmt ./...


### util
mkbin:
	mkdir -p bin

### debug
pprof-cpu:
	go tool pprof -http=:8090 ./bin/gamesync cpu.prof
pprof-mem:
	go tool pprof -http=:8091 ./bin/gamesync mem.prof
trace:
	go tool trace trace.out

.PHONY: all build man install uninstall clean
