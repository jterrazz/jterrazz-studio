---
name: jterrazz-stack
description: Overview of the @jterrazz ecosystem: shared npm packages, project naming ({product}-{role}, package-*), required files, CI/CD workflows, testing and architecture conventions. Use when working on any jterrazz project, choosing a package, scaffolding decisions, or understanding how the stack composes.
---

# @jterrazz Stack

The map of the `@jterrazz` ecosystem: composable packages, one set of conventions, interchangeable repos. Canonical home: `jterrazz-studio` `docs/concepts/stack.md`.

## Mental model

- **Shared packages, one concern each**: `@jterrazz/typescript` (toolchain) and `@jterrazz/test` (testing) under every project; `telemetry`, `intelligence`, `analytics`, `broadcast`, `manifest`, `attestation` where a project needs them. CI/CD is shared workflows from `jterrazz-actions`; infra is `jterrazz-infrastructure`.
- **Two project shapes**: libraries (`package-*`, published `@jterrazz/*`, ports & adapters) and applications (`{product}-{role}`, hexagonal). Both build/lint/test through `@jterrazz/typescript` and `@jterrazz/test`.
- **A project declares what it is once**: `jterrazz.config.ts` (kind, platform, presets, exceptions with reasons); the editor files are one-line derivations, the Makefile answers the six gestures (`install`, `build`, `lint`, `fix`, `test`, `check`), the validate workflow is shared, and the root files are the repo-structure doctrine's.

## Where to look

| Task                                                | Read                                       |
| --------------------------------------------------- | ------------------------------------------ |
| The whole picture: packages, naming, required files | `jterrazz-studio` `docs/concepts/stack.md` |
| Build / lint / format / API docs                    | `jterrazz-typescript` skill                |
| Testing conventions and structure                   | `jterrazz-test` skill                      |
| Scaffolding a new project                           | `jterrazz-new-project` skill               |
| CI/CD workflows                                     | `jterrazz-workflows` skill                 |
| Where knowledge lives in a repo                     | `jterrazz-repo-structure` skill            |
| What is installed on the machine                    | `jterrazz-toolbelt` skill                  |

## Always

- `npm` only (never pnpm/yarn), Node.js 24, ESM only, `.js` extensions in Node imports.
- Author: `Jean-Baptiste Terrazzoni <contact@jterrazz.com>`.
- Never restate another skill's territory here; route to it.
