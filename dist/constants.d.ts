export declare const NATS_SUBJECT = "ratatoskr.messages";
export declare const DETACHED_BRANCH = "(detached)";
export declare const HARNESS_CLAUDE_CODE = "claude-code";
export declare const HARNESS_PI = "pi";
export declare const ROLE_USER = "user";
export declare const ROLE_ASSISTANT = "assistant";
export declare const EVENTS_SUBJECT_ROOT = "events";
export declare const EVENTS_STREAM = "EVENTS";
export declare const EVENT_SPEC_VERSION = "1.0";
export declare const EVENT_TOKEN_PATTERN: RegExp;
export declare const EVENT_TYPE_PATTERN: RegExp;
/** Returns the NATS subject for an event: events.<source>.<subject>.<type>. Throws on values the schema rejects. */
export declare function eventSubject(event: {
    source: string;
    subject: string;
    type: string;
}): string;
//# sourceMappingURL=constants.d.ts.map