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
| `work` | work of the current claim (float64, `0..256`) | `0` |
| `claimedAt` | server time of the current claim, Unix seconds (int64) | `0` |
| `generation` | number of accepted claims on this cell (uint32) | `0` |

## 2. Proof-of-work

A proof is a `nonce` (uint64) chosen by the client.

```text
preimage = "hp1"                3 bytes, ASCII
        || uint16be(x)          2 bytes
        || uint16be(y)          2 bytes
        || uint8(color)         1 byte
        || uint32be(generation) 4 bytes
        || uint64be(nonce)      8 bytes
                                = 20 bytes
hash = SHA-256(preimage)
H    = hash read as a 256-bit unsigned big-endian integer
work = 256 - log2(H)            (256 if H = 0)
```

Work is measured in **bits** from the actual value of the hash, so it is fractional: the smaller
the hash, the more work. A hash with `work >= w` takes `2^w` attempts on average, so each extra
bit doubles the expected work. Work is computed in IEEE 754 float64.

The proof is bound to the cell, the color and the cell's current `generation`. Once any claim is
accepted on the cell, its generation changes and every other proof for that cell becomes invalid.
Proofs cannot be replayed or reused on another cell or color.

There is no server challenge: proofs never expire while the cell is untouched, so players may
mine ahead of time.

## 3. Claim rule

| Parameter | Value |
|-----------|-------|
| `minWork` | 22 bits (about 4M hashes, about 0.5 s on a laptop) |
| `decayPeriod` | 300 s (a claim loses 1 bit every 5 minutes) |

A claim's strength decays linearly in bits, so the work needed to retake it halves every
`decayPeriod`:

```text
strength(now) = work - (now - claimedAt) / decayPeriod
```

A claim `(x, y, color, generation, nonce)` with work `w` received at server time `now` is
**accepted** if and only if:

1. `x`, `y` and `color` are in range,
2. `generation == cell.generation`,
3. `w >= minWork`,
4. `w > strength(now)`, evaluated in float64 as
   `(cell.work - w) * decayPeriod < now - cell.claimedAt`.

On acceptance, atomically for the cell:

```text
color = color, work = w, claimedAt = now, generation = generation + 1
```

Consequences:

- An empty or fully decayed cell costs `minWork`.
- The price of a painted cell is set by whoever painted it last: mining more buys more
  protection. Each extra bit doubles the work and adds 5 minutes of protection.
- Lucky proofs count: the stored `work` comes from the hash actually found, not a target.
- Repainting your own cell (same color) is an ordinary claim. There is no special case.
- At `claimedAt`, retaking needs strictly more work than the claim. After that, the threshold
  drops by 1 bit every `decayPeriod`.

For clients: a proof wins if its work is `>= minWork` and `> cell.work - (now - cell.claimedAt) /
decayPeriod`. Mining against the hash target `H < 2^(256 - threshold)` finds such a proof.

### Claim rule examples

| `cell.work` | `now - claimedAt` | `w` | Accepted | Why |
|-------------|-------------------|-----|----------|-----|
| 0 | any | 21.9 | no | below `minWork` |
| 0 | any | 22 | yes | empty cell |
| 30 | 0 | 30 | no | needs `w > 30` |
| 30 | 0 | 30.01 | yes | |
| 30 | 1 | 30 | yes | `0 < 1` |
| 30 | 150 | 29.5 | no | `0.5 * 300 = 150 < 150` is false |
| 30 | 151 | 29.5 | yes | |
| 30 | 300 | 29 | no | `300 < 300` is false |
| 40 | 3601 | 28 | yes | `12 * 300 = 3600 < 3601` |

## 4. Time

The server clock is the only clock. Times are whole Unix seconds. Client-supplied times are never
used.

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
- **Stockpiling**: players can pre-mine proofs for untouched cells and submit them at once.
- **State loss**: proofs are not bound to a canvas instance. If cell state is lost and generations
  restart at 0, old proofs become valid again. Persistence must keep generations.
