# Machines

`workstation machine` manages a small registry of the machines you own (typically a client box, your laptop, and one or more servers) and runs status checks, remote actions, and server-only configuration.

## The registry & config model

Every machine has an alias, a role (`client` or `server`), and an optional SSH endpoint. The registry lives in `~/.studio/config.json` and is the **single source of truth**: adding a machine also writes a managed `Host` block in `~/.ssh/config`.

```jsonc
{
    "remote": { "mode": "userspace", "auth_method": "oauth" },
    "self": "macbook",
    "machines": {
        "macbook": { "role": "client" },
        "mac-mini": { "role": "server", "ssh": "agent@192.168.1.106" },
    },
}
```

The role decides what `workstation machine status` reports and which items `workstation config` exposes for this box (see [Configuration](configuration.md)).

```sh
workstation machine init                                                 # Bootstrap THIS machine (interactive)
workstation machine list                                                 # Table of registered machines (* marks self)
workstation machine add mac-mini --role server --ssh agent@192.168.1.106 # Add a remote
workstation machine add macbook  --role client                           # Add a local-only entry
workstation machine remove mac-mini                                      # Refuses if alias is self
```

## Inspect

```sh
workstation machine status              # FileVault, SSH, plus services (server role only)
workstation machine probe <alias>       # ping + ssh + OpenClaw gateway port + console owner
workstation machine restart <alias> -y  # FileVault-aware authrestart, waits for SSH to come back
workstation machine unlock <alias>      # Pre-boot SSH session to enter the FileVault password
```

`status` runs locally and adapts to the role:

- **client**: Machine state only: FileVault, SSH (port 22).
- **server**: Machine state + Services: OpenClaw runtime, OpenClaw config, channel health (Slack/Telegram/BlueBubbles), OrbStack.

`probe`/`restart`/`unlock` resolve the SSH endpoint from the registry. They refuse to act on the alias marked as self.

## Remote (Tailscale)

```sh
workstation remote up       # Connect (userspace mode, SSH enabled, keep-awake)
workstation remote down     # Disconnect and stop daemon
workstation remote status   # Show connection state
```

Supports `auto`/`userspace` mode and `oauth`/`authkey` authentication. Daemon state lives under `~/.studio/tailscale/`. To change the endpoint settings, open `workstation config` and switch to the Remote tab.

## Related

- [Release](../guides/release.md): the `~/.studio/` footprint and how `workstation` gets installed.
- [Configuration](configuration.md): role-gated `workstation config` items.
