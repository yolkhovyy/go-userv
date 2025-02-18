ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

## Mocks:

.PHONY: check-mockery
check-mockery:
ifeq (, $(shell command -v mockery 2> /dev/null))
	@echo "❗ mockery not installed, run `make install`"
	@exit 1
endif

.PHONY: remove-mocks
remove-mocks: ## Remove generated mock files
	@echo "🗑 Removing generated mock files"
	@find . -type f -name "mock_*.go" -delete

.PHONY: generate-mocks
generate-mocks: check-mockery remove-mocks ## Generate mock implementations for interfaces using mockery
	@echo "⚙️ Generating mocks"
	@mockery --config ${ROOT_DIR}/.mockery.yml --all --quiet
