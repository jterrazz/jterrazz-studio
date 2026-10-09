package commands

import "github.com/jterrazz/jterrazz-studio/src/internal/config"

func init() {
	config.RegisterServerActions(config.ServerActions{
		AutologinInstall:        installAutologin,
		AutologinUninstall:      uninstallAutologin,
		AutologinCheck:          checkAutologinInstalled,
		PowerInstall:            installPower,
		PowerUninstall:          uninstallPower,
		PowerCheck:              checkPowerInstalled,
		LockAfterLoginInstall:   installLockAfterLogin,
		LockAfterLoginUninstall: uninstallLockAfterLogin,
		LockAfterLoginCheck:     checkLockAfterLoginInstalled,
		SshdInstall:             installSshd,
		SshdUninstall:           uninstallSshd,
		SshdCheck:               checkSshdInstalled,
	})
}
