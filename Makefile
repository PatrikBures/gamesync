BIN_NAME := gamesync
PREFIX ?= /usr/local
BIN_DIR := $(DESTDIR)$(PREFIX)/bin
MAN1_DIR := $(DESTDIR)$(PREFIX)/share/man/man1

all: build man

build:
	@echo "Building $(BIN_NAME)..."
	cd client && go build -o ../$(BIN_NAME) ./cmd/gamesync

man: build
	@echo "Generating man pages..."
	mkdir -p manpages
	./$(BIN_NAME) gen-man

install: all
	@echo "Installing binary to $(BIN_DIR)..."
	install -d $(BIN_DIR)
	install -m 755 $(BIN_NAME) $(BIN_DIR)/$(BIN_NAME)

	@echo "Installing man pages..."
	install -d $(MAN1_DIR)
	install -m 644 manpages/*.1 $(MAN1_DIR)/

	@echo "Updating man database..."
	-mandb >/dev/null 2>&1 || true

uninstall:
	@echo "Removing $(BIN_NAME)..."
	rm -f $(BIN_DIR)/$(BIN_NAME)
	rm -f $(MAN1_DIR)/gamesync*.1

clean:
	@echo "Cleaning up..."
	rm -f $(BIN_NAME)
	rm -rf manpages

go-install:
	@echo "installing..."
	cd client && go install ./cmd/gamesync

.PHONY: all build man install uninstall clean
