---
name: pow-protocol
description: Rules and invariants for hashplace claims, proof-of-work and decay. Load before designing, implementing, testing or reviewing anything that creates, verifies, stores or displays a claim.
---

# Proof-of-work claims

The **normative** source is [`docs/spec/claims.md`](../../../docs/spec/claims.md). This skill
summarizes it. If the two disagree, the spec wins: report the mismatch.

## In one paragraph

A proof is a `nonce`. `work` is the number of leading zero bits of
`SHA-256("hp1" || x || y || color || generation || nonce)` (20 bytes, big-endian). A claim with
work `w` wins a cell if `w >= 22` and `(cell.bits - w) * 300 < now - cell.claimedAt`. On success
the cell takes the new color, `bits = w`, `claimedAt = now`, and `generation` goes up by one.

## Invariants

1. The server is authoritative for time and state. Never use client time.
2. Validate ranges before hashing. Verification is one SHA-256 and integer comparisons only:
   no floating point anywhere in the rule.
3. The preimage encoding is byte-exact. Go and TypeScript both pass
   `docs/spec/vectors/pow-v1.json`.
4. Accepting a claim and updating the cell is atomic per cell. With two valid claims for the same
   generation, exactly one wins.
5. Decay is computed from `(bits, claimedAt, now)` when needed, with an injected clock. Nothing
   rewrites cells in the background.

## Tests

- Table-driven tests for the claim rule, including every example in spec section 3.
- Fake clock only. No `time.Sleep`.
- The golden vectors file in both Go and TypeScript.

## Changing the rules

Change the spec first, in its own PR, and bump the protocol version if the encoding changes.
Never change a rule only in code.
