# Agent brief — `j` (jterrazz-studio)

A Go CLI (Cobra + Bubble Tea TUIs) that bootstraps and manages a macOS dev machine. This file **routes**; it does not restate what the corpus already says.

## Mental model

- **Registry-driven.** The tools it installs, the config items, and the curated agent skills are declared as data in `src/internal/config/` — the same catalogue powers `j install`, `j status`, and `j config`. Add to the registry, not to bespoke command code.
- **The machine registry is the source of truth.** `~/.jterrazz/config.json` (aliases, roles, ssh) drives everything; adding a machine also writes a managed `~/.ssh/config` block. Role (`client`/`server`) gates what `status` reports and which `config` items appear.
- **Dotfiles are versioned.** `dotfiles/applications/*` are installed onto the machine by `j config`.
- **Specs drive the real binary.** `specs/cli/` uses `@jterrazz/test` (`specification.cli`) against a freshly built `j`; the runner rebuilds when any `src/**/*.go` is newer than the test binary. A scenario is a `<case>.spec.yaml` document, not code — see `docs/03-testing.md`.

## Where knowledge lives (route here first)

The corpus is `docs/` + `README.md`, mapped by `docs/README.md`. Do not duplicate it — link to it.

| Working on…                               | Read                          |
| ----------------------------------------- | ----------------------------- |
| Layers, registries, the Cobra root        | `docs/01-architecture.md`     |
| Building, and which file a change opens   | `docs/02-developing.md`       |
| The two suites, the spec documents        | `docs/03-testing.md`          |
| The release, install, the footprint       | `docs/04-operating.md`        |
| status / install / upgrade / clean / run  | `docs/05-commands.md`         |
| Machine registry, `config.json`, remote   | `docs/06-machines.md`         |
| The `j config` TUI, items, categories     | `docs/07-configuration.md`    |
| Tool + skill registries (`config/`)       | `docs/08-tools-and-skills.md` |
| Dotfiles (`dotfiles/applications/`)       | `docs/09-dotfiles.md`         |
| `@jterrazz` stack conventions             | `docs/10-stack.md`            |
| Repo doctrine (corpus/injection/compiler) | `docs/11-repo-structure.md`   |

The repo-structure doctrine is authored **here** (`docs/11-repo-structure.md`) and ships to agents as the `jterrazz-repo-structure` skill. This repo is an application: it adopts the corpus + injection layers and never generates a `docs/reference/` projection. Its one compiler is `make skills` (`src/cmd/skillsgen/`), which projects the config registries into the `jterrazz-toolbelt` skill rosters — see `docs/08-tools-and-skills.md`.

## Setup & commands

```bash
make build     # Build .artifacts/go/j
make test      # Go unit tests
make test-e2e  # npm install + rebuild j + vitest --run (the specs/cli suite)
make lint      # golangci-lint (pinned) + typescript check
```

The e2e specs need `@jterrazz/test` (npm) and rebuild `j` via an mtime check — see `specs/cli/cli.specification.ts`. `make lint` installs the npm half itself, so both linters run from one gesture.

## Standing rules

- **Add to a registry, not to command code.** Which registry file a change opens is the table in `docs/02-developing.md`.
- **Four things land in the same commit as the change that makes them true** — the `make skills` projection, the regenerated spec documents, the chapter the behaviour falsified, and the skill a doctrine chapter feeds. Each is stated once, in `docs/02-developing.md` § What a change owes.
- **Never edit a generated section by hand** — the blocks between the `GENERATED` markers of `skills/jterrazz-toolbelt/SKILL.md` come from `make skills`.
