# Operating

This repository ships one artefact — the `j` binary — and it reaches a machine two ways: a tagged GitHub release, or a build from a checkout.

| Section                   | Answers                                         |
| ------------------------- | ----------------------------------------------- |
| The tagged release        | What a `v*` tag builds and publishes            |
| How a machine installs it | The one-line install, and what it leaves behind |
| From a checkout           | `make install`, `make check`, `make uninstall`  |
| The footprint             | Everything `j` writes outside the repository    |

## The tagged release

There is no deploy: the release IS the delivery. Pushing a version tag is the whole gesture.

```sh
git tag v1.0.0
git push --tags
```

`.github/workflows/release.yaml` fires on `push: tags: v*` and calls the shared `jterrazz-actions` `release-go.yaml` workflow with the binary name `j`, the build path `./src/cmd/j`, Go 1.24, and four targets: `darwin/arm64`, `darwin/amd64`, `linux/arm64`, `linux/amd64`. Nothing publishes on a merge to `main` — a merge only runs the validate workflow ([Testing](03-testing.md)).

Each target is published as an asset named `j-<os>-<arch>`, which is the name the install script asks for. A release that renames those assets breaks every machine's install path.

## How a machine installs it

```sh
xcode-select --install
curl -fsSL https://raw.githubusercontent.com/jterrazz/jterrazz-studio/main/scripts/install.sh | sh
source ~/.zshrc
```

`scripts/install.sh` reads the latest release tag from the GitHub API, downloads `j-<os>-<arch>` for the detected platform into `~/.jterrazz/bin/`, and makes it executable. It then does two things beyond the binary, because the binary alone is not the product: it clones the repository to `~/Developer/jterrazz-studio` when that directory does not exist — the dotfiles and the skills are files on disk, not compiled in — and appends a `source …/dotfiles/applications/zsh/zshrc.sh` line to `~/.zshrc` when one is not already there. That clone is what [Dotfiles](09-dotfiles.md) and the `skills/` folder are served from.

Upgrading is re-running the same script: it overwrites the binary and leaves the clone and the shell line alone.

## From a checkout

```sh
git clone https://github.com/jterrazz/jterrazz-studio.git ~/Developer/jterrazz/jterrazz-studio
cd ~/Developer/jterrazz/jterrazz-studio
make install
source ~/.zshrc
```

`make install` builds and copies the binary to `~/.jterrazz/bin/j`, appends the `zshrc.sh` source line if it is missing, and migrates what an older layout left behind — `~/.config/jterrazz/jrc.json` becomes `~/.jterrazz/config.json`, and the Tailscale state directory moves with it, both only when the new location is empty. It also removes a stale `/usr/local/bin/j` from the install that predates `~/.jterrazz/bin`, and says what to run with `sudo` when it cannot.

`make check` answers whether the machine is equipped: `j` on `PATH`, and `~/.jterrazz` present. `make uninstall` removes the binary from both locations and touches nothing else — the user data survives.

## The footprint

Everything `j` owns on a machine lives under `~/.jterrazz/`:

```
~/.jterrazz/
├── bin/           # the binary
├── config.json    # the machine registry and the remote settings
├── tailscale/     # userspace daemon state
└── dns/           # generated DNS profiles
```

Outside it, `j` writes exactly two things a user did not create: the managed `Host` block in `~/.ssh/config` that follows the machine registry, and the `source` line in `~/.zshrc`. The config model itself is [Machines](06-machines.md)'s; what each `j config` item writes onto the machine is [Configuration](07-configuration.md)'s.

## After the install

```sh
j machine init     # bootstrap THIS machine (interactive)
j status           # the dashboard
j config           # configure the local machine
```

## Related

- [Machines](06-machines.md) — the registry `config.json` holds, and remote access.
- [Commands](05-commands.md) — every verb the installed binary answers to.
