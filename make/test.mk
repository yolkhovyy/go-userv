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
	@$(MAKE) dc-build-up
	@./make/scripts/integration-test.sh
	@$(MAKE) dc-stop

.PHONY: remove-generated
remove-generated: ## Remove generated folders and files
	@echo "🗑 Removing generated folders and files"
	@rm -rf coverage/ results/ docs/dep-*.*

## Coverage:

.PHONY: gocovmerge-check
gocovmerge-check:
ifeq (, $(shell command -v gocovmerge 2> /dev/null))
	@echo "❌ gocovmerge not installed, run 'make install'"
	@exit 1
endif

.PHONY: gocover-cobertura-check
gocover-cobertura-check:
ifeq (, $(shell command -v gocover-cobertura 2> /dev/null))
	@echo "❌ gocover-cobertura not installed, run 'make install'"
	@exit 1
endif

.PHONY: coverage
coverage: gocovmerge-check gocover-cobertura-check ## Make coverage report
	@if [ ! -d "coverage/" ]; then \
		${MAKE} test; \
	fi
	@for file in coverage/*.cov; do \
		cp $$file $$file.tmp; \
		cat $$file.tmp | grep -v -e "mock_" -v -e "test" -v -e "pb.go" > $$file; \
	done
	@gocovmerge coverage/*.cov > coverage/total.cov
	@go tool cover -html=coverage/total.cov -o coverage/total.html
	@go tool cover -func coverage/total.cov > coverage/total.txt
	@gocover-cobertura < coverage/total.cov > coverage/total.xml
