# Agent brief: `workstation` (studio-workstation)

A Go CLI (Cobra + Bubble Tea TUIs) that bootstraps and manages a macOS dev machine. This file **routes**; it does not restate what the corpus already says.

## Mental model

- **Registry-driven.** The tools it installs, the config items, and the curated agent skills are declared as data in `src/internal/config/`, and the same catalogue powers `workstation install`, `workstation status`, and `workstation config`. Add to the registry, not to bespoke command code.
- **The machine registry is the source of truth.** `~/.studio/config.json` (aliases, roles, ssh) drives everything; adding a machine also writes a managed `~/.ssh/config` block. Role (`client`/`server`) gates what `status` reports and which `config` items appear.
- **Dotfiles are versioned.** `dotfiles/applications/*` are installed onto the machine by `workstation config`.
- **Specs drive the real binary.** `specs/cli/` uses `@jterrazz/test` (`specification.cli`) against a freshly built `workstation`; the runner rebuilds when any `src/**/*.go` is newer than the test binary. A scenario is a `<case>.spec.yaml` document, not code (see `docs/guides/test.md`).

## Where knowledge lives (route here first)

The corpus is `docs/` + `README.md`, mapped by `docs/README.md`. Do not duplicate it; link to it.

| Working on                                      | Read                                 |
| ----------------------------------------------- | ------------------------------------ |
| A fresh machine: install, init, configure       | `docs/start/getting-started.md`      |
| Layers, registries, the Cobra root              | `docs/concepts/architecture.md`      |
| Building, and which file a change opens         | `docs/guides/develop.md`             |
| The two suites, the spec documents              | `docs/guides/test.md`                |
| The tagged release, `make provision`            | `docs/guides/release.md`             |
| What `workstation` writes on a machine          | `docs/reference/footprint.md`        |
| status / install / upgrade / clean / run        | `docs/reference/commands.md`         |
| Machine registry, `config.json`, remote         | `docs/reference/machines.md`         |
| The `workstation config` TUI, items, categories | `docs/reference/configuration.md`    |
| Tool + skill registries (`config/`)             | `docs/reference/tools-and-skills.md` |
| Dotfiles (`dotfiles/applications/`)             | `docs/reference/dotfiles.md`         |
| `@jterrazz` stack conventions                   | `docs/concepts/stack.md`             |
| Repo doctrine (corpus/injection/compiler)       | `docs/concepts/repo-structure.md`    |

The repo-structure doctrine is authored **here** (`docs/concepts/repo-structure.md`) and ships to agents as the `studio-repo-structure` skill. This repo is an application: it adopts the corpus + injection layers, and generates no API reference. Its one compiler is `make skills` (`src/cmd/skillsgen/`), which projects the config registries into the `studio-toolbelt` skill rosters (see `docs/reference/tools-and-skills.md`).

## Setup & commands

```bash
make install    # Go modules + the npm half
make build      # Build .artifacts/go/workstation
make check      # lint (gofmt, golangci-lint pinned, typescript check) then test (go test, the specs/cli suite)
make fix        # gofmt -w, then typescript fix
make provision  # put this build on the machine as the real workstation
```

The specs rebuild `workstation` via an mtime check (see `specs/cli/cli.specification.ts`). `make install` never touches the machine; `make provision` does.

## Standing rules

- **Add to a registry, not to command code.** Which registry file a change opens is the table in `docs/guides/develop.md`.
- **Four things land in the same commit as the change that makes them true**: the `make skills` projection, the regenerated spec documents, the page the behaviour falsified, and the skill a doctrine page feeds. Each is stated once, in `docs/guides/develop.md` § What a change owes.
- **Never edit a generated section by hand**: the blocks between the `GENERATED` markers of `skills/studio-toolbelt/SKILL.md` come from `make skills`.
