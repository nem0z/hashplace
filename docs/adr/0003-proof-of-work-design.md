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

1. **Bits decay**: strength in bits, `256 - log2(hash)`, minus a fixed number of bits per period.
   Linear decay in bits is exponential in work, since 1 bit is a factor of 2.
   - Pros: retake cost halves every period, so even a huge proof is cheap again after a bounded
     time. Hardware advantage is capped (1000x is only about 10 bits).
   - Cons: protection grows only logarithmically with work.
2. **Work decays** (strength in expected hashes, minus a fixed number of hashes per second)
   - Pros: protection is proportional to effort.
   - Cons: GPU owners can lock cells for a very long time.

**Unit of work**

1. **Actual hash value** (`256 - log2(hash)`, fractional): every hash counts for exactly its
   value, so there are no factor-of-2 jumps between thresholds. Cons: float64 math.
2. **Count of leading zero bits** (integer): integer-only math. Cons: work moves in factor-of-2
   steps, and a lucky hash gets no credit beyond its last zero bit.

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

**Claim time**

1. **Chosen by the client and hashed**, checked only against the server clock (`ts <= now`):
   the stored hash contains its own time, so a cell's history can be re-verified from stored
   claims. Backdating only hurts the claimer. Cons: clients need server time.
2. **Server time at receipt**: no client clock involved. Cons: the time is not part of the hash,
   so stored claims cannot be re-verified on their own.

**Decay shape**: continuous (`bits - elapsed / period`) or step (whole bits at each period).

## Decision

Bits decay, continuous, SHA-256, and proofs bound to cell, color, generation and a client-chosen
`ts`, with no challenge. Work is the actual hash value expressed in bits (`256 - log2(hash)`,
float64). The rule is `claimed_work - (new_ts - claimed_ts) / decayPeriod < new_work`, with
elapsed time measured claim to claim and `new_ts <= server clock`. Parameters: 256 x 256 canvas,
16 colors, `minWork = 22`, `decayPeriod = 300 s`. Players are anonymous. Details are in
[`docs/spec/claims.md`](../spec/claims.md).

## Consequences

- Verification is one SHA-256 over 28 bytes, one `log2` and a comparison, so it is safe to run on
  every request.
- Stockpiling and GPU advantage are accepted and bounded by decay. Revisit with a challenge or a
  different hash if they hurt the game.
- Parameters are config values and should be tuned after a real browser hash-rate benchmark.
