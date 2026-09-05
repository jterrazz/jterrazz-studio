# Setup of Binaries

## Unified Command System (`j`) - Modular

A unified, modular command interface for all jterrazz tools and aliases. Each category is a separate module that auto-loads.

### Usage

```shell
j <category> <command> [args...]
```

### Available Commands

**Docker:**

```shell
j docker rm        # Remove all containers
j docker rmi       # Remove all images
j docker clean     # Clean up system (prune)
j docker reset     # Remove containers + images
j docker ps        # List containers
j docker images    # List images
```

**Dock:**

```shell
j dock add-spacer  # Add spacer to macOS Dock
j dock reset       # Reset dock to defaults
```

**Finder:**

```shell
j finder format ~/Documents            # Format with icon view (default)
j finder format ~/Downloads --view list # Format with list view
j finder reset ~/Desktop               # Reset to system defaults
```

**Git:**

```shell
j git feat "msg"   # Add all and commit with feat: prefix
j git fix "msg"    # Add all and commit with fix: prefix
j git chore "msg"  # Add all and commit with chore: prefix
j git push         # Push current branch to origin
j git sync         # Fetch and pull from remote
j git wip          # Quick WIP commit
j git unwip        # Undo WIP commit
j git status       # Git status
j git log          # Recent commits
j git branches     # List branches
```

**System:**

```shell
j system update    # Update system packages (Homebrew + npm global)
j system install   # Install development tools (brew, ohmyzsh, nvm, git-ssh)
j system clean     # Clean system (Homebrew, Docker, Multipass, trash)
```

**Help:**

```shell
j help             # Show all categories
j docker help      # Show docker commands
j dock help        # Show dock commands
j finder help      # Show finder commands
j git help         # Show git commands
j system help      # Show system commands
```

### Modular Architecture

The system automatically discovers modules from `configuration/binaries/j/*/commands.sh`:

```
configuration/binaries/
├── EXTENDING.md            # Module development guide
├── README.md               # This file
├── j/                      # Unified command system
│   ├── j.sh               # Main orchestrator
│   ├── docker/
│   │   └── commands.sh    # Docker module
│   ├── dock/
│   │   └── commands.sh    # Dock module
│   ├── finder/
│   │   └── commands.sh    # Finder module
│   ├── git/
│   │   └── commands.sh    # Git module
│   └── system/
│       └── commands.sh    # System module
└── zsh/
    └── zshrc.sh           # Shell integration & autocomplete
```

### Adding New Modules

1. Create `configuration/binaries/j/mymodule/commands.sh`
2. Implement the required functions (see `EXTENDING.md`)
3. Make it executable: `chmod +x commands.sh`
4. Reload: the module auto-loads on next `j` command

### Autocomplete Support

The system includes full zsh autocomplete support for all commands and nested subcommands:

```shell
j <TAB>                    # Shows: docker dock finder git system help
j system <TAB>             # Shows: update install clean help
j system install <TAB>     # Shows: brew ohmyzsh nvm git-ssh all help
j git <TAB>                # Shows: feat fix chore push sync wip unwip status log branches help
j docker <TAB>             # Shows: rm rmi clean reset ps images help
j dock <TAB>               # Shows: add-spacer reset help
j finder <TAB>             # Shows: format reset help
```

### Setup

#### Option 1: Automated Installation (Recommended)

Use the Makefile from the project root:

```shell
cd ~/Developer/jterrazz-configuration
make install
```

This will automatically add the source line to your `~/.zshrc` and provide helpful feedback.

**Additional commands:**

- `make check-installed` - Check if already installed
- `make uninstall` - Remove from ~/.zshrc
- `make help` - Show all options

#### Option 2: Manual Installation

Add to your `~/.zshrc`:

```shell
source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh
```

## Automated Installation

Instead of manual setup, use the automated installation commands:

```shell
# Full development environment setup
j system install all

# Or install components individually
j system install brew      # Homebrew + dev tools (ansible, terraform, kubectl, multipass, biome, bun, python)
j system install ohmyzsh   # Oh My Zsh + shell configuration
j system install nvm       # Node Version Manager + Node.js
j system install git-ssh   # Git SSH key generation + configuration
```

All installation commands are idempotent and safe to run multiple times.
