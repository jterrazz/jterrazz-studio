# Getting started

From a fresh Mac to a machine `j` knows: the one-line install, then the three commands that register the machine, show it and configure it.

## 1. Install `j`

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/jterrazz-studio/main/scripts/install.sh | sh
source ~/.zshrc
```

`scripts/install.sh` reads the latest release tag from the GitHub API, downloads `j-<os>-<arch>` for the detected platform into `~/.jterrazz/bin/`, and makes it executable. It then does two things beyond the binary, because the binary alone is not the product: it clones the repository to `~/Developer/jterrazz-studio` when that directory does not exist (the dotfiles and the skills are files on disk, not compiled in), and appends a `source .../dotfiles/applications/zsh/zshrc.sh` line to `~/.zshrc` when one is not already there. That clone is what [Dotfiles](../reference/dotfiles.md) and the `skills/` folder are served from.

Upgrading is re-running the same script: it overwrites the binary and leaves the clone and the shell line alone.

## 2. Register this machine

```sh
j machine init
```

It asks for this machine's alias and role (`client` for a laptop, `server` for a box other machines reach) and writes them to `~/.jterrazz/config.json`. The role decides which checks `j machine status` runs and which items `j config` offers ([Machines](../reference/machines.md)).

## 3. Look at it, then configure it

```sh
j status    # the dashboard: system, workspace, applications, configuration
j config    # the configurator: terminal, security, editor, skills, remote
```

Everything `j` wrote on the way lives under `~/.jterrazz/`, listed in [Footprint](../reference/footprint.md); every verb is in [Commands](../reference/commands.md).
