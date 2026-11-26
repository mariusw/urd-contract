import type { Event } from "urd-contract/event";
import type { Message } from "urd-contract/message";
import { DETACHED_BRANCH, EVENT_SPEC_VERSION, HARNESS_PI, NATS_SUBJECT, ROLE_USER } from "urd-contract/constants";

const _msg: Message = {
  id: "AAAAA",
  session_id: "sess-1",
  harness: HARNESS_PI,
  cwd: "/tmp/demo",
  git_branch: DETACHED_BRANCH,
  role: ROLE_USER,
  content: NATS_SUBJECT,
  timestamp: new Date().toISOString(),
};

const _event: Event = {
  specversion: EVENT_SPEC_VERSION,
  id: "0193f4c2-7a10-7c3e-9b1d-5e2f8a6c4d01",
  source: "heimdall",
  subject: "exo",
  type: "status.changed",
  time: new Date().toISOString(),
  dataversion: 1,
  data: { from: "up", to: "degraded" },
};

export { _event, _msg };
