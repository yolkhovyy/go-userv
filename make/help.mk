.PHONY: display-help
display-help:
	@echo 'Usage:'
	@echo '  $(MAKE) <target>'
	@echo 'Targets:'
	@awk \
		'BEGIN {FS = ":.*?## "} { \
			if (/^## .*$$/) {\
				printf "  %s\n", substr($$1, 4) \
			} else if (/^[a-zA-Z\-_%]+:.*?##.*$$/) {\
				printf "    %-22s%s\n", $$1, $$2 \
			} \
		}' $(MAKEFILE_LIST)
