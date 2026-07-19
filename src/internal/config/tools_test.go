package config

import "testing"

func TestToolUninstallCustomFnWins(t *testing.T) {
	called := false
	tool := Tool{
		Name:        "x",
		Method:      InstallManual, // would otherwise error
		UninstallFn: func() error { called = true; return nil },
	}
	if err := tool.Uninstall(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("UninstallFn was not invoked")
	}
}

func TestToolUninstallErrorsOnCheckOnlyMethods(t *testing.T) {
	for _, method := range []InstallMethod{InstallManual, InstallXcode, InstallMAS, InstallNvm, InstallMethod("")} {
		tool := Tool{Name: "x", Method: method}
		if err := tool.Uninstall(); err == nil {
			t.Errorf("method %q: expected error, got nil", method)
		}
	}
}

func TestToolUninstallable(t *testing.T) {
	cases := []struct {
		name string
		tool Tool
		want bool
	}{
		{"custom UninstallFn", Tool{UninstallFn: func() error { return nil }}, true},
		{"brew formula", Tool{Method: InstallBrewFormula}, true},
		{"brew cask", Tool{Method: InstallBrewCask}, true},
		{"npm", Tool{Method: InstallNpm}, true},
		{"bun", Tool{Method: InstallBun}, true},
		{"uv", Tool{Method: InstallUV}, true},
		{"mas", Tool{Method: InstallMAS}, false},
		{"xcode", Tool{Method: InstallXcode}, false},
		{"manual", Tool{Method: InstallManual}, false},
		{"nvm", Tool{Method: InstallNvm}, false},
		{"no method, no fn", Tool{}, false},
	}
	for _, c := range cases {
		if got := c.tool.Uninstallable(); got != c.want {
			t.Errorf("%s: Uninstallable() = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestToolInstallable(t *testing.T) {
	cases := []struct {
		name string
		tool Tool
		want bool
	}{
		{"custom InstallFn", Tool{InstallFn: func() error { return nil }}, true},
		{"brew formula", Tool{Method: InstallBrewFormula}, true},
		{"brew cask", Tool{Method: InstallBrewCask}, true},
		{"npm", Tool{Method: InstallNpm}, true},
		{"bun", Tool{Method: InstallBun}, true},
		{"uv", Tool{Method: InstallUV}, true},
		{"mas", Tool{Method: InstallMAS}, false},
		{"xcode", Tool{Method: InstallXcode}, false},
		{"manual", Tool{Method: InstallManual}, false},
		{"nvm", Tool{Method: InstallNvm}, false},
		{"no method, no fn", Tool{}, false},
	}
	for _, c := range cases {
		if got := c.tool.Installable(); got != c.want {
			t.Errorf("%s: Installable() = %v, want %v", c.name, got, c.want)
		}
	}
}
