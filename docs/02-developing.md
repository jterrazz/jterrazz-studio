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
git clone https://github.com/jterrazz/jterrazz-studio.git ~/Developer/jterrazz/jterrazz-studio
cd ~/Developer/jterrazz/jterrazz-studio
make build                # .artifacts/go/j
./.artifacts/go/j status  # run the build without installing it
```

Go 1.24+ is the only requirement for the Go half; the end-to-end specs additionally need npm ([Testing](03-testing.md)). `make install` puts the build on the machine as the real `j` — that gesture, and everything else about a machine that has `j` on it, is [Operating](04-operating.md)'s.

## The make targets

| Target           | Does                                                                |
| ---------------- | ------------------------------------------------------------------- |
| `make build`     | Build `.artifacts/go/j`                                             |
| `make test`      | `go test ./src/...`                                                 |
| `make test-e2e`  | npm install, rebuild the test binary, run the `specs/cli` suite     |
| `make fmt`       | `gofmt -w ./src/`                                                   |
| `make vet`       | `go vet ./src/...`                                                  |
| `make lint`      | Both halves: `golangci-lint run ./src/...`, then `typescript check` |
| `make skills`    | Regenerate the toolbelt skill's rosters from the registries         |
| `make install`   | Build and install to `~/.jterrazz/bin`                              |
| `make uninstall` | Remove the installed binary                                         |
| `make check`     | Verify the installation                                             |
| `make clean`     | Remove `.artifacts/`                                                |

`make lint` is the whole gate, and it fails on the first finding of either half. The Go half runs the `golangci-lint` version the Makefile PINS — it installs that exact one when it is missing, because a linter resolved at `@latest` makes a green run a fact about the day it ran. Which linters judge the Go sources, and the four exclusions that carry a reason, are `.golangci.yml`. The TypeScript half is `typescript check` from `@jterrazz/typescript`, whose passes are the toolchain's own [Quality checks](https://github.com/jterrazz/package-typescript/blob/main/docs/06-quality-checks.md).

## Where a change goes

Most changes are a registry entry, not new code — that is the shape [Architecture](01-architecture.md) describes. The file to open:

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

Four things land in the same commit as the change that makes them true:

- **A registry change regenerates the projection.** Run `make skills`; the sync test in `make test` fails otherwise, and the sections between the `GENERATED` markers are never edited by hand ([Tools and skills](08-tools-and-skills.md)).
- **A change to a command's output regenerates the specs.** `TEST_UPDATE=1 make test-e2e` rewrites the `exit:` and the streams of the affected documents, and nothing else ([Testing](03-testing.md)).
- **A change to behaviour updates its chapter.** The corpus is where this repository's knowledge is authored; a page left saying what is no longer true is a bug that shipped.
- **A change to a doctrine chapter updates its skill.** [The stack](10-stack.md) ships as `jterrazz-stack` and [Repo structure](11-repo-structure.md) as `jterrazz-repo-structure`; skills route, they never author.

## Related

- [Architecture](01-architecture.md) — the shape a change lands in.
- [Testing](03-testing.md) — the two suites, and the document form of a spec.
- [Operating](04-operating.md) — installing, and the tagged release.
