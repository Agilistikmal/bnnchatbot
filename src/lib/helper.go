package lib

import (
	"fmt"
	"os"
	"os/exec"
)

func RestartApp() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	cmd := exec.Command(exePath)
	cmd.Start()

	os.Exit(0)
}
