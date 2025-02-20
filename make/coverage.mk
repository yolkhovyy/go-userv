## Coverage:

.PHONY: gocovmerge-check
gocovmerge-check:
ifeq (, $(shell command -v gocovmerge 2> /dev/null))
	@echo "❗ gocovmerge not installed, run 'make install'"
	@exit 1
endif

.PHONY: gocover-cobertura-check
gocover-cobertura-check:
ifeq (, $(shell command -v gocover-cobertura 2> /dev/null))
	@echo "❗ gocover-cobertura not installed, run 'make install'"
	@exit 1
endif

.PHONY: coverage
coverage: gocovmerge-check gocover-cobertura-check ## Make coverage report
	@if [ ! -d "coverage/" ]; then \
		${MAKE} test; \
	fi
	@for file in coverage/*.cov; do \
		bname=$$(basename $$file .cov); \
		name=$$bname.out; \
		cat $$file | grep -v -e "mock_" -v -e "test" -v -e "pb.go" > coverage/$$name; \
	done
	@gocovmerge coverage/*.out > coverage/system.out
	@go tool cover -html=coverage/system.out -o coverage/system.html
	@go tool cover -func coverage/system.out > coverage/system.txt
	@gocover-cobertura < coverage/system.out > coverage/system.xml
