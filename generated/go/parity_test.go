package contract

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// schema mirrors the parts of message.schema.json this test asserts against.
type schema struct {
	Required   []string                   `json:"required"`
	Properties map[string]json.RawMessage `json:"properties"`
}

func loadSchema(t *testing.T) schema {
	t.Helper()
	path := filepath.Join("..", "..", "message.schema.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var s schema
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	return s
}

// structTags maps each Message field's JSON name to whether it carries omitempty.
func structTags(t *testing.T) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	rt := reflect.TypeOf(Message{})
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			t.Fatalf("field %s has no json tag", rt.Field(i).Name)
		}
		parts := strings.Split(tag, ",")
		name := parts[0]
		omitempty := false
		for _, p := range parts[1:] {
			if p == "omitempty" {
				omitempty = true
			}
		}
		out[name] = omitempty
	}
	return out
}

// TestFieldParity asserts the Message struct has exactly the schema's properties.
func TestFieldParity(t *testing.T) {
	s := loadSchema(t)
	tags := structTags(t)

	for name := range s.Properties {
		if _, ok := tags[name]; !ok {
			t.Errorf("schema property %q missing from Message struct", name)
		}
	}
	for name := range tags {
		if _, ok := s.Properties[name]; !ok {
			t.Errorf("struct field %q has no corresponding schema property", name)
		}
	}
}

// TestRequiredOmitempty enforces the wire contract: required fields must always
// be present (no omitempty), optional fields must be omitted when absent (omitempty).
func TestRequiredOmitempty(t *testing.T) {
	s := loadSchema(t)
	tags := structTags(t)

	required := map[string]bool{}
	for _, r := range s.Required {
		required[r] = true
	}

	for name, omitempty := range tags {
		if required[name] && omitempty {
			t.Errorf("required field %q must not have omitempty", name)
		}
		if !required[name] && !omitempty {
			t.Errorf("optional field %q must have omitempty", name)
		}
	}
}
