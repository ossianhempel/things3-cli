BIN_NAME := things
PREFIX ?= /usr/local
BIN_DIR := $(PREFIX)/bin
MAN_DIR := $(PREFIX)/share/man/man1
BUILD_DIR := bin

.PHONY: build test install uninstall check-things-skill test-things-skill-sync sync-things-skill

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/$(BIN_NAME)

test:
	go test ./...

check-things-skill:
	./scripts/check-things-skill-sync.sh

test-things-skill-sync:
	./scripts/test-sync-things-skill.sh

sync-things-skill:
	./scripts/sync-things-skill.sh

install: build
	@mkdir -p $(BIN_DIR)
	cp $(BUILD_DIR)/$(BIN_NAME) $(BIN_DIR)/$(BIN_NAME)
	@mkdir -p $(MAN_DIR)
	@if [ -f share/man/man1/things.1 ]; then \
		cp share/man/man1/things.1 $(MAN_DIR)/things.1; \
	fi

uninstall:
	rm -f $(BIN_DIR)/$(BIN_NAME)
	rm -f $(MAN_DIR)/things.1
