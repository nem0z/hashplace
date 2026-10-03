---
applyTo: "frontend/**"
---

# Frontend standards (React + TypeScript)

**Simplest option first.** Simple, short and clear code beats clever or "future-proof" code.
Add complexity (state libraries, abstractions, memoization) only when the simple version does not
work or the maintainer asks for it.

## Toolchain

- Node **LTS** and **npm**. Commit `package-lock.json`. Use Vite, React and TypeScript with
  `strict: true`.
- ESLint (typescript-eslint, react-hooks) and Prettier. Vitest and React Testing Library.
- Required `package.json` scripts: `dev`, `build`, `lint`, `typecheck` (`tsc --noEmit`),
  `test` (`vitest run`, non-watch), `format`.

## Architecture

- The client is a **view and a miner**. It never decides whether a claim is valid. It shows what
  the server says.
- Folder layout under `frontend/src/`:
  - `api/`: typed HTTP and stream client. The only place that calls `fetch`.
  - `pow/`: preimage encoding and the mining Web Worker. It must match the Go encoding byte for
    byte.
  - `canvas/`: rendering (a single `<canvas>` element, never one DOM node per pixel), pan and zoom.
  - `components/`: presentational components.
  - `hooks/`: state and effects glue.
- Mining runs **only** in a Web Worker and can be cancelled. The UI shows progress and expected time.
- Keep server state in one place (the `api` layer plus hooks). No ad-hoc `fetch` calls in
  components.

## Code

- Function components and hooks only. No `any`, no non-null `!` without a comment, and no
  `@ts-ignore` (use `@ts-expect-error` with a reason if you have to).
- Mirror the backend's API types in `api/types.ts`. Do not hand-wave them with `unknown` casts.
- Accessibility: interactive elements are keyboard-reachable and labeled. The color picker has
  visible names.
- No new dependency without a justification in the PR. Prefer the platform (Canvas, Workers,
  `crypto.subtle` where it fits).

## Tests

- Unit-test logic (encoding, state reducers, hooks). Test components through user-visible
  behavior with Testing Library.
- No snapshot tests for large component trees.

## Commands

```sh
cd frontend
npm ci
npm run lint
npm run typecheck
npm test
npm run build
```
