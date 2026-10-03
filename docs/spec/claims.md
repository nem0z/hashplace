# Claims, proof-of-work and decay

- **Status**: DRAFT v0. Not final. Open questions are tracked in the "Finalize claim spec v1" issue.
- **Protocol version**: `hp1` (provisional)

This document is **normative**. Code must match it. Change it in a PR (via the `architect` agent)
before changing behavior.

## 1. Concepts

| Term | Definition |
|------|------------|
| Canvas | A `W x H` grid of cells. Coordinates `0 <= x < W`, `0 <= y < H`. |
| Palette | A fixed list of `P` colors, addressed by index `0 <= c < P`. |
| Cell state | `color`, `strength0` (millibits), `claimedAt` (server time), `generation` (uint64, incremented on every accepted claim). |
| Work | The number of leading zero bits of a proof's hash, as an integer `0..256`. |
| Strength | How much work it currently takes to overwrite a cell. Decays over time. |

## 2. Proof-of-work

### 2.1 Challenge

The server publishes a **challenge** that rotates every `epoch` (for example 10 minutes):

```text
challenge = HMAC-SHA256(serverSecret, "hp1/epoch" || uint64be(epochIndex))
```

The server accepts proofs against the **current and previous** epochs only. Old proofs expire, so
work cannot be precomputed far in advance. The challenge is stateless to verify.

### 2.2 Preimage (byte-exact)

```text
preimage = "hp1"                      (3 bytes, ASCII)
        || challenge                  (32 bytes)
        || uint16be(x) || uint16be(y) (4 bytes)
        || uint8(color)               (1 byte)
        || uint64be(generation)       (8 bytes)
        || nonce                      (16 bytes, chosen by the client)
hash    = SHA-256(preimage)
work    = leadingZeroBits(hash)
```

Binding the proof to `(x, y, color, generation)` means it is valid for exactly one claim attempt.
It cannot be replayed after a decay, reused on another cell or used for another color. If someone
else claims the cell first, the generation changes and the proof becomes worthless. This is
intended: it creates a race, just like mining.

### 2.3 Verification

1. Decode and check bounds (`x`, `y`, `color`, nonce length) and epoch freshness. Reject before
   hashing.
2. Recompute the hash and `work`.
3. Apply the claim rule (section 3) atomically for the cell.

## 3. Claim rule and decay

All values are integers. Strength is stored in **millibits** (1 bit = 1000) to avoid floating point.

```text
elapsed  = now - claimedAt                       (whole seconds, server clock)
strength = max(0, strength0 - decayRate * elapsed)    (millibits; decayRate in millibits/s)

accept iff  work >= minWork
       and  work * 1000 > strength
       and  generation == cell.generation

on accept:  color = c, strength0 = work * 1000, claimedAt = now, generation += 1
```

- Decay is **linear in bits**, so the work needed to retake a cell **halves** every
  `1000 / decayRate` seconds. This is the "initial pow - time delta" idea from the project brief,
  expressed in the log domain of work.
- A cell that has never been claimed has `strength = 0`, so `minWork` is the floor price of any
  pixel.
- Ties lose: you must strictly beat the current strength.
- The client may submit any hash it found. Lucky extra zero bits count as extra strength.

### Illustrative parameters (to be calibrated)

Assuming a browser mines about 2^22 SHA-256 hashes per second (about 4 MH/s):

| Parameter | Value | Effect |
|-----------|-------|--------|
| `minWork` | 22 bits | about 1 s to paint an unclaimed or fully decayed cell |
| `decayRate` | 1 bit per 120 s | work to retake halves every 2 minutes |
| A 30-bit claim | about 4 minutes of mining | needs more than 30 bits immediately; drops to `minWork` after about 16 minutes |

## 4. Open questions (to resolve before v1)

1. **Hash function.** SHA-256 gives GPU and ASIC owners a 1000x or greater advantage over browsers.
   Alternatives: a memory-hard function (Argon2id, scrypt) with a higher verification cost, or a
   per-identity cap. Which matters more: fairness or cheap verification?
2. **Work granularity.** Use integer leading zero bits (simple), or fractional work
   `-log2(hash / 2^256)` in fixed point (smoother, more complex)?
3. **Generation binding vs. contention.** Binding to `generation` prevents replay but wastes work
   on busy cells. The alternative is binding only to the epoch plus storing spent proofs.
4. **Identity and abuse.** Is the game anonymous? Do we need per-IP limits on verification
   requests (verification is cheap but not free)?
5. **Parameters.** Canvas size, palette size, `epoch`, `minWork`, `decayRate`, and whether they
   adapt to global activity.
6. **Stacking.** Can the current owner "reinforce" their own cell, and is work additive or replacing?
7. **Real-time updates.** WebSocket or Server-Sent Events for broadcasting cell changes
   (transport ADR).
