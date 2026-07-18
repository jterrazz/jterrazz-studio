# Machines

`j machine` manages a small registry of the machines you own — typically a client box (your laptop) and one or more servers — and runs status checks, remote actions, and server-only configuration.

## The registry & config model

Every machine has an alias, a role (`client` or `server`), and an optional SSH endpoint. The registry lives in `~/.jterrazz/config.json` and is the **single source of truth** — adding a machine also writes a managed `Host` block in `~/.ssh/config`.

```jsonc
{
  "remote":    { "mode": "userspace", "auth_method": "oauth" },
  "self":      "macbook",
  "machines": {
    "macbook":  { "role": "client" },
    "mac-mini": { "role": "server", "ssh": "agent@192.168.1.106" }
  }
}
```

The role decides what `j machine status` reports and which items `j config` exposes for this box (see [Configuration](04-configuration.md)).

```sh
j machine init                                                 # Bootstrap THIS machine (interactive)
j machine list                                                 # Table of registered machines (* marks self)
j machine add mac-mini --role server --ssh agent@192.168.1.106 # Add a remote
j machine add macbook  --role client                           # Add a local-only entry
j machine remove mac-mini                                      # Refuses if alias is self
```

## Inspect

```sh
j machine status              # FileVault, SSH, plus services (server role only)
j machine probe <alias>       # ping + ssh + OpenClaw gateway port + console owner
j machine restart <alias> -y  # FileVault-aware authrestart, waits for SSH to come back
j machine unlock <alias>      # Pre-boot SSH session to enter the FileVault password
```

`status` runs locally and adapts to the role:

- **client** — Machine state only: FileVault, SSH (port 22).
- **server** — Machine state + Services: OpenClaw runtime, OpenClaw config, channel health (Slack/Telegram/BlueBubbles), OrbStack.

`probe`/`restart`/`unlock` resolve the SSH endpoint from the registry. They refuse to act on the alias marked as self.

## Remote (Tailscale)

```sh
j remote up       # Connect (userspace mode, SSH enabled, keep-awake)
j remote down     # Disconnect and stop daemon
j remote status   # Show connection state
```

Supports `auto`/`userspace` mode and `oauth`/`authkey` authentication. Daemon state lives under `~/.jterrazz/tailscale/`. To change the endpoint settings, open `j config` and switch to the Remote tab.

## Related

- [Getting started](01-getting-started.md) — `~/.jterrazz/` layout.
- [Configuration](04-configuration.md) — role-gated `j config` items.
