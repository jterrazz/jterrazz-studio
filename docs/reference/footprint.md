# Footprint

Everything `j` writes on a machine. Everything it owns lives under `~/.jterrazz/`:

```text
~/.jterrazz/
├── bin/           # the binary
├── config.json    # the machine registry and the remote settings
├── tailscale/     # userspace daemon state
└── dns/           # generated DNS profiles
```

Outside it, and outside what a `j config` item the user chose writes onto the machine (each item's files are [Configuration](configuration.md)'s), three things land that a user did not create:

- the clone of this repository at `~/Developer/jterrazz-studio`, which the install script makes when that directory does not exist ([Getting started](../start/getting-started.md)); `j` reads its dotfiles and scripts from it, and `J_REPO_PATH` names another clone instead;
- the `source` line in `~/.zshrc`, which points at that clone's `dotfiles/applications/zsh/zshrc.sh`;
- the managed `Host` block in `~/.ssh/config` per registered machine with an ssh endpoint, which follows the registry: removing the machine gives the file back as it was.

The config model itself is [Machines](machines.md)'s.

## Related

- [Release](../guides/release.md): how the binary reaches `~/.jterrazz/bin/`.
- [Machines](machines.md): the registry `config.json` holds.
