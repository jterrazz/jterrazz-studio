# Getting started

From a fresh Mac to a machine `workstation` knows: the one-line install, then the three commands that register the machine, show it and configure it.

## 1. Install `workstation`

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/studio-workstation/main/scripts/install.sh | sh
source ~/.zshrc
```

`scripts/install.sh` reads the latest release tag from the GitHub API, downloads `workstation-<os>-<arch>` for the detected platform into `~/.studio/bin/`, and makes it executable. It then does two things beyond the binary, because the binary alone is not the product: it clones the repository to `~/Developer/studio-workstation` when that directory does not exist (the dotfiles and the skills are files on disk, not compiled in), and appends a `source .../dotfiles/applications/zsh/zshrc.sh` line to `~/.zshrc` when one is not already there. That clone is what [Dotfiles](../reference/dotfiles.md) and the `skills/` folder are served from.

Upgrading is re-running the same script: it overwrites the binary and leaves the clone and the shell line alone.

## 2. Register this machine

```sh
workstation machine init
```

It asks for this machine's alias and role (`client` for a laptop, `server` for a box other machines reach) and writes them to `~/.studio/config.json`. The role decides which checks `workstation machine status` runs and which items `workstation config` offers ([Machines](../reference/machines.md)).

## 3. Look at it, then configure it

```sh
workstation status    # the dashboard: system, workspace, applications, configuration
workstation config    # the configurator: terminal, security, editor, skills, remote
```

Everything `workstation` wrote on the way lives under `~/.studio/`, listed in [Footprint](../reference/footprint.md); every verb is in [Commands](../reference/commands.md).
