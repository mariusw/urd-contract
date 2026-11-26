# Session ledger: CI hardening

_Last reconciled: 2026-11-27_

## Context

Follow-up to the event envelope (v0.3.0). Closes the CI holes that work opened or exposed: unpinned generator, hand-copied patterns, unchecked doc example, no formatting check. The wider CI plan for all repositories is mariusw/homelab#7, which names this repository as the reference for the contract type.

## Current state

- Done on `chore/ci-hardening` (worktree `~/worktrees/urd-contract/ci-hardening`): pinned `go-jsonschema` v0.23.1 with `make tools` and a version check in `make generate`; `gofmt -l` in `make verify`; Go and TypeScript pattern-parity checks; doc example validation; BUILD-HISTORY v7 marked shipped and v8 added. Both parity checks were mutation-tested.
- Next: push, PR, CI, merge on Marius's word, then remove worktree and branch.

## Decisions

- Pin in the Makefile, not only in CI, so local and CI generation cannot disagree.
- Validate doc examples through the generated `Event` decoder instead of adding a JSON Schema validator dependency.
- Out of scope here, planned in homelab#7: breaking-change detection, conformance fixtures, consumer builds, tag and `package.json` version agreement, Node version in CI.

## Next steps

1. PR, CI, merge, clean up.
2. Tag only if Marius wants a release for this; it changes no contract, but `EVENT_TOKEN_PATTERN` and `EVENT_TYPE_PATTERN` are newly exported from `urd-contract/constants`.

## References

- mariusw/homelab#7 (CI plan), mariusw/homelab#6 (event log)
- `Makefile`, `schema_parity_test.go`, `docs_examples_test.go`, `smoke/runtime-import.mjs`
