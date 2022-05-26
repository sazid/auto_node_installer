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

	tryInstall := func(args []string) error {
		cmd := exec.Command(pythonInstallerPath, args...)

		err := cmd.Run()
		return err
	}

	if err := tryInstall([]string{
		"/passive",
		"InstallAllUsers=0",
		"Include_launcher=1",
		"Include_test=0",
		"Include_pip=1",
		"Include_tools=1",
		"Include_exe=1",
		"Include_dev=1",
		"Include_tcltk=1",
		"PrependPath=1",
		"AssociateFiles=0",
		fmt.Sprintf("DefaultJustForMeTargetDir=%s", pythonInstallDir),
	}); err != nil {
		log.Println("failed to obtain admin rights, installing without .py file association")
		log.Printf("to associate .py files, right click on any .py file, select Open With and select 'python.exe' from this folder: %s", pythonInstallDir)
		if err := tryInstall([]string{
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
		}); err != nil {
			log.Fatalf("failed to install python: %v", err)
		}
	}
}
