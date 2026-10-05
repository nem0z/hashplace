# 0005. Persistence

- **Status**: Accepted
- **Date**: 2026-10-05
- **Issue/PR**: #9

## Context

The canvas lives in memory. If it is lost, old proofs become valid again (see
[`docs/spec/claims.md`](../spec/claims.md) section 7), so **every accepted claim must be durable
before it is acknowledged**. We also want to keep the full claim history, for timelapses and so
anyone can re-verify it. There are 65,536 cells, and claims arrive at most a few per second.

## Considered options

| Topic | Chosen | Rejected alternatives |
|-------|--------|-----------------------|
| Storage | **SQLite file**: embedded, one file, transactions, easy to query. | **Append-only log**: standard library only, but needs custom replay and compaction. **Periodic snapshot**: loses recent claims on a crash, which allows replays. **Postgres**: an extra service, overkill for this size. |
| Driver | **`modernc.org/sqlite`** (pure Go): no C toolchain, builds the same on Windows arm64, Linux and in CI. About 2-3x slower than cgo, still thousands of writes per second. | **`mattn/go-sqlite3`** (cgo): faster, but needs a C compiler everywhere. |
| Data kept | **Full history**: every accepted claim is appended to `claims`, and the current state is in `cells`. | **Current cells only**: smaller, but no timelapse and no history to re-verify. |
| Write path | **Synchronous**: the claim is written in a transaction before the cell is updated in memory. | **Batched writes**: faster, but a crash loses acknowledged claims. |

## Decision

SQLite through `modernc.org/sqlite`, in WAL mode, in one file (path from `HASHPLACE_DB`, default
`hashplace.db`). Schema:

```sql
CREATE TABLE claims (
  id    INTEGER PRIMARY KEY,   -- insertion order
  x     INTEGER NOT NULL,
  y     INTEGER NOT NULL,
  color INTEGER NOT NULL,
  ts    INTEGER NOT NULL,
  nonce INTEGER NOT NULL,      -- uint64 stored as its int64 bit pattern
  hash  BLOB    NOT NULL,      -- 32 bytes
  work  REAL    NOT NULL
);

CREATE TABLE cells (
  x        INTEGER NOT NULL,
  y        INTEGER NOT NULL,
  claim_id INTEGER NOT NULL REFERENCES claims(id),
  PRIMARY KEY (x, y)
);
```

Flow:

- **Startup**: load every row of `cells` joined with `claims` into the in-memory canvas.
- **Claim**: the canvas checks the claim under its lock. If it is accepted, the canvas calls its
  store, which in one transaction inserts into `claims` and upserts `cells`. Only if that succeeds
  is the cell updated in memory and the claim acknowledged. A store error rejects the claim
  (`500`).
- Reads (`GET /api/canvas`, `GET /api/cells/...`) are served from memory, never from SQLite.

Package boundaries (ADR 0002): `canvas` declares a small `Store` interface
(`SaveClaim(Cell) error`), `internal/store` implements it with SQLite, and `main` wires them. In
tests, the canvas uses an in-memory fake.

## Consequences

- One new dependency, `modernc.org/sqlite`, and still a single binary with no extra service.
- Every claim costs one SQLite transaction (about a millisecond) while the canvas lock is held.
  This is fine for a few claims per second. Revisit (for example per-cell locks) if needed.
- The `claims` table grows without limit, at about 100 bytes per claim. Pruning or archiving can
  come later if it becomes a problem.
- Backups are a copy of one file (`sqlite3 .backup` or a copy while the server is stopped).
