gofiles := $(filter-out ./vendor/%, $(shell find . -iname "*.go" -type f))
gosrcfiles := $(filter-out %_test.go, $(gofiles))

.PHONY: all
all: build

.PHONY: build
build: dist/$(shell basename "$$PWD")-$(shell go env "GOOS")-$(shell go env "GOARCH")

.PHONY: build-all
build-all:
	$(MAKE) -s build-android-arm

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

.PHONY: ci
ci: clean all

