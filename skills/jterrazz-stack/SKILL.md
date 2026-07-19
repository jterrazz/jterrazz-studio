---
name: jterrazz-stack
description: Overview of the @jterrazz ecosystem — shared npm packages, project naming ({product}-{role}, package-*), required files, CI/CD workflows, testing and architecture conventions. Use when working on any jterrazz project, choosing a package, scaffolding decisions, or understanding how the stack composes.
---

# @jterrazz Stack

The map of the `@jterrazz` ecosystem: composable packages, one set of conventions, interchangeable repos. Canonical home: `jterrazz-studio` `docs/07-stack.md`.

## Mental model

- **Five shared packages, one concern each** — `@jterrazz/typescript` (toolchain), `@jterrazz/test` (testing), `@jterrazz/logger`, `@jterrazz/intelligence`, `@jterrazz/broadcast`. CI/CD is shared workflows from `jterrazz-actions`; infra is `jterrazz-infra`.
- **Two project shapes** — libraries (`package-*`, published `@jterrazz/*`, ports & adapters) and applications (`{product}-{role}`, hexagonal). Both build/lint/test through `@jterrazz/typescript` and `@jterrazz/test`.
- **Same required files everywhere** — Makefile (`build`/`lint`/`test`), tsconfig + lint configs extending the shared presets, the validate workflow, and the root files the repo-structure doctrine mandates.

## Where to look

| Task                                                | Read                                            |
| --------------------------------------------------- | ----------------------------------------------- |
| The whole picture: packages, naming, required files | `jterrazz-studio` `docs/07-stack.md`            |
| Build / lint / format / API docs                    | `jterrazz-typescript` skill                     |
| Testing conventions and structure                   | `jterrazz-test` skill                           |
| Scaffolding a new project                           | `jterrazz-new-project` skill                    |
| CI/CD workflows                                     | `jterrazz-workflows` skill                      |
| Where knowledge lives in a repo                     | `jterrazz-repo-structure` skill                 |
| What is installed on the machine                    | `jterrazz-toolbelt` skill                       |

## Always

- `npm` only (never pnpm/yarn), Node.js 24, ESM only, `.js` extensions in Node imports.
- Author: `Jean-Baptiste Terrazzoni <contact@jterrazz.com>`.
- Never restate another skill's territory here — route to it.
