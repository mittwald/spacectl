package helper

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