export const NATS_SUBJECT = "ratatoskr.messages";
export const DETACHED_BRANCH = "(detached)";
export const HARNESS_CLAUDE_CODE = "claude-code";
export const HARNESS_PI = "pi";
export const ROLE_USER = "user";
export const ROLE_ASSISTANT = "assistant";
// Event constants. Build a subject with eventSubject rather than by hand.
export const EVENTS_SUBJECT_ROOT = "events";
export const EVENTS_STREAM = "EVENTS";
export const EVENT_SPEC_VERSION = "1.0";
// Mirror event.schema.json. A dot or wildcard in a token would change routing.
export const EVENT_TOKEN_PATTERN = /^[a-z][a-z0-9-]*$/;
export const EVENT_TYPE_PATTERN = /^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$/;
/** Returns the NATS subject for an event: events.<source>.<subject>.<type>. Throws on values the schema rejects. */
export function eventSubject(event) {
    if (!EVENT_TOKEN_PATTERN.test(event.source))
        throw new Error(`event source ${JSON.stringify(event.source)}: must match ${EVENT_TOKEN_PATTERN}`);
    if (!EVENT_TOKEN_PATTERN.test(event.subject))
        throw new Error(`event subject ${JSON.stringify(event.subject)}: must match ${EVENT_TOKEN_PATTERN}`);
    if (!EVENT_TYPE_PATTERN.test(event.type))
        throw new Error(`event type ${JSON.stringify(event.type)}: must match ${EVENT_TYPE_PATTERN}`);
    return `${EVENTS_SUBJECT_ROOT}.${event.source}.${event.subject}.${event.type}`;
}
