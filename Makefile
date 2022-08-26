# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

platform="linux/amd64,linux/arm64"
docker: ## build multi-platforms and push docker image
	docker buildx build --platform=$(platform) --push --tag=witjem/feedpls:main --progress=plain .
.PHONY: docker

tests: ## run go tests
	go test ./...
.PHONY: tests

tests-e2e: ## run end-to-end  tests
	./test/run.sh
.PHONY: tests-e2e

tests-all: tests tests-e2e ## run end-to-end tests
.PHONY: tests-all

lint: ## check by golangci-lint linters
	golangci-lint run
.PHONY: lint
