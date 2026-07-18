---
name: jterrazz-stack
description: Overview of the @jterrazz ecosystem — shared npm packages, naming conventions, project patterns, and how everything composes together. Activates when working on any jterrazz project, choosing packages, or understanding the stack.
---

# @jterrazz Stack

The @jterrazz ecosystem is a set of composable packages that define how every project builds, lints, tests, logs, and deploys. Each package owns one concern and all projects follow the same conventions.

## Packages

| Package | Purpose | npm script |
|---------|---------|------------|
| `@jterrazz/typescript` | Toolchain — build (tsdown) + quality (tsgo, oxlint, oxfmt, knip) + docs | `typescript build`, `typescript bundle`, `typescript dev`, `typescript start`, `typescript check`, `typescript fix`, `typescript docs` |
| `@jterrazz/test` | Testing framework — conventions, structure, mocking | `vitest --run` |
| `@jterrazz/logger` | Structured logging (pino) | — |
| `@jterrazz/intelligence` | AI toolkit (OpenRouter, Langfuse) | — |
| `@jterrazz/broadcast` | Multi-channel announcements (App Store, push) | — |

## Project types

**Library** (`package-*`):
```json
{
  "build": "typescript bundle",
  "lint": "typescript check",
  "lint:fix": "typescript fix",
  "test": "vitest --run"
}
```

**Application** (`signews-api`, `signews-broadcast`, etc.):
```json
{
  "build": "typescript build",
  "start": "typescript start",
  "dev": "typescript dev",
  "lint": "typescript check",
  "lint:fix": "typescript fix",
  "test": "vitest --run"
}
```

## Naming conventions

```
{product}-{role}
├── signews-api          # Backend API
├── signews-web          # Web client
├── signews-mobile       # iOS/Android app
├── signews-broadcast    # Event broadcaster
├── signews-blueprint    # Architecture docs
└── package-{name}       # Shared @jterrazz/* packages
```

Roles: `-api`, `-web`, `-mobile`, `-broadcast`, `-blueprint`, `-landing`

## Required files

Every project must have:
- `Makefile` with `build`, `lint`, `test` targets
- `tsconfig.json` extending `@jterrazz/typescript/tsconfig/node`
- `oxlint.config.ts` + `oxfmt.config.ts` importing presets from `@jterrazz/typescript`
- `.github/workflows/validate.yaml` using shared workflow

## Repo structure

Where knowledge lives — README vitrine, `AGENTS.md` router + `CLAUDE.md` symlink, the `docs/` corpus, and (for packages) the generated `docs/reference/` — is one shared doctrine. Do not restate it here: follow the **`jterrazz-repo-structure`** skill / `package-typescript` `docs/06-repo-structure.md`.

## CI/CD

Shared workflows from `jterrazz/jterrazz-actions`:
- `validate.yaml` — runs `make build`, `make lint`, `make test`
- `release-docker.yaml` — Docker build + Helm deploy
- `release-npm.yaml` — npm publish with OIDC provenance

## Testing convention

Defined by `@jterrazz/test` — all projects follow the same structure:

- **Unit** (`thing.test.ts`) — colocated next to source, no I/O
- **Integration** (`thing.integration.test.ts`) — testcontainers + in-process app
- **E2E** (`thing.e2e.test.ts`) — full docker compose, real HTTP
- Docker required — `docker/compose.test.yaml` defines test infrastructure
- Service factories: `postgres({ compose: "db" })`, `redis({ compose: "cache" })`
- One folder per test, data colocated (`seeds/`, `requests/`, `responses/`, `mock/`, `expected/`)
- `tests/setup/` for infrastructure, `tests/fixtures/` for sample apps, `tests/helpers/` for utilities

## Architecture pattern

Libraries use **ports & adapters**:
- `src/ports/` — interfaces
- `src/adapters/` — implementations
- `src/index.ts` — public API exports

Apps use **hexagonal architecture**:
- `src/domain/` — pure business logic
- `src/application/` — use cases, ports
- `src/infrastructure/` — adapters (HTTP, DB, external APIs)

## Always

- Use `npm` (not pnpm/yarn) — all repos have `package-lock.json`
- Node.js 24
- ESM only (`"type": "module"`)
- `.js` extensions in imports for Node.js projects
- Author: `Jean-Baptiste Terrazzoni <contact@jterrazz.com>`
