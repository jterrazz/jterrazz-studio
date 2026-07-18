# Getting started

`j` is a single CLI to bootstrap and manage a macOS development machine — tools, configs, templates, and remote access. No sudo required.

## Install

**Fresh machine** (no Go needed):

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/jterrazz-cli/main/scripts/install.sh | sh
source ~/.zshrc
```

**From source** (requires Go 1.24+):

```sh
git clone https://github.com/jterrazz/jterrazz-cli.git ~/Developer/jterrazz-cli
cd ~/Developer/jterrazz-cli
make install
source ~/.zshrc
```

The binary lives at `~/.jterrazz/bin/j`. All user data goes under `~/.jterrazz/` — see [Machines](03-machines.md) for the config model.

## First steps

```sh
j machine init     # Bootstrap THIS machine (interactive)
j status           # Full-screen dashboard
j install          # List tracked tools with status
j config           # Configure the local machine (TUI)
```

## User data

Everything lives under `~/.jterrazz/`:

```
~/.jterrazz/
├── bin/           # CLI binary
├── config.json    # Runtime config (remote/Tailscale, machine registry)
├── tailscale/     # Userspace daemon state
└── dns/           # Generated DNS profiles
```

The `config.json` schema and its role as the single source of truth are covered in [Machines](03-machines.md).

## Development

```sh
make build     # Build ./j
make test      # Run Go unit tests
make install   # Build + install to ~/.jterrazz/bin
make check     # Verify installation
```

End-to-end specs (written with `@jterrazz/test`) drive the built binary:

```sh
make test-e2e  # npm install + rebuild j + vitest --run
```

The runner rebuilds `j` whenever any `src/**/*.go` is newer than the test binary (an mtime check), so editing the CLI and re-running the suite always exercises the change. See [`specs/cli/cli.specification.ts`](../specs/cli/cli.specification.ts).

### Releasing

Push a version tag to build and publish binaries via GitHub Actions:

```sh
git tag v1.0.0
git push --tags
```

Builds for `darwin/arm64`, `darwin/amd64`, `linux/arm64`, `linux/amd64`.

## Project structure

```
src/
├── cmd/j/main.go             # Entry point
└── internal/
    ├── commands/             # CLI commands (Cobra)
    ├── config/               # Tool, script, command, and skill registries
    ├── domain/               # Version parsing, status loading, skills
    └── presentation/         # TUI views, components, theme
dotfiles/
└── applications/             # App configs (ghostty, starship, tmux, vscode, zed, zsh)
docs/                         # This corpus
skills/                       # Claude Code skills for the @jterrazz workflow
specs/cli/                    # End-to-end specs (@jterrazz/test)
```

## Related

- [Commands](02-commands.md) — the everyday verbs.
- [Machines](03-machines.md) — the registry & config model.
- [Tools and skills](05-tools-and-skills.md) — the curated registries `install`/`config` draw from.
