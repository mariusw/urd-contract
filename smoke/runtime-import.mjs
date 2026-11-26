import { DETACHED_BRANCH, HARNESS_PI, NATS_SUBJECT, eventSubject } from "urd-contract/constants";

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
