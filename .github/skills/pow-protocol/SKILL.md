---
name: pow-protocol
description: Rules and invariants for hashplace claims, proof-of-work and decay. Load before designing, implementing, testing or reviewing anything that creates, verifies, stores or displays a claim.
---

# Proof-of-work claims

The **normative** source is [`docs/spec/claims.md`](../../../docs/spec/claims.md). This skill
summarizes it. If the two disagree, the spec wins: report the mismatch.

## In one paragraph

A proof is a client-chosen `ts` (Unix seconds) and `nonce`. `work = 256 - log2(H)` in bits
(float64), where `H` is the value of
`SHA-256("hp1" || x || y || color || generation || ts || nonce)` (28 bytes, big-endian) read as a
256-bit integer. Work is **not** the count of leading zero bits. A claim with work `w` wins a cell
if `ts <= server clock`, `w >= 22` and `(cell.work - w) * 300 < ts - cell.ts`. Elapsed time is
measured claim to claim. On success the cell stores the new claim (color, `ts`, nonce, hash,
work), and `generation` goes up by one.

## Invariants

1. The server clock is used for one check only: reject `ts > now`. Decay uses the claims' own
   `ts`. Never trust anything else from the client.
2. Validate ranges before hashing. Verification is one SHA-256, one `log2` and a comparison,
   exactly as written in the spec.
3. The preimage encoding is byte-exact. Go and TypeScript both pass
   `docs/spec/vectors/pow-v1.json`.
4. Accepting a claim and updating the cell is atomic per cell. With two valid claims for the same
   generation, exactly one wins.
5. Decay is computed from the stored `(work, ts)` and the new claim's `ts`. Nothing rewrites cells
   in the background. The server clock is injected so tests can control it.

## Tests

- Table-driven tests for the claim rule, including every example in spec section 3.
- Fake clock only. No `time.Sleep`.
- The golden vectors file in both Go and TypeScript.

## Changing the rules

Change the spec first, in its own PR, and bump the protocol version if the encoding changes.
Never change a rule only in code.
