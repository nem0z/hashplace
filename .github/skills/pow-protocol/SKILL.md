---
name: pow-protocol
description: Domain rules and invariants for hashplace claims, proof-of-work and decay. Load before designing, implementing, testing or reviewing anything that creates, verifies, stores or displays a claim.
---

# Proof-of-work claims: working rules

The **normative** source is [`docs/spec/claims.md`](../../../docs/spec/claims.md). This skill
summarizes it and explains how to work on it safely. If the two disagree, the spec wins. Report
the mismatch.

## Mental model

- A cell holds a color and a **claim**: `(strength0, claimedAt)`.
- The current strength decays linearly in **bits** over time:
  `strength(now) = max(0, strength0 - decayRate * (now - claimedAt))`. Each bit lost halves the
  work needed to take the cell back.
- A proof's **work** is measured in bits from its hash. A new claim wins only if
  `work >= minWork` and `work > strength(now)`.
- Proofs are bound to the target cell, the color, the server challenge and the cell's current
  generation. A proof cannot be reused elsewhere or replayed.

## Invariants (must hold in code and tests)

1. The server is authoritative for time, challenges and state. Never trust client time.
2. Verification is O(1) and cheap: decode, a bounds check, one hash, compare. Reject oversized or
   malformed input **before** hashing.
3. The preimage encoding is byte-exact and versioned. Go and TypeScript implementations must
   pass the same **golden test vectors** (`docs/spec/vectors/`, created with the first
   implementation).
4. Claim acceptance and the state update are atomic per cell. With two concurrent valid claims,
   exactly one wins against a given generation.
5. Decay is computed lazily from `(strength0, claimedAt, now)` with an injected clock. No
   background goroutine rewrites cells to apply decay.
6. Never use floating point in consensus-relevant comparisons. Use the fixed-point
   representation defined in the spec.

## Testing requirements

- Table-driven tests for acceptance and rejection at the boundaries: `work == strength`,
  `work == minWork`, fully decayed cell, generation mismatch, expired challenge.
- Fake clock only. No `time.Sleep` and no wall clock in tests.
- Fuzz the request decoder and the verifier (`go test -fuzz`).
- Benchmark verification and keep it allocation-free on the hot path.

## Changing the rules

A rule change means a spec PR first (via the `architect` agent): bump the protocol version,
update the vectors, then write the implementation issues. Never change a rule only in code.
