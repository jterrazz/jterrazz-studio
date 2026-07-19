# The stack

How every `@jterrazz` project composes: the shared packages, the naming scheme, the required files, and the conventions that make repos interchangeable. Deep knowledge stays in each package's own corpus — this chapter is the map of the whole.

## Shared packages

| Package                  | Owns                                                        | Corpus                        |
| ------------------------ | ----------------------------------------------------------- | ----------------------------- |
| `@jterrazz/typescript`   | Toolchain — build, quality checks, lint presets, API docs   | `package-typescript/docs/`    |
| `@jterrazz/test`         | Testing — conventions, structure, mocking, service factories | `package-test/docs/`          |
| `@jterrazz/logger`       | Structured logging (pino)                                   | `package-logger`              |
| `@jterrazz/intelligence` | AI toolkit (OpenRouter, Langfuse)                           | `package-intelligence`        |
| `@jterrazz/broadcast`    | Multi-channel announcements (App Store, push)               | `package-broadcast/docs/`     |

Shared CI/CD lives in `jterrazz/jterrazz-actions` (validate, release-npm, release-docker, release-go); infrastructure in `jterrazz/jterrazz-infra` (K3s, Helm, Traefik).

## Project types and naming

```
{product}-{role}             # applications: signews-api, signews-web, signews-mobile…
package-{name}               # shared libraries, published as @jterrazz/*
```

Roles: `-api`, `-web`, `-mobile`, `-broadcast`, `-blueprint`, `-landing`.

**Library** npm scripts: `build: typescript bundle`, `lint: typescript check`, `lint:fix: typescript fix`, `test: vitest --run`.
**Application** adds `start: typescript start` and `dev: typescript dev`, and builds with `typescript build` instead of `bundle`.

## Required files

Every project carries:

- `Makefile` with `build`, `lint`, `test` targets.
- `tsconfig.json` extending `@jterrazz/typescript/tsconfig/node`.
- `oxlint.config.ts` + `oxfmt.config.ts` importing the `@jterrazz/typescript` presets.
- `.github/workflows/validate.yaml` using the shared workflow (runs `make build`, `make lint`, `make test`).
- The root files mandated by the [repo-structure doctrine](08-repo-structure.md): README vitrine, `AGENTS.md` (+ `CLAUDE.md` symlink), a `docs/` corpus.

## Testing

The convention is defined by `@jterrazz/test`, and every project follows it: colocated `*.test.ts` units (no I/O), `*.integration.test.ts` against testcontainers, `*.e2e.test.ts` against a real compose stack, with data colocated per test. The full convention lives in `package-test`'s corpus — route there, don't restate it.

## Architecture

Libraries use **ports & adapters** (`src/ports/` interfaces, `src/adapters/` implementations, `src/index.ts` barrel). Applications use **hexagonal** (`src/domain/`, `src/application/`, `src/infrastructure/`) — enforced by the `hexagonal` oxlint preset, documented in `package-typescript/docs/04-lint-presets.md`.

## Standing conventions

- `npm` only — never pnpm or yarn; every repo commits `package-lock.json`.
- Node.js 24, ESM only (`"type": "module"`), `.js` extensions in Node imports.
- Author: `Jean-Baptiste Terrazzoni <contact@jterrazz.com>`.

## Related

- [Repo structure](08-repo-structure.md) — where knowledge lives in every repo.
- [Tools and skills](05-tools-and-skills.md) — the machine-side registries that distribute the skills.
