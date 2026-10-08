package config

import (
	"testing"
)

func TestSetSelfAndSelfMachine(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := SetSelf("macbook"); err == nil {
		t.Fatal("SetSelf on unknown alias should fail, got nil")
	}

	_ = AddMachine("macbook", Machine{Role: RoleClient})
	if err := SetSelf("macbook"); err != nil {
		t.Fatalf("SetSelf(macbook) error = %v", err)
	}

	alias, m, ok := SelfMachine()
	if !ok || alias != "macbook" || m.Role != RoleClient {
		t.Fatalf("SelfMachine() = (%q, %+v, %v)", alias, m, ok)
	}
}
