---
name: address-review-feedback
description: Work through review comments on a hashplace pull request you own. Make one commit per comment (or group of related comments), reply to and resolve every thread, then hand back to the maintainer. Use when a PR has new review comments or requested changes.
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

If two comments conflict, ask.

## 3. One commit per comment

- Make **one commit per review comment**, or one per group of comments that ask for the same
  change. Each thread then maps to exactly one commit.
- The commit message follows the commit standard in `CONTRIBUTING.md` and describes the change,
  for example `fix (backend): reject negative coordinates in claim request`.
- Push new commits on top of the branch, so the maintainer can see what changed since their last
  review.
- Re-run the full checks for the touched areas before pushing.

## 4. Reply to and resolve every thread

After pushing, go through each thread: reply inside it (not as a new top-level comment) with what
you did and the commit SHA, then **resolve it**. Example reply: "Fixed in a1b2c3d: bounds are now
validated in `canvas.Cell`". A comment answered without a code change gets the answer, then is
resolved too.

```sh
gh api repos/nem0z/hashplace/pulls/<n>/comments/<root-comment-id>/replies -f body='...'
gh api graphql -F id=<thread-id> -f query='
mutation($id:ID!){ resolveReviewThread(input:{threadId:$id}){ thread{ isResolved } } }'
```

Leave a thread unresolved only when you need a decision from the maintainer (a disagreement or an
open question). Say so explicitly in your reply.

## 5. Turn general feedback into rules

If a maintainer comment expresses a preference that applies beyond this PR (naming, formatting,
structure, API style), make sure future work follows it without being told again:

1. Prefer enforcing it with a tool (a linter rule, a formatter setting) in the current PR when
   the PR already owns that config.
2. Otherwise, propose the rule in the matching standards file (`AGENTS.md` > Code style,
   `.github/instructions/*.instructions.md` or `CONTRIBUTING.md`) in a separate small PR.
3. Mention the new rule or linter in your thread reply.

## 6. Hand back

Post one short summary comment ("Round N addressed: ...") listing what changed and what is
still open for discussion. Then stop and wait for the maintainer.
