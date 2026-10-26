package commands

import (
	"strings"

	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/print"
)

func installSshd() error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	print.Header("install sshd", machineContext())
	print.Category("Before")
	dumpSshdState()
	print.Empty()

	print.Category("Applying")
	if out, err := runQuiet("/usr/sbin/systemsetup", "-setremotelogin", "on"); err != nil {
		print.Warning("systemsetup -setremotelogin on: " + oneLineOrDash(out))
	} else {
		print.Success("Remote Login enabled")
	}

	print.Empty()
	print.Category("After")
	dumpSshdState()
	return nil
}

func uninstallSshd() error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	print.Header("uninstall sshd", machineContext())
	print.Category("Before")
	dumpSshdState()
	print.Empty()

	print.Category("Applying")
	if out, err := runQuiet("/usr/sbin/systemsetup", "-setremotelogin", "off"); err != nil {
		print.Warning("systemsetup -setremotelogin off: " + oneLineOrDash(out))
	} else {
		print.Success("Remote Login disabled")
	}

	print.Empty()
	print.Category("After")
	dumpSshdState()
	return nil
}

// outside: the TUI runs as the user, and `systemsetup -getremotelogin` then exits 0 with
// its error on stdout; `launchctl print-disabled system` is readable without root and
// carries the sshd state systemsetup sets.
func checkSshdInstalled() config.CheckResult {
	out, err := runQuiet("/bin/launchctl", "print-disabled", "system")
	if err != nil {
		return config.CheckResult{}
	}
	if !strings.Contains(out, `"com.openssh.sshd" => enabled`) {
		return config.CheckResult{}
	}
	return config.InstalledWithDetail("Remote Login on")
}

func dumpSshdState() {
	if out, err := runQuiet("/usr/sbin/systemsetup", "-getremotelogin"); err == nil {
		print.Linef("  systemsetup: %s", oneLineOrDash(out))
	} else {
		print.Linef("  systemsetup: error %s", err.Error())
	}
}
