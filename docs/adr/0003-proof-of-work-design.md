# 0003. Proof-of-work design

- **Status**: Accepted
- **Date**: 2026-10-04
- **Issue/PR**: #2

## Context

Painting a cell costs proof-of-work instead of a cooldown timer. A claim must weaken over time so
that any cell can eventually be retaken by a normal computer. The design must be simple, cheap to
verify on the server and fast enough to mine in a browser.

## Considered options

**Decay**

1. **Bits decay** (strength in leading zero bits, minus a fixed number of bits per period)
   - Pros: retake cost halves every period, so even a huge proof is cheap again after a bounded
     time. Hardware advantage is capped (1000x is only about 10 bits).
   - Cons: protection grows only logarithmically with work.
2. **Work decays** (strength in expected hashes, minus a fixed number of hashes per second)
   - Pros: protection is proportional to effort.
   - Cons: GPU owners can lock cells for a very long time.

**Hash function**

1. **SHA-256**: native in Go, fast in browsers, verification costs about 1 µs.
   Cons: GPUs are about 1000x faster than browsers.
2. **Memory-hard (Argon2id)**: smaller GPU advantage. Cons: tens of hashes per second in a
   browser (noisy solve times), and milliseconds of server CPU and memory per verification.

**Replay protection**

1. **Bind to cell, color and generation** (no server challenge): stateless, nothing to rotate.
   Cons: proofs can be mined ahead of time.
2. **Add a rotating server challenge**: limits stockpiling. Cons: a server secret, an extra endpoint
   and expired work.

**Decay shape**: continuous (`bits - elapsed / period`) or step (whole bits at each period).

## Decision

Bits decay, continuous, SHA-256, and proofs bound to cell, color and generation with no
challenge. Integer leading zero bits are the unit of work. The claim rule is evaluated with
integers only. Parameters: 256 x 256 canvas, 16 colors, `minBits = 22`, `decayPeriod = 300 s`.
Players are anonymous. Details are in [`docs/spec/claims.md`](../spec/claims.md).

## Consequences

- Verification is one SHA-256 over 20 bytes plus a few integer comparisons, so it is safe to run
  on every request.
- Stockpiling and GPU advantage are accepted and bounded by decay. Revisit with a challenge or a
  different hash if they hurt the game.
- Parameters are config values and should be tuned after a real browser hash-rate benchmark.
