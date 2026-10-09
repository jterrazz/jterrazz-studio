# Developing

How to work on `j` itself: build it, run it from the checkout, and know which file a change belongs in.

| Section             | Answers                                     |
| ------------------- | ------------------------------------------- |
| Working from source | The clone, Go 1.24+, and the loop           |
| The make targets    | Every gesture the Makefile offers           |
| Where a change goes | Which file to open for which kind of change |
| What a change owes  | What must land in the same commit           |

## Working from source

```sh
git clone https://github.com/jterrazz/jterrazz-studio.git ~/Developer/jterrazz-studio
cd ~/Developer/jterrazz-studio
make build                # .artifacts/go/j
./.artifacts/go/j status  # run the build without installing it
```

The path is the one the install script clones to, because the machine serves one clone: the `~/.zshrc` line sources its `zshrc.sh`, and `j` reads its dotfiles and scripts from it unless `J_REPO_PATH` names another ([Footprint](../reference/footprint.md)). A second clone elsewhere is a checkout to build, not the one the shell runs.

Go 1.24+ and npm are the two requirements: `make install` fetches the Go modules and the npm half the specs and the toolchain run on ([Test](test.md)). `make provision` puts the build on the machine as the real `j` — that gesture, and everything else about a machine that has `j` on it, is [Release](release.md)'s.

## The make targets

| Target             | Does                                                                               |
| ------------------ | ---------------------------------------------------------------------------------- |
| `make install`     | Fetch the Go modules and install the npm half (`npm ci`)                           |
| `make build`       | Build `.artifacts/go/j`                                                            |
| `make lint`        | Both halves: `gofmt -l` and `golangci-lint run ./src/...`, then `typescript check` |
| `make fix`         | Every fixer: `gofmt -w ./src/`, then `typescript fix`                              |
| `make test`        | `go test ./src/...`, then the `specs/cli` suite against a fresh test binary        |
| `make check`       | `lint` then `test` — the gate a change passes                                      |
| `make skills`      | Regenerate the toolbelt skill's rosters from the registries                        |
| `make provision`   | Build and put `j` on this machine: `~/.jterrazz/bin`, the `zshrc.sh` line          |
| `make deprovision` | Remove the provisioned binary                                                      |
| `make clean`       | Remove `.artifacts/`                                                               |

The first six are the gestures every repository of the estate answers, so a session that never opened this one installs, proves and repairs it the same way; `make install` touches nothing outside the checkout, and putting `j` on the machine is `make provision`'s alone. `make check` is the whole gate, and it fails on the first finding of either half. The Go half runs the `golangci-lint` version the Makefile PINS — it installs that exact one when it is missing, because a linter resolved at `@latest` makes a green run a fact about the day it ran. Which linters judge the Go sources, and the four exclusions that carry a reason, are `.golangci.yml`. The TypeScript half is `typescript check` from `@jterrazz/typescript`, whose passes are the toolchain's own [Quality checks](https://github.com/jterrazz/package-typescript/blob/main/docs/reference/quality-checks.md).

## Where a change goes

Most changes are a registry entry, not new code — that is the shape [Architecture](../concepts/architecture.md) describes. The file to open:

| Change                               | File                                                                              |
| ------------------------------------ | --------------------------------------------------------------------------------- |
| A tool `j install` offers            | `src/internal/config/tools_catalog.go`                                            |
| A `j config` item                    | `src/internal/config/scripts.go` (server items: `server_*.go` beside the command) |
| A `j run` shortcut                   | `src/internal/config/commands.go`                                                 |
| Something `j clean` reclaims         | `src/internal/config/cleanables.go`                                               |
| A package manager `j upgrade` knows  | `src/internal/config/upgraders.go`                                                |
| A curated agent skill                | `src/internal/config/skills.go`                                                   |
| A new verb                           | A file under `src/internal/commands/`, registering itself in `init()`             |
| What a TUI shows                     | `src/internal/presentation/views/<view>/`                                         |
| A shell shortcut or an app's dotfile | `dotfiles/applications/`                                                          |

## What a change owes

Five things land in the same commit as the change that makes them true:

- **A registry change regenerates the projection.** Run `make skills`; the sync test in `make test` fails otherwise, and the sections between the `GENERATED` markers are never edited by hand ([Tools and skills](../reference/tools-and-skills.md)).
- **A change to a command's output regenerates the specs.** `TEST_UPDATE=1 npm test` rewrites the `exit:` and the streams of the affected documents, and nothing else ([Test](test.md)).
- **A change to behaviour updates its page.** The corpus is where this repository's knowledge is authored; a page left saying what is no longer true is a bug that shipped.
- **A comment states its kind, in both halves.** The default is none; a comment that stays opens with `invariant:`, `hazard:`, `outside:` or `value:` ([A comment states its kind](https://github.com/jterrazz/package-typescript/blob/main/docs/reference/lint-presets.md#a-comment-states-its-kind)). The Go half holds the same law in Go's shapes: a package comment (staticcheck's ST1000 asks for it), and one line of doc above an exported identifier or member, its hover. A kinded comment about an exported function goes inside its body, since ST1020 wants that function's doc to open with its name. `typescript check` reads the configs and scripts, not `.go` files, so on the Go side the review is the gate.
- **A change to a doctrine page updates its skill.** [The stack](../concepts/stack.md) ships as `jterrazz-stack` and [Repo structure](../concepts/repo-structure.md) as `jterrazz-repo-structure`; skills route, they never author.

## Related

- [Architecture](../concepts/architecture.md) — the shape a change lands in.
- [Test](test.md) — the two suites, and the document form of a spec.
- [Release](release.md) — installing, and the tagged release.
