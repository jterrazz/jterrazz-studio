package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

const (
	autologinTargetUser  = "jterrazz.agent"
	autologinPasswordEnv = "AGENT_PASSWORD"
)

func installAutologin(values config.InputValues) error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	if _, err := exec.LookPath("sysadminctl"); err != nil {
		return fmt.Errorf("sysadminctl not found: %w", err)
	}
	if _, err := os.Stat("/Users/" + autologinTargetUser); err != nil {
		return fmt.Errorf("user %s does not exist on this Mac", autologinTargetUser)
	}

	password := values.Get("password")
	if password == "" {
		password = os.Getenv(autologinPasswordEnv)
	}
	if password == "" {
		return errors.New("empty password: refusing to set /etc/kcpassword with an empty value")
	}

	print.Header("install autologin", machineContext())
	print.Category("Before")
	dumpAutologinState()
	print.Empty()

	print.Category("Applying")
	if err := run("/usr/bin/defaults", "write", "/Library/Preferences/com.apple.loginwindow", "DisableFDEAutoLogin", "-bool", "NO"); err != nil {
		return err
	}
	print.Success("DisableFDEAutoLogin = NO")

	// outside: on macOS 26 `sysadminctl -autoLogin set` refuses on a FileVault disk
	// whatever DisableFDEAutoLogin says, so the login is written as /etc/kcpassword (the
	// XOR cipher loginwindow has read since 10.4) plus autoLoginUser; loginwindow at boot
	// still honors DisableFDEAutoLogin.
	if err := os.WriteFile("/etc/kcpassword", encodeKCPassword(password), 0o600); err != nil {
		return fmt.Errorf("write /etc/kcpassword: %w", err)
	}
	if err := os.Chmod("/etc/kcpassword", 0o600); err != nil {
		return err
	}
	// outside: /etc/kcpassword must be owned by root:wheel.
	_ = os.Chown("/etc/kcpassword", 0, 0)
	print.Success("/etc/kcpassword written (root:wheel 0600)")

	if err := run("/usr/bin/defaults", "write", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser", "-string", autologinTargetUser); err != nil {
		return fmt.Errorf("defaults write autoLoginUser: %w", err)
	}
	print.Success("autoLoginUser = " + autologinTargetUser)

	if _, err := os.Stat("/etc/kcpassword"); err != nil {
		return errors.New("/etc/kcpassword missing after write")
	}
	out, err := runQuiet("/usr/bin/defaults", "read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser")
	if err != nil || strings.TrimSpace(out) != autologinTargetUser {
		return fmt.Errorf("autoLoginUser cross-check failed (got %q)", strings.TrimSpace(out))
	}

	// value: clearing the screen-locked-after-resume flag keeps the auto-login session
	// interactive; lock-after-login owns the lock itself.
	_ = exec.Command("/usr/bin/defaults", "delete", "/Library/Preferences/com.apple.loginwindow", "autoLoginUserScreenLocked").Run()

	print.Empty()
	print.Category("After")
	dumpAutologinState()
	print.Empty()
	print.Dim("Verify end-to-end: `sudo fdesetup authrestart -delayminutes 0` (or `j machine restart` from the MacBook).")
	return nil
}

func uninstallAutologin() error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	print.Header("uninstall autologin", machineContext())
	_ = run("/usr/sbin/sysadminctl", "-autologin", "off")
	_ = os.Remove("/etc/kcpassword")
	for _, key := range []string{"autoLoginUser", "oneTimeAutoLogin", "autoLoginUserScreenLocked"} {
		_ = exec.Command("/usr/bin/defaults", "delete", "/Library/Preferences/com.apple.loginwindow", key).Run()
	}
	print.Empty()
	print.Category("After")
	dumpAutologinState()
	return nil
}

func checkAutologinInstalled() config.CheckResult {
	out, err := runQuiet("/usr/bin/defaults", "read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser")
	if err != nil || strings.TrimSpace(out) != autologinTargetUser {
		return config.CheckResult{}
	}
	if _, err := os.Stat("/etc/kcpassword"); err != nil {
		return config.CheckResult{}
	}
	return config.InstalledWithDetail(autologinTargetUser)
}

func dumpAutologinState() {
	if out, err := runQuiet("/usr/sbin/sysadminctl", "-autologin", "status"); err == nil || out != "" {
		print.Linef("  sysadminctl: %s", oneLineOrDash(out))
	} else {
		print.Linef("  sysadminctl: %s", err.Error())
	}

	for _, key := range []string{"autoLoginUser", "DisableFDEAutoLogin"} {
		out, err := runQuiet("/usr/bin/defaults", "read", "/Library/Preferences/com.apple.loginwindow", key)
		if err != nil {
			print.Linef("  %s: (unset)", key)
			continue
		}
		print.Linef("  %s: %s", key, oneLineOrDash(out))
	}

	if info, err := os.Stat("/etc/kcpassword"); err == nil {
		print.Linef("  /etc/kcpassword: present (%d bytes, mode %v)", info.Size(), info.Mode().Perm())
	} else {
		print.Linef("  /etc/kcpassword: absent")
	}

	if owner, err := runQuiet("/usr/bin/stat", "-f", "%Su", "/dev/console"); err == nil {
		print.Linef("  /dev/console owner: %s", oneLineOrDash(owner))
	}
}

func oneLineOrDash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return strings.Join(strings.Fields(strings.Split(s, "\n")[0]), " ")
}
