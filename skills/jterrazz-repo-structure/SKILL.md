---
name: jterrazz-repo-structure
description: Use when organizing or restructuring a repository's documentation, skills, AGENTS.md/CLAUDE.md, README, or generated docs — where knowledge lives, what gets committed, corpus vs injection layers vs compiler, projections vs presentations.
---

# Repo structure

The doctrine every `@jterrazz` repo follows for where knowledge lives. Canonical home: `jterrazz-studio` `docs/08-repo-structure.md`.

## Mental model

A repo is three distinct roles — conflating them is the failure mode this skill guards against:

1. **The corpus** — `docs/` (numbered chapters) + `README.md`. Where knowledge is _authored_.
2. **The injection layers** — `AGENTS.md`/`CLAUDE.md` (agent brief), `skills/` entries. They **route**, they never author.
3. **The compiler** — a repo-specific projector, when there is a source layer worth projecting: `typescript docs` → `docs/reference/` in TypeScript packages, `make skills` → the toolbelt rosters in jterrazz-studio. No source layer, no compiler.

Scope: an API-reference projection applies to **packages** (a public API surface to project). An **application** adopts roles 1–2 in full and never generates `docs/reference/` — though it may run its own compiler over another source layer. Non-TypeScript repos: same doctrine, the compiler half substitutes the language's own generator or is absent.

Three golden rules, compressed:

- **Written once** — a piece of knowledge lives in one chapter; everything else links to it. No duplication between README, a skill, and a chapter.
- **Route, don't retell** — AGENTS.md/CLAUDE.md and skills are maps: a mental model + a routing table into `docs/`. When they start explaining, knowledge has leaked out of the corpus.
- **The compiler projects, never authors** — generated files are projections, stamped `DO-NOT-EDIT`. Change the source and regenerate; never hand-edit a projection.

Projection vs presentation — the line is which layers a derived file spans. A **projection crosses a layer** (source code → docs): it is **committed** and sync-checked, diffable and agent-readable from the tree. A **presentation re-packages within the docs layer** (an `llms.txt` concatenation, a rendered HTML site): it authors nothing new, is never committed, and is built in CI only if a delivery target exists — with none, no presentation is produced and the corpus is read directly.

## Canonical root-file layout

| File           | Role                                                                                  |
| -------------- | ------------------------------------------------------------------------------------- |
| `README.md`    | The vitrine — what the repo is, install, a pointer into `docs/`. Not a manual.        |
| `AGENTS.md`    | The agent brief — mental model + routing table into `docs/`. Routes, does not retell. |
| `CLAUDE.md`    | A symlink to `AGENTS.md`.                                                             |
| `CHANGELOG.md` | The release history.                                                                  |
| `TODO.md`      | The working backlog, when one is kept.                                                |
| `docs/`        | The corpus: numbered chapters, plus any committed projections.                        |
| `skills/`      | Injection layer for agents — one skill per capability, routes into `docs/`.           |

## Where to look

| Task                                                              | Read                                              |
| ----------------------------------------------------------------- | ------------------------------------------------- |
| The doctrine itself (three roles, golden rules, layout)           | `jterrazz-studio` `docs/08-repo-structure.md`     |
| TypeScript compiler mechanics (`typescript docs`, sync-checking)  | `package-typescript` `docs/05-docs-pipeline.md`   |
| The studio's own compiler (`make skills`, the toolbelt rosters)   | `jterrazz-studio` `docs/05-tools-and-skills.md`   |

How `@jterrazz` projects compose (packages, naming, CI) is a separate capability: see the `jterrazz-stack` skill.

## Always

- Never restate a chapter's content in a skill or in AGENTS.md/CLAUDE.md — link to it instead.
- Never hand-edit a projection — regenerate it (`typescript docs`, `make skills`).
- Regenerate a projection in the same change that touches its source, so the sync check stays green.
- One skill per capability, split by genuinely distinct trigger conditions — not one skill restating two chapters.
