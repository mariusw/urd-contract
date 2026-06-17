.PHONY: generate verify

# The TypeScript interface is generated directly from the schema.
# The Go struct (generated/go/message.go) is hand-authored from the schema
# because go-jsonschema's output diverges from the required wire shape
# (ID/CWD casing, plain string enums, no extra tags). Its correctness against
# the schema is guaranteed by the field-parity test in generated/go, not by
# regeneration. See message.schema.json Task 3 note.
generate:
	npx json-schema-to-typescript message.schema.json -o generated/ts/message.ts

verify:
	go test ./...
	@git diff --exit-code generated/ || (echo "generated files are out of sync with schema; run make generate" && exit 1)
