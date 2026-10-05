# Claims, proof-of-work and decay

- **Status**: v1
- **Protocol version**: `hp1`

This document is **normative**. Code must match it. To change a rule, change this file first (in
its own PR), then the code.

## 1. Canvas

| Parameter | Value |
|-----------|-------|
| Size | 256 x 256 cells, coordinates `0 <= x, y < 256` |
| Palette | 16 colors, index `0 <= color < 16` (see section 6) |
| Identity | None. Players are anonymous. |

Each cell stores:

| Field | Type | Empty cell |
|-------|------|------------|
| `color` | palette index (uint8) | `0` (white) |
| `ts` | timestamp of the current claim, chosen by the claimer, Unix seconds (int64) | `0` |
| `nonce` | nonce of the current claim (uint64) | `0` |
| `hash` | hash of the current claim (32 bytes) | none |
| `work` | work of `hash` (float64, `0..256`) | `0` |

A claim can be re-verified by anyone from the stored fields alone.

## 2. Proof-of-work

A proof is a timestamp `ts` (int64, Unix seconds) and a `nonce` (uint64), both chosen by the
client.

```text
preimage = "hp1"                3 bytes, ASCII
        || uint16be(x)          2 bytes
        || uint16be(y)          2 bytes
        || uint8(color)         1 byte
        || int64be(ts)          8 bytes
        || uint64be(nonce)      8 bytes
                                = 24 bytes
hash = SHA-256(preimage)
H    = hash read as a 256-bit unsigned big-endian integer
work = 256 - log2(H)            (256 if H = 0)
```

Work is measured in **bits** from the actual value of the hash, so it is fractional: the smaller
the hash, the more work. A hash with `work >= w` takes `2^w` attempts on average, so each extra
bit doubles the expected work. Work is computed in IEEE 754 float64.

The proof is bound to the cell, the color and its own `ts`, so it cannot be reused on another cell
or color. It is not bound to the cell's current state: a proof wins against whatever claim is
stored when it arrives, as long as it beats it (section 3).

There is no server challenge: proofs never expire, so players may mine ahead of time.

## 3. Claim rule

| Parameter | Value |
|-----------|-------|
| `minWork` | 22 bits (about 4M hashes, about 0.5 s on a laptop) |
| `decayPeriod` | 300 s (a claim loses 1 bit every 5 minutes) |

A claim's strength decays linearly in bits, so the work needed to retake it halves every
`decayPeriod`:

```text
strength(t) = work - (t - ts) / decayPeriod
```

A claim `(x, y, color, ts, nonce)` with work `w`, processed when the server clock reads `now`, is
**accepted** if and only if:

1. `x`, `y` and `color` are in range,
2. `ts <= now` (no claims from the future),
3. `w >= minWork`,
4. `w > strength(ts)` of the stored claim, evaluated in float64 as
   `(cell.work - w) * decayPeriod < ts - cell.ts`.

Elapsed time is measured **claim to claim**, from the stored claim's `ts` to the new claim's
`ts`. The server clock is only used for check 2.

On acceptance, atomically for the cell:

```text
color = color, ts = ts, nonce = nonce, hash = hash, work = w
```

**Replays are impossible.** Check 4 is equivalent to
`w - ts / decayPeriod > cell.work - cell.ts / decayPeriod`: every claim has a score
`work - ts / decayPeriod`, and a new claim must have a strictly higher score than the cell's.
A cell's score therefore only goes up, so a proof that was already accepted, or already beaten,
can never be accepted again.

Consequences:

- An empty or fully decayed cell costs `minWork`.
- The price of a painted cell is set by whoever painted it last: mining more buys more
  protection. Each extra bit doubles the work and adds 5 minutes of protection.
- Lucky proofs count: the stored `work` comes from the hash actually found, not a target.
- Repainting your own cell (same color) is an ordinary claim. There is no special case.
- Backdating `ts` is allowed but only hurts the claimer: the claim is less likely to beat the
  stored one, and it decays sooner.

For clients: set `ts` to the current server time and refresh it while mining. A proof wins if its
work is `>= minWork` and `> cell.work - (ts - cell.ts) / decayPeriod`. Mining against the hash
target `H < 2^(256 - threshold)` finds such a proof.

### Claim rule examples

| `cell.work` | `ts - cell.ts` | `w` | Accepted | Why |
|-------------|----------------|-----|----------|-----|
| 0 | any `>= 0` | 21.9 | no | below `minWork` |
| 0 | any `>= 0` | 22 | yes | empty cell |
| 30 | 0 | 30 | no | needs `w > 30` |
| 30 | 0 | 30.01 | yes | |
| 30 | 1 | 30 | yes | `0 < 1` |
| 30 | 150 | 29.5 | no | `0.5 * 300 = 150 < 150` is false |
| 30 | 151 | 29.5 | yes | |
| 30 | 300 | 29 | no | `300 < 300` is false |
| 40 | 3601 | 28 | yes | `12 * 300 = 3600 < 3601` |
| 30 | -300 | 31 | no | backdated: `-300 < -300` is false |

Any claim with `ts > now` is rejected, whatever its work.

## 4. Time

Times are whole Unix seconds (int64). Claim timestamps are chosen by the client and are part of
the hash, so a cell's history can be re-verified without the server clock. The server clock only
rejects timestamps from the future. Clients should sync with server time: if a client's clock is
ahead, its claims are rejected; if it is behind, its claims are backdated.

## 5. Test vectors

[`vectors/pow-v1.json`](vectors/pow-v1.json) lists preimages, hashes and work values. The Go
verifier and the browser miner must both pass them, comparing `work` with a tolerance of `1e-9`.

## 6. Palette

| Index | Color | Index | Color |
|-------|-------|-------|-------|
| 0 | `#FFFFFF` | 8 | `#E5D900` |
| 1 | `#E4E4E4` | 9 | `#94E044` |
| 2 | `#888888` | 10 | `#02BE01` |
| 3 | `#222222` | 11 | `#00D3DD` |
| 4 | `#FFA7D1` | 12 | `#0083C7` |
| 5 | `#E50000` | 13 | `#0000EA` |
| 6 | `#E59500` | 14 | `#CF6EE4` |
| 7 | `#A06A42` | 15 | `#820080` |

## 7. Known limitations

- **Hardware advantage**: a GPU hashes about 1000x faster than a browser (about 10 bits). With
  decay this buys about 50 extra minutes of protection, not permanent ownership.
- **Stockpiling**: players can pre-mine proofs for any cell, using a future `ts`, and submit them
  all once that time arrives. A pre-mined proof stays usable as long as it beats the cell's
  current claim.
- **State loss**: proofs are not bound to a canvas instance. If cell state is lost, old proofs
  become valid again. Persistence must keep the current claim of every cell.
