# Testing

Two suites prove a change: Go unit tests over the packages, and end-to-end specs that drive the real `j` binary and state what it printed.

| Section                  | Answers                                               |
| ------------------------ | ----------------------------------------------------- |
| The two suites           | What `make test` runs, and in which order             |
| A scenario is a document | The `<case>.spec.yaml` form and how it is regenerated |
| What the host decides    | The tokens that stand for a host's verdicts           |
| What CI runs             | The gate a pull request passes                        |

## The two suites

```sh
make test      # go test ./src/..., then the specs against a fresh build of j
npm test       # the specs alone
```

The Go module tests come first: `*_test.go` files next to the code they judge — the registry's parsing and validation (`src/internal/config/`), version parsing and the skill lock (`src/internal/domain/`), the TUI's section and update logic (`src/internal/presentation/views/`). One of them is a sync test over the projection: `src/cmd/skillsgen/main_test.go` fails when `skills/jterrazz-toolbelt/SKILL.md` no longer matches the registry, which is why a registry change and its `make skills` land together.

Then `make test` runs the CLI specs under `specs/cli/` on Vitest, through `@jterrazz/test`. There is exactly one runner: `specs/cli/cli.specification.ts` builds `.artifacts/go/j-test` before any spec runs, so the whole suite shares one compilation, and every spec drives that binary. It rebuilds when the binary is missing, when `J_FORCE_REBUILD=1` is set (what `make test` does), or when any `src/**/*.go` is newer than the binary — an mtime check, so editing the CLI and re-running the suite always exercises the change instead of a stale build.

`vitest.config.ts` is the one line `@jterrazz/test` derives from `jterrazz.config.ts`: its `cli` project collects `specs/cli/`, and every `*.spec.yaml` there becomes a one-test module driving that runner. The runner also gives every run a `HOME` of its own — the spec's working directory — so `j` never reads the developer's `~/.jterrazz/` registry.

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

A document may stand on ground of its own: `fixture: ./__fixtures__/server-registry/` spreads `specs/cli/machine/__fixtures__/server-registry/` into the working directory, which is the run's `HOME`, so the binary reads the registry inside it — which is how the server-role rows are exercised deterministically on a client laptop. Ground sits in a special folder spelt with double underscores; a subject folder never starts with one.

`TEST_UPDATE=1 npm test` rewrites the `exit:` and the streams of every document from what the binary actually printed — deliberately, after a change to a command's output, and never as a way to make a red suite go green. Nothing else in the file is touched.

## What the host decides

Some of what `j` prints comes from the HOST, not from the binary: the `✓`/`✗` install-state column of the catalog, and the machine-status verdicts — FileVault, sshd, running services, OrbStack. A document still states the whole stream; the cell the host decides is the token `{{string}}` (`{{string}} FileVault {{string}}`), so each row's label and its place are asserted and the machine's answer is not. What the binary decides on its own stays literal: the catalog's families and methods, the config paths under `{{workdir}}`, the rows a server role adds.

`TEST_UPDATE=1` keeps a token wherever it still matches, but writes a NEW row with the host's literal glyph: after adding a tool, replace its `✓` or `✗` with `{{string}}` before committing, or the document holds this machine's answer.

The full grammar is `@jterrazz/test`'s [Spec documents](https://github.com/jterrazz/package-test/blob/main/docs/guides/spec-documents.md) and [CLI specs](https://github.com/jterrazz/package-test/blob/main/docs/guides/cli-specs.md).

## What CI runs

`.github/workflows/validate.yaml` calls the shared `jterrazz-actions` validate workflow on every push and every pull request to `main`; what that workflow runs against a repository's `make` targets is [The stack](../concepts/stack.md)'s. Locally the same gate is `make build` then `make check`.

The npm half is linted by the same toolchain as every other repository: `npm run lint` is `typescript check`, and `make lint` runs it after the Go half ([Develop](develop.md)).

## Related

- [Develop](develop.md) — the change that these suites judge.
- [Release](release.md) — what happens after a green run: the tagged release.
