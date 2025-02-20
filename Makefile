## Top level:

.PHONY: help
help: display-help ## Display this help

.PHONY: all
all: install lint test coverage build run ## Install, lint, test, coverage, build and run the project

.PHONY: install
install: install-env install-tools install-git-hooks ## Install project

.PHONY: lint ## Lint the project
lint: lint-go

.PHONY: test ## Run tests
test: unit-tests

.PHONY: build
build: dc-build ## Build docker image

.PHONY: run
run: dc-up ## Run the project

.PHONY: stop
stop: dc-stop ## Stop the project

.PHONY: clean
clean: dc-down remove-mocks remove-generated lint-go-clean ## Clean the project

-include make/*.mk
