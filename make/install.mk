## Install:

.PHONY: install-tools
install-tools: ## Install required project tools
	@echo "🛠️ Installing tools"
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/jstemmer/go-junit-report@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/boumenot/gocover-cobertura@latest
	@go install github.com/wadey/gocovmerge@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest


.PHONY: install-git-hooks
install-git-hooks: ## Install git hooks
	@if [ -d ".git/hooks" ]; then \
		echo "🛠️ Installing git hooks"; \
		cp make/pre-commit.sh .git/hooks/pre-commit; \
		cp make/pre-push.sh .git/hooks/pre-push; \
	fi

.PHONY: install-env ## Install .env
install-env:
	@echo "🛠️ Installing .env"
	@cp .env.local .env
