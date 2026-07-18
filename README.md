# j

A single CLI to bootstrap and manage a macOS development machine — tools, configs, templates, and remote access. No sudo required.

## Install

**Fresh machine** (no Go needed):

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/jterrazz-cli/main/scripts/install.sh | sh
source ~/.zshrc
```

The binary lives at `~/.jterrazz/bin/j`; all user data goes under `~/.jterrazz/`. From-source install and the full setup are in [Getting started](docs/01-getting-started.md).

## What it does

```sh
j status      # Full-screen dashboard (system, workspace, apps, config)
j machine     # Machine registry, remote actions, server config
j install     # Install from a 100+ tool catalogue
j config      # Configure the local machine (terminal, security, editor, skills)
j remote      # Tailscale connect / disconnect
j run         # git + docker shortcuts
```

## Documentation

The corpus lives in [`docs/`](docs/):

- [Getting started](docs/01-getting-started.md) — install, user data, development.
- [Commands](docs/02-commands.md) — status, install, upgrade, clean, run, shortcuts.
- [Machines](docs/03-machines.md) — the machine registry & `config.json` model, remote.
- [Configuration](docs/04-configuration.md) — the `j config` TUI, items, and categories.
- [Tools and skills](docs/05-tools-and-skills.md) — the curated tool and skill registries.
- [Dotfiles](docs/06-dotfiles.md) — the versioned application configs.

Agents start at [`AGENTS.md`](AGENTS.md).

## License

MIT © [Jean-Baptiste Terrazzoni](https://github.com/jterrazz)
