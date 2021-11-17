DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps
	@mkdir -p bin
	go build -o bin/nrped
	@bash --norc -i ./scripts/build.sh

deps:
	go mod tidy

test: deps
	go test ./...

clean:
	rm -rf bin
