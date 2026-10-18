# Release

This repository ships one artefact, the `j` binary, and it reaches a machine two ways: a tagged GitHub release, or a build from a checkout.

| Section                | Answers                               |
| ---------------------- | ------------------------------------- |
| The tagged release     | What a `v*` tag builds and publishes  |
| How a machine takes it | The one-line install, and the upgrade |
| From a checkout        | `make provision`, `make deprovision`  |

## The tagged release

There is no deploy: the release IS the delivery. Pushing a version tag is the whole gesture.

```sh
git tag v1.0.0
git push --tags
```

`.github/workflows/release.yaml` fires on `push: tags: v*` and calls the shared `jterrazz-actions` `release-go.yaml` workflow with the binary name `j`, the build path `./src/cmd/j`, Go 1.24, and four targets: `darwin/arm64`, `darwin/amd64`, `linux/arm64`, `linux/amd64`. Nothing publishes on a merge to `main`: a merge only runs the validate workflow ([Test](test.md)).

Each target is published as an asset named `j-<os>-<arch>`, which is the name the install script asks for. A release that renames those assets breaks every machine's install path.

## How a machine takes it

A machine takes the release through the one-line install of [Getting started](../start/getting-started.md), which downloads the asset named for its platform; re-running it is the upgrade.

## From a checkout

```sh
git clone https://github.com/jterrazz/jterrazz-studio.git ~/Developer/jterrazz-studio
cd ~/Developer/jterrazz-studio
make provision
source ~/.zshrc
```

`make provision` builds and copies the binary to `~/.jterrazz/bin/j`, appends the `zshrc.sh` source line if it is missing, and migrates what an older layout left behind: `~/.config/jterrazz/jrc.json` becomes `~/.jterrazz/config.json`, and the Tailscale state directory moves with it, both only when the new location is empty. It also removes a stale `/usr/local/bin/j` from the install that predates `~/.jterrazz/bin`, and says what to run with `sudo` when it cannot.

`make deprovision` removes the binary from both locations and touches nothing else; the user data survives.

## Related

- [Footprint](../reference/footprint.md): everything the binary writes on the machine it reaches.
- [Machines](../reference/machines.md): the registry `config.json` holds, and remote access.
- [Commands](../reference/commands.md): every verb the installed binary answers to.
