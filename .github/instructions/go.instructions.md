---
applyTo: "backend/**"
---

# Go standards (backend)

**Simplest option first.** Simple, short and clear code beats clever or "future-proof" code.
Add complexity only when the simple version does not work or the maintainer asks for it.

## Toolchain

- Go **1.27** (the `go` directive in `backend/go.mod` is the source of truth). Module path:
  `github.com/nem0z/hashplace/backend`.
- Format with `gofmt` (imports grouped by `goimports`). Lint with **golangci-lint v2** using
  `backend/.golangci.yml`. No `//nolint` without a `// reason:` comment on the same line.

## Layout and dependencies

- `cmd/hashplaced/` holds wiring only: config, construction, start and graceful shutdown.
- `internal/<domain>/` packages follow ADR 0002. Domain packages (`pow`, `canvas`) are pure: no
  HTTP, no SQL, no `os`, no global clock.
- Dependencies point inward: `api -> canvas -> pow`. Domain packages never import transport or
  storage packages.
- Define interfaces where they are **consumed**, only when you need one (a second implementation
  or a test fake). Keep them small and return concrete types.
- Prefer the standard library: `net/http` routing patterns, `log/slog`, `encoding/json`,
  `crypto/*`. A new module needs a justification in the PR; significant ones need an ADR.

## Code

Follow `AGENTS.md` > Code style. In Go specifically:

- Blank lines between blocks are enforced by the `wsl_v5` linter, configured to match
  `AGENTS.md` > Code style. Fix its findings and never add `//nolint` for it. If a finding
  conflicts with the style rules, change the linter settings instead.
- Defaults are `const` declarations at the top of the package that owns them
  (`const defaultAddr = ":8080"`), never literals inside functions.
- Goroutines: `go func() { errChan <- srv.ListenAndServe() }()` is fine on one line; anything
  longer goes in a named function (`go worker(ctx, jobs)`).
- Shared HTTP behavior (for example the `Hashplace-Now` header) is a middleware wrapping the mux.
- `context.Context` is the first parameter of anything that does I/O or may block. Never store
  it in a struct.
- Errors: wrap only with context the **current** function owns: what *it* was doing and with
  which inputs, for example `fmt.Errorf("load cell (%d,%d): %w", x, y, err)`. Context about the
  operation that failed belongs in the function that returned the error, not in its callers. If
  the current function has nothing relevant to add, return `err` as is.
- Expose sentinel or typed errors from domain packages (`canvas.ErrClaimTooWeak`) and check them
  with `errors.Is` and `errors.As`. Handle an error once: log it **or** return it, not both.
- No panics outside `main` startup and genuine programmer errors.
- No mutable package-level state. Configuration is a struct populated in `main` from flags or env.
- **Inject time only where logic depends on it.** Pass a clock (`now func() time.Time`) to code
  whose behavior depends on time: claim checks, decay, expiry. Code that only reports the time
  (for example the `Hashplace-Now` header) calls `time.Now()` directly. Randomness that affects
  behavior is injected as an `io.Reader`.
- Concurrency: a mutex is owned by the type it protects and documented next to the field. Every
  goroutine has a clear owner and a stop path (context cancellation). No goroutine leaks.
- HTTP: set server timeouts, cap bodies with `http.MaxBytesReader`, validate every field
  (bounds, palette index, lengths) before doing any work, and return JSON errors with stable
  machine-readable codes.
- Logs: structured `slog` with stable keys. Never log secrets or full request bodies.
- Exported identifiers have doc comments. Comment *why*, not *what*.

## Tests

- Tests sit next to the code (`_test.go`). Use table-driven tests with `t.Run`, and `t.Parallel()`
  when there is no shared state.
- Deterministic only: fake clock, seeded or fixed randomness, no network, no `time.Sleep`.
- Use `net/http/httptest` for handlers and test through the public package API.
- Add fuzz tests for decoders and the proof-of-work verifier, and benchmarks for hot paths
  (verification, canvas snapshot and diff).
- `go test -race ./...` must pass. Target >= 90% coverage for `internal/pow` and `internal/canvas`.
- Bug fixes start with a failing regression test.

## Commands

```sh
cd backend
gofmt -l .            # must print nothing
go vet ./...
golangci-lint run
go test -race ./...
```
