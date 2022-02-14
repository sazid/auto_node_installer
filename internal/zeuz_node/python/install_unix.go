//go:build !windows
// +build !windows

package python

import (
	"log"
)

// installPython installs the extracted Python installation file.
func installPython(payloadDir, pythonInstallerFilename string) {
	log.Println("we're running on a platform other than, python installation will be skipped.")
}
