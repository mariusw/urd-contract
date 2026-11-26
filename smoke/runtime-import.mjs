import { readFileSync } from "node:fs";
import {
  DETACHED_BRANCH,
  EVENT_TOKEN_PATTERN,
  EVENT_TYPE_PATTERN,
  HARNESS_PI,
  NATS_SUBJECT,
  eventSubject,
} from "urd-contract/constants";

// The TypeScript patterns are copies of the schema's; keep them identical.
const schema = JSON.parse(readFileSync(new URL("../event.schema.json", import.meta.url), "utf8"));
for (const [property, pattern] of [
  ["source", EVENT_TOKEN_PATTERN],
  ["subject", EVENT_TOKEN_PATTERN],
  ["type", EVENT_TYPE_PATTERN],
]) {
  if (pattern.source !== schema.properties[property].pattern) {
    throw new Error(`pattern for ${property}: ts/constants.ts has ${pattern.source}, event.schema.json has ${schema.properties[property].pattern}`);
  }
}

if (NATS_SUBJECT !== "ratatoskr.messages") {
  throw new Error(`unexpected NATS_SUBJECT: ${NATS_SUBJECT}`);
}
if (DETACHED_BRANCH !== "(detached)") {
  throw new Error(`unexpected DETACHED_BRANCH: ${DETACHED_BRANCH}`);
}
if (HARNESS_PI !== "pi") {
  throw new Error(`unexpected HARNESS_PI: ${HARNESS_PI}`);
}

const subject = eventSubject({ source: "heimdall", subject: "exo", type: "status.changed" });
if (subject !== "events.heimdall.exo.status.changed") {
  throw new Error(`unexpected eventSubject: ${subject}`);
}
let rejected = false;
try {
  eventSubject({ source: "heim.dall", subject: "exo", type: "status.changed" });
} catch {
  rejected = true;
}
if (!rejected) {
  throw new Error("eventSubject accepted a source containing a dot");
}

console.log("runtime import smoke passed");
