//go:build windows
// +build windows

package python

import (
	"fmt"
	"log"
	"os/exec"
)

// installPython installs the extracted Python installation file.
func installPython(pythonInstallerPath, pythonInstallDir string) {
	log.Println("installing python")

	installerArgs := []string{
		"/passive",
		"InstallAllUsers=0",
		"Include_launcher=0",
		"Include_test=0",
		"Include_pip=1",
		"Include_tools=1",
		"Include_exe=1",
		"Include_dev=1",
		"Include_tcltk=1",
		"PrependPath=1",
		"AssociateFiles=1",
		fmt.Sprintf("DefaultJustForMeTargetDir=%s", pythonInstallDir),
	}
	cmd := exec.Command(pythonInstallerPath, installerArgs...)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to install python: %v", err)
	}
}
