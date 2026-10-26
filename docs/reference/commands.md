# Commands

The everyday verbs. Machine management is in [Machines](machines.md); the interactive configurator in [Configuration](configuration.md).

## `workstation status`

Full-screen TUI dashboard, organised into 4 tabs (`←/→` to cycle, `1..4` to jump directly):

- **System**: live CPU/Memory/GPU/Network sparklines, top processes, network, Tailscale peers, and system health (firewall, DNS, etc.)
- **Workspace**: tracked git repos, Docker containers, project dependencies
- **Applications**: every catalogue tool with its version, by category
- **Configuration**: every `workstation config` item with its current state, grouped by category (Terminal / Security / Editor / System / Server / Network / Identity). Server subsection only shows on a server-registered machine.

Everything loads in parallel with a progress bar; the System tab's live readings refresh every second.

## `workstation install [tool...]`

With no arguments in an interactive terminal, opens a full-screen TUI with the same architecture as `workstation config`. The 18 tool families are grouped into 4 pages: **Dev** (Package Managers, Runtimes, Shell & Terminal, CLI Tools, Git, Editors & IDEs), **Infra** (Containers & VMs, Deploy, Remote Access, Security), **AI** (AI Agents, AI Tooling, AI Apps), **Apps** (Browsers, Communication, Productivity, Media, System Utilities). Each is split into collapsible family sections, rows in registry order. Press `i` to install, `u` to uninstall (with a Yes/No confirm first), `space` for a detail panel (description, replaces, method, formula, dependencies, post-install scripts), `tab` to fold a section, `←/→`/`1..4` to switch pages. Entries that are neither installable nor uninstallable through the registry (App Store apps, Xcode, and other manual/check-only entries) show a muted `(read-only)` marker instead of `i`/`u` hints. Installing a tool with unmet `Dependencies` refuses with "install `<dep>` first" instead of running.

```sh
workstation install                 # Open the interactive install TUI (TTY only)
workstation install --list          # Print the plain-text tool catalog instead
workstation install homebrew go node         # Install specific tools directly
workstation install claude codex ollama rtk  # AI tools
workstation install ghostty tmux zed         # Terminal + editor
```

Falls back to the plain-text catalog listing (the classic `✓`/`✗` table by category) when stdout isn't a terminal (scripts, CI, piped output) or when `--list` is passed. `workstation install <tool> [tool...]` (the direct, non-interactive form) is unchanged either way.

The tool catalogue and its categories live in [Tools and skills](tools-and-skills.md).

## `workstation upgrade [package...]`

```sh
workstation upgrade --all          # Upgrade all package managers (brew, npm, bun)
workstation upgrade --brew         # Upgrade Homebrew only
workstation upgrade node claude    # Upgrade specific packages
```

## `workstation clean [item...]`

```sh
workstation clean --all            # Clean everything (brew cache, docker, multipass, trash)
workstation clean docker trash     # Clean specific items
```

## `workstation run`

```sh
workstation run git feat "message"    # git add . && commit "feat: message"
workstation run git fix "message"     # git add . && commit "fix: message"
workstation run git wip               # git add --all && commit "WIP"
workstation run git unwip             # Undo last commit
workstation run git push              # Push current branch
workstation run git sync              # Fetch + pull
workstation run docker reset          # Remove all containers + images
workstation run docker clean          # System prune
```

## Shell shortcuts

Sourced via `dotfiles/applications/zsh/zshrc.sh` (see [Dotfiles](dotfiles.md)):

| Command | Action                     |
| ------- | -------------------------- |
| `jj`    | Attach tmux session `main` |
| `jc`    | Open Claude in tmux        |
| `jo`    | Open Codex in tmux         |
| `jg`    | Open Gemini in tmux        |

## Related

- [Machines](machines.md): `workstation machine` and `workstation remote`.
- [Configuration](configuration.md): `workstation config`.
- [Tools and skills](tools-and-skills.md): what `install`/`upgrade`/`clean` act on.
