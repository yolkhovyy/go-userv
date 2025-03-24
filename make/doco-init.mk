
export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1
export COMPOSE_BAKE = true

DOCO = docker compose -f docker-compose.yml
export DEPENDENCIES = otel-collector postgres zookeeper kafka kafka-initializer

# Telemetry
ifdef NR		#----- New Relic
	export OTEL_COLLECTOR_CONFIG = config-newrelic.yml
else ifdef GLT	#----- grafana loki, tempo, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-tempo.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki tempo prometheus promtail
	export OTEL_COLLECTOR_CONFIG = config-grafana-loki-tempo.yml
else			#----- grafana loki, jaeger, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-jaeger.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki jaeger prometheus promtail
	export OTEL_COLLECTOR_CONFIG = config-grafana-loki-jaeger.yml
endif

export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DOCO} config --services))
