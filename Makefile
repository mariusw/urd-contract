.PHONY: tools generate verify

# The Go generator is pinned. An unpinned release could fail the drift guard on
# an unrelated change, or alter the contract without a schema edit. Bump it on
# purpose, regenerate, and commit the result in the same change.
GO_JSONSCHEMA_VERSION := v0.23.1
GO_JSONSCHEMA := $(shell go env GOPATH)/bin/go-jsonschema

tools:
	go install github.com/atombender/go-jsonschema@$(GO_JSONSCHEMA_VERSION)

# Each contract's bindings are generated from its schema (message.schema.json,
# event.schema.json), which is its single source of truth. The Go generator is
# given -t so the root type takes its name from the schema's "title" (Message,
# Event) rather than the filename; no other shape is imposed. Whatever the
# generators emit IS the contract.
generate:
	@go version -m $(GO_JSONSCHEMA) | grep -q 'go-jsonschema[[:space:]]*$(GO_JSONSCHEMA_VERSION)' || (echo "go-jsonschema is not $(GO_JSONSCHEMA_VERSION); run make tools" && exit 1)
	$(GO_JSONSCHEMA) -p contract -t -o message.go message.schema.json
	npx json-schema-to-typescript message.schema.json -o ts/message.ts
	$(GO_JSONSCHEMA) -p contract -t -o event.go event.schema.json
	npx json-schema-to-typescript event.schema.json -o ts/event.ts

verify:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && echo "the files above are not gofmt-formatted; run gofmt -w" && exit 1)
	go build ./...
	go vet ./...
	go test ./...
	npx tsc --noEmit --strict ts/message.ts ts/event.ts ts/constants.ts
	npm run build:ts
	node smoke/runtime-import.mjs
	npx tsc -p tsconfig.smoke.json
	@git diff --exit-code message.go ts/message.ts event.go ts/event.ts dist/ || (echo "generated or built files are out of sync; run make generate and npm run build:ts, then commit dist/" && exit 1)
