## Top level:

.PHONY: help
help: display-help ## Display this help

.PHONY: all
all: install lint test build run ## Install, lint, test, coverage, build and run the project

.PHONY: install
install: install-env install-tools install-git-hooks ## Install project

.PHONY: lint
lint: lint-go ## Lint the project

.PHONY: test
test: unit-test integration-test coverage ## Run tests

.PHONY: build
build: doco-build ## Build docker image

.PHONY: run
run: doco-up ## Run the project

.PHONY: stop
stop: doco-stop ## Stop the project

.PHONY: clean
clean: doco-down remove-mocks remove-generated lint-go-clean ## Clean the project

-include make/*.mk
