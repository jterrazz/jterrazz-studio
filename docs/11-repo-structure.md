# Repo structure

A repository is a written corpus, some thin injection layers, and a compiler that projects source into machine-readable forms. Keep those three roles distinct and knowledge stays in exactly one place.

This is the doctrine every `@jterrazz` repo follows, whatever the language. It is the canonical home for it.

## The corpus belongs to the repository

A repository's own knowledge lives in its `docs/`. Anything a change in this repository would make false — its architecture, its layers, its conventions, its runbook, its quirks, the decisions it alone took — is authored here and nowhere else.

An outside corpus holds two things about a repository and no more: a pointer to it, and the contracts it has with other repositories. The OS wiki (`jterrazz-os`, `home/<brand>/wiki/`) is the outside corpus in practice — it carries one row per repository and the agreements that span several of them. When a sentence there could be falsified by one repository alone, it is in the wrong place.

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

## Canonical layout

| Path              | Role                                                                                          |
| ----------------- | --------------------------------------------------------------------------------------------- |
| `README.md`       | The vitrine — what the repo is, install, a pointer into `docs/`. Not a manual.                |
| `AGENTS.md`       | The agent brief — mental model + routing table into `docs/`. Routes, does not retell.         |
| `CLAUDE.md`       | A symlink to `AGENTS.md`.                                                                     |
| `TODO.md`         | The working backlog, when one is kept.                                                        |
| `docs/README.md`  | The map of the corpus — one sentence, then a table of the chapters.                           |
| `docs/NN-*.md`    | The chapters, numbered, one subject each, plus any committed projections.                     |
| `docs/decisions/` | The repo-local decision records, `NNN-title.md`, beside the `_template.md` they are cut from. |
| `docs/reference/` | Committed projections — a compiler's output, never authored by hand.                          |
| `docs/_assets/`   | Files a chapter embeds. Underscored, because it is ground rather than a subject.              |
| `skills/`         | Injection layer for agents — one skill per capability, routes into `docs/`.                   |

The list is exhaustive for the root: a file that is not on it does not belong there, and a repository carries no `CHANGELOG.md`. Those three subfolders are exhaustive too — a fourth folder under `docs/` means a subject is hiding from the numbering.

`docs/README.md` is the one file an outside corpus points at. It is the entry every reader lands on, so a chapter can be renumbered or split without breaking a link held somewhere else.

### The spine

Every repository's `docs/` opens on the same four numbers, whatever the language and whatever the product. A reader landing on an unfamiliar repository knows where the shape is before opening a file, and so does an agent.

| Chapter              | Owns                                                                                      |
| -------------------- | ----------------------------------------------------------------------------------------- |
| `01-architecture.md` | The shape of the thing — its parts, their boundaries, and why the lines are drawn there.  |
| `02-developing.md`   | How a change is made — the toolchain, the loop, which file a change opens, what it owes.  |
| `03-testing.md`      | What proves a change — the suites, the ground they stand on, how a golden is regenerated. |
| `04-operating.md`    | How it ships and how it runs — the release or the deploy, and the footprint it leaves.    |

The repository's own subjects follow, numbered contiguously from `05`: one subject per chapter, in whatever order the repository reads best. A chapter's name is a subject, never a moment — `exploration`, `review`, `notes`, `draft`, `wip` name a date rather than a thing, and what happened on a date is a decision record or git's history, not a chapter.

`01`, `02` and `03` are required of every repository: everything has a shape, a way of being changed, and something that proves the change — a repository that cannot fill one of them has found a real hole, not an exemption.

`04-operating.md` is required when the repository ships something that runs — an image it deploys, a package it publishes, a binary it releases, a platform it provisions. The rule only ever requires it; it never forbids it. Whether a repository ships is not reliably readable from its tree — a tagged release or a deploy workflow lives in `.github/`, in a dialect a checker has no business parsing — so a repository that ships that way writes its `04` and nothing has to ask it to, and a library that writes a hollow one gets a review comment, not a red gate.

This chapter owns the spine. Where a checker holds part of it mechanically, the rule ids and the sentence each one prints belong to that checker, so that the message and the rule cannot drift apart — this page never lists them.

## Decision records

A decision this repository alone took is recorded in `docs/decisions/NNN-title.md`, numbered in the order the decisions were taken. A record carries a status and a date, then three sections: the context that forced the decision, what was decided, and the consequences that follow. The mold sits where the record is written — `docs/decisions/_template.md`, copied for each new record.

The status is `Proposed` until the repository's owner accepts it. A record that replaces an earlier one says so, and the record it replaces is marked superseded by it — the link closes both ways, so a reader who lands on the old one is told where the live decision is.

A decision that spans two or more repositories is not recorded here. It belongs to the corpus that spans them, and this repository links to it.

## Packages vs applications

A `docs/reference/` API projection exists for **packages** — a published library has a public API surface that others consume, and the compiler keeps its readable form in lockstep with the code (see `package-typescript` [Docs pipeline](https://github.com/jterrazz/package-typescript/blob/main/docs/08-docs-pipeline.md)). An **application** (an API server, a product CLI, a web app) has no API consumers: it adopts the doctrine in full — the corpus, the routing layers, the single-home rule — and never generates `docs/reference/`. An application may still have its own compiler when it holds another source layer worth projecting, as this repo does with its install registry.

## Language specifics

The doctrine applies in full to every repo: a written corpus, thin routing layers, and knowledge that lives in exactly one place. Only the _compiler_ half is language-specific — TypeScript packages use `typescript docs`; a Go, Rust, or infra repo substitutes its own generator or has none. The projection-vs-presentation criterion is unchanged: a cross-layer compile is committed, a same-layer re-packaging is a presentation built only when a delivery target exists.

## Related

- [The stack](10-stack.md) — how `@jterrazz` projects compose.
- [Tools and skills](08-tools-and-skills.md) — this repo's own compiler, `make skills`.
