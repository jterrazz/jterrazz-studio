package config

import (
	"testing"
)

func TestParseSSHEndpoint(t *testing.T) {
	tests := []struct {
		in       string
		wantUser string
		wantHost string
		wantErr  bool
	}{
		{"agent@192.168.1.106", "agent", "192.168.1.106", false},
		{" admin@host ", "admin", "host", false},
		{"justhost", "", "", true},
		{"@host", "", "", true},
		{"user@", "", "", true},
		{"", "", "", true},
	}
	for _, tt := range tests {
		u, h, err := parseSSHEndpoint(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseSSHEndpoint(%q) err=%v wantErr=%v", tt.in, err, tt.wantErr)
		}
		if u != tt.wantUser || h != tt.wantHost {
			t.Errorf("parseSSHEndpoint(%q) = (%q, %q), want (%q, %q)", tt.in, u, h, tt.wantUser, tt.wantHost)
		}
	}
}

func TestRemoveSSHAliasMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := RemoveSSHAlias("nope"); err != nil {
		t.Fatalf("expected no-op on missing file, got %v", err)
	}
}
