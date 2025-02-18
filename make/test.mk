## Tests:

.PHONY: unit-test
unit-test: generate-mocks ## Run unit tests
	@echo "⚙️ Running unit tests"
	@mkdir -p results coverage
	@go test -v -count=1 -coverpkg=./... -coverprofile=coverage/unit-test.cov ./... | tee results/unit-test.0
	@go-junit-report -set-exit-code < results/unit-test.0 > results/unit-test.xml

.PHONY: integration-test
integration-test: export BUILD_TARGET = test
integration-test: ## Run integration tests
	@$(MAKE) dc-build-up; sleep 7
	@./test/integration.sh
	@$(MAKE) dc-stop
