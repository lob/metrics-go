DIRS     ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)
PKG_NAME ?= metrics-go

BFLAGS ?=
LFLAGS ?=
TFLAGS ?=

COVERAGE_PROFILE ?= coverage.out

default: test enforce

.PHONY: clean
clean:
	@echo "---> Cleaning"
	rm -rf ./vendor

.PHONY: enforce
enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: lint
lint:
	@echo "---> Linting"
	gometalinter --vendor --tests --deadline=2m $(LFLAGS) $(DIRS)

.PHONY: setup
setup:
	@echo "--> Installing linter"
	go get -u -v github.com/alecthomas/gometalinter github.com/golang/dep/cmd/dep
	gometalinter --install

.PHONY: install
install:
	@echo "---> Installing dependencies"
	dep ensure

.PHONY: test
test:
	@echo "---> Testing"
	ENVIRONMENT=test go test ./... -race -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)
