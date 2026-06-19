import { DETACHED_BRANCH, HARNESS_PI, NATS_SUBJECT } from "urd-contract/constants";

if (NATS_SUBJECT !== "ratatoskr.messages") {
  throw new Error(`unexpected NATS_SUBJECT: ${NATS_SUBJECT}`);
}
if (DETACHED_BRANCH !== "(detached)") {
  throw new Error(`unexpected DETACHED_BRANCH: ${DETACHED_BRANCH}`);
}
if (HARNESS_PI !== "pi") {
  throw new Error(`unexpected HARNESS_PI: ${HARNESS_PI}`);
}

console.log("runtime import smoke passed");
