# j — documentation

`j` is the workstation product: one Go CLI that bootstraps and manages a macOS development machine. It installs from a registry-driven tool catalogue, configures the machine and its dotfiles, keeps a registry of machines reachable over the tailnet, and curates the agent skills an operator wants on top of what the estate's own repositories already supply.

This corpus is where that knowledge is authored. `AGENTS.md` and the `skills/` entries route into it; they never restate it.

## Table of contents

| Chapter                                          | Covers                                                                                       |
| ------------------------------------------------ | -------------------------------------------------------------------------------------------- |
| [01 — Getting started](01-getting-started.md)    | Install, user data under `~/.jterrazz/`, development, the spec documents, release, the layout |
| [02 — Commands](02-commands.md)                  | `status`, `install`, `upgrade`, `clean`, `run`, and the shell shortcuts                      |
| [03 — Machines](03-machines.md)                  | The machine registry, the `config.json` model, roles, remote access                          |
| [04 — Configuration](04-configuration.md)        | The `j config` TUI: items, categories, and what each one writes                              |
| [05 — Tools and skills](05-tools-and-skills.md)  | The tool and skill registries, the `skills` CLI integration, the toolbelt projection         |
| [06 — Dotfiles](06-dotfiles.md)                  | The versioned application configs under `dotfiles/applications/`                             |
| [07 — The stack](07-stack.md)                    | How `@jterrazz` projects compose: packages, naming, required files, conventions              |
| [08 — Repo structure](08-repo-structure.md)      | Where knowledge lives in every `@jterrazz` repo — the canonical home of that doctrine        |

## How this documentation is organized

- **Chapters 01–02** are the entry: install the binary, then learn what it does.
- **Chapters 03–06** follow the four things `j` manages — machines, configuration, the registries, and the dotfiles.
- **Chapters 07–08** are doctrine rather than product. They describe every `@jterrazz` repository, not this one, and each ships to agents as a skill: `jterrazz-stack` and `jterrazz-repo-structure`. A change to either chapter updates its skill in the same change.
