# Agent brief — `j` (jterrazz-cli)

A Go CLI (Cobra + Bubble Tea TUIs) that bootstraps and manages a macOS dev machine. This file **routes**; it does not restate what the corpus already says.

## Mental model

- **Registry-driven.** The tools it installs, the config items, and the curated agent skills are declared as data in `src/internal/config/` — the same catalogue powers `j install`, `j status`, and `j config`. Add to the registry, not to bespoke command code.
- **The machine registry is the source of truth.** `~/.jterrazz/config.json` (aliases, roles, ssh) drives everything; adding a machine also writes a managed `~/.ssh/config` block. Role (`client`/`server`) gates what `status` reports and which `config` items appear.
- **Dotfiles are versioned.** `dotfiles/applications/*` are installed onto the machine by `j config`.
- **Specs drive the real binary.** `specs/cli/` uses `@jterrazz/test` (`specification.cli`) against a freshly built `j`; the runner rebuilds when any `src/**/*.go` is newer than the test binary.

## Where knowledge lives (route here first)

The corpus is `docs/` + `README.md`. Do not duplicate it — link to it.

| Working on…                              | Read                          |
| ---------------------------------------- | ----------------------------- |
| Install, user data, dev/release, layout  | `docs/01-getting-started.md`  |
| status / install / upgrade / clean / run | `docs/02-commands.md`         |
| Machine registry, `config.json`, remote  | `docs/03-machines.md`         |
| The `j config` TUI, items, categories    | `docs/04-configuration.md`    |
| Tool + skill registries (`config/`)      | `docs/05-tools-and-skills.md` |
| Dotfiles (`dotfiles/applications/`)      | `docs/06-dotfiles.md`         |

Repo doctrine (corpus / injection layers / compiler — where knowledge lives) is the shared `jterrazz-repo-structure` skill / `package-typescript` `docs/06-repo-structure.md`. This is an application: it adopts the corpus + injection layers and never generates a `docs/reference/` projection.

## Setup & commands

```bash
make build     # Build ./j
make test      # Go unit tests
make test-e2e  # npm install + rebuild j + vitest --run (the specs/cli suite)
make lint      # golangci-lint
```

The e2e specs need `@jterrazz/test` (npm) and rebuild `j` via an mtime check — see `specs/cli/cli.specification.ts`.

## Repo layout

```
src/cmd/j/main.go            # entry point
src/internal/commands/       # Cobra commands
src/internal/config/         # tool, script, command, and SKILL registries (skills.go)
src/internal/domain/         # version parsing, status loading, skill install (shells to the `skills` CLI)
src/internal/presentation/   # Bubble Tea TUI views, components, theme
dotfiles/applications/       # versioned app configs installed by `j config`
docs/                        # the corpus (this brief routes into it)
skills/                      # Claude Code skills for the @jterrazz workflow
specs/cli/                   # end-to-end specs (@jterrazz/test)
```

## Standing rules

- The curated skills list lives in **`src/internal/config/skills.go`** (`StudioSkills` + `StudioRepos`). It is the single source of truth — no lockfile, no golden; a new `{ repo, skill }` entry is all that's needed.
- A change to product behaviour or a command's output updates the matching `docs/` chapter in the same change, and regenerates any affected `specs/cli/**/expected/*` golden with `TEST_UPDATE=1` (deliberately).
- Never author repo-layout doctrine here or in a skill — route to `jterrazz-repo-structure`.
