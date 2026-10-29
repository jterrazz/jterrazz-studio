# workstation

A single CLI to bootstrap and manage a macOS development machine: tools, configs, templates, and remote access. No sudo required.

## Install

**Fresh machine** (no Go needed):

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/studio-workstation/main/scripts/install.sh | sh
source ~/.zshrc
```

The binary lives at `~/.studio/bin/workstation`; all user data goes under `~/.studio/`. The first steps after it are [Getting started](docs/start/getting-started.md); what the install leaves on a machine is [Footprint](docs/reference/footprint.md), and the release it comes from [Release](docs/guides/release.md).

## What it does

```sh
workstation status      # Full-screen dashboard (system, workspace, apps, config)
workstation machine     # Machine registry, remote actions, server config
workstation install     # Install from the tool catalogue
workstation config      # Configure the local machine (terminal, security, editor, skills)
workstation remote      # Tailscale connect / disconnect
workstation run         # git + docker shortcuts
```

## Documentation

The corpus lives in [`docs/`](docs/README.md); start at its map, which lists every page.

For agents: start at [`CLAUDE.md`](CLAUDE.md). The [`skills/`](skills/) directory ships the shared skills: `studio-toolbelt` (the machine's tool roster, generated from the install registry), `studio-stack`, `studio-new-project` and `studio-repo-structure`. The manifest and SEO domain skills ship with [`package-manifest`](https://github.com/jterrazz/package-manifest).

## License

MIT © [Jean-Baptiste Terrazzoni](https://github.com/jterrazz)
