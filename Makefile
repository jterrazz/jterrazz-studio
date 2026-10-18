.PHONY: help install build lint fix test check provision deprovision skills clean

BINARY          := j
BUILD_PATH      := .artifacts/go/$(BINARY)
JTERRAZZ_DIR    := $(HOME)/.jterrazz
BIN_DIR         := $(JTERRAZZ_DIR)/bin
INSTALL_PATH    := $(BIN_DIR)/$(BINARY)
ZSHRC_SOURCE    := dotfiles/applications/zsh/zshrc.sh
OLD_CONFIG_DIR  := $(HOME)/.config/jterrazz
OLD_INSTALL     := /usr/local/bin/$(BINARY)

# hazard: an unpinned `@latest` makes a green run a fact about the day it ran.
# Bump this line, read what the new release finds, land both together.
GOLANGCI_VERSION := v2.13.2
GOLANGCI         := $(shell go env GOPATH)/bin/golangci-lint

CYAN  := \033[36m
GREEN := \033[32m
DIM   := \033[2m
RESET := \033[0m

help: ## Show available targets
	@printf "$(CYAN)jterrazz-studio$(RESET)\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(CYAN)%-12s$(RESET) %s\n", $$1, $$2}'

node_modules/.install: package.json package-lock.json
	@npm ci --silent
	@touch node_modules/.install

install: node_modules/.install ## Install the dependencies of both halves
	@go mod download

build: ## Build the binary
	@go build -o $(BUILD_PATH) ./src/cmd/j
	@printf "$(GREEN)✓$(RESET) Built $(BUILD_PATH)\n"

lint: node_modules/.install ## Both halves: gofmt and golangci-lint (pinned) over Go, typescript check over the rest
	@unformatted=$$(gofmt -l ./src/); \
		if [ -n "$$unformatted" ]; then printf "gofmt would rewrite:\n%s\n" "$$unformatted"; exit 1; fi
	@$(GOLANGCI) --version 2>/dev/null | grep -q "$(GOLANGCI_VERSION:v%=%)" || \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	@$(GOLANGCI) run ./src/...
	@printf "$(GREEN)✓$(RESET) Go lint passed\n"
	@npm run --silent lint

fix: node_modules/.install ## Every fixer: gofmt over Go, then typescript fix
	@status=0; \
	gofmt -w ./src/ || status=1; \
	npx typescript fix || status=1; \
	exit $$status

test: node_modules/.install ## The Go module tests, then the CLI specs against a fresh build of j
	@go test ./src/...
	@J_FORCE_REBUILD=1 npm test

check: lint test ## Every gate a change passes

provision: build ## Put this build on the machine: ~/.jterrazz/bin, and the zshrc line
	@mkdir -p $(BIN_DIR)
	@cp $(BUILD_PATH) $(INSTALL_PATH)
	@chmod +x $(INSTALL_PATH)
	@printf "$(GREEN)✓$(RESET) Installed $(INSTALL_PATH)\n"
	@if [ -f "$(OLD_INSTALL)" ]; then \
		rm "$(OLD_INSTALL)" 2>/dev/null || printf "$(DIM)Cannot remove $(OLD_INSTALL); run: sudo rm $(OLD_INSTALL)$(RESET)\n"; \
	fi
	@if [ -f "$(OLD_CONFIG_DIR)/jrc.json" ] && [ ! -f "$(JTERRAZZ_DIR)/config.json" ]; then \
		cp "$(OLD_CONFIG_DIR)/jrc.json" "$(JTERRAZZ_DIR)/config.json"; \
		printf "$(GREEN)✓$(RESET) Migrated config.json\n"; \
	fi
	@if [ -d "$(OLD_CONFIG_DIR)/tailscale" ] && [ ! -d "$(JTERRAZZ_DIR)/tailscale" ]; then \
		cp -r "$(OLD_CONFIG_DIR)/tailscale" "$(JTERRAZZ_DIR)/tailscale"; \
		printf "$(GREEN)✓$(RESET) Migrated tailscale state\n"; \
	fi
	@if [ -f "$$HOME/.zshrc" ]; then \
		if ! grep -q "$(ZSHRC_SOURCE)" "$$HOME/.zshrc"; then \
			printf '\n# jterrazz-studio\nsource $(CURDIR)/$(ZSHRC_SOURCE)\n' >> "$$HOME/.zshrc"; \
			printf "$(GREEN)✓$(RESET) Added source to ~/.zshrc\n"; \
		fi \
	else \
		printf "$(DIM)~/.zshrc not found; add it by hand: source $(CURDIR)/$(ZSHRC_SOURCE)$(RESET)\n"; \
	fi
	@printf "$(GREEN)✓$(RESET) Done. Run $(DIM)source ~/.zshrc$(RESET) then $(DIM)j help$(RESET)\n"

deprovision: ## Remove the binary from ~/.jterrazz/bin; the user data stays
	@if [ -f "$(INSTALL_PATH)" ]; then \
		rm "$(INSTALL_PATH)"; \
		printf "$(GREEN)✓$(RESET) Removed $(INSTALL_PATH)\n"; \
	else \
		printf "$(DIM)$(BINARY) not found at $(INSTALL_PATH)$(RESET)\n"; \
	fi
	@if [ -f "$(OLD_INSTALL)" ]; then \
		rm "$(OLD_INSTALL)" 2>/dev/null || printf "$(DIM)Cannot remove $(OLD_INSTALL); run: sudo rm $(OLD_INSTALL)$(RESET)\n"; \
	fi

skills: ## Regenerate the toolbelt skill rosters
	@go run ./src/cmd/skillsgen
	@printf "$(GREEN)✓$(RESET) Regenerated skills/jterrazz-toolbelt/SKILL.md\n"

clean: ## Remove build artifacts
	@rm -rf .artifacts
	@printf "$(GREEN)✓$(RESET) Cleaned\n"
