package browser

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	cmd      = "url.dll,FileProtocolHandler"
	runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
)

func OpenURL(url string) error {
	return WinOpen(url)
}

func OpenURLFork(url string) error {
	return OpenURL(url)
}

func WinOpen(input string) error {
	cmd := exec.Command(runDll32, cmd, input)
	return cmd.Run()
}