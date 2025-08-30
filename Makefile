SHELL := /bin/bash

.PHONY: all lint test vet vuln fmt tidy tools check ci

GO := go
PKGS := ./...

GOLANGCI_LINT := $(shell command -v golangci-lint 2>/dev/null)
GOVULNCHECK := $(shell command -v govulncheck 2>/dev/null)

all: check

lint:
ifndef GOLANGCI_LINT
	@echo "golangci-lint not found. Run 'make tools' first." && exit 1
endif
	golangci-lint run ./...

test:
	$(GO) test $(PKGS) -race -coverprofile=coverage.out

vet:
	$(GO) vet $(PKGS)

vuln:
ifndef GOVULNCHECK
	@echo "govulncheck not found. Run 'make tools' first." && exit 1
endif
	govulncheck $(PKGS)

fmt:
	$(GO) fmt $(PKGS)

tidy:
	$(GO) mod tidy

tools:
	@echo "Installing local dev tools..."
	# Use a golangci-lint version compatible with newer Go (e.g., 1.24)
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Tools installed to \"$$(go env GOPATH)/bin\". Ensure it is on your PATH."

.PHONY: cyclo
cyclo:
	@command -v gocyclo >/dev/null 2>&1 || { echo "Installing gocyclo..."; $(GO) install github.com/fzipp/gocyclo/cmd/gocyclo@latest; }
	gocyclo -over 15 ./...

check: fmt vet lint test

ci: tidy check vuln
