# AGENTS.md

Instructions for every AI agent working in this repository. Read this file first, then the
documents it links to for the area you are touching.

## Project in one paragraph

**hashplace** is a collaborative pixel canvas (like r/place) where painting a cell costs
proof-of-work instead of a cooldown timer. A claim's strength comes from the work behind it and
decays over time, so any cell can eventually be reclaimed. The backend (Go, in `backend/`) is the
product's core and owns all rules. The frontend (React + TypeScript, in `web/`) is a thin client
that renders the canvas and mines proofs in a Web Worker.

## Non-negotiable rules

1. **Never push to `main`, never merge a PR, never approve a PR.** Only the maintainer (@nem0z)
   merges. Branch protection enforces most of this; do not try to work around it.
2. **One issue -> one branch -> one PR.** Keep PRs small and focused (target < 400 changed lines,
   excluding lockfiles and generated code). If scope grows, stop and propose a split.
3. **Do not weaken the safety net.** Never delete or skip tests, disable linters, add `//nolint`
   or `eslint-disable` without a justification comment, lower coverage, or edit CI, branch
   protection, `AGENTS.md`, `.github/agents/` or `.github/skills/` unless the issue explicitly
   asks for it.
4. **The spec is the source of truth for game rules.** `docs/spec/` defines how claims,
   proof-of-work and decay behave. Code must match it. If the code needs a rule change, change
   the spec in the same PR and call it out in the PR description.
5. **Decisions are written down.** Any architectural decision (new dependency, new package
   boundary, storage choice, protocol change) needs an ADR in `docs/adr/` (see the
   `write-adr` skill).
6. **Ask instead of guessing.** If an issue is ambiguous, ask the maintainer (in the session or as
   an issue comment) before writing code. Record the answer in the issue or PR.
7. **No secrets** in code, tests, fixtures, logs or PR descriptions.

## Repository map

```text
backend/            Go module (github.com/nem0z/hashplace/backend)
  cmd/hashplaced/   server entrypoint
  internal/         application packages (see docs/adr/0002-repository-layout.md)
web/                React + TypeScript client (Vite)
docs/
  workflow.md       how work flows from issue to merged PR
  spec/             normative game/protocol rules
  adr/              architecture decision records
.github/
  agents/           custom agent definitions (architect, backend, frontend)
  skills/           reusable procedures (open-pull-request, address-review-feedback, ...)
  instructions/     path-scoped coding standards (auto-applied by Copilot)
  workflows/        CI
```

`backend/` and `web/` are created by their scaffold issues; until then CI skips them.

## Commands

Run the checks for every area you touched before opening or updating a PR.

| Area    | Command (from repo root)                                                                 |
|---------|-------------------------------------------------------------------------------------------|
| Backend | `cd backend && gofmt -l . && go vet ./... && golangci-lint run && go test -race ./...`   |
| Web     | `cd web && npm ci && npm run lint && npm run typecheck && npm test && npm run build`     |

## Standards

- Go: [`.github/instructions/go.instructions.md`](.github/instructions/go.instructions.md)
- React/TypeScript: [`.github/instructions/web.instructions.md`](.github/instructions/web.instructions.md)
- Git, commits and PRs: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- End-to-end workflow: [`docs/workflow.md`](docs/workflow.md)

## Agents and skills

| Agent       | Use it for                                                              |
|-------------|-------------------------------------------------------------------------|
| `architect` | Turning ideas into specs, ADRs and agent-ready issues. Writes no product code. |
| `backend`   | Implementing Go issues in `backend/`.                                   |
| `frontend`  | Implementing React issues in `web/`.                                    |

| Skill                     | Use it when                                              |
|---------------------------|----------------------------------------------------------|
| `open-pull-request`       | Your change is ready to be proposed.                     |
| `address-review-feedback` | A PR you own received review comments.                   |
| `write-adr`               | You are making or proposing an architectural decision.   |
