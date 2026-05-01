go := "/usr/local/go/bin/go"

build:
    {{go}} build -o ./bin/yo ./cmd/yo

install:
    {{go}} install ./cmd/yo

test:
    {{go}} test ./internal/...

test-integration:
    {{go}} test ./internal/... -tags integration

coverage:
    {{go}} test ./internal/... -tags integration \
        -coverpkg=github.com/simskij/yo/internal/config,github.com/simskij/yo/internal/dots,github.com/simskij/yo/internal/ui \
        -coverprofile=coverage.out
    {{go}} tool cover -func=coverage.out | grep total

lint:
    {{go}} vet ./...

tidy:
    {{go}} mod tidy
