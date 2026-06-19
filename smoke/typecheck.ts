import type { Message } from "urd-contract/message";
import { DETACHED_BRANCH, HARNESS_PI, NATS_SUBJECT, ROLE_USER } from "urd-contract/constants";

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

export { _msg };
