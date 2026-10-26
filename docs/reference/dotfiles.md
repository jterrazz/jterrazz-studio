# Dotfiles

`dotfiles/applications/` holds the versioned configuration `workstation config` installs onto a machine, one folder per application. Keeping them in the repo means every machine converges on the same setup, and changes are diffable and reviewed.

| Application | Contents                                                                                                     |
| ----------- | ------------------------------------------------------------------------------------------------------------ |
| `ghostty`   | Terminal `config` + `themes/`                                                                                |
| `starship`  | `starship.toml` prompt                                                                                       |
| `tmux`      | `tmux.conf`                                                                                                  |
| `vscode`    | `Default.code-profile`                                                                                       |
| `zed`       | `settings.json` + `.oxfmtrc.json`: the formatter values Zed applies when it saves its own settings           |
| `zsh`       | `zshrc.sh`, sourced into `~/.zshrc` by `make install`; defines the `jj` / `jc` / `jo` / `jg` shell shortcuts |

The relevant `workstation config` items (Terminal, Editor categories) point at these files (see [Configuration](configuration.md)). The shell shortcuts are listed in [Commands](commands.md).

## Related

- [Configuration](configuration.md): the items that install these configs.
- [Release](../guides/release.md): how `zshrc.sh` gets sourced onto a machine.
