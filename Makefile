.PHONY: build
build: clean ## Build executable.
	go build -v -o main ./cmd/app

.PHONY: test
test: ## Run unit tests.
	go test -v -race -timeout 30s ./..

.PHONY: clean
clean: ## Cleanup.
	rm -f main

.PHONY: run
run: build ## Run executable.
	./main

.PHONY: help
help: ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help