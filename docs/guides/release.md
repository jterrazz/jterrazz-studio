# Release

This repository ships one artefact, the `workstation` binary, and it reaches a machine two ways: a tagged GitHub release, or a build from a checkout.

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

`.github/workflows/release.yaml` fires on `push: tags: v*` and calls the shared `studio-actions` `release-go.yaml` workflow with the binary name `workstation`, the build path `./src/cmd/workstation`, Go 1.26, and four targets: `darwin/arm64`, `darwin/amd64`, `linux/arm64`, `linux/amd64`. Nothing publishes on a merge to `main`: a merge only runs the validate workflow ([Test](test.md)).

Each target is published as an asset named `workstation-<os>-<arch>`, which is the name the install script asks for. A release that renames those assets breaks every machine's install path.

## How a machine takes it

A machine takes the release through the one-line install of [Getting started](../start/getting-started.md), which downloads the asset named for its platform; re-running it is the upgrade.

## From a checkout

```sh
git clone https://github.com/jterrazz/studio-workstation.git ~/Developer/studio-workstation
cd ~/Developer/studio-workstation
make provision
source ~/.zshrc
```

`make provision` builds and copies the binary to `~/.studio/bin/workstation` and appends the `zshrc.sh` source line if it is missing.

`make deprovision` removes the binary and touches nothing else; the user data survives.

## Related

- [Footprint](../reference/footprint.md): everything the binary writes on the machine it reaches.
- [Machines](../reference/machines.md): the registry `config.json` holds, and remote access.
- [Commands](../reference/commands.md): every verb the installed binary answers to.
