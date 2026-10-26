# The manual of workstation

`workstation` is the workstation product: one Go CLI that bootstraps and manages a macOS development machine. It installs from a registry-driven tool catalogue, configures the machine and its dotfiles, keeps a registry of machines reachable over the tailnet, and curates the agent skills an operator wants. This repository also authors two doctrines every `@jterrazz` repository follows, shipped to agents as skills. This page lists every page of the manual, by what the reader came to do.

## Start

| Page                                        | Answers                                                        |
| ------------------------------------------- | -------------------------------------------------------------- |
| [Getting started](start/getting-started.md) | From a fresh Mac to a registered, configured machine, in order |

## Guides

| Page                         | Answers                                                                          |
| ---------------------------- | -------------------------------------------------------------------------------- |
| [Develop](guides/develop.md) | Working from source: the gestures, which file a change opens, what a change owes |
| [Test](guides/test.md)       | The Go module tests, the `specs/cli` documents, how a host's verdict is written  |
| [Release](guides/release.md) | The tagged release, how a machine takes it, provisioning from a checkout         |

## Concepts

| Page                                         | Answers                                                                                       |
| -------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [Architecture](concepts/architecture.md)     | The binary, the four packages, why a feature is a registry entry, the one compiler            |
| [The stack](concepts/stack.md)               | How `@jterrazz` projects compose: packages, naming, required files; shipped as `studio-stack` |
| [Repo structure](concepts/repo-structure.md) | Where knowledge lives in every `@jterrazz` repository; shipped as `studio-repo-structure`     |

## Reference

| Page                                              | Answers                                                                             |
| ------------------------------------------------- | ----------------------------------------------------------------------------------- |
| [Commands](reference/commands.md)                 | `status`, `install`, `upgrade`, `clean`, `run`, and the shell shortcuts             |
| [Machines](reference/machines.md)                 | The machine registry, the `config.json` model, roles, remote access                 |
| [Configuration](reference/configuration.md)       | The `workstation config` TUI: its tabs, items, categories, and what each one writes |
| [Tools and skills](reference/tools-and-skills.md) | The tool and skill registries, the `skills` CLI, the toolbelt projection            |
| [Dotfiles](reference/dotfiles.md)                 | The versioned application configs under `dotfiles/applications/`                    |
| [Footprint](reference/footprint.md)               | Everything `workstation` writes on a machine                                        |

## Decisions

The records of decisions this repository alone took go in [`decisions/`](decisions/), numbered in the order they were taken, the mold `_template.md` beside them. A decision spanning several repositories is recorded by the corpus that spans them.
