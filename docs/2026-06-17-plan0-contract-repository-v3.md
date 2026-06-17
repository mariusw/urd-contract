# Plan 0 — Contract Repository

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish the `mariusw/urd-contract` repository as the single source of truth for the Canonical Message wire format. All other plans depend on this one. Nothing in ratatoskr or vault defines the Message type locally.

**Architecture:** One JSON Schema file defines the format. A Makefile drives codegen for Go and TypeScript. Generated files are checked in. CI validates that committed generated output always matches the schema.

**Prerequisite:** None. This plan must be complete before Plan 1 can add the contract as a Go dependency.

---

## File Map

| File | Responsibility |
|------|----------------|
| `message.schema.json` | Canonical JSON Schema definition of Message and constants |
| `generated/go/message.go` | Generated Go struct — do not edit by hand |
| `generated/go/constants.go` | Handwritten Go constants (NATSSubject, DetachedBranch, typed harness/role values) |
| `generated/ts/message.ts` | Generated TypeScript interface — do not edit by hand |
| `generated/ts/constants.ts` | Handwritten TypeScript constants (NATS_SUBJECT, DETACHED_BRANCH) |
| `go.mod` | Module declaration: `github.com/mariusw/urd-contract` |
| `Makefile` | `make generate` runs both codegen tools and overwrites generated files |
| `.github/workflows/validate.yml` | CI: regenerate and fail if output differs from committed files |

---

## Task 1: Repository initialisation

- [ ] **Step 1: Initialise git and Go module**

```bash
cd ~/localcode/github.com/mariusw/urd-contract
git init
go mod init github.com/mariusw/urd-contract
```

Expected: `go.mod` created with correct module path.

- [ ] **Step 2: Create directory structure**

```bash
mkdir -p generated/go generated/ts
```

- [ ] **Step 3: Commit scaffold**

```bash
git add go.mod
git commit -m "chore: initialise contract module"
```

---

## Task 2: JSON Schema definition

**Files:**
- Create: `message.schema.json`

The schema is the authoritative definition of every field in the Canonical Message. All fields from CONTEXT.md must be present.

- [ ] **Step 1: Write the schema**

Create `message.schema.json`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Message",
  "description": "Canonical event published by every Ratatoskr bridge.",
  "type": "object",
  "required": [
    "id",
    "session_id",
    "harness",
    "cwd",
    "git_branch",
    "role",
    "content",
    "timestamp",
    "sequence"
  ],
  "properties": {
    "id": {
      "type": "string",
      "description": "Idempotency key. Consumers must treat this as the deduplication key."
    },
    "session_id": {
      "type": "string"
    },
    "parent_session_id": {
      "type": "string",
      "description": "Set only for subagent sessions. Omitted (not null) when absent."
    },
    "harness": {
      "type": "string",
      "enum": ["claude-code", "pi"]
    },
    "cwd": {
      "type": "string"
    },
    "git_branch": {
      "type": "string",
      "description": "Always present on the wire. Detached HEAD is encoded as the sentinel '(detached)'."
    },
    "project": {
      "type": "string",
      "description": "Derived from cwd at bridge time. May be empty for sessions outside a git repo."
    },
    "role": {
      "type": "string",
      "enum": ["user", "assistant"]
    },
    "content": {
      "type": "string"
    },
    "timestamp": {
      "type": "string",
      "format": "date-time"
    },
    "sequence": {
      "type": "integer",
      "description": "Per-session message order. Assigned by the bridge."
    }
  },
  "additionalProperties": false
}
```

- [ ] **Step 2: Verify it is valid JSON**

```bash
python3 -m json.tool message.schema.json
```

- [ ] **Step 3: Commit**

```bash
git add message.schema.json
git commit -m "feat: add canonical Message JSON Schema"
```

---

## Task 3: Go types

**Files:**
- Create: `generated/go/message.go`
- Create: `generated/go/constants.go`

**Note on codegen tooling:** `go-jsonschema` (`github.com/atombender/go-jsonschema`) can generate Go structs from JSON Schema. If the generated output requires significant manual post-processing to add correct JSON tags or typed constants, prefer writing `message.go` by hand (clearly marked as generated-from-schema) and enforce correctness via a CI field-parity check rather than fighting the codegen tool.

- [ ] **Step 1: Attempt codegen**

```bash
go install github.com/atombender/go-jsonschema/cmd/go-jsonschema@latest
go-jsonschema --package contract message.schema.json -o generated/go/message.go
```

Inspect the output. If the struct fields and JSON tags match the schema cleanly, use it. If significant manual correction is needed, write `message.go` by hand instead and add a `// Code generated from message.schema.json. Do not edit.` header.

- [ ] **Step 2: Ensure the Go struct matches the schema exactly**

Required Go struct shape regardless of whether it was generated or hand-written:

```go
// Code generated from message.schema.json. Do not edit.

package contract

import "time"

type Message struct {
    ID              string    `json:"id"`
    SessionID       string    `json:"session_id"`
    ParentSessionID string    `json:"parent_session_id,omitempty"`
    Harness         string    `json:"harness"`
    CWD             string    `json:"cwd"`
    GitBranch       string    `json:"git_branch"`
    Project         string    `json:"project,omitempty"`
    Role            string    `json:"role"`
    Content         string    `json:"content"`
    Timestamp       time.Time `json:"timestamp"`
    Sequence        int       `json:"sequence"`
}
```

- [ ] **Step 3: Write typed constants**

Create `generated/go/constants.go`:

```go
package contract

const (
    NATSSubject    = "ratatoskr.messages"
    DetachedBranch = "(detached)"

    HarnessClaudeCode = "claude-code"
    HarnessPi         = "pi"

    RoleUser      = "user"
    RoleAssistant = "assistant"
)
```

- [ ] **Step 4: Verify**

```bash
go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add generated/go/
git commit -m "feat: add generated Go types and constants"
```

---

## Task 4: TypeScript types

**Files:**
- Create: `generated/ts/message.ts`
- Create: `generated/ts/constants.ts`

- [ ] **Step 1: Install codegen**

```bash
npm init -y
npm install --save-dev json-schema-to-typescript typescript
```

- [ ] **Step 2: Generate the TypeScript interface**

```bash
npx json-schema-to-typescript message.schema.json -o generated/ts/message.ts
```

- [ ] **Step 3: Write constants**

Create `generated/ts/constants.ts`:

```typescript
export const NATS_SUBJECT = "ratatoskr.messages";
export const DETACHED_BRANCH = "(detached)";

export const HARNESS_CLAUDE_CODE = "claude-code";
export const HARNESS_PI = "pi";

export const ROLE_USER = "user";
export const ROLE_ASSISTANT = "assistant";
```

- [ ] **Step 4: Verify**

```bash
npx tsc --noEmit --strict generated/ts/message.ts generated/ts/constants.ts
```

- [ ] **Step 5: Commit**

```bash
git add generated/ts/ package.json package-lock.json
git commit -m "feat: add generated TypeScript types and constants"
```

---

## Task 5: Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Write the Makefile**

```makefile
.PHONY: generate verify

generate:
	go-jsonschema --package contract message.schema.json -o generated/go/message.go || \
	  echo "go-jsonschema failed; update generated/go/message.go manually"
	npx json-schema-to-typescript message.schema.json -o generated/ts/message.ts

verify:
	@git diff --exit-code generated/ || (echo "generated files are out of sync with schema; run make generate" && exit 1)
```

- [ ] **Step 2: Commit**

```bash
git add Makefile
git commit -m "chore: add generate and verify Makefile targets"
```

---

## Task 6: CI validation

**Files:**
- Create: `.github/workflows/validate.yml`

- [ ] **Step 1: Write the workflow**

```yaml
name: validate

on: [push, pull_request]

jobs:
  generated-in-sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - run: npm ci
      - run: go install github.com/atombender/go-jsonschema/cmd/go-jsonschema@latest
      - run: make generate
      - run: make verify
```

- [ ] **Step 2: Commit**

```bash
git add .github/
git commit -m "ci: validate generated files match schema on every PR"
```

---

## Self-Review

- [ ] `message.schema.json` contains all fields from CONTEXT.md Canonical Message
- [ ] `git_branch` is in `required` (never omitted from wire)
- [ ] `parent_session_id` is NOT in `required` (optional for non-subagent sessions)
- [ ] Go struct JSON tags match schema field names exactly
- [ ] TypeScript interface compiles without errors
- [ ] `make generate && make verify` exits 0
- [ ] CI workflow runs and passes
