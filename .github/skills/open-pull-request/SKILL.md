---
name: open-pull-request
description: Prepare, verify and open a hashplace pull request that meets the repo's PR standards. Use when a change is ready to be proposed for review.
---

# Open a pull request

Follow every step. Do not skip verification to save time.

## 1. Sync and check the branch

```sh
git fetch origin
git rebase origin/main          # resolve conflicts, re-run tests afterwards
git branch --show-current       # must match <type>/<issue>-<slug>
```

If the branch was already pushed, a rebase needs a force-push. Always use
`git push --force-with-lease`, never plain `--force`: it refuses to overwrite commits you have not
seen.

## 2. Verify

Run the checks for every area you touched (see `AGENTS.md` > Commands). Everything must pass
locally. If a check cannot run in your environment (a missing tool, for example), say so in the
PR's **Verification** section. Never claim a check passed if you did not run it.

## 3. Self-review the diff

```sh
git diff --stat origin/main...
git diff origin/main...
```

- Remove debug output, commented-out code, TODOs without an issue link and unrelated changes.
- Check the size. More than about 250 changed lines (excluding lockfiles and generated code) means
  you must break the work down into smaller PRs.
- Check that spec, ADR and doc updates are included where needed.

## 4. Push and open the PR

Open a regular PR, ready for review (not a draft).

```sh
git push -u origin HEAD
gh pr create --base main \
  --title "<type>(<scope>): <summary>" \
  --body-file <filled-template.md> \
  --label "<type:...>" --label "<area:...>"
```

- Build the body from `.github/pull_request_template.md` and fill in **every** section. Write
  "None" or "N/A" rather than deleting a section.
- The first line must be `Closes #<issue>`.
- The title must be a Conventional Commit (CI enforces this): lower-case summary, imperative mood,
  no trailing period.

## 5. Wait for CI

Run `gh pr checks <n> --watch`. Fix failures with new commits until everything is green.

## 6. Hand off

Tell the maintainer the PR is ready and give its URL, a one-line summary and any open questions.
Then stop. **Never merge, approve or enable auto-merge.**
