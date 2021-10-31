GOLANGCILINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

.PHONY: build
build: ## Compiles CLI binary.
	@go build -o ./fleetlockctl main.go

.PHONY: test-working-tree-clean
test-working-tree-clean: ## Checks if working directory is clean.
	@test -z "$$(git status --porcelain)" || (echo "Commit all changes before running this target"; exit 1)

.PHONY: update-linters
update-linters: ## Updates list of linters in .golangci.yml file based on currently installed golangci-lint binary.
	# Remove all enabled linters.
	sed -i '/^  enable:/q0' .golangci.yml
	# Then add all possible linters to config.
	$(GOLANGCILINT) linters | grep -E '^\S+:' | cut -d: -f1 | sort | sed 's/^/    - /g' | grep -v -E "($$(grep '^  disable:' -A 100 .golangci.yml  | grep -E '    - \S+$$' | awk '{print $$2}' | tr \\n '|' | sed 's/|$$//g'))" >> .golangci.yml

.PHONY: test-update-linters
test-update-linters: test-working-tree-clean update-linters ## Verifies that list of linters in .golangci.yml file is up to date.
	@test -z "$$(git status --porcelain)" || (echo "Linter configuration outdated. Run 'make update-linters' and commit generated changes to fix."; exit 1)

.PHONY: help
help: ## Prints help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
