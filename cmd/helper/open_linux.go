package helper

import (
	"os"
	"syscall"
)

func OpenURL(url string) error {
	var err error

	_, err = os.Stat("/usr/bin/xdg-open")
	if err == nil {
		return syscall.Exec("/usr/bin/xdg-open", []string{"/usr/bin/xdg-open", url}, os.Environ())
	}

	_, err = os.Stat("/usr/bin/sensible-browser")
	if err == nil {
		return syscall.Exec("/usr/bin/sensible-browser", []string{"/usr/bin/sensible-browser", url}, os.Environ())
	}

	return nil
}