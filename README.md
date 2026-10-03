# hashplace

A collaborative pixel canvas in the spirit of r/place, but **paid for in hashes instead of a cooldown timer**.

To paint a cell you solve a proof-of-work. The harder the work, the stronger your claim. Claims decay over time, so every pixel eventually becomes cheap to take back.

> Status: early design.

## Stack

- **Backend**: Go (the main focus)
- **Frontend**: React + TypeScript

## How this project is built

All code is written by AI agents. Every change lands through a GitHub pull request that the maintainer (@nem0z) reviews and approves.

- [`AGENTS.md`](AGENTS.md): rules and entry point for agents
- [`CONTRIBUTING.md`](CONTRIBUTING.md): branch, commit and PR standards
- [`docs/workflow.md`](docs/workflow.md): from idea to merged PR
- [`docs/adr/`](docs/adr/): architecture decisions
