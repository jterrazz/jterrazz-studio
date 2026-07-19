# Repo structure

A repository is a written corpus, some thin injection layers, and a compiler that projects source into machine-readable forms. Keep those three roles distinct and knowledge stays in exactly one place.

This is the doctrine every `@jterrazz` repo follows, whatever the language. It is the canonical home for it.

## The three roles

1. **The corpus** — `docs/` (numbered chapters) plus `README.md`. This is where knowledge is _authored_: prose a human writes and maintains.
2. **The injection layers** — `AGENTS.md` (agent brief, with `CLAUDE.md` as a symlink) and the `skills/` entries. They **route**, they never author. A skill points at a chapter; it does not restate what the chapter says.
3. **The compiler** — a repo-specific command that projects a source layer into docs. In `package-typescript`, `typescript docs` projects the source barrel into `docs/reference/`; in `jterrazz-studio`, `make skills` projects the install registry into the toolbelt skill rosters. A repo with nothing to project simply has no compiler.

## Three golden rules

- **A piece of knowledge is written once.** It lives in one chapter. Everything else links to it. No paragraph is duplicated between the README, a skill, and a chapter.
- **AGENTS.md and skills route without retelling.** They are maps, not territory: a mental model plus a routing table into `docs/`. When they start explaining, the knowledge has leaked out of the corpus.
- **The compiler projects; it never authors.** Generated files are projections of a source layer, stamped `DO-NOT-EDIT`. You change the source and regenerate — you never edit a projection.

## Committed projections vs presentations

The dividing line is which layers a derived file spans:

- **A projection crosses a layer boundary** — it compiles one layer into another (source code → docs). It is **committed**, so it is diffable, greppable, and agent-readable straight from the tree, and kept honest by a sync check (`typescript docs --check`, the skillsgen sync test).
- **A presentation re-packages within one layer** — it re-arranges docs into another docs shape (an `llms.txt` concatenation, a rendered HTML site). It authors no new knowledge, so it is never committed; it is **built in CI only if a delivery target exists** (a hosted site, an agent-ingestion endpoint). With no such target, no presentation is produced at all — the committed corpus is read directly.

## Root file conventions

| File           | Role                                                                                   |
| -------------- | -------------------------------------------------------------------------------------- |
| `README.md`    | The vitrine — what the repo is, install, a pointer into `docs/`. Not a manual.          |
| `AGENTS.md`    | The agent brief — mental model + routing table into `docs/`. Routes, does not retell.   |
| `CLAUDE.md`    | A symlink to `AGENTS.md`.                                                               |
| `CHANGELOG.md` | The release history.                                                                    |
| `TODO.md`      | The working backlog, when one is kept.                                                  |
| `docs/`        | The corpus: numbered chapters, plus any committed projections.                          |
| `skills/`      | Injection layer for agents — one skill per capability, routes into `docs/`.             |

## Packages vs applications

A `docs/reference/` API projection exists for **packages** — a published library has a public API surface that others consume, and the compiler keeps its readable form in lockstep with the code (see `package-typescript` [Docs pipeline](https://github.com/jterrazz/package-typescript/blob/main/docs/05-docs-pipeline.md)). An **application** (an API server, a product CLI, a web app) has no API consumers: it adopts the doctrine in full — the corpus, the routing layers, the single-home rule — and never generates `docs/reference/`. An application may still have its own compiler when it holds another source layer worth projecting, as this repo does with its install registry.

## Language specifics

The doctrine applies in full to every repo: a written corpus, thin routing layers, and knowledge that lives in exactly one place. Only the _compiler_ half is language-specific — TypeScript packages use `typescript docs`; a Go, Rust, or infra repo substitutes its own generator or has none. The projection-vs-presentation criterion is unchanged: a cross-layer compile is committed, a same-layer re-packaging is a presentation built only when a delivery target exists.

## Related

- [The stack](07-stack.md) — how `@jterrazz` projects compose.
- [Tools and skills](05-tools-and-skills.md) — this repo's own compiler, `make skills`.
