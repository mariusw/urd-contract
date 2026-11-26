# Session ledger: event envelope contract

_Last reconciled: 2026-11-27_

## Context

First step of the stack-wide event log (mariusw/homelab#6): persist metadata from every service in one replayable JetStream stream, `EVENTS`, starting as soon as possible. The contract must exist before the first producer. Phase 1 of #6 is contract, backups (mariusw/homelab#5), NATS access control, and one long-retention stream; Heimdall republishing its observations comes after Heimdall's frontend redesign ships.

## Current state

- Done on `feat/event-envelope` (worktree `~/worktrees/urd-contract/event-envelope`, branched from `origin/main`): `event.schema.json` with generated `event.go` and `ts/event.ts`, `subject.go` (`EventSubject`) with tests, TypeScript `eventSubject` and constants, smoke checks, `docs/events.md`, BUILD-HISTORY v7 entry. `make verify` passes locally.
- Merged in mariusw/urd-contract#1 and tagged v0.3.0. Worktree and branch removed. CI follow-ups continue on `chore/ci-hardening`.

## Decisions

- **Lives in urd-contract**, chosen by Marius over a new neutral package or a document-only start: existing schema-to-bindings tooling and CI, speed first. The repo name no longer fits; renaming is deferred.
- **CloudEvents 1.0 structured JSON mode as prior art**, with a stricter profile: `time`, `subject`, `dataversion` required; `source`/`subject`/`type` restricted to NATS-safe tokens; no reverse-DNS `type` prefix; bare-token `source`; `id` globally unique; no unknown top-level attributes.
- **Subject `events.<source>.<subject>.<type>`.** `source` is the publisher (NATS user limited to `events.<source>.>`, so the broker enforces provenance); `subject` is the service the event is about, which lets Heimdall report on third-party services under its own source.
- **Payload versioning by `dataversion` per type**, not by a version in `type` or the subject. Adding a field needs no bump.
- **Metadata only** in `data`; prompt and response content is already persisted as JSONL elsewhere.
- **Stream created by the broker admin, never by publishers**; discard policy `new` so a full stream rejects rather than silently deleting history.
- Generated Go does not reject unknown attributes (same as Message); documented, not worked around.

## Next steps

1. Done: review, PR #1, merge, tag v0.3.0.
2. Done: CI hardening moved to its own branch (see docs/sessions/chore/ci-hardening.md).
3. Homelab: apply NATS access control from `docs/nats.md`, add an `events-heimdall`-style user per publisher, create the `EVENTS` stream (size decision against dozer's disk and the 10 GB file-store limit), record it in `docs/nats.md`.
4. Deferred: Heimdall publishes `status.changed` and `metric.observed` (after the redesign); service-published metric descriptors (mariusw/heimdall#9); analytical store choice.

## References

- mariusw/homelab#6 (project), mariusw/homelab#5 (backups), mariusw/heimdall#9 (descriptors)
- `docs/events.md`, `event.schema.json`, `subject.go`
- CloudEvents spec: https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md
