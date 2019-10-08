gofiles := $(filter-out ./vendor/%, $(shell find . -iname "*.go" -type f))

.PHONY: all
all: build

.PHONY: build
build: dist/$(shell basename "$$PWD")-$(shell go env "GOOS")-$(shell go env "GOARCH")

.PHONY: build-all
build-all:
	# $(MAKE) -s build-android-arm
	# $(MAKE) -s build-android-arm64
	$(MAKE) -s build-linux-386
	$(MAKE) -s build-linux-amd64
	$(MAKE) -s build-linux-arm
	$(MAKE) -s build-linux-arm64

build-%:
	@\
		GOOS="$(shell echo "$*" | cut -d "-" -f 1)" \
		GOARCH="$(shell echo "$*" | cut -sd "-" -f 2)" \
		$(MAKE) -s build

.PHONY: clean
clean:
	rm -rf dist/

dist/%: $(shell go list -f "{{range .GoFiles}}{{.}} {{end}}")
	go build -o "dist/$*" .

# Development

coverage_file := coverage.txt

.PHONY: ca
ca:
	golangci-lint run

.PHONY: ci
ci: clean-dev test-race lint ca coverage build-all

.PHONY: clean-dev
clean-dev: clean
	rm -rf $(coverage_file)

.PHONY: coverage
coverage: $(coverage_file)
	go tool cover -func $<

.PHONY: coverage-web
coverage-web: $(coverage_file)
	go tool cover -html $<

.PHONY: format
format:
	gofmt -s -w -l $(gofiles)

.PHONY: lint
lint:
	gofmt -d -e -s $(gofiles)

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race -v ./...

$(coverage_file): $(shell go list -f "{{range .GoFiles}}{{.}} {{end}}") $(shell go list -f "{{range .TestGoFiles}}{{.}} {{end}}")
	go test -coverprofile $(coverage_file) ./...

