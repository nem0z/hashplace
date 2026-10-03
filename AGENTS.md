# AGENTS.md

Instructions for every AI agent working in this repository. Read this file first, then the
documents it links to for the area you are touching.

## Project in one paragraph

**hashplace** is a collaborative pixel canvas (like r/place) where painting a cell costs
proof-of-work instead of a cooldown timer. A claim's strength comes from the work behind it and
decays over time, so any cell can eventually be reclaimed. The backend (Go, in `backend/`) is the
product's core and owns all rules. The frontend (React + TypeScript, in `frontend/`) is a thin client
that renders the canvas and mines proofs in a Web Worker.

## Core principle: simplest option first

**Simple, short and clear code is the primary goal of every coding agent.**

- Always pick the **simplest option that works**. Move to something more complex only when the
  simple option demonstrably does not work (a failing test, a measured problem) or when the
  maintainer explicitly asks for it. If you do, say why in the PR.
- Fewer lines, fewer types, fewer layers. No speculative abstractions, interfaces nobody needs
  yet, generic helpers, configuration knobs or optimizations "for later".
- Straightforward control flow: early returns, small functions, obvious names. No clever tricks.
- Reuse the standard library and existing code before adding dependencies or new patterns.
- When two options are equally simple, pick the one that is easier to read.

## Non-negotiable rules

1. **Never push to `main`, never merge a PR, never approve a PR.** Only the maintainer (@nem0z)
   merges. Branch protection enforces most of this; do not try to work around it.
2. **One issue -> one branch -> one PR.** Keep PRs small and focused (< 250 changed lines,
   excluding lockfiles and generated code). If scope grows, break it down into smaller PRs.
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
frontend/           React + TypeScript client (Vite)
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

`backend/` and `frontend/` are created by their scaffold issues; until then CI skips them.

## Commands

Run the checks for every area you touched before opening or updating a PR.

| Area     | Command (from repo root)                                                                  |
|----------|-------------------------------------------------------------------------------------------|
| Backend  | `cd backend && gofmt -l . && go vet ./... && golangci-lint run && go test -race ./...`    |
| Frontend | `cd frontend && npm ci && npm run lint && npm run typecheck && npm test && npm run build` |

## Commit messages

Every commit (and every PR title) follows this format:

```text
<type>(<scope>): <description>
```

- **type**: `feat`, `fix`, `test`, `docs`, `chore`, `refactor`, `perf`, `build`, `ci`, `revert`.
- **scope** (the context): `backend`, `frontend`, `protocol`, `ci`, `docs`, `agents`, `deps`.
  Optional, but use it whenever the change belongs to one area.
- **description**: starts with an **action verb** in the imperative (`add`, `fix`, `remove`,
  `rename`, ...), lower case, no trailing period, at most 72 characters.

Examples:

```text
feat(backend): add route to get current map
fix(backend): reject claims outside the canvas
test(frontend): cover color picker keyboard navigation
chore(deps): bump golangci-lint to v2.14
```

Details (breaking changes, trailers) are in [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Standards

- Go: [`.github/instructions/go.instructions.md`](.github/instructions/go.instructions.md)
- React/TypeScript: [`.github/instructions/frontend.instructions.md`](.github/instructions/frontend.instructions.md)
- Git, commits and PRs: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- End-to-end workflow: [`docs/workflow.md`](docs/workflow.md)

## Agents and skills

| Agent       | Use it for                                                              |
|-------------|-------------------------------------------------------------------------|
| `architect` | Turning ideas into specs, ADRs and agent-ready issues. Writes no product code. |
| `backend`   | Implementing Go issues in `backend/`.                                   |
| `frontend`  | Implementing React issues in `frontend/`.                               |

| Skill                     | Use it when                                              |
|---------------------------|----------------------------------------------------------|
| `open-pull-request`       | Your change is ready to be proposed.                     |
| `address-review-feedback` | A PR you own received review comments.                   |
| `write-adr`               | You are making or proposing an architectural decision.   |
