export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1

DC = docker compose
export DEPENDENCIES = postgres zookeeper kafka kafka-initializer
export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DC} config --services))

## Docker compose:

.PHONY: dc-up-dependencies
dc-up-dependencies: ## Start project dependency containers
	@echo "🐳 Starting project dependencies in Docker containers"
	@${DC} up --detach --quiet-pull ${DEPENDENCIES}

.PHONY: dc-build
dc-build: lint ## Build Docker images for the project
	@echo "🐳 Building Docker images"
	@${DC} build ${SERVICES}

.PHONY: dc-build-up
dc-build-up: lint dc-up-dependencies ## Build Docker images and start the services
	@echo "🐳 Building and starting services in Docker containers"
	@${DC} up --build --detach --force-recreate --remove-orphans ${SERVICES}

.PHONY: dc-up
dc-up: dc-up-dependencies ## Start project services in Docker containers
	@echo "🐳 Starting project services in Docker containers"
	@${DC} up --detach --remove-orphans ${SERVICES}

.PHONY: dc-stop
dc-stop: ## Stop running Docker containers
	@echo "🐳 Stopping running Docker containers"
	@${DC} stop

.PHONY: dc-down
dc-down: ## Stop and remove Docker containers and associated network
	@echo "🗑 Stopping and removing Docker containers and associated network"
	@${DC} down

.PHONY: dc-watch
dc-watch: ## Display the status of running Docker containers
	@watch ${DC} ps

