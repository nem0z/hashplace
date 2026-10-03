# Development workflow

hashplace is built entirely by AI agents. The maintainer (@nem0z) owns every decision and every
merge. This document describes how work moves from an idea to `main`.

```mermaid
flowchart LR
    I[Idea] -->|architect agent| D[Spec / ADR PR]
    D -->|maintainer merges| S[Agent-ready issues]
    S -->|backend / frontend agent| B[Branch + tests + code]
    B -->|open-pull-request skill| P[Draft PR]
    P -->|CI green| R[Ready for review]
    R -->|maintainer review| F{Approved?}
    F -->|changes requested| A[address-review-feedback skill]
    A --> R
    F -->|yes| M[Maintainer squash-merges]
```

## Roles

| Who              | Does                                                            | Never does                     |
|------------------|-----------------------------------------------------------------|--------------------------------|
| Maintainer       | Prioritizes, decides, reviews, merges                           | -                              |
| `architect`      | Specs, ADRs, splits work into agent-ready issues                | Product code                   |
| `backend`        | Implements `backend/` issues, test-first                        | Merges, touches `frontend/`    |
| `frontend`       | Implements `frontend/` issues                                   | Merges, touches `backend/`     |

## 1. From idea to issues

1. The maintainer opens an issue or describes the idea to the `architect` agent.
2. The architect asks questions, then opens a docs PR that changes `docs/spec/` and/or adds an ADR.
3. When the maintainer merges it, the architect creates **agent-ready** issues: one per PR-sized
   step, with acceptance criteria, scope and dependencies.

## 2. From issue to PR

1. Open a session on the issue in the Copilot desktop app (agents run on the maintainer's local
   machine) and select the `backend` or `frontend` agent.
2. The agent branches (`<type>/<issue>-<slug>`), writes tests first, implements, and runs all checks.
3. The agent opens a **draft** PR with the `open-pull-request` skill, waits for CI to be green
   and marks the PR ready.

## 3. Review and merge

1. The maintainer reviews on GitHub, leaving inline comments and a "Request changes" or
   "Comment" review.
2. The agent runs the `address-review-feedback` skill. It makes one commit per comment (or group
   of related comments), replies to and resolves each thread, and posts a round summary.
3. The maintainer re-reviews (and can reopen any thread), then **squash-merges** when CI is green
   and all threads are resolved. The PR title becomes the commit on `main`.

## Guardrails on `main`

A repository ruleset on `main` enforces:

- changes only through pull requests (no direct pushes, no force-pushes, no deletion),
- required checks: `ci-ok` (backend and frontend CI) and `conventional-title`,
- all review threads resolved before merging,
- squash merge only, linear history.

> Agents running locally use the maintainer's GitHub identity, so GitHub cannot tell them apart
> from the maintainer. The "never merge, never approve" rule is enforced by `AGENTS.md`, not by
> GitHub. For hard enforcement, run agents under a separate bot account or GitHub App and require
> one code-owner approval.

## Labels

| Label                    | Meaning                                         |
|--------------------------|-------------------------------------------------|
| `type:feature`, `type:bug`, `type:chore`, `type:docs` | Kind of work          |
| `area:backend`, `area:frontend`, `area:protocol`, `area:infra` | Part of the system |
| `status:needs-design`    | Needs a spec or ADR before implementation       |
| `status:agent-ready`     | Fully specified; any agent can pick it up       |
| `status:blocked`         | Waiting on another issue or a decision          |
