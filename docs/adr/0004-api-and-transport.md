# 0004. HTTP API and real-time transport

- **Status**: Accepted
- **Date**: 2026-10-04
- **Issue/PR**: #8

## Context

Clients need to load the canvas, submit claims, see other players' claims in real time and know
the server time, because a claim `ts` must not be in the future. The canvas has 65,536 cells. We
want the simplest design that works with the Go standard library.

## Considered options

| Topic | Chosen | Rejected alternatives |
|-------|--------|-----------------------|
| Real-time | **Server-Sent Events** for server-to-client updates, plain `POST` for claims: standard library only, browsers reconnect automatically. | **WebSocket**: two-way, but needs a third-party library and manual reconnection. **Polling**: laggy and wasteful. |
| Canvas load | **Colors-only snapshot** (65,536 bytes) plus `GET /api/cells/{x}/{y}` on demand: small for every visitor, and full details only for the cell being mined. | **Full JSON snapshot** (about 4 MB). **Full binary snapshot** (about 1.4 MB, custom decoder). |
| Server time | **`now` in every response and event**, plus `GET /api/time`. | **HTTP `Date` header**: not available on stream events. |
| Nonce encoding | **JSON number** (any uint64): the rule does not depend on client technology. Clients that cannot represent every uint64 mine a smaller range. | **Decimal or hex string**: exact in every language, but not a natural JSON type for a number. |
| Event replay | **None**: a reconnecting client reloads the 64 KB snapshot. | **Event ids with a replay buffer**: more server state for little gain. |

## Decision

Implement the API exactly as described in [`docs/spec/api.md`](../spec/api.md).

## Consequences

- One HTTP server and the standard library cover everything (`net/http`, `encoding/json`).
- Every accepted claim fans out to all open streams. Slow clients are disconnected instead of
  buffered without limit.
- There is no rate limiting in v1. Revisit if abuse shows up.
