---
name: jterrazz-new-project
description: Scaffolds a new @jterrazz project with all conventions — jterrazz.config.ts, package.json, the derived editor files, the Makefile's six gestures, CI workflows. Use when creating a new app or library from scratch.
---

# Create a new @jterrazz project

A project is born declared: one `jterrazz.config.ts`, the derived editor files, a Makefile answering the six gestures. How each is written — the steps, the files whole, the first `make check` — is `@jterrazz/typescript`'s [Getting started](https://github.com/jterrazz/package-typescript/blob/main/docs/start/getting-started.md), and where a test lives is `@jterrazz/test`'s [Set up a project](https://github.com/jterrazz/package-test/blob/main/docs/start/set-up-a-project.md). Follow them; this skill holds only what the estate adds.

## For a new library (package-{name})

```bash
mkdir package-{name} && cd package-{name}
git init
npm init -y
```

The identity `package.json` carries beside what Getting started writes:

```json
{
    "name": "@jterrazz/{name}",
    "version": "0.1.0",
    "author": "Jean-Baptiste Terrazzoni <contact@jterrazz.com>",
    "type": "module",
    "files": ["dist"],
    "publishConfig": { "registry": "https://registry.npmjs.org/" },
    "repository": { "type": "git", "url": "https://github.com/jterrazz/package-{name}" }
}
```

`jterrazz.config.ts` declares `kind: 'library'`, the platform it runs on, and `presets: [test()]`.

## For a new application ({product}-{role})

`"private": true`, no `exports` or `publishConfig`; `jterrazz.config.ts` declares the application's kind — `api`, `cli`, `website` or `mobile` — and its platform.

**The composition root.** An application is hexagonal — `typescript: { architecture: { hexagonal: true } }` in `jterrazz.config.ts` — and is built in one place, `src/composition/compose.ts`: `compose(config, overrides?: Partial<Ports>)` calls the constructors and returns the typed graph with its `[Symbol.asyncDispose]`. No DI container library. Only the entry and the spec runners import it:

```ts
// src/index.ts
import { compose } from './composition/compose.js';

await using app = compose();
await app.server.start(app.configuration.http);
```

The law and its why are `@jterrazz/typescript`'s [The composition root](https://github.com/jterrazz/package-typescript/blob/main/docs/concepts/the-composition-root.md); how a spec runner builds the app is `@jterrazz/test`'s [API specs](https://github.com/jterrazz/package-test/blob/main/docs/guides/api-specs.md#how-the-runner-builds-the-app).

## CI workflows

**.github/workflows/validate.yaml:**

```yaml
name: Validate
on:
    push:
        branches: [main]
    pull_request:
        branches: [main]
jobs:
    validate:
        uses: jterrazz/jterrazz-actions/.github/workflows/validate.yaml@main
        with:
            node-version: '24'
```

For libraries, also add **.github/workflows/publish.yaml:**

```yaml
name: Publish
on:
    release:
        types: [created]
permissions:
    contents: read
    id-token: write
jobs:
    publish:
        uses: jterrazz/jterrazz-actions/.github/workflows/release-npm.yaml@main
        with:
            node-version: '24'
        secrets: inherit
```

## Repo structure — born structured

A new project starts with the corpus + injection layout, not just code. Do not author layout advice here — follow the **`jterrazz-repo-structure`** skill, whose canonical page is `jterrazz-studio` `docs/concepts/repo-structure.md`:

- `README.md` — the vitrine (what it is, install, a pointer into `docs/`), not a manual.
- `AGENTS.md` — the agent brief (mental model + routing table into `docs/`); `CLAUDE.md` is a **symlink** to it (`ln -s AGENTS.md CLAUDE.md`).
- `docs/` — the manual cut by mood: `start/`, `guides/`, `concepts/`, `reference/`, `decisions/`; `concepts/architecture.md`, `guides/develop.md` and `guides/test.md` from day one.

## Never

- Never use pnpm or yarn
- Never skip the Makefile
- Never put business logic in `src/index.ts` — a library's only exports, an application's only composes and serves
- Never author repo-layout advice in a skill — route to `jterrazz-repo-structure`
