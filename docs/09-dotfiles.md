# Dotfiles

`dotfiles/applications/` holds the versioned configuration `j config` installs onto a machine — one folder per application. Keeping them in the repo means every machine converges on the same setup, and changes are diffable and reviewed.

| Application | Contents                                                                                                      |
| ----------- | ------------------------------------------------------------------------------------------------------------- |
| `ghostty`   | Terminal `config` + `themes/`                                                                                 |
| `starship`  | `starship.toml` prompt                                                                                        |
| `tmux`      | `tmux.conf`                                                                                                   |
| `vscode`    | `Default.code-profile`                                                                                        |
| `zed`       | `settings.json`                                                                                               |
| `zsh`       | `zshrc.sh` — sourced into `~/.zshrc` by `make install`; defines the `jj` / `jc` / `jo` / `jg` shell shortcuts |

The relevant `j config` items (Terminal, Editor categories) point at these files — see [Configuration](07-configuration.md). The shell shortcuts are listed in [Commands](05-commands.md).

## Related

- [Configuration](07-configuration.md) — the items that install these configs.
- [Operating](04-operating.md) — how `zshrc.sh` gets sourced onto a machine.
