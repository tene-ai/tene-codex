.PHONY: build test vet check plugin-validate clean

build:
	go build -trimpath -o dist/tene-workflow ./cmd/tene-workflow

test:
	go test ./...

vet:
	go vet ./...

check: test vet
	python3 -m json.tool .codex-plugin/plugin.json >/dev/null
	python3 -m json.tool hooks/hooks.json >/dev/null
	for schema in schemas/*.json; do python3 -m json.tool "$$schema" >/dev/null; done
	python3 -m unittest discover -s tests -p '*_test.py'

clean:
	go clean
