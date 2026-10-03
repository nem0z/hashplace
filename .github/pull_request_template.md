Closes #

## What

<!-- One or two sentences: what this PR changes. -->

## Why

<!-- The problem or issue goal this addresses. -->

## How

<!-- Key implementation choices and trade-offs. Point the reviewer at the files that matter most. -->

## Verification

<!-- Commands run and their outcome. New or changed tests. Manual checks, if any. -->

- [ ] `backend`: `gofmt`, `go vet`, `golangci-lint`, `go test -race ./...`
- [ ] `web`: `lint`, `typecheck`, `test`, `build`
- [ ] Not applicable (docs / CI only)

## Checklist

- [ ] Title is a Conventional Commit (`type(scope): summary`)
- [ ] Scoped to a single issue; < 400 changed lines or the split is justified below
- [ ] Tests added or updated for the new behavior
- [ ] `docs/spec/claims.md` updated if game/protocol rules changed
- [ ] ADR added if an architectural decision was made
- [ ] No new dependencies, or each one is justified below
- [ ] `reviewer` agent pass done and findings addressed

## Notes for the reviewer

<!-- Open questions, follow-ups, anything you are unsure about. Write "None" if empty. -->
