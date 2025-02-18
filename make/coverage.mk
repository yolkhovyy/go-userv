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
		cp $$file $$file.tmp; \
		cat $$file.tmp | grep -v -e "mock_" -v -e "test" -v -e "pb.go" > $$file; \
	done
	@gocovmerge coverage/*.cov > coverage/total.cov
	@go tool cover -html=coverage/total.cov -o coverage/total.html
	@go tool cover -func coverage/total.cov > coverage/total.txt
	@gocover-cobertura < coverage/total.cov > coverage/total.xml
