## Dependency:

.PHONY: deps
deps: check-goda check-graphviz ## Generate dependency graph/tree
	@echo "🔍 Generating dependency graph"
	@goda graph ./... > docs/dep-graph.dot
	@dot -Tpng docs/dep-graph.dot -o docs/dep-graph.png
	@echo "🔍 Generating dependency tree"
	@goda tree ./... > docs/dep-tree.txt

.PHONY: check-graphviz
check-graphviz:
ifeq (, $(shell command -v dot 2> /dev/null))
	@echo "❗ graphviz not found, please install"
	@exit 1
endif

.PHONY: check-goda
check-goda:
ifeq (, $(shell command -v goda 2> /dev/null))
	@echo "❗ goda not found, run `make install`"
	@exit 1
endif
