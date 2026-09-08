# j — documentation

`j` is the workstation product: one Go CLI that bootstraps and manages a macOS development machine. It installs from a registry-driven tool catalogue, configures the machine and its dotfiles, keeps a registry of machines reachable over the tailnet, and curates the agent skills an operator wants on top of what the estate's own repositories already supply.

This corpus is where that knowledge is authored. `AGENTS.md` and the `skills/` entries route into it; they never restate it.

## Table of contents

| Chapter                                         | Covers                                                                                |
| ----------------------------------------------- | ------------------------------------------------------------------------------------- |
| [01 — Architecture](01-architecture.md)         | The binary, the four packages, why a feature is a registry entry                      |
| [02 — Developing](02-developing.md)             | Working from source, the make targets, which file a change opens                      |
| [03 — Testing](03-testing.md)                   | The Go units, the `specs/cli` documents, `make test-e2e`                              |
| [04 — Operating](04-operating.md)               | The tagged release, how a machine installs `j`, its footprint                         |
| [05 — Commands](05-commands.md)                 | `status`, `install`, `upgrade`, `clean`, `run`, and the shell shortcuts               |
| [06 — Machines](06-machines.md)                 | The machine registry, the `config.json` model, roles, remote access                   |
| [07 — Configuration](07-configuration.md)       | The `j config` TUI: items, categories, and what each one writes                       |
| [08 — Tools and skills](08-tools-and-skills.md) | The tool and skill registries, the `skills` CLI integration, the toolbelt projection  |
| [09 — Dotfiles](09-dotfiles.md)                 | The versioned application configs under `dotfiles/applications/`                      |
| [10 — The stack](10-stack.md)                   | How `@jterrazz` projects compose: packages, naming, required files, conventions       |
| [11 — Repo structure](11-repo-structure.md)     | Where knowledge lives in every `@jterrazz` repo — the canonical home of that doctrine |

## How this documentation is organized

- **Chapters 01–04** are the spine every `@jterrazz` repository carries in the same four numbers: the shape, how it is changed, what proves the change, and how it ships.
- **Chapters 05–09** follow the five things `j` manages — the verbs, the machines, the configuration, the registries, and the dotfiles.
- **Chapters 10–11** are doctrine rather than product. They describe every `@jterrazz` repository, not this one, and each ships to agents as a skill: `jterrazz-stack` and `jterrazz-repo-structure`. A change to either chapter updates its skill in the same change.
