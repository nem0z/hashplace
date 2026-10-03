---
name: address-review-feedback
description: Work through review comments on a hashplace pull request you own. Fix or answer every thread, push incremental commits and hand back to the maintainer. Use when a PR has new review comments or requested changes.
---

# Address review feedback

## 1. Collect every open thread

```sh
gh pr view <n> --json reviews,comments,headRefName
gh api graphql -F owner=nem0z -F repo=hashplace -F pr=<n> -f query='
query($owner:String!,$repo:String!,$pr:Int!){
  repository(owner:$owner,name:$repo){
    pullRequest(number:$pr){
      reviewThreads(first:100){
        nodes{ id isResolved path line
          comments(first:20){ nodes{ databaseId author{login} body } } }
      }
    }
  }
}'
```

Work only on unresolved threads, plus top-level review bodies and PR comments from the maintainer.

## 2. Triage each comment

| Kind                          | Action                                                      |
|-------------------------------|-------------------------------------------------------------|
| Clear defect or request       | Fix it.                                                     |
| Question                      | Answer it. Change the code or docs if the answer shows they were unclear. |
| Suggestion you disagree with  | Reply with a concise technical argument. Do not silently ignore it. |
| Out of scope                  | Propose a follow-up issue. Create it only if the maintainer agrees. |

Maintainer comments always win over reviewer-agent comments. If two comments conflict, ask.

## 3. Fix with incremental commits

- **Push new commits. Do not force-push during review**, so the maintainer can see what changed
  since their last review.
- Use one commit per logical fix or one per review round, for example
  `fix(backend): reject negative coordinates in claim request`.
- Re-run the full checks for the touched areas before pushing.

## 4. Reply in every thread

Reply inside the thread (not as a new top-level comment) with what you did and the commit SHA,
for example "Fixed in a1b2c3d: bounds are now validated in `canvas.Cell`".

```sh
gh api repos/nem0z/hashplace/pulls/<n>/comments/<root-comment-id>/replies -f body='...'
```

**Do not resolve threads.** The maintainer resolves them once satisfied. Branch protection
requires every thread to be resolved before merging.

## 5. Hand back

Post one short summary comment ("Round N addressed: ...") listing what changed and what is
still open for discussion. Then stop and wait for the maintainer.
