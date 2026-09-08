# Architecture

`j` is one Go binary whose behaviour is declared as data: a registry of tools, configuration items, run shortcuts and curated skills, read by a thin Cobra command layer and rendered by Bubble Tea views.

| Section                         | Answers                                                  |
| ------------------------------- | -------------------------------------------------------- |
| The binary                      | What runs when you type `j`                              |
| The four packages               | Where each kind of code lives                            |
| Registry-driven                 | Why a feature is usually a data entry, not a new command |
| The two sources of truth        | The registry in the repo, the config on the machine      |
| The compiler                    | The one source layer this repo projects                  |

## The binary

`src/cmd/j/main.go` does one thing: call `commands.Execute()` and exit non-zero on error. `Execute` runs the single Cobra root declared in `src/internal/commands/root.go` — `j`, with Cobra's own `help` command and default completion command hidden.

Every verb is a `*cobra.Command` in its own file under `src/internal/commands/`, registered on the root by that file's `init()`. A verb is small on purpose: `status.go` and `config.go` hand straight over to a view (`statusview.RunOrExit()`, `configview.RunOrExit()`), and `install.go` holds only the branch between the interactive TUI and the plain-text catalogue — the decision it makes is whether stdout is a terminal, or whether `--list` was passed.

Two commands ship as a second binary each: `src/cmd/skillsgen/` is the projection generator run by `make skills`, and the CLI specs build their own copy of `j` (see [Testing](03-testing.md)).

## The four packages

| Package                        | Holds                                                                                                  |
| ------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `src/internal/commands/`       | One Cobra command per file, plus the server-side actions `j config` items call (`server_*.go`)          |
| `src/internal/config/`         | The registries — tools, scripts, run commands, cleanables, upgraders, skills — and the machine registry |
| `src/internal/domain/`         | Behaviour that is not a registry entry: version parsing, status loading, the `skills` CLI integration   |
| `src/internal/presentation/`   | Bubble Tea views (`views/status`, `views/install`, `views/config`), shared `components/`, `theme/`, and `print/` for plain output |

The registry package is the centre rather than a bottom layer: `commands`, `domain/status` and the views all import `config`, and `config` reaches back into `presentation/print` for the lines its install actions emit. Only `domain/tool` — version detection and command probing — is imported by others without importing anything of its own.

The four server items are the one place the direction is inverted. Their implementations need `commands`, which already imports `config`, so `server.go` hands them over at startup: its `init()` calls `config.RegisterServerActions(...)` with the check, install and uninstall functions defined in `server_*.go`. The registry stays the thing the TUI reads; the cycle never forms.

## Registry-driven

Each registry is a slice of structs that declares both the data and the functions that act on it. A `Script` (`src/internal/config/scripts.go`) carries its category, its help text, an optional `Role` gate, a `CheckFn` returning a `CheckResult`, an `InstallFn` taking the values a modal form collected, and an optional `UninstallFn` — and it is the presence of that last function that makes the item toggleable in the TUI. A `Tool` declares its install method, its dependencies, its version detection and its post-install scripts; a `Cleanable` declares where reclaimable storage lives and how big it is; a `RunCommand` declares a group of `j run` subcommands.

The consequence is the rule this repo is built on: **add to the registry, not to bespoke command code**. One entry gives the same tool a row in `j install`, a line in the Applications tab of `j status`, and an uninstall path — because all three read the same slice. A change that needs new command code is the exception, and worth questioning.

What each registry contains, and how the catalogue is organised, is [Tools and skills](08-tools-and-skills.md); what the items look like on screen is [Configuration](07-configuration.md).

## The two sources of truth

The repository holds the catalogue: what `j` knows how to install and configure. The machine holds its own registry: `~/.jterrazz/config.json`, which names this box, the machines it can reach, and their roles. Nothing else on the machine is consulted for that — adding a machine also writes the managed block in `~/.ssh/config`, and the role decides which checks run and which `j config` items exist at all (`Script.Role`, `Machine.Validate`). The model is [Machines](06-machines.md)'s.

## The compiler

This repo is an application in the sense of the [repo-structure doctrine](11-repo-structure.md): it authors a corpus, routes to it from `AGENTS.md` and `skills/`, and generates no API reference. It does hold one source layer worth projecting — the tool registry — and `make skills` (`src/cmd/skillsgen/`) renders it between the `GENERATED` markers of `skills/jterrazz-toolbelt/SKILL.md`. The projection is committed and a Go test fails when it drifts; see [Tools and skills](08-tools-and-skills.md).

## Related

- [Developing](02-developing.md) — the gestures that change any of the above.
- [Commands](05-commands.md) — what each verb does for a user.
- [Repo structure](11-repo-structure.md) — the corpus / injection / compiler doctrine this shape follows.
