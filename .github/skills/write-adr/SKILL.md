---
name: write-adr
description: Write an Architecture Decision Record for hashplace. Use when introducing a dependency, a storage or transport choice, a package boundary, a protocol change, or any decision a future contributor would ask "why?" about.
---

# Write an ADR

ADRs live in `docs/adr/` as `NNNN-kebab-title.md`, numbered sequentially. They are short (1-2
pages) and immutable once accepted. To change a decision, write a new ADR that supersedes the old
one and update the old one's status line only.

## When an ADR is required

- Adding a third-party Go module or a significant npm package
- Choosing storage, transport (REST, WebSocket, SSE), serialization or deployment
- Creating or changing boundaries between `internal/` packages
- Any change to proof-of-work, claim or decay rules (together with a `docs/spec/` change)
- Trade-offs that reviewers are likely to question

## Steps

1. Find the next number: list `docs/adr/` and add 1 to the highest number.
2. Copy `docs/adr/0000-template.md` to `docs/adr/NNNN-<title>.md`.
3. Fill in every section. **Considered options** must contain at least two real alternatives
   with honest pros and cons. Quantify where possible (latency, memory, hashes per second, cost).
4. Set the status to `Proposed`. Change it to `Accepted` in the same PR once the maintainer
   agrees in review. Merging the PR is the acceptance.
5. Link the ADR from the PR description and from any spec section it affects.
