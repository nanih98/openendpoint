SHELL:=/bin/sh
.PHONY: all

help: ## this help
	@awk 'BEGIN {FS = ":.*?## ";  printf "Usage:\n  make \033[36m<target> \033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

doctoc: ## Create table of contents with doctoc
	doctoc .

goreleaser: ## Generate go binaries using goreleaser (brew install goreleaser)
	goreleaser release --snapshot --rm-dist

golangci-lint: ## Lint Golang code (brew install golangci-lint)
	golangci-lint

##https://github.com/moovweb/gvm

