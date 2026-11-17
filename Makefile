.PHONY: run lint cloc help

run: ## run main program
	go run .

lint: ## check for code quality
	go fmt ./...
	go test ./...
	go vet ./...
	go fix ./...

cloc: ## count lines of code
	cloc . --vcs git

help: ## help message
	@echo "usage:"
	@echo
	@sed -n 's/^\([A-Za-z0-9_.-]*\):.*## \(.*\)$$/\t\1: \2/p' Makefile | sort | column -t -s ':'
	@echo
