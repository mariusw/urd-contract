package contract

import (
	"fmt"
	"regexp"
)

// Patterns mirror event.schema.json. They are repeated here because EventSubject
// must reject a value before it becomes part of a NATS subject, where a dot or a
// wildcard character would silently change the routing.
var (
	eventTokenPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	eventTypePattern  = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`)
)

// EventSubject returns the NATS subject an event is published on:
// events.<source>.<subject>.<type>. It is the only place that subject is built.
func EventSubject(e Event) (string, error) {
	if !eventTokenPattern.MatchString(e.Source) {
		return "", fmt.Errorf("event source %q: must match %s", e.Source, eventTokenPattern)
	}
	if !eventTokenPattern.MatchString(e.Subject) {
		return "", fmt.Errorf("event subject %q: must match %s", e.Subject, eventTokenPattern)
	}
	if !eventTypePattern.MatchString(e.Type) {
		return "", fmt.Errorf("event type %q: must match %s", e.Type, eventTypePattern)
	}
	return EventsSubjectRoot + "." + e.Source + "." + e.Subject + "." + e.Type, nil
}
