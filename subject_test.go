package contract

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEventSubject(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		subject string
		typ     string
		want    string
		wantErr bool
	}{
		{name: "collector reporting on another service", source: "heimdall", subject: "exo", typ: "status.changed", want: "events.heimdall.exo.status.changed"},
		{name: "service reporting on itself", source: "eir", subject: "eir", typ: "session.usage.recorded", want: "events.eir.eir.session.usage.recorded"},
		{name: "hyphenated subject", source: "heimdall", subject: "blog-staging", typ: "metric.observed", want: "events.heimdall.blog-staging.metric.observed"},
		{name: "dot in source would add a token", source: "heim.dall", subject: "exo", typ: "status.changed", wantErr: true},
		{name: "wildcard in subject", source: "heimdall", subject: "*", typ: "status.changed", wantErr: true},
		{name: "upper case", source: "Heimdall", subject: "exo", typ: "status.changed", wantErr: true},
		{name: "single-segment type", source: "heimdall", subject: "exo", typ: "changed", wantErr: true},
		{name: "empty source", source: "", subject: "exo", typ: "status.changed", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EventSubject(Event{Source: tc.source, Subject: tc.subject, Type: tc.typ})
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got subject %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

const sampleEvent = `{
  "specversion": "1.0",
  "id": "0193f4c2-7a10-7c3e-9b1d-5e2f8a6c4d01",
  "source": "heimdall",
  "subject": "exo",
  "type": "status.changed",
  "time": "2026-11-26T19:58:31.000Z",
  "dataversion": 1,
  "data": {"from": "up", "to": "degraded", "reason": "2 of 3 nodes ready"}
}`

func TestEventUnmarshal(t *testing.T) {
	var e Event
	if err := json.Unmarshal([]byte(sampleEvent), &e); err != nil {
		t.Fatalf("valid event rejected: %v", err)
	}
	if e.Data["to"] != "degraded" {
		t.Fatalf("data not decoded: %#v", e.Data)
	}

	missingTime := strings.Replace(sampleEvent, `"time": "2026-11-26T19:58:31.000Z",`, "", 1)
	if err := json.Unmarshal([]byte(missingTime), &e); err == nil {
		t.Fatal("event without time accepted")
	}

	wrongVersion := strings.Replace(sampleEvent, `"specversion": "1.0"`, `"specversion": "0.3"`, 1)
	if err := json.Unmarshal([]byte(wrongVersion), &e); err == nil {
		t.Fatal("event with specversion 0.3 accepted")
	}

	badSource := strings.Replace(sampleEvent, `"source": "heimdall"`, `"source": "heim.dall"`, 1)
	if err := json.Unmarshal([]byte(badSource), &e); err == nil {
		t.Fatal("event with a dot in source accepted")
	}
}
