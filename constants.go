package contract

// These constants are not expressible in the Message schema (they are not
// message fields), so they are maintained here by hand. The harness and role
// values live in message.go as the generated MessageHarness / MessageRole
// typed enums.
const (
	NATSSubject    = "ratatoskr.messages"
	DetachedBranch = "(detached)"
)

// Event constants. Build a subject with EventSubject rather than by hand.
const (
	// EventsSubjectRoot is the first token of every event subject; the EVENTS
	// stream captures EventsSubjectRoot + ".>".
	EventsSubjectRoot = "events"
	// EventsStream is the JetStream stream that holds every event. It is created
	// by the broker's administrator, never by a publisher.
	EventsStream = "EVENTS"
)
