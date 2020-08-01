# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

all: build

.PHONY: build
build:
	@go build -o gpkg cmd/gpkg/main.go

.PHONY: clean
clean:
	@rm -f gpkg
	@rm -f coverage.txt

# ------------------------------------------------------------------------------
# Test
# ------------------------------------------------------------------------------

.PHONY: test
test: test.unit

.PHONY: test.all
test.all: test.unit test.integration

.PHONY: test.unit
test.unit:
	@go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic ./pkg/... ./internal/...

.PHONY: test.integration
test.integration:
	@go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic ./test/...
