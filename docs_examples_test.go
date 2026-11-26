package contract

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

// jsonCodeBlock matches a fenced ```json block in Markdown.
var jsonCodeBlock = regexp.MustCompile("(?s)```json\n(.*?)\n```")

// Every JSON example in docs/events.md must be a valid event: it decodes
// against the generated contract and yields a valid subject. This keeps the
// documentation from drifting away from the schema.
func TestEventsDocExamplesAreValid(t *testing.T) {
	doc, err := os.ReadFile("docs/events.md")
	if err != nil {
		t.Fatalf("read docs/events.md: %v", err)
	}
	examples := jsonCodeBlock.FindAllSubmatch(doc, -1)
	if len(examples) == 0 {
		t.Fatal("docs/events.md has no ```json examples; the check would pass vacuously")
	}
	for i, example := range examples {
		var event Event
		if err := json.Unmarshal(example[1], &event); err != nil {
			t.Errorf("example %d is not a valid event: %v", i+1, err)
			continue
		}
		if _, err := EventSubject(event); err != nil {
			t.Errorf("example %d does not yield a valid subject: %v", i+1, err)
		}
	}
}
