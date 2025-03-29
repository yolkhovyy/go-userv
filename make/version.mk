
export GIT_TAG = $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

ifeq ($(shell git diff-index --quiet HEAD --; echo $$?), 1)
    export TIMESTAMP = $(shell date +%Y%m%d%H%M%S)
    export APP_VERSION = $(GIT_TAG)-$(TIMESTAMP)
else
    export APP_VERSION = $(GIT_TAG)
endif
