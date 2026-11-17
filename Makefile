.PHONY: run cloc fmt

run:
	go run .

fmt:
	go fmt ./...

cloc:
	cloc . --vcs git
