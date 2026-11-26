package contract

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

// eventSchemaPatterns reads the regex patterns for source, subject and type
// from event.schema.json, the single source of truth they are copied from.
func eventSchemaPatterns(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("event.schema.json")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var schema struct {
		Properties map[string]struct {
			Pattern string `json:"pattern"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	patterns := map[string]string{}
	for _, name := range []string{"source", "subject", "type"} {
		pattern := schema.Properties[name].Pattern
		if pattern == "" {
			t.Fatalf("schema property %q has no pattern", name)
		}
		patterns[name] = pattern
	}
	return patterns
}

// EventSubject repeats the schema's patterns so it can reject a value before it
// becomes a subject token. This keeps the copies identical to the schema.
func TestEventSubjectPatternsMatchSchema(t *testing.T) {
	patterns := eventSchemaPatterns(t)
	checks := []struct {
		property string
		inCode   *regexp.Regexp
	}{
		{"source", eventTokenPattern},
		{"subject", eventTokenPattern},
		{"type", eventTypePattern},
	}
	for _, check := range checks {
		if check.inCode.String() != patterns[check.property] {
			t.Errorf("pattern for %s: subject.go has %q, event.schema.json has %q", check.property, check.inCode.String(), patterns[check.property])
		}
	}
}
