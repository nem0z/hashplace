# HTTP API

- **Status**: v1

This document is **normative**. Game rules are in [`claims.md`](claims.md).

## Conventions

- Base path: `/api`. The frontend is served from the same origin, so there is no CORS.
- Bodies are JSON (`application/json`) unless stated otherwise.
- Times are Unix seconds (integer). `nonce` is a JSON number: any uint64 is valid. Clients that
  cannot represent every uint64 (JavaScript, for example) simply mine a smaller range. `hash` is a
  lowercase hex string, or `null` for an empty cell.
- Every JSON response includes `now`, the server time, so clients can keep a clock offset for
  their claim `ts`.

A **cell** object, here after the claim from the test vectors (`nonce` `9620748`) was accepted:

```json
{"x": 12, "y": 34, "color": 5, "ts": 1791072000, "nonce": 9620748,
 "hash": "000002b3b357ccc62fe5625cbc469abaec00ed5b43380c2a519e66ef73544c96",
 "work": 22.56599618898906}
```

An empty cell has `ts` `0`, `nonce` `0`, `hash` `null` and `work` `0`.

## Endpoints

### `GET /api/time`

```json
{"now": 1791072000}
```

### `GET /api/canvas`

The colors of all cells as `application/octet-stream`: 65,536 bytes, one palette index per cell.
The color of `(x, y)` is at byte `y * 256 + x`.

### `GET /api/cells/{x}/{y}`

The full state of one cell, which is what a client needs to start mining it.

```json
{"cell": { ... }, "now": 1791072000}
```

`404 {"error": "out_of_range", "now": ...}` if the coordinates are outside the canvas.

### `POST /api/claims`

```json
{"x": 12, "y": 34, "color": 5, "ts": 1791072000, "nonce": 9620748}
```

The body is limited to 1 KiB. The server checks the claim as described in
[`claims.md`](claims.md) section 3.

| Status | `error` | Meaning |
|--------|---------|---------|
| 200 | - | Accepted. Body: `{"cell": { ... }, "now": ...}` with the updated cell |
| 400 | `invalid_request` | Malformed JSON, missing or invalid field, body too large |
| 400 | `out_of_range` | `x`, `y` or `color` out of range |
| 422 | `future_timestamp` | `ts` is after the server clock |
| 422 | `work_too_low` | Work below `minWork` |
| 422 | `claim_too_weak` | Work does not beat the current claim. The body also contains the current `cell`, so the miner can see what it has to beat |

Error bodies are `{"error": "<code>", "now": ...}`, plus `cell` for `claim_too_weak`. Codes are
stable and machine-readable.

### `GET /api/events`

A Server-Sent Events stream (`text/event-stream`).

| Event | `data` | When |
|-------|--------|------|
| `cell` | `{"cell": { ... }, "now": ...}` | After every accepted claim |

Events have no `id` and there is no replay. A client that reconnects reloads the canvas.

If a client is too slow to read its events, the server closes its stream rather than buffering
without limit. The browser reconnects automatically.

## Client loading sequence

1. Open `GET /api/events` and buffer `cell` events.
2. Fetch `GET /api/canvas`.
3. Apply the snapshot, then the buffered events in order, then live events.

Events carry the full cell and arrive in order, so this converges to the server state even if the
snapshot already contains some buffered events.

## Limits

- Request bodies: 1 KiB.
- No rate limiting in v1. Verifying a claim costs about 1 µs, so revisit only if abuse shows up.
