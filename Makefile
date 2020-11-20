## test: Run test
test:
	go test ./... -coverprofile cp.out && go tool cover -func=cp.out

## coverage: Run test with html coverage report
coverage:
	go test ./... -coverprofile cp.out && go tool cover -html=cp.out

## benchmark: Run benchmark test
benchmark:
	go test -bench=.

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
