package browser

import (
	"os"
	"syscall"
)

func OpenURL(url string) error {
	var err error

	_, err = os.Stat("/usr/bin/open")
	if err == nil {
		return syscall.Exec("/usr/bin/open", []string{"/usr/bin/open", url}, os.Environ())
	}

	return nil
}

func OpenURLFork(url string) error {
	var err error

	_, err = os.Stat("/usr/bin/open")
	if err == nil {
		_, err := syscall.ForkExec("/usr/bin/open", []string{"/usr/bin/open", url}, &syscall.ProcAttr{Env: os.Environ()})
		return err
	}

	return nil
}
