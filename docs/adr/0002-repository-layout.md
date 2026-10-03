# 0002. Repository and backend package layout

- **Status**: Accepted
- **Date**: 2026-10-04
- **Issue/PR**: dev workflow PR

## Context

The project has a Go backend (the focus) and a React client. Agents need an unambiguous place for
every kind of code and clear dependency rules, so changes stay small and reviewable.

## Considered options

1. **Monorepo with `backend/` and `frontend/` top-level folders**, each with its own toolchain
   - Pros: atomic cross-cutting PRs (for example a protocol change with vectors), one place for
     the spec, CI, agents and skills.
   - Cons: CI has to handle two toolchains.
2. **Two repositories**
   - Pros: independent CI.
   - Cons: protocol changes span repositories, and agent instructions and the spec are duplicated.
3. **Go module at the repository root** (with `frontend/` nested)
   - Pros: shorter import paths.
   - Cons: Go tooling scans `node_modules` and frontend files, and the boundary is blurrier.

## Decision

We use a monorepo (option 1). The backend is the Go module `github.com/nem0z/hashplace/backend`
with this layout:

```text
backend/
  cmd/hashplaced/      main: config, wiring, graceful shutdown
  internal/pow/        proof-of-work encoding and verification (pure)
  internal/canvas/     canvas state, claim rules and decay (pure, injected clock)
  internal/store/      persistence behind interfaces consumed by canvas/api
  internal/api/        HTTP and streaming transport, DTOs, validation
  internal/config/     configuration loading
```

Dependency direction: `cmd -> api -> canvas -> pow`, and `store` implements interfaces declared
by its consumers. `pow` and `canvas` never import `api` or `store`.

The client lives in `frontend/` (Vite, React, TypeScript). The layout under `frontend/src/` is described in
`.github/instructions/frontend.instructions.md`.

## Consequences

- CI runs separate `backend` and `frontend` jobs and skips each one until its folder exists.
- New top-level `internal/` packages need an ADR update or a new ADR.
