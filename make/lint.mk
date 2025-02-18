## Lint:

.PHONY: lint-go
lint-go: check-lint-go generate-mocks ## Run golangci-lint to check Go code for style and potential errors
	@echo "🔍 Linting Go files"
	@golangci-lint run

.PHONY: lint-go-clean
lint-go-clean: check-lint-go ## Clean golangci-lint cache
	@echo "🗑 Cleaning golangci-lint cache"
	@golangci-lint cache clean

.PHONY: check-lint-go
check-lint-go:
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@echo "❗ golangci-lint not found, run `make install`"
	@exit 1
endif

