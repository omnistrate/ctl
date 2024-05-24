GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
TAG?=latest

.GIT_COMMIT=$(shell git rev-parse HEAD)
.GIT_VERSION=$(shell git describe --tags 2>/dev/null || echo "$(.GIT_COMMIT)")
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	.GIT_COMMIT := $(.GIT_COMMIT)-dirty
endif

export GOFLAGS=-mod=vendor
GIT_USER?=$(shell gh api user -q ".login") # gets current user using github cli if the variable is not already set
GIT_TOKEN?=$(shell gh config get -h github.com oauth_token) # gets current user using github cli if the variable is not already set
PROJECT_NAME=omnistrate-ctl
DOCKER_PLATFORM=linux/arm64
TESTCOVERAGE_THRESHOLD=0
REPO_ROOT=$(shell git rev-parse --show-toplevel)

# Build info
BUILD_INFO_PKG=github.com/omnistrate/ctl/build
BUILD_VERSION=0.5
BUILD_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIMESTAMP=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS=-ldflags "-X $(BUILD_INFO_PKG).CommitID=$(BUILD_COMMIT) -X $(BUILD_INFO_PKG).Timestamp=$(BUILD_TIMESTAMP) -X $(BUILD_INFO_PKG).Version=$(BUILD_VERSION)"
CGO_ENABLED=0
GOPRIVATE=github.com/omnistrate

export ROOT_DOMAIN=omnistrate.dev

.PHONY: all
all: tidy build unit-test lint sec

.PHONY: tidy
tidy:
	echo "Tidy dependency modules"
	go mod tidy
	go mod vendor

.PHONY: unit-test
unit-test:
	echo "Running unit tests for service"
	go test ./... -skip ./test/... $(ARGS) -cover -coverprofile coverage.out -covermode count
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/[%]//g' | awk 'current=$$1; {if (current < ${TESTCOVERAGE_THRESHOLD}) {print "\033[31mTest coverage is " current " which is below threshold\033[0m"; exit 1} else {print "\033[32mTest coverage is above threshold\033[0m"}}'

.PHONY: smoke-test
smoke-test:
	echo "Running smoke tests for service"
	export ENABLE_SMOKE_TEST=true && \
	export ROOT_DOMAIN=omnistrate.dev && \
	export LOG_LEVEL=debug && \
	export LOG_FORMAT=pretty && \
	export SMOKE_TEST_EMAIL=xzhang+customer-hosted@omnistrate.com && \
	go clean -testcache && \
	go test ./... -skip ./test/... $(ARGS) 

.PHONY: build
build:
	echo "Building go binaries for omnistrate cli"
	go build -mod=mod ${BUILD_FLAGS} -o omnistrate-ctl main.go

.PHONY: test-coverage-report
test-coverage-report:
	go test ./... -skip ./test/... -cover -coverprofile coverage.out -covermode count
	go tool cover -html=coverage.out

lint-install:
	echo "Installing golangci-lint"
	brew install golangci-lint
	brew upgrade golangci-lint

.PHONY: lint
lint:
	echo "Running checks for service"
	golangci-lint run ./...

.PHONY: sec-install
sec-install:
	echo "Installing gosec"
	go install github.com/securego/gosec/v2/cmd/gosec@v2.18.0
.PHONY: sec
sec:
	echo "Security scanning for service"
	gosec -tests --quiet ./...
.PHONY: sec-verbose
sec-verbose:
	echo "Security scanning for service"
	gosec -tests ./...

.PHONY: update-dependencies
update-dependencies:
	echo "Updating dependencies"
	go get -t -u ./...
	go mod tidy

.PHONY: docker
docker: docker-build
.PHONY: docker-build
docker-build:
	docker build --platform=${DOCKER_PLATFORM} --build-arg GIT_USER=${GIT_USER} --build-arg GIT_TOKEN=${GIT_TOKEN} -f ./build/Dockerfile  -t ${PROJECT_NAME}:latest .

.PHONY: docker-run
docker-run:
	docker run --platform=${DOCKER_PLATFORM} -t ${PROJECT_NAME}:latest

# Other
.PHONY: clean
clean:
	echo "Cleaning up"
	rm ./omnistrate-cli
	rm ./omnistrate-ctl
	rm ./coverage.out
	rm ./coverage-report.html
	rm ./coverage-report.txt
	rm ./test-report.json
	rm ./security-report.html