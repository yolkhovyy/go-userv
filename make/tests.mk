## Tests:

.PHONY: unit-tests
unit-tests: generate-mocks ## Run unit tests
	@echo "⚙️ Running unit tests"
	@mkdir -p results coverage
	@go test -v -count=1 -coverpkg=./... -coverprofile=coverage/unit-tests.cov ./... | tee results/unit-tests.0
	@go-junit-report -set-exit-code < results/unit-tests.0 > results/unit-tests.xml
