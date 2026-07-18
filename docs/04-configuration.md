# Configuration

`j config` is an interactive TUI for configuring the local machine, organised into 3 tabs (`←/→` to cycle, `1..3` to jump directly):

- **Configuration** — installable items grouped by category. Sections are collapsible; items show their current state.
- **Skills** — install / list / remove AI agent skills (requires the `skills` CLI on PATH — see [Tools and skills](05-tools-and-skills.md)).
- **Remote** — read-only summary of the Tailscale endpoint; press `i` to open a form that rewrites `~/.jterrazz/config.json`.

```
 j config                                                self: mac-mini · server
 [Configuration]  Skills  Remote
 ──────────────────────────────────────────────────────────────────────────────
 ▾ Terminal               3/3
   ✓ ghostty
   ✓ tmux
   ✓ hushlogin

 ▸ Security               4/5
 ▾ Editor                 1/1
   ✓ zed

 ▸ System                 2/4
 ▾ Server                2/4
   ✓ autologin
 ▶ ✗ power
   ✓ lock-after-login
   ✗ sshd
 ──────────────────────────────────────────────────────────────────────────────
 ▶ power |  i install   space details
```

## Categories

The Configuration tab groups items by category (Server only appears when the current machine is registered as `server` — see [Machines](03-machines.md)):

- **Terminal** — ghostty, tmux, hushlogin
- **Security** — GPG commit signing, SSH keygen, GitHub CLI auth, encrypted DNS (Quad9), Spotlight exclusion
- **Editor** — Zed config
- **System** — JAVA_HOME, nvm, dock reset/spacer
- **Server** — autologin, power policy, lock-after-login, sshd

Items that need extra inputs (e.g. autologin's password) open a modal form before installing — built on [Charm's huh](https://github.com/charmbracelet/huh). Set `AGENT_PASSWORD` in your environment to pre-fill the autologin password field.

## Keys

| Key | Action |
|---|---|
| `←` `→` `1..3` | switch tab |
| `↑` `↓` `j` `k` | navigate |
| `tab` | collapse/expand current section |
| `space` | toggle the inline detail panel (Configuration tab) |
| `i` | install the current item (or open the reconfigure form on the Remote tab) |
| `u` | uninstall (only for toggleable items that are currently installed) |
| `q` `esc` | quit |

## Related

- [Tools and skills](05-tools-and-skills.md) — the Skills tab and the underlying registries.
- [Machines](03-machines.md) — how the machine role gates which items appear.
