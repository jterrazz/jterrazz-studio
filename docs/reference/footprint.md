# Footprint

Everything `j` writes on a machine. Everything it owns lives under `~/.jterrazz/`:

```text
~/.jterrazz/
├── bin/           # the binary
├── config.json    # the machine registry and the remote settings
├── tailscale/     # userspace daemon state
└── dns/           # generated DNS profiles
```

Outside it, `j` writes exactly two things a user did not create: the managed `Host` block in `~/.ssh/config` that follows the machine registry, and the `source` line in `~/.zshrc`. The config model itself is [Machines](machines.md)'s; what each `j config` item writes onto the machine is [Configuration](configuration.md)'s.

## Related

- [Release](../guides/release.md) — how the binary reaches `~/.jterrazz/bin/`.
- [Machines](machines.md) — the registry `config.json` holds.
