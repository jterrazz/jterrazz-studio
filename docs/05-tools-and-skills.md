# Tools and skills

`j` is registry-driven: the tools it installs and the AI agent skills it curates are declared as data in `src/internal/config/`, so the same catalogue powers `j install`, `j status`, and `j config`.

## The tools registry

89 tools across 18 families — package managers, runtimes, shell & terminal, CLI tools, git, editors & IDEs, containers & VMs, deploy, AI agents/tooling/apps, browsers, communication, productivity, media, remote access, security, and system utilities. Each tool declares its install method (brew, cask, npm, bun, manual), dependencies, version detection, and any post-install script. Both install and uninstall dispatch per method straight from the registry (brew formula/cask, npm, bun, uv), so a tool is one entry away from being toggleable in the `j install` TUI without bespoke command code. `j install`, `j upgrade`, and `j clean` (see [Commands](02-commands.md)) all act on this catalogue, and the Applications tab of `j status` renders it with live versions.

## The skills registry

The curated Claude Code / agent skills live in `src/internal/config/skills.go`, split into two lists:

- **StudioSkills** — the `@jterrazz` skills, the foundation of every project (toolbelt, stack, new-project, repo-structure, reach + its seo/geo/structured-data/content-reach domain skills, infra, typescript, broadcast, test, workflows).
- **CommunitySkills** — third-party skills worth having.

Each entry is a `{ repo, skill }` pair; `StudioRepos` / `CommunityRepos` carry the human-readable repo descriptions shown in the UI. `skills.go` is the single source of truth — adding an entry there is all it takes (plus a `make skills` to refresh the projected roster below); there is no lockfile.

The `j config` **Skills** tab installs, lists, and removes these by shelling out to the standalone **`skills`** CLI (it must be on `PATH`):

```sh
skills add <repo> -g -y --skill <name>   # install one skill globally
skills list -g                           # list installed
skills remove -g -y <name>               # remove
```

See `src/internal/domain/skill/` for the integration.

> **Note:** the toolchain skill set follows the `@jterrazz/typescript@6` consolidation — `@jterrazz/codestyle` (lint/format) was absorbed into `@jterrazz/typescript`, so a single `jterrazz-typescript` skill now covers build, lint, format, and docs, alongside `jterrazz-repo-structure` for repo doctrine.

## The toolbelt skill — a committed projection

The registries are also projected into an agent-facing skill: [`skills/jterrazz-toolbelt/`](../skills/jterrazz-toolbelt/SKILL.md) carries the full tool roster and the modern-replacement table, so any agent starting work on a jterrazz machine knows what is available. (The curated skills list is deliberately not projected — an installed skill is already injected into the agent.) The sections between its `GENERATED` markers are compiled from `src/internal/config/` by **`make skills`** (`src/cmd/skillsgen/`) — never edited by hand — and a Go sync test run by `make test` fails when the registry and the skill drift. Change a registry, run `make skills`, commit both in the same change.

This is the repo's one compiler in the sense of the repo-structure doctrine: the registry is the source layer, the skill roster is its committed, diffable projection.

## Related

- [Commands](02-commands.md) — `j install` / `j upgrade` / `j clean`.
- [Configuration](04-configuration.md) — the Skills tab.
