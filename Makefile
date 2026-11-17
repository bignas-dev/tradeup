.PHONY: run cloc

run:
	go run .

cloc:
	cloc . --vcs git
