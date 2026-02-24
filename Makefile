.PHONY: build test run lint clean build-web build-go ensure-dist

build-web:
	cd web && npm install && npm run build

# Ensure web/dist exists with at least a placeholder (for go build without frontend)
ensure-dist:
	@mkdir -p web/dist
	@test -f web/dist/index.html || echo '<!doctype html><html><body>Run make build to include frontend</body></html>' > web/dist/index.html

build: build-web
	go build -o ghist .

build-go: ensure-dist
	go build -o ghist .

test: ensure-dist
	go test -race ./...

run: ensure-dist
	go run main.go

lint:
	go vet ./...

clean:
	rm -f ghist
	rm -rf web/dist web/node_modules

release: build-web
	goreleaser release --snapshot --clean
