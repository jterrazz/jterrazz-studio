package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

func requireDarwin() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this command only runs on macOS (got %s)", runtime.GOOS)
	}
	return nil
}

func requireRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	return errors.New("requires root — re-run with: sudo " + strings.Join(os.Args, " "))
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runQuiet(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimRight(string(out), "\n"), err
}

func failOn(err error) {
	if err == nil {
		return
	}
	print.Error(err.Error())
	os.Exit(1)
}
