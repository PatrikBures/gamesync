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
	go build -ldflags "-X gamesync/internal/vars.Version=$(VERSION)" -o bin/$(BIN_NAME) ./cmd/gamesync

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

up-pg:
	docker compose -f ./docker-compose-pg.yml up --build --remove-orphans -d
up-pg-gen:
	docker compose -f ./docker-compose-pg-gen.yml up --build --remove-orphans -d
down-pg:
	docker compose -f ./docker-compose-pg.yml down -v

psql:
	docker compose exec -it db psql $(DB_NAME) $(DB_USER) 

up:
	docker compose up --build --remove-orphans -d
down:
	docker compose down -v


### generate code
gen-all: gen-pg gen-api gen-strings
gen-most: gen-api gen-strings

gen-pg: down-pg up-pg-gen
	GAMESYNC_DB_TYPE=postgres GAMESYNC_DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME) go run ./cmd/gen/main.go
gen-sqlite: up # should not be used, just the gen-pg one 
	GAMESYNC_DB_TYPE=sqlite GAMESYNC_DB_URL=./data/sqlite_db/gamesync.sqlite go run ./cmd/gen/main.go

gen-api:
	go generate .

gen-strings:
	go generate ./internal/server/permissions


### linting
lint-go:
	golangci-lint run
lint-api:
	vacuum --ext-refs lint --min-score 100 api/*
lint-apid:
	vacuum --ext-refs dashboard api/*



### util
mkbin:
	mkdir -p bin

.PHONY: all build man install uninstall clean
