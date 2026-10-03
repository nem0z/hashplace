# 0001. Record architecture decisions

- **Status**: Accepted
- **Date**: 2026-10-04
- **Issue/PR**: dev workflow PR

## Context

hashplace is written entirely by AI agents across many independent sessions. Agents start each
session without memory of past discussions, so the reasoning behind decisions has to be written
down where both agents and the maintainer can find it.

## Considered options

1. **ADRs in the repository** (`docs/adr/`)
   - Pros: versioned with the code, reviewed in PRs, readable by agents.
   - Cons: needs discipline to keep writing them.
2. **Decisions recorded only in issue and PR comments**
   - Pros: no extra files.
   - Cons: hard to discover, and agents rarely read closed threads.

## Decision

We use lightweight ADRs in `docs/adr/`, following `0000-template.md` and the `write-adr` skill.
Normative game and protocol rules live separately in `docs/spec/`, because they describe *what*
the system does, while ADRs record *why* it was built that way.

## Consequences

PRs that make an architectural decision must include an ADR. The reviewer should ask for one
when it is missing.
