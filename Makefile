BIN_NAME := gamesync
BIN_NAME_DEV := $(BIN_NAME)-dev
BIN_STATE_NAME := gamesync-state
CONTAINER_NAME := $(BIN_NAME)
VERSION ?= dev
PREFIX ?= /usr/local
BIN_DIR := $(DESTDIR)$(PREFIX)/bin
MAN1_DIR := $(DESTDIR)$(PREFIX)/share/man/man1
LICENSE_DIR := $(DESTDIR)$(PREFIX)/share/licenses/$(BIN_NAME)

all: build man

build:
	@echo "Building $(BIN_NAME)..."
	mkdir -p bin
	go build -ldflags "-X gamesync/internal/vars.Version=$(VERSION)" -o bin/$(BIN_NAME) ./cmd/gamesync

man: build
	@echo "Generating man pages..."
	mkdir -p manpages
	bin/$(BIN_NAME) gen-man

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

go-test:
	@echo "testing..."
	go test -v ./...


build-state:
	@echo "building $(BIN_STATE_NAME)..."
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/$(BIN_STATE_NAME) ./cmd/gamesync-state

build-container: build-state
	@echo "building container..."
	docker build ./ -t $(CONTAINER_NAME):$(VERSION)

.PHONY: all build man install uninstall clean
