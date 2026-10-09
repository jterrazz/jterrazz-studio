package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// value: the alias sits in both markers, so each block names itself to a human who opens the file.
const (
	sshAliasMarkerPrefix = "# >>> jterrazz machine "
	sshAliasMarkerSuffix = " >>>"
	sshAliasEndPrefix    = "# <<< jterrazz machine "
	sshAliasEndSuffix    = " <<<"
)

func sshAliasBegin(alias string) string { return sshAliasMarkerPrefix + alias + sshAliasMarkerSuffix }
func sshAliasEnd(alias string) string   { return sshAliasEndPrefix + alias + sshAliasEndSuffix }

func sshConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".ssh", "config")
}

// WriteSSHAlias writes or refreshes alias's managed Host block, refusing a foreign Host line.
func WriteSSHAlias(alias string, m Machine) error {
	if strings.TrimSpace(alias) == "" {
		return errors.New("alias is required")
	}
	user, host, err := parseSSHEndpoint(m.SSH)
	if err != nil {
		return err
	}

	path := sshConfigPath()
	if path == "" {
		return errors.New("cannot resolve ~/.ssh/config: HOME unset")
	}
	if mkdirErr := os.MkdirAll(filepath.Dir(path), 0o700); mkdirErr != nil {
		return mkdirErr
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	block := buildSSHAliasBlock(alias, host, user, m.Identity)
	updated, conflict := upsertSSHAliasBlock(existing, alias, block)
	if conflict != "" {
		return fmt.Errorf("refusing to modify %s: %s", path, conflict)
	}
	if bytes.Equal(updated, existing) {
		return nil
	}
	return os.WriteFile(path, updated, 0o600)
}

// RemoveSSHAlias removes alias's managed block; a missing file or block is a no-op.
func RemoveSSHAlias(alias string) error {
	path := sshConfigPath()
	if path == "" {
		return nil
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	updated := removeSSHAliasBlock(existing, alias)
	if bytes.Equal(updated, existing) {
		return nil
	}
	return os.WriteFile(path, updated, 0o600)
}

func parseSSHEndpoint(endpoint string) (user, host string, err error) {
	endpoint = strings.TrimSpace(endpoint)
	at := strings.Index(endpoint, "@")
	if at <= 0 || at == len(endpoint)-1 {
		return "", "", fmt.Errorf("ssh endpoint %q must be user@host", endpoint)
	}
	return endpoint[:at], endpoint[at+1:], nil
}

func buildSSHAliasBlock(alias, host, user, identity string) string {
	if identity == "" {
		identity = "~/.ssh/id_ed25519"
	}
	return strings.Join([]string{
		sshAliasBegin(alias),
		"Host " + alias,
		"  HostName " + host,
		"  User " + user,
		"  IdentityFile " + identity,
		"  ServerAliveInterval 30",
		"  ServerAliveCountMax 3",
		sshAliasEnd(alias),
		"",
	}, "\n")
}

func upsertSSHAliasBlock(existing []byte, alias, block string) (updated []byte, conflict string) {
	text := string(existing)
	begin := sshAliasBegin(alias)
	end := sshAliasEnd(alias)

	beginIdx := strings.Index(text, begin)
	endIdx := strings.Index(text, end)
	if beginIdx >= 0 && endIdx > beginIdx {
		endLine := endIdx + len(end)
		if endLine < len(text) && text[endLine] == '\n' {
			endLine++
		}
		return []byte(text[:beginIdx] + block + text[endLine:]), ""
	}

	if conflict := findForeignHostAlias(text, alias); conflict != "" {
		return existing, conflict
	}

	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	if text != "" {
		text += "\n"
	}
	return []byte(text + block), ""
}

func removeSSHAliasBlock(existing []byte, alias string) []byte {
	text := string(existing)
	begin := sshAliasBegin(alias)
	end := sshAliasEnd(alias)

	beginIdx := strings.Index(text, begin)
	endIdx := strings.Index(text, end)
	if beginIdx < 0 || endIdx <= beginIdx {
		return existing
	}
	endLine := endIdx + len(end)
	if endLine < len(text) && text[endLine] == '\n' {
		endLine++
	}
	// invariant: an append always writes exactly one blank line before its block
	// (upsertSSHAliasBlock), even after a blank line the user wrote, so removal
	// takes exactly one and an add then remove gives the file back byte for byte.
	before, after := text[:beginIdx], text[endLine:]
	switch {
	case strings.HasSuffix(before, "\n\n"):
		before = before[:len(before)-1]
	case before == "" && strings.HasPrefix(after, "\n"):
		after = after[1:]
	}
	return []byte(before + after)
}

// invariant: any managed block counts, not only this alias's: a neighbour's block is
// not a foreign Host.
func findForeignHostAlias(text, alias string) string {
	inManaged := false
	for i, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, sshAliasMarkerPrefix) && strings.HasSuffix(trimmed, sshAliasMarkerSuffix) {
			inManaged = true
			continue
		}
		if strings.HasPrefix(trimmed, sshAliasEndPrefix) && strings.HasSuffix(trimmed, sshAliasEndSuffix) {
			inManaged = false
			continue
		}
		if inManaged {
			continue
		}
		if !strings.HasPrefix(trimmed, "Host ") {
			continue
		}
		for _, name := range strings.Fields(trimmed)[1:] {
			if name == alias {
				return fmt.Sprintf("Host %s already defined at line %d (outside any managed block)", alias, i+1)
			}
		}
	}
	return ""
}
