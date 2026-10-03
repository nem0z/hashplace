# Contributing

All code in hashplace is written by AI agents and reviewed by the maintainer (@nem0z) through
GitHub pull requests. These rules apply to agents and humans alike. The full lifecycle is in
[`docs/workflow.md`](docs/workflow.md).

## Issues

Work starts from an issue. An issue is **agent-ready** (label `status:agent-ready`) when it has:

- a clear goal and the reasoning behind it,
- acceptance criteria that can be checked,
- explicit scope and out-of-scope lists,
- an `area:*` label and a `type:*` label.

Issues that still need a decision get `status:needs-design` and go to the `architect` agent first.

## Branches

```text
<type>/<issue-number>-<short-slug>
```

Examples: `feat/12-claim-endpoint`, `fix/31-decay-rounding`, `docs/4-pow-spec`.

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(<scope>): <imperative summary, lower case, no trailing period>
```

| Types | `feat`, `fix`, `perf`, `refactor`, `test`, `docs`, `build`, `ci`, `chore`, `revert` |
|-------|----------------------------------------------------------------------------------------|
| Scopes | `backend`, `web`, `protocol`, `ci`, `docs`, `agents`, `deps`                          |

Breaking changes use `!` (`feat(protocol)!: bind proofs to cell generation`) and a
`BREAKING CHANGE:` footer. AI-authored commits keep the tool's `Co-authored-by` trailer.

## Pull requests

- **Title** is a Conventional Commit. It becomes the squash-commit message on `main` (CI checks it).
- **Body** follows the PR template. Link the issue with `Closes #<n>`.
- **Size**: target < 400 changed lines (excluding lockfiles and generated code). Split larger work
  into stacked PRs.
- **One concern per PR.** Drive-by refactors go in their own PR.
- Open as **draft**, wait for CI to be green, then mark ready.
- During review, **push new commits** instead of force-pushing, so the reviewer can see what
  changed. Rebasing on `main` before the final review is fine.
- The agent replies to every review thread. **Only the maintainer resolves threads and merges.**
- Merge method: **squash only**. The branch is deleted automatically.

## Definition of done

- [ ] CI is green (`ci-ok` and `conventional-title` checks).
- [ ] New behavior has tests; bug fixes start with a failing test.
- [ ] `docs/spec/` is updated if game or protocol rules changed.
- [ ] An ADR is added if an architectural decision was made.
- [ ] The PR template is fully filled in, including how it was verified.
