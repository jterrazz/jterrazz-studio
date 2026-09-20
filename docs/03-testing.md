# Testing

Two suites prove a change: Go unit tests over the packages, and end-to-end specs that drive the real `j` binary and state what it printed.

| Section                  | Answers                                               |
| ------------------------ | ----------------------------------------------------- |
| The two suites           | What `make test` and `make test-e2e` each run         |
| A scenario is a document | The `<case>.spec.yaml` form and how it is regenerated |
| When a spec stays code   | The one reason to write TypeScript instead            |
| What CI runs             | The gate a pull request passes                        |

## The two suites

```sh
make test      # go test ./src/...
make test-e2e  # npm install + rebuild j + vitest --run
```

`make test` is the fast one: colocated `*_test.go` files next to the code they judge — the registry's parsing and validation (`src/internal/config/`), version parsing and the skill lock (`src/internal/domain/`), the TUI's section and update logic (`src/internal/presentation/views/`). One of them is a sync test over the projection: `src/cmd/skillsgen/main_test.go` fails when `skills/jterrazz-toolbelt/SKILL.md` no longer matches the registry, which is why a registry change and its `make skills` land together.

`make test-e2e` runs the CLI specs under `specs/cli/` on Vitest, through `@jterrazz/test`. There is exactly one runner: `specs/cli/cli.specification.ts` builds `.artifacts/go/j-test` before any spec runs, so the whole suite shares one compilation, and every spec drives that binary. It rebuilds when the binary is missing, when `J_FORCE_REBUILD=1` is set (what the Makefile target does), or when any `src/**/*.go` is newer than the binary — an mtime check, so editing the CLI and re-running the suite always exercises the change instead of a stale build.

`vitest.config.ts` names that runner once, through `defineSpecConfig({ test: { projects: [cli()] } })`: the `cli()` facet collects `specs/cli/**/*.spec.ts` and wires the literate plugin onto `specs/cli/cli.specification.ts` by default, so every `*.spec.yaml` under `specs/cli/` becomes a one-test module driving it.

## A scenario is a document

Almost every spec here is a terminal session — a command, its exit code, what it printed — so it is written as one: a `<case>.spec.yaml` in the command's folder, in the literate format of `@jterrazz/test`. The file IS the test; `description:` is its title in the runner, and each entry of `runs:` states the command, its `exit:` and its `stdout:`/`stderr:` byte-exact.

```yaml
# specs/cli/status/help-surface.spec.yaml
description: documents the status command under --help
runs:
    - command: status --help
      exit: 0
      stdout: |
          Show comprehensive system status

          Usage:
            j status [flags]
```

A spec may mount ground of its own: `cli.fixture("server-registry/")` mounts `specs/cli/machine/_fixtures/server-registry/` as the working directory, and `.env({ HOME: "$WORKDIR" })` points the binary at the registry inside it — which is how the server-role rows are exercised deterministically on a client laptop. Ground carries a leading underscore; the spec's own folder never does.

`TEST_UPDATE=1 make test-e2e` rewrites the `exit:` and the streams of every document from what the binary actually printed — deliberately, after a change to a command's output, and never as a way to make a red suite go green. Nothing else in the file is touched.

## When a spec stays code

A spec stays a `*.spec.ts` — code under `specs/`, not a document — for one reason: the output's text comes from the HOST, and a byte-exact stream cannot promise it. That is the `✓`/`✗` install-state column and the machine-status verdicts — macOS version, FileVault, sshd, running services. Those specs (`specs/cli/install/install.spec.ts`, `specs/cli/machine/machine.spec.ts`) probe the rows the binary always emits rather than the text the machine happened to fill them with.

The full grammar, and the rest of the reasons to reach for code, are `@jterrazz/test`'s [`docs/12-cli.md`](https://github.com/jterrazz/package-test/blob/main/docs/12-cli.md).

## What CI runs

`.github/workflows/validate.yaml` calls the shared `jterrazz-actions` validate workflow on every push and every pull request to `main`; what that workflow runs against a repository's `make` targets is [The stack](10-stack.md)'s. Locally the same gate is `make build`, `make lint`, `make test`, `make test-e2e`.

The npm half is linted by the same toolchain as every other repository: `npm run lint` is `typescript check`, and `make lint` runs it after the Go half ([Developing](02-developing.md)).

## Related

- [Developing](02-developing.md) — the change that these suites judge.
- [Operating](04-operating.md) — what happens after a green run: the tagged release.
