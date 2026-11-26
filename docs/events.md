# Events

The contract for publishing events to NATS: every service in the stack publishes what happens to it, and everything is kept. The JetStream stream `EVENTS` is the system of record; any database built from it is a projection that can be dropped and rebuilt by replaying the stream. The schema is `event.schema.json`, with generated bindings in `event.go` and `ts/event.ts`. Background and phasing: mariusw/homelab#6.

## The envelope

An event is a [CloudEvents 1.0](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md) event in structured JSON mode, with a stricter profile.

| Attribute | Required | Meaning |
|---|---|---|
| `specversion` | yes | Always `"1.0"`. |
| `id` | yes | Globally unique event id and idempotency key. See "Identity and duplicates". |
| `source` | yes | The component that published the event: `heimdall`, `urd`, `eir`. |
| `subject` | yes | The service the event is about. Equal to `source` when a service reports on itself; different when a collector reports on something else (`source: heimdall`, `subject: exo`). |
| `type` | yes | What happened, as dotted lower-case segments ending in a past-tense verb: `request.served`, `status.changed`, `metric.observed`. |
| `time` | yes | When it happened, not when it was published. RFC 3339, UTC, upper-case `T` and `Z`, millisecond precision. |
| `dataversion` | yes | Version of the `data` shape for this `type`, starting at 1. |
| `correlationid` | no | Joins related events across services: an agent session id, an upstream request id. |
| `host` | no | The machine the event was produced on, by its short role name. |
| `data` | yes | The payload, an object. See "Payload rules". |

`source` and `subject` match `^[a-z][a-z0-9-]*$`. Each `type` segment matches `^[a-z][a-z0-9_]*$`, and there are at least two.

### Where this profile departs from plain CloudEvents

- `time`, `subject` and `dataversion` are required; CloudEvents makes the first two optional and has no `dataversion`.
- `source` is a bare token, which is a valid relative URI-reference. CloudEvents recommends an absolute URI; a token is what lets `source` be a subject token and a NATS permission boundary.
- `type` has no reverse-DNS prefix. CloudEvents recommends one for global uniqueness; here `source` plus `type` is already unique, and the prefix would repeat in every subject.
- `id` is unique across all sources, not only within one. That lets a single JetStream duplicate window cover the whole stream.
- `dataversion`, `correlationid` and `host` are CloudEvents extension attributes, named to its rules (lower-case letters and digits only).
- No other top-level attributes are allowed. The schema says so; the generated Go bindings do not enforce it and silently drop unknown attributes when decoding. New attributes are added here and in the schema first.

## Subjects

```
events.<source>.<subject>.<type>
```

`type`'s dots become further tokens. Build subjects with `EventSubject` (Go) or `eventSubject` (TypeScript), never by hand: both reject a value that would add a token or a wildcard.

| Subject | Meaning |
|---|---|
| `events.heimdall.exo.status.changed` | Heimdall saw exo change status. |
| `events.eir.eir.session.usage.recorded` | Eir reports its own usage. |

Useful filters:

| Filter | Selects |
|---|---|
| `events.>` | Everything. |
| `events.heimdall.>` | Everything Heimdall published. |
| `events.*.exo.>` | Everything about exo, from any source. |
| `events.*.*.status.changed` | Every status change. |

Each publisher has its own NATS user, permitted to publish only to `events.<source>.>`. The broker therefore guarantees that `source` is who it claims to be. A collector that reports on third-party services publishes under its own source and names the service in `subject`.

## Identity and duplicates

- When an event has a natural identity, `id` is a hash of the fields that define it (for example source, subject, type and the upstream request id), so a retried publish yields the same id.
- Otherwise `id` is a UUIDv7, which also sorts by time.
- Publish through JetStream and wait for the acknowledgement, never with a plain core-NATS publish. Set the `Nats-Msg-Id` header to `id` (`nats.MsgId(e.Id)` with the legacy JetStream context, `jetstream.WithMsgID(e.Id)` with the `jetstream` package) so the stream drops duplicates inside its window.
- Consumers still deduplicate by `id`. The window is short, and replays deliver everything again.

## Payload rules

- **Metadata only.** Never prompt or response bodies, message content, credentials, tokens or keys. Content is already persisted elsewhere.
- Keys are `snake_case`.
- Quantities carry their unit as a suffix: `latency_ms`, `size_bytes`, `cpu_pct`, `duration_s`. Counts are plain integers named for what they count: `input_tokens`, `retries`.
- Timestamps inside `data` follow the same format as `time`.
- Enumerated values are lower-case strings: `up`, `degraded`.
- Omit a field that has no value; do not send `null`.
- Bump `dataversion` when a field is removed, renamed, or changes type, unit or meaning. Adding a field does not need a bump; consumers ignore fields they do not know.

## The stream

`EVENTS` is created and configured by the broker's administrator (see `docs/nats.md` in the homelab repo), never by a publisher. Publishers that call `AddStream` with their own config, as Urd does for its streams, would fail or fight over it.

| Setting | Value | Why |
|---|---|---|
| Subjects | `events.>` | One stream for everything. |
| Storage | file | Durable across restarts. |
| Retention | limits, no maximum age | Keep everything. |
| Maximum size | sized to the disk and the backup destination | The broker's file store is limited to 10 GB today. |
| Discard policy | new | When full, reject new events, which publishers notice, instead of silently deleting history. Alert well before the limit. |
| Duplicate window | 2 minutes | Covers publisher retries. |

## Event catalogue

Register a type here before its first publish. An entry names its source, subject, type, `dataversion`, and the fields of `data`.

Entries marked draft can still change without a `dataversion` bump until their first producer ships.

### `heimdall` · any service · `status.changed` · v1 (draft)

Heimdall observed a service's status change.

| Field | Type | Meaning |
|---|---|---|
| `from` | string | Previous status: `up`, `degraded`, `down`, `unknown`. |
| `to` | string | New status, same values. |
| `reason` | string | Why, when Heimdall has a reason. |

### `heimdall` · any service · `metric.observed` · v1 (draft)

Heimdall observed one metric value. The `key` matches Heimdall's metric registry. This is the one deliberate exception to the unit-suffix rule: a generic sample carries many units, so the unit is a field.

| Field | Type | Meaning |
|---|---|---|
| `key` | string | Metric key, e.g. `cliproxyapi.accounts.usable`. |
| `value` | number | The value, in `unit`. |
| `unit` | string | Unit as shown in the registry, e.g. `ms`, `%`, `h`. Omitted for dimensionless values. |
| `window` | string | Aggregation window for rates and percentiles, e.g. `5m`. Omitted for gauges. |

## Example

```json
{
  "specversion": "1.0",
  "id": "0193f4c2-7a10-7c3e-9b1d-5e2f8a6c4d01",
  "source": "heimdall",
  "subject": "exo",
  "type": "status.changed",
  "time": "2026-11-26T19:58:31.000Z",
  "dataversion": 1,
  "data": { "from": "up", "to": "degraded", "reason": "2 of 3 nodes ready" }
}
```

Published on `events.heimdall.exo.status.changed` with `Nats-Msg-Id: 0193f4c2-7a10-7c3e-9b1d-5e2f8a6c4d01`.
