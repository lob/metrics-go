GOTOOLS := \
	github.com/alecthomas/gometalinter \
	github.com/git-chglog/git-chglog/cmd/git-chglog \
	github.com/golang/dep/cmd/dep \

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

.PHONY: release
release: ## Creates a new release with the given tag
	@echo "---> Creating new release"
ifndef tag
	$(error tag must be specified)
endif
	git-chglog --output CHANGELOG.md --next-tag $(tag)
	git add CHANGELOG.md
	git commit -m $(tag)
	git tag $(tag)
	git push origin master --tags

.PHONY: setup
setup:
	@echo "--> Installing linter"
	go get -u -v $(GOTOOLS)
	gometalinter --install

.PHONY: install
install:
	@echo "---> Installing dependencies"
	dep ensure

.PHONY: test
test:
	@echo "---> Testing"
	ENVIRONMENT=test go test ./... -race -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)
