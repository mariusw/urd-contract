.PHONY: generate verify

# Each contract's bindings are generated from its schema (message.schema.json,
# event.schema.json), which is its single source of truth. The Go generator is
# given -t so the root type takes its name from the schema's "title" (Message,
# Event) rather than the filename; no other shape is imposed. Whatever the generators emit IS the contract.
generate:
	$(shell go env GOPATH)/bin/go-jsonschema -p contract -t -o message.go message.schema.json
	npx json-schema-to-typescript message.schema.json -o ts/message.ts
	$(shell go env GOPATH)/bin/go-jsonschema -p contract -t -o event.go event.schema.json
	npx json-schema-to-typescript event.schema.json -o ts/event.ts

verify:
	go build ./...
	go vet ./...
	go test ./...
	npx tsc --noEmit --strict ts/message.ts ts/event.ts ts/constants.ts
	npm run build:ts
	node smoke/runtime-import.mjs
	npx tsc -p tsconfig.smoke.json
	@git diff --exit-code message.go ts/message.ts event.go ts/event.ts dist/ || (echo "generated or built files are out of sync; run make generate and npm run build:ts, then commit dist/" && exit 1)
