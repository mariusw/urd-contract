.PHONY: generate verify

# Both bindings are generated from message.schema.json, which is the single
# source of truth. The Go generator is given -t so the root type takes its name
# from the schema's "title" (Message) rather than the filename; no other shape
# is imposed. Whatever the generators emit IS the contract.
generate:
	$(shell go env GOPATH)/bin/go-jsonschema -p contract -t -o generated/go/message.go message.schema.json
	npx json-schema-to-typescript message.schema.json -o generated/ts/message.ts

verify:
	go build ./...
	@git diff --exit-code generated/ || (echo "generated files are out of sync with schema; run make generate" && exit 1)
