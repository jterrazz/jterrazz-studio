package tool

import "testing"

func TestVersionParsers(t *testing.T) {
	tests := []struct {
		name     string
		parser   func(string) string
		given    string
		expected string
	}{
		{"TrimVersion basic", TrimVersion, "1.2.3", "1.2.3"},
		{"TrimVersion with v prefix", TrimVersion, "v1.2.3", "1.2.3"},
		{"TrimVersion with newline", TrimVersion, "v1.2.3\n", "1.2.3"},

		{"ParseBrewVersion multiline", ParseBrewVersion, "Homebrew 4.2.0\nHomebrew/homebrew-core", "4.2.0"},
		{"ParseBrewVersion single line", ParseBrewVersion, "Homebrew 5.0.7", "5.0.7"},

		{"ParseGitVersion basic", ParseGitVersion, "git version 2.39.0", "2.39.0"},
		{"ParseGitVersion apple", ParseGitVersion, "git version 2.39.0 (Apple Git-145)", "2.39.0"},

		{"ParseTmuxVersion basic", ParseTmuxVersion, "tmux 3.6a", "3.6a"},
		{"ParseTmuxVersion with newline", ParseTmuxVersion, "tmux 3.5a\n", "3.5a"},

		{"ParseTailscaleVersion single line", ParseTailscaleVersion, "1.86.2", "1.86.2"},
		{"ParseTailscaleVersion multi line", ParseTailscaleVersion, "1.86.2\ntailscale commit: abc", "1.86.2"},

		{"ParseGoVersion darwin", ParseGoVersion, "go version go1.21.0 darwin/arm64", "1.21.0"},
		{"ParseGoVersion linux", ParseGoVersion, "go version go1.22.5 linux/amd64", "1.22.5"},

		{"ParsePythonVersion basic", ParsePythonVersion, "Python 3.12.0", "3.12.0"},
		{"ParsePythonVersion with newline", ParsePythonVersion, "Python 3.11.5\n", "3.11.5"},

		{"ParseRustVersion", ParseRustVersion, "rustc 1.84.0 (9fc6b4312 2025-01-07)", "1.84.0"},

		{"ParseTerraformVersion", ParseTerraformVersion, "Terraform v1.5.7\non darwin_arm64", "1.5.7"},

		{"ParseAnsibleVersion core", ParseAnsibleVersion, "ansible [core 2.15.0]", "2.15.0"},
		{"ParseAnsibleVersion simple", ParseAnsibleVersion, "ansible 2.14.0", "2.14.0"},

		{"ParseMultipassVersion", ParseMultipassVersion, "multipass 1.12.0+mac\nmultipassd 1.12.0+mac", "1.12.0+mac"},

		{"ParseCodexVersion with prefix", ParseCodexVersion, "codex 0.1.0", "0.1.0"},
		{"ParseCodexVersion simple", ParseCodexVersion, "0.1.0", "0.1.0"},

		{"ParseMoleVersion with leading newline", ParseMoleVersion, "\nMole version 1.14.5\nmacOS: 14.0", "1.14.5"},
		{"ParseMoleVersion no leading newline", ParseMoleVersion, "Mole version 1.13.0\nmacOS: 13.0", "1.13.0"},

		{"ParseClaudeVersion with suffix", ParseClaudeVersion, "2.0.76 (Claude Code)", "2.0.76"},
		{"ParseClaudeVersion simple", ParseClaudeVersion, "1.0.0", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.given

			result := tt.parser(input)

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStripAnsi(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		{"plain text", "plain text", "plain text"},
		{"red color", "\x1b[31mred\x1b[0m", "red"},
		{"bold green", "\x1b[1;32mbold green\x1b[0m", "bold green"},
		{"mixed text and color", "no\x1b[33mcolor\x1b[0mhere", "nocolorhere"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.given

			result := StripAnsi(input)

			if result != tt.expected {
				t.Errorf("StripAnsi(%q) = %q, want %q", input, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		given    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"500 bytes", 500, "500 B"},
		{"1 KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"1 MB", 1048576, "1.0 MB"},
		{"1 GB", 1073741824, "1.0 GB"},
		{"1.5 GB", 1610612736, "1.5 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.given)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.given, result, tt.expected)
			}
		})
	}
}
