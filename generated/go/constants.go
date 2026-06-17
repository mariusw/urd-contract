package contract

// These constants are not expressible in the Message schema (they are not
// message fields), so they are maintained here by hand. The harness and role
// values live in message.go as the generated MessageHarness / MessageRole
// typed enums.
const (
	NATSSubject    = "ratatoskr.messages"
	DetachedBranch = "(detached)"
)
