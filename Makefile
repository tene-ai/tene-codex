.PHONY: build test vet check routing-eval plugin-validate clean

build:
	go build -trimpath -o dist/tene-workflow ./cmd/tene-workflow

test:
	go test ./...

vet:
	go vet ./...

check: test vet routing-eval
	python3 -m json.tool .codex-plugin/plugin.json >/dev/null
	python3 -m json.tool hooks/hooks.json >/dev/null
	python3 -m json.tool .agents/plugins/marketplace.json >/dev/null
	for schema in schemas/*.json; do python3 -m json.tool "$$schema" >/dev/null; done
	python3 -m unittest discover -s tests -p '*_test.py'
	./scripts/release-smoke.sh
	./scripts/requirements-audit.py >/dev/null

routing-eval:
	go run ./cmd/tene-routing-eval evals/routing-corpus.json >/dev/null

clean:
	go clean
