# 0003. Proof-of-work design

- **Status**: Accepted
- **Date**: 2026-10-04
- **Issue/PR**: #2

## Context

Painting a cell costs proof-of-work instead of a cooldown timer. A claim must weaken over time so
that any cell can eventually be retaken by a normal computer. The design must be simple, cheap to
verify on the server and fast enough to mine in a browser.

## Considered options

| Topic | Chosen | Rejected alternative |
|-------|--------|----------------------|
| Decay | **Bits decay**: strength in bits minus a fixed number of bits per period. This is exponential in work, so the retake cost halves every period and a 1000x hardware advantage is only about 10 bits. | **Work decays** (expected hashes minus a fixed rate): protection proportional to effort, but GPU owners could lock cells for a very long time. |
| Unit of work | **Actual hash value**, `256 - log2(hash)` (fractional, float64): every hash counts for exactly its value. | **Count of leading zero bits** (integer): integer-only math, but factor-of-2 jumps and no credit for lucky hashes. |
| Hash | **SHA-256**: native in Go, fast in browsers, about 1 µs to verify. GPUs are about 1000x faster. | **Argon2id**: smaller GPU advantage, but tens of hashes per second in browsers and milliseconds of server CPU and memory per verification. |
| Replay | **Bind to cell, color and generation**: stateless. Proofs can be mined ahead of time. | **Rotating server challenge**: limits stockpiling, but needs a secret, an endpoint and expires work. |
| Claim time | **Client-chosen `ts` in the hash**, only checked `ts <= server clock`: stored claims are verifiable on their own, and backdating only hurts the claimer. | **Server time at receipt**: no client clock, but the time is not in the hash. |
| Decay shape | **Continuous**: `work - elapsed / period`. | **Steps**: whole bits at each period. |

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
