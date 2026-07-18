# Tools and skills

`j` is registry-driven: the tools it installs and the AI agent skills it curates are declared as data in `src/internal/config/`, so the same catalogue powers `j install`, `j status`, and `j config`.

## The tools registry

100+ tools across 7 categories — package managers, runtimes, devops, AI, terminal, GUI apps, and Mac App Store. Each tool declares its install method (brew, cask, npm, bun, manual), dependencies, version detection, and any post-install script. `j install`, `j upgrade`, and `j clean` (see [Commands](02-commands.md)) all act on this catalogue, and the Applications tab of `j status` renders it with live versions.

## The skills registry

The curated Claude Code / agent skills live in `src/internal/config/skills.go`, split into two lists:

- **StudioSkills** — the `@jterrazz` skills, the foundation of every project (stack, new-project, infra, typescript, repo-structure, broadcast, test, actions).
- **CommunitySkills** — third-party skills worth having.

Each entry is a `{ repo, skill }` pair; `StudioRepos` / `CommunityRepos` carry the human-readable repo descriptions shown in the UI. `skills.go` is the single source of truth — adding an entry there is all it takes; there is no lockfile or golden to keep in sync.

The `j config` **Skills** tab installs, lists, and removes these by shelling out to the standalone **`skills`** CLI (it must be on `PATH`):

```sh
skills add <repo> -g -y --skill <name>   # install one skill globally
skills list -g                           # list installed
skills remove -g -y <name>               # remove
```

See `src/internal/domain/skill/` for the integration.

> **Note:** the toolchain skill set follows the `@jterrazz/typescript@6` consolidation — `@jterrazz/codestyle` (lint/format) was absorbed into `@jterrazz/typescript`, so a single `jterrazz-typescript` skill now covers build, lint, format, and docs, alongside `jterrazz-repo-structure` for repo doctrine.

## Related

- [Commands](02-commands.md) — `j install` / `j upgrade` / `j clean`.
- [Configuration](04-configuration.md) — the Skills tab.
