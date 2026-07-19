# Commands

The everyday verbs. Machine management is in [Machines](03-machines.md); the interactive configurator in [Configuration](04-configuration.md).

## `j status`

Full-screen TUI dashboard, organised into 4 tabs (`←/→` to cycle, `1..4` to jump directly):

- **System** — live CPU/Memory/GPU/Network sparklines, top processes, network, Tailscale peers, and system health (firewall, DNS, etc.)
- **Workspace** — tracked git repos, Docker containers, project dependencies
- **Applications** — 100+ tracked tools with versions, by category
- **Configuration** — every `j config` item with its current state, grouped by category (Terminal / Security / Editor / System / Server / Network / Identity). Server subsection only shows on a server-registered machine.

Everything loads in parallel with a progress bar; the System tab's live readings refresh every second.

## `j install [tool...]`

With no arguments in an interactive terminal, opens a full-screen TUI (same architecture as `j config`): the 18 tool families are grouped into 4 pages — **Dev** (Package Managers, Runtimes, Shell & Terminal, CLI Tools, Git, Editors & IDEs), **Infra** (Containers & VMs, Deploy, Remote Access, Security), **AI** (AI Agents, AI Tooling, AI Apps), **Apps** (Browsers, Communication, Productivity, Media, System Utilities) — each split into collapsible family sections, rows in registry order. Press `i` to install, `u` to uninstall (with a Yes/No confirm first), `space` for a detail panel (description, replaces, method, formula, dependencies, post-install scripts), `tab` to fold a section, `←/→`/`1..4` to switch pages. Entries that are neither installable nor uninstallable through the registry (App Store apps, Xcode, and other manual/check-only entries) show a muted `(read-only)` marker instead of `i`/`u` hints. Installing a tool with unmet `Dependencies` refuses with "install `<dep>` first" instead of running.

```sh
j install                 # Open the interactive install TUI (TTY only)
j install --list          # Print the plain-text tool catalog instead
j install homebrew go node         # Install specific tools directly
j install claude codex ollama rtk  # AI tools
j install ghostty tmux zed         # Terminal + editor
```

Falls back to the plain-text catalog listing (the classic `✓`/`✗` table by category) when stdout isn't a terminal — scripts, CI, piped output — or when `--list` is passed. `j install <tool> [tool...]` (the direct, non-interactive form) is unchanged either way.

The tool catalogue and its categories live in [Tools and skills](05-tools-and-skills.md).

## `j upgrade [package...]`

```sh
j upgrade --all          # Upgrade all package managers (brew, npm, bun)
j upgrade --brew         # Upgrade Homebrew only
j upgrade node claude    # Upgrade specific packages
```

## `j clean [item...]`

```sh
j clean --all            # Clean everything (brew cache, docker, multipass, trash)
j clean docker trash     # Clean specific items
```

## `j run`

```sh
j run git feat "message"    # git add . && commit "feat: message"
j run git fix "message"     # git add . && commit "fix: message"
j run git wip               # git add --all && commit "WIP"
j run git unwip             # Undo last commit
j run git push              # Push current branch
j run git sync              # Fetch + pull
j run docker reset          # Remove all containers + images
j run docker clean          # System prune
```

## Shell shortcuts

Sourced via `dotfiles/applications/zsh/zshrc.sh` (see [Dotfiles](06-dotfiles.md)):

| Command | Action |
|---------|--------|
| `jj` | Attach tmux session `main` |
| `jc` | Open Claude in tmux |
| `jo` | Open Codex in tmux |
| `jg` | Open Gemini in tmux |

## Related

- [Machines](03-machines.md) — `j machine` and `j remote`.
- [Configuration](04-configuration.md) — `j config`.
- [Tools and skills](05-tools-and-skills.md) — what `install`/`upgrade`/`clean` act on.
