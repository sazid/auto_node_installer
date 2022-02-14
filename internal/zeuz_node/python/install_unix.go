//go:build !windows
// +build !windows

package python

import (
	"log"
)

// installPython installs the extracted Python installation file.
func installPython(pythonInstallerPath, pythonInstallDir string) {
	log.Println("we're running on a *nix platform, python installation will be skipped.")
}
